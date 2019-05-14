/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package helper

import (
	"encoding/json"
	"errors"
)

const JsonMarshalFuncFailureMessage = "jsonMarshalFuncReturnsFailure"

// FactoryJsonMarshalFuncReturnsFailureOnFirstCall function returns a function that conforms to Marshaller contract.
// It returns an error on the first call and the results from json.Marshal() on subsequent calls.
func FactoryJsonMarshalFuncReturnsFailureOnFirstCall() func(v interface{}) ([]byte, error) {
	callCount := 0
	return func(v interface{}) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return nil, errors.New(JsonMarshalFuncFailureMessage)
		}
		return json.Marshal(v)
	}
}
