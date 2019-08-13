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

package impl

import (
	"context"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/michaelestrin/cloudmqtt/internal/cloudmqtt/contract"
)

// notify is a receiver wrapping a metadata query-and-forward implementation.
type notify struct {
	loggingClient  logger.LoggingClient
	send           contract.Sender
	marshal        contract.Marshaller
	metadataClient contract.MetadataClient
}

// NewNotifier is a constructor that returns an instance of notify configured to communicate with a
// EdgeX core-metadata instance.
func NewNotifier(
	loggingClient logger.LoggingClient,
	send contract.Sender,
	marshal contract.Marshaller,
	metadataClient contract.MetadataClient) *notify {

	return &notify{
		loggingClient:  loggingClient,
		send:           send,
		marshal:        marshal,
		metadataClient: metadataClient,
	}
}

// deviceCallFailedLogMessage function formats and returns the log message for when a device call fails.
func deviceCallFailedLogMessage(eventId string, errorMessage string) string {
	return fmt.Sprintf("device call failed for %s (%s)", eventId, errorMessage)
}

// marshalFailedLogMessage function formats and returns the log message for when an attempt to marshal a type fails.
func marshalFailedLogMessage(eventId string, errorMessage string) string {
	return fmt.Sprintf("marshal failed for %s (%s)", eventId, errorMessage)
}

// Notify method implements Notifier contract; it queries an EdgeX core-metadata instance for a specific device's
// metadata and forwards the result northbound.
func (n *notify) Notify(event *models.Event) bool {
	result, err := n.metadataClient.DeviceForName(event.Device, context.Background())
	if err != nil {
		n.loggingClient.Error(deviceCallFailedLogMessage(event.ID, err.Error()))
		return false
	}

	bytes, err := n.marshal(result)
	if err != nil {
		n.loggingClient.Error(marshalFailedLogMessage(event.ID, err.Error()))
		return false
	}

	return n.send(bytes)
}
