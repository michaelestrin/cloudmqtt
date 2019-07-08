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

package contract

import (
	"context"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

// Sender defines function contract for transmitting bytes to Cloud.
type Sender func(data []byte) bool

// Notifier defines function contract for notifying Cloud of newly added device's metadata.
type Notifier func(event *models.Event) bool

// Receiver defines function contract for handling southbound command received from Cloud.
type Receiver func(command string)

// Marshaller defines function contract for marshalling type to []byte; supports unit testing.
type Marshaller func(v interface{}) ([]byte, error)

// CleanUp defines function contract for end of execution clean up activity.
type CleanUp func()

// MetadataClient defines interface for interacting with EdgeX core-metadata service; defined to facilitate
// unit testing.
type MetadataClient interface {
	// DeviceForName loads the device for the specified name
	DeviceForName(name string, ctx context.Context) (models.Device, error)
}

// EdgeXContext defines interface for interacting with Applications Functions SDK's edgexcontext; defined to facilitate
// unit testing.
type EdgeXContext interface {
	MarkAsPushed() error
}
