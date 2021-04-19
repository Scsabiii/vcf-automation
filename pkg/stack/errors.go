/******************************************************************************
*
*  Copyright 2021 SAP SE
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*
******************************************************************************/

package controller

import "errors"

var ErrStringEmpty = errors.New("string empty")
var ErrNotImplemented = errors.New("not implemented")
var ErrNotSupported = errors.New("not supported")
var ErrNodeExists = errors.New("node exists")
var ErrStackNotInitialized = errors.New("stack not initialized")
var ErrBackendURLNotSet = errors.New("env variable PULUMI_BACKEND_URL not set")
var ErrKeypairNotSet = errors.New("Config.Props.Keypair not set")
var ErrBadFormat = errors.New("bad format")
