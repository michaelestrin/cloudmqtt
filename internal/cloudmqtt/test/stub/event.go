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

package stub

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/google/uuid"
)

// NewEventForDevice returns a stub models.Event struct with a unique ID and the specified deviceName
func NewEventForDevice(deviceName string) models.Event {
	return models.Event{
		ID:     uuid.New().String(),
		Device: deviceName,
	}
}

// NewEvent returns an stub models.Event struct with a unique ID
func NewEvent() models.Event {
	return NewEventForDevice("device")
}
