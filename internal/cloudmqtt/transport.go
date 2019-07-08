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

package cloudmqtt

import (
	"fmt"
	"github.com/dell-iot/cloudmqtt/internal/cloudmqtt/contract"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"sync"
	"time"
)

// transport is a receiver wrapping a generic event and metadata export adapter.
type transport struct {
	loggingClient                logger.LoggingClient
	sendFailureWaitInNanoseconds time.Duration
	send                         contract.Sender
	notify                       contract.Notifier
	marshal                      contract.Marshaller
	cleanUp                      contract.CleanUp
	wg                           sync.WaitGroup
	events                       chan *models.Event
}

// NewTransport is a constructor that returns a configured transport receiver whose Run() method can be included in
// a call to the EdgeX Applications Functions SDK's SetFunctionsPipeline() method.
func NewTransport(
	loggingClient logger.LoggingClient,
	sendFailureWaitInNanoseconds time.Duration,
	send contract.Sender,
	notify contract.Notifier,
	marshal contract.Marshaller,
	cleanUp contract.CleanUp) *transport {

	t := &transport{
		loggingClient:                loggingClient,
		sendFailureWaitInNanoseconds: sendFailureWaitInNanoseconds,
		send:                         send,
		notify:                       notify,
		marshal:                      marshal,
		cleanUp:                      cleanUp,
		events:                       make(chan *models.Event, 16),
	}
	t.wg.Add(1)
	go t.newDeviceHandler()
	return t
}

// detectedNewDeviceLogMessage function formats and returns the log message for when a new device is detected.
func detectedNewDeviceLogMessage(deviceName string) string {
	return fmt.Sprintf("detected new device %s", deviceName)
}

// newDeviceHandler method is executed as goroutine by constructor and is responsible for tracking known devices and
// calling Notifier implementation for new devices.
func (t *transport) newDeviceHandler() {
	defer t.wg.Done()

	devices := make(map[string]bool)
	for event := range t.events {
		_, ok := devices[event.Device]
		if !ok {
			if t.notify(event) {
				devices[event.Device] = true
				t.loggingClient.Debug(detectedNewDeviceLogMessage(event.Device))
			}
		}
	}
}

// newDeviceHandler function formats and returns the log message for when an attempt to marshal a type fails.
func marshalFailedLogMessage(eventId string, errorMessage string) string {
	return fmt.Sprintf("marshal failed for %s (%s)", eventId, errorMessage)
}

// sentLogMessage function formats and returns the log message for when an event has been successfully sent northbound.
func sentLogMessage(eventId string) string {
	return fmt.Sprintf("sent for %s", eventId)
}

// handleEvent method transmits an event northbound.
func (t *transport) handleEvent(EdgeXContext contract.EdgeXContext, event *models.Event) {
	bytes, err := t.marshal(event)
	if err != nil {
		t.loggingClient.Warn(marshalFailedLogMessage(event.ID, err.Error()))
		return
	}

	for {
		if t.send(bytes) {
			t.loggingClient.Debug(sentLogMessage(event.ID))

			err = EdgeXContext.MarkAsPushed()
			if err != nil {
				t.loggingClient.Error(err.Error())
			}

			return
		}
		time.Sleep(t.sendFailureWaitInNanoseconds)
	}
}

// run method is internal implementation delegated to by publicly accessible Run(); implemented to facilitate
// unit testing
func (t *transport) run(EdgeXContext contract.EdgeXContext, params ...interface{}) (bool, interface{}) {
	for _, param := range params {
		if event, ok := param.(models.Event); ok {
			t.events <- &event
			t.handleEvent(EdgeXContext, &event)
		}
	}
	return true, params
}

// Run method is an EdgeX Applications Function SDK-compatible function that can be included in a call to its
// SetFunctionsPipeline() method.
func (t *transport) Run(EdgeXContext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	return t.run(EdgeXContext, params[0])
}

// CleanUp method ensures the newDeviceHandler() goroutine has completed.
func (t *transport) CleanUp() {
	close(t.events)
	t.wg.Wait()
	t.cleanUp()
}
