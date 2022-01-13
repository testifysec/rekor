// Code generated by go-swagger; DO NOT EDIT.

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
//

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DsseV001Schema dsse v0.0.1 Schema
//
// Schema for dsse object
//
// swagger:model dsseV001Schema
type DsseV001Schema struct {

	// payload of the envelope
	// Required: true
	// Format: byte
	Payload *strfmt.Base64 `json:"payload"`

	// payload hash
	PayloadHash *DsseV001SchemaPayloadHash `json:"payloadHash,omitempty"`

	// type descriping the payload
	// Required: true
	PayloadType *string `json:"payloadType"`

	// collection of all signatures of the envelope's payload
	// Required: true
	// Min Items: 1
	Signatures []*DsseV001SchemaSignaturesItems0 `json:"signatures"`
}

// Validate validates this dsse v001 schema
func (m *DsseV001Schema) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePayload(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePayloadHash(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePayloadType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSignatures(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DsseV001Schema) validatePayload(formats strfmt.Registry) error {

	if err := validate.Required("payload", "body", m.Payload); err != nil {
		return err
	}

	return nil
}

func (m *DsseV001Schema) validatePayloadHash(formats strfmt.Registry) error {
	if swag.IsZero(m.PayloadHash) { // not required
		return nil
	}

	if m.PayloadHash != nil {
		if err := m.PayloadHash.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("payloadHash")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("payloadHash")
			}
			return err
		}
	}

	return nil
}

func (m *DsseV001Schema) validatePayloadType(formats strfmt.Registry) error {

	if err := validate.Required("payloadType", "body", m.PayloadType); err != nil {
		return err
	}

	return nil
}

func (m *DsseV001Schema) validateSignatures(formats strfmt.Registry) error {

	if err := validate.Required("signatures", "body", m.Signatures); err != nil {
		return err
	}

	iSignaturesSize := int64(len(m.Signatures))

	if err := validate.MinItems("signatures", "body", iSignaturesSize, 1); err != nil {
		return err
	}

	for i := 0; i < len(m.Signatures); i++ {
		if swag.IsZero(m.Signatures[i]) { // not required
			continue
		}

		if m.Signatures[i] != nil {
			if err := m.Signatures[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("signatures" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("signatures" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this dsse v001 schema based on the context it is used
func (m *DsseV001Schema) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePayloadHash(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSignatures(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DsseV001Schema) contextValidatePayloadHash(ctx context.Context, formats strfmt.Registry) error {

	if m.PayloadHash != nil {
		if err := m.PayloadHash.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("payloadHash")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("payloadHash")
			}
			return err
		}
	}

	return nil
}

func (m *DsseV001Schema) contextValidateSignatures(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Signatures); i++ {

		if m.Signatures[i] != nil {
			if err := m.Signatures[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("signatures" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("signatures" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *DsseV001Schema) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DsseV001Schema) UnmarshalBinary(b []byte) error {
	var res DsseV001Schema
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DsseV001SchemaPayloadHash hash of the envelope's payload after being PAE encoded
//
// swagger:model DsseV001SchemaPayloadHash
type DsseV001SchemaPayloadHash struct {

	// The hasing function used to compue the hash value
	// Enum: [sha256]
	Algorithm string `json:"algorithm,omitempty"`

	// The hash value of the PAE encoded payload
	Value string `json:"value,omitempty"`
}

// Validate validates this dsse v001 schema payload hash
func (m *DsseV001SchemaPayloadHash) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlgorithm(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var dsseV001SchemaPayloadHashTypeAlgorithmPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["sha256"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dsseV001SchemaPayloadHashTypeAlgorithmPropEnum = append(dsseV001SchemaPayloadHashTypeAlgorithmPropEnum, v)
	}
}

const (

	// DsseV001SchemaPayloadHashAlgorithmSha256 captures enum value "sha256"
	DsseV001SchemaPayloadHashAlgorithmSha256 string = "sha256"
)

// prop value enum
func (m *DsseV001SchemaPayloadHash) validateAlgorithmEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, dsseV001SchemaPayloadHashTypeAlgorithmPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *DsseV001SchemaPayloadHash) validateAlgorithm(formats strfmt.Registry) error {
	if swag.IsZero(m.Algorithm) { // not required
		return nil
	}

	// value enum
	if err := m.validateAlgorithmEnum("payloadHash"+"."+"algorithm", "body", m.Algorithm); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this dsse v001 schema payload hash based on the context it is used
func (m *DsseV001SchemaPayloadHash) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *DsseV001SchemaPayloadHash) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DsseV001SchemaPayloadHash) UnmarshalBinary(b []byte) error {
	var res DsseV001SchemaPayloadHash
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DsseV001SchemaSignaturesItems0 a signature of the envelope's payload along with the public key for the signature
//
// swagger:model DsseV001SchemaSignaturesItems0
type DsseV001SchemaSignaturesItems0 struct {

	// optional id of the key used to create the signature
	Keyid string `json:"keyid,omitempty"`

	// public key that corresponds to this signature
	// Read Only: true
	// Format: byte
	PublicKey strfmt.Base64 `json:"publicKey,omitempty"`

	// signature of the payload
	// Format: byte
	Sig strfmt.Base64 `json:"sig,omitempty"`
}

// Validate validates this dsse v001 schema signatures items0
func (m *DsseV001SchemaSignaturesItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this dsse v001 schema signatures items0 based on the context it is used
func (m *DsseV001SchemaSignaturesItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePublicKey(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DsseV001SchemaSignaturesItems0) contextValidatePublicKey(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "publicKey", "body", strfmt.Base64(m.PublicKey)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DsseV001SchemaSignaturesItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DsseV001SchemaSignaturesItems0) UnmarshalBinary(b []byte) error {
	var res DsseV001SchemaSignaturesItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
