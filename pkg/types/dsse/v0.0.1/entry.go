//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dsse

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"

	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/secure-systems-lab/go-securesystemslib/dsse"
	"github.com/spf13/viper"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/sigstore/rekor/pkg/generated/models"
	"github.com/sigstore/rekor/pkg/log"
	"github.com/sigstore/rekor/pkg/pki/x509"
	"github.com/sigstore/rekor/pkg/types"
	rekordsse "github.com/sigstore/rekor/pkg/types/dsse"
	"github.com/sigstore/sigstore/pkg/signature"
)

const (
	APIVERSION = "0.0.1"
)

func init() {
	if err := rekordsse.VersionMap.SetEntryFactory(APIVERSION, NewEntry); err != nil {
		log.Logger.Panic(err)
	}
}

type V001Entry struct {
	DsseObj models.DsseV001Schema
}

func (v V001Entry) APIVersion() string {
	return APIVERSION
}

func NewEntry() types.EntryImpl {
	return &V001Entry{}
}

func (v V001Entry) IndexKeys() ([]string, error) {
	var result []string
	payloadBytes := v.DsseObj.Payload
	payloadType := *v.DsseObj.PayloadType
	h := sha256.Sum256(payloadBytes)
	payloadKey := "sha256:" + hex.EncodeToString(h[:])
	result = append(result, payloadKey)
	result = append(result, v.DsseObj.PayloadHash.Algorithm+":"+v.DsseObj.PayloadHash.Value)

	for _, sig := range v.DsseObj.Signatures {
		// decode the key to a string, then hash it
		keyHash := sha256.Sum256(sig.PublicKey)
		result = append(result, "sha256:"+hex.EncodeToString(keyHash[:]))
	}

	switch payloadType {
	case in_toto.PayloadType:
		statement, err := parseIntotoStatement(payloadBytes)
		if err != nil {
			return result, err
		}

		for _, s := range statement.Subject {
			for alg, ds := range s.Digest {
				result = append(result, alg+":"+ds)
			}
		}
	default:
		log.Logger.Infof("Cannot index payload of type: %s", payloadType)
	}

	return result, nil
}

func parseIntotoStatement(p []byte) (*in_toto.Statement, error) {
	ps := in_toto.Statement{}
	if err := json.Unmarshal(p, &ps); err != nil {
		return nil, err
	}

	return &ps, nil
}

func (v *V001Entry) Unmarshal(pe models.ProposedEntry) error {
	dsseModel, ok := pe.(*models.Dsse)
	if !ok {
		return errors.New("cannot unmarshal non DSSE v0.0.1 type")
	}

	if err := types.DecodeEntry(dsseModel.Spec, &v.DsseObj); err != nil {
		return err
	}

	// field validation
	if err := v.DsseObj.Validate(strfmt.Default); err != nil {
		return err
	}

	env := &dsse.Envelope{}
	// this weird juggling is because dsse.Envelope expects env.Payload to be a base64 string,
	// while v.DsseObj.Payload is a base64 byte array... but casting it to a string decodes it.
	// so... cast it to a string to decode it, and then re-encode it as a base64 encoded string....
	payload := string(v.DsseObj.Payload)
	env.Payload = base64.StdEncoding.EncodeToString([]byte(payload))
	env.PayloadType = *v.DsseObj.PayloadType
	allPubKeyBytes := make([][]byte, 0)
	for _, sig := range v.DsseObj.Signatures {
		env.Signatures = append(env.Signatures, dsse.Signature{
			KeyID: sig.Keyid,
			Sig:   string(sig.Sig),
		})

		allPubKeyBytes = append(allPubKeyBytes, sig.PublicKey)
	}

	_, err := verifyEnvelope(allPubKeyBytes, env)
	if err != nil {
		return fmt.Errorf("could not verify envelope: %w", err)
	}

	return nil
}

func (v *V001Entry) Canonicalize(ctx context.Context) ([]byte, error) {
	canonicalEntry := models.DsseV001Schema{
		PayloadHash: v.DsseObj.PayloadHash,
		Signatures:  v.DsseObj.Signatures,
		PayloadType: v.DsseObj.PayloadType,
	}

	model := models.Dsse{}
	model.APIVersion = swag.String(APIVERSION)
	model.Spec = canonicalEntry
	return json.Marshal(&model)
}

func (v *V001Entry) Attestation() []byte {
	payload := v.DsseObj.Payload
	if len(payload) > viper.GetInt("max_attestation_size") {
		log.Logger.Infof("Skipping attestation storage, size %d is greater than max %d", len(payload), viper.GetInt("max_attestation_size"))
		return nil
	}

	return payload
}

type verifier struct {
	v        signature.Verifier
	pub      crypto.PublicKey
	keyBytes []byte
	id       string
}

func (v *verifier) KeyID() (string, error) {
	return v.id, nil
}

func (v *verifier) Public() crypto.PublicKey {
	return v.pub
}

func (v *verifier) Verify(data, sig []byte) error {
	if v.v == nil {
		return errors.New("nil verifier")
	}
	return v.v.VerifySignature(bytes.NewReader(sig), bytes.NewReader(data))
}

func (v V001Entry) CreateFromArtifactProperties(_ context.Context, props types.ArtifactProperties) (models.ProposedEntry, error) {
	returnVal := models.Dsse{}
	re := V001Entry{}

	var err error
	artifactBytes := props.ArtifactBytes
	if artifactBytes == nil {
		if props.ArtifactPath == nil {
			return nil, errors.New("path to artifact file must be specified")
		}
		if props.ArtifactPath.IsAbs() {
			return nil, errors.New("dsse envelopes cannot be fetched over HTTP(S)")
		}
		artifactBytes, err = ioutil.ReadFile(filepath.Clean(props.ArtifactPath.Path))
		if err != nil {
			return nil, err
		}
	}

	env := dsse.Envelope{}
	if err := json.Unmarshal(artifactBytes, &env); err != nil {
		return nil, fmt.Errorf("payload must be a valid dsse envelope: %w", err)
	}

	allPubKeyBytes := make([][]byte, 0)
	if props.PublicKeyBytes != nil {
		allPubKeyBytes = append(allPubKeyBytes, props.PublicKeyBytes)
	}

	allPubKeyBytes = append(allPubKeyBytes, props.PublicKeysBytes...)
	allPubKeyPaths := make([]*url.URL, 0)
	if props.PublicKeyPath != nil {
		allPubKeyPaths = append(allPubKeyPaths, props.PublicKeysPaths...)
	}

	for _, path := range allPubKeyPaths {
		if path.IsAbs() {
			return nil, errors.New("dsse public keys cannot be fetched over HTTP(S)")
		}

		publicKeyBytes, err := ioutil.ReadFile(filepath.Clean(path.Path))
		if err != nil {
			return nil, fmt.Errorf("error reading public key file: %w", err)
		}

		allPubKeyBytes = append(allPubKeyBytes, publicKeyBytes)
	}

	keysBySig, err := verifyEnvelope(allPubKeyBytes, &env)
	if err != nil {
		return nil, err
	}

	decodedPayload, err := base64.StdEncoding.DecodeString(env.Payload)
	if err != nil {
		return nil, fmt.Errorf("could not decode envelope payload: %w", err)
	}

	paeEncodedPayload := dsse.PAE(env.PayloadType, decodedPayload)
	h := sha256.Sum256(paeEncodedPayload)
	re.DsseObj.Payload = decodedPayload
	re.DsseObj.PayloadType = &env.PayloadType
	re.DsseObj.PayloadHash = &models.DsseV001SchemaPayloadHash{
		Algorithm: models.DsseV001SchemaPayloadHashAlgorithmSha256,
		Value:     hex.EncodeToString(h[:]),
	}

	for _, sig := range env.Signatures {
		key, ok := keysBySig[sig.Sig]
		if !ok {
			return nil, errors.New("all signatures must have a key that verifies it")
		}

		canonKey, err := key.CanonicalValue()
		if err != nil {
			return nil, fmt.Errorf("could not canonicize key: %w", err)
		}

		keyBytes := strfmt.Base64(canonKey)
		sigBytes := strfmt.Base64([]byte(sig.Sig))
		re.DsseObj.Signatures = append(re.DsseObj.Signatures, &models.DsseV001SchemaSignaturesItems0{
			Keyid:     sig.KeyID,
			Sig:       sigBytes,
			PublicKey: keyBytes,
		})
	}

	returnVal.APIVersion = swag.String(re.APIVersion())
	returnVal.Spec = re.DsseObj
	return &returnVal, nil
}

// verifyEnvelope takes in an array of possible key bytes and attempts to parse them as x509 public keys.
// it then uses these to verify the envelope and makes sure that every signature on the envelope is verified.
// it returns a map of verifiers indexed by the signature the verifier corresponds to.
func verifyEnvelope(allPubKeyBytes [][]byte, env *dsse.Envelope) (map[string]*x509.PublicKey, error) {
	// generate a fake id for these keys so we can get back to the key bytes and match them to their corresponding signature
	verifierBySig := make(map[string]*x509.PublicKey)

	for _, pubKeyBytes := range allPubKeyBytes {
		key, err := x509.NewPublicKey(bytes.NewReader(pubKeyBytes))
		if err != nil {
			return nil, fmt.Errorf("could not parse public key as x509: %w", err)
		}

		vfr, err := signature.LoadVerifier(key.CryptoPubKey(), crypto.SHA256)
		if err != nil {
			return nil, fmt.Errorf("could not load verifier: %w", err)
		}

		dsseVfr, err := dsse.NewEnvelopeVerifier(&verifier{
			v:        vfr,
			keyBytes: pubKeyBytes,
		})

		if err != nil {
			return nil, fmt.Errorf("could not use public key as a dsse verifier: %w", err)
		}

		accepted, err := dsseVfr.Verify(env)
		if err != nil {
			return nil, fmt.Errorf("could not verify envelope: %w", err)
		}

		for _, accept := range accepted {
			verifierBySig[accept.Sig.Sig] = key
		}
	}

	return verifierBySig, nil
}
