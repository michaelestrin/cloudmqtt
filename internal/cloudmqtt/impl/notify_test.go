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
	"encoding/json"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/google/uuid"
	"github.com/michaelestrin/cloudmqtt/internal/cloudmqtt/contract"
	"github.com/michaelestrin/cloudmqtt/internal/cloudmqtt/test/helper"
	"github.com/michaelestrin/cloudmqtt/internal/cloudmqtt/test/stub"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//  test stubs
//

type deviceForNameCalledInstance struct {
	Name    string
	Context context.Context
}

type deviceForNameResult struct {
	Device models.Device
	Err    error
}

type metadataClientImpl struct {
	DeviceForNameCalledCount     int
	DeviceForNameCalledInstances []deviceForNameCalledInstance
	DeviceForNameResult          deviceForNameResult
}

func newMetadataClientImpl(device models.Device, err error) *metadataClientImpl {
	return &metadataClientImpl{
		DeviceForNameCalledCount: 0,
		DeviceForNameResult: deviceForNameResult{
			Device: device,
			Err:    err,
		},
	}
}

func (c *metadataClientImpl) DeviceForName(name string, ctx context.Context) (models.Device, error) {
	c.DeviceForNameCalledCount++
	c.DeviceForNameCalledInstances = append(
		c.DeviceForNameCalledInstances,
		deviceForNameCalledInstance{
			Name:    name,
			Context: ctx,
		})
	return c.DeviceForNameResult.Device, c.DeviceForNameResult.Err
}

func newDevice(deviceName string) models.Device {
	return models.Device{
		Id:   uuid.New().String(),
		Name: deviceName,
	}

}

func newMetadataClientImplReturnFailure(errorMessage string) contract.MetadataClient {
	return newMetadataClientImpl(newDevice("device"), errors.New(errorMessage))
}

func newMetadataClientImplReturnSuccess() contract.MetadataClient {
	return newMetadataClientImpl(newDevice("device"), nil)
}

//
//  SUT factory
//

func newNotifierSUT(
	loggingClient logger.LoggingClient,
	sender contract.Sender,
	marshal contract.Marshaller,
	metadataClient contract.MetadataClient) *notify {

	return NewNotifier(loggingClient, sender, marshal, metadataClient)
}

//
//  unit tests
//

func TestNotifyCallToMetadataClientFailureReturnsFalse(t *testing.T) {
	sut := newNotifierSUT(
		stub.NewLoggerStub(),
		stub.NewSenderImpl().Send,
		json.Marshal,
		newMetadataClientImplReturnFailure(uuid.New().String()))
	event := stub.NewEvent()

	result := sut.Notify(&event)

	assert.False(t, result)
}

func TestNotifyCallToMetadataClientFailureLogsError(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	errorMessage := uuid.New().String()
	sut := newNotifierSUT(
		loggingClient,
		stub.NewSenderImpl().Send,
		json.Marshal,
		newMetadataClientImplReturnFailure(errorMessage))
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.True(t, loggingClient.SpecificErrorOccurred(deviceCallFailedLogMessage(event.ID, errorMessage)))
}

func TestSendNotCalledWhenCallToMetadataClientFails(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newNotifierSUT(
		stub.NewLoggerStub(),
		sender.Send,
		json.Marshal,
		newMetadataClientImplReturnFailure(uuid.New().String()))
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.Equal(t, 0, sender.SendCalledCount)
}

func TestNotifyCallToMetadataClientSuccessDoesNotLogError(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newNotifierSUT(loggingClient, stub.NewSenderImpl().Send, json.Marshal, newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.False(t, loggingClient.ErrorsOccurred())
}

func TestNotifyCallToMetadataClientSuccessReturnsTrue(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newNotifierSUT(loggingClient, stub.NewSenderImpl().Send, json.Marshal, newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	result := sut.Notify(&event)

	assert.True(t, result)
}

func TestNotifyCallToMarshalFailureReturnsFalse(t *testing.T) {
	sut := newNotifierSUT(
		stub.NewLoggerStub(),
		stub.NewSenderImpl().Send,
		helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall(),
		newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	result := sut.Notify(&event)

	assert.False(t, result)
}

func TestNotifyCallToMarshalFailureLogsError(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newNotifierSUT(
		loggingClient,
		stub.NewSenderImpl().Send,
		helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall(),
		newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.True(t, loggingClient.SpecificErrorOccurred(marshalFailedLogMessage(event.ID, helper.JsonMarshalFuncFailureMessage)))
}

func TestSendNotCalledWhenCallToMarshalFails(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newNotifierSUT(
		stub.NewLoggerStub(),
		sender.Send,
		helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall(),
		newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.Equal(t, 0, sender.SendCalledCount)
}

func TestNotifyCallToMarshalSuccessDoesNotLogError(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newNotifierSUT(loggingClient, stub.NewSenderImpl().Send, json.Marshal, newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.False(t, loggingClient.ErrorsOccurred())
}

func TestNotifyCallToMarshalSuccessReturnsTrue(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newNotifierSUT(loggingClient, stub.NewSenderImpl().Send, json.Marshal, newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	result := sut.Notify(&event)

	assert.True(t, result)
}

func TestNotifyCallsSenderOnce(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newNotifierSUT(stub.NewLoggerStub(), sender.Send, json.Marshal, newMetadataClientImplReturnSuccess())
	event := stub.NewEvent()

	sut.Notify(&event)

	assert.Equal(t, 1, sender.SendCalledCount)
}
