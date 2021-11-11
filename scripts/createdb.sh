#!/bin/bash
#
# Copyright 2021 The Sigstore Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


DB="rekor"
USER="user"
PASS="password"
ROOTPASS="admin"

echo -e "Creating $DB database and $USER user account"


mysql <<MYSQL_SCRIPT
DROP DATABASE IF EXISTS $DB;
CREATE DATABASE $DB;
CREATE USER IF NOT EXISTS 'user;@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON rekor.* TO 'user'@'%';
FLUSH PRIVILEGES;
MYSQL_SCRIPT


echo -e "Loading table data.."

mysql -u $USER -p$PASS -D $DB < ./storage.sql
