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
	"encoding/json"
	"github.com/dell-iot/cloudmqtt/internal/cloudmqtt/contract"
	"github.com/dell-iot/cloudmqtt/internal/cloudmqtt/test/helper"
	"github.com/dell-iot/cloudmqtt/internal/cloudmqtt/test/stub"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const sendFailureWaitInNanosecondsForTesting = 50000000

//
//  test stubs
//

type notifierImpl struct {
	NotifyCalledCount int
	Notified          []models.Event
	notifyResult      bool
}

func newNotifierImplWithSpecificResult(notifyResult bool) *notifierImpl {
	return &notifierImpl{
		NotifyCalledCount: 0,
		notifyResult:      notifyResult,
	}
}

func newNotifierImpl() *notifierImpl {
	return newNotifierImplWithSpecificResult(true)
}

func (n *notifierImpl) notify(event *models.Event) bool {
	n.NotifyCalledCount++
	n.Notified = append(n.Notified, *event)
	return n.notifyResult
}

type edgeXContextImpl struct {
	MarkAsPushedCalledCount int
	markAsPushedResult      error
}

func newEdgeXContextImplWithSpecificResult(result error) *edgeXContextImpl {
	return &edgeXContextImpl{
		MarkAsPushedCalledCount: 0,
		markAsPushedResult:      result,
	}
}

func newEdgeXContextImpl() *edgeXContextImpl {
	return newEdgeXContextImplWithSpecificResult(nil)
}

func (e *edgeXContextImpl) MarkAsPushed() error {
	e.MarkAsPushedCalledCount++
	return e.markAsPushedResult
}

//
//  SUT factory
//

func newSUT(
	loggingClient logger.LoggingClient,
	sender contract.Sender,
	notifier contract.Notifier,
	marshal contract.Marshaller) *transport {

	return NewTransport(
		loggingClient,
		sendFailureWaitInNanosecondsForTesting,
		sender,
		notifier,
		marshal)
}

//
//  utility and helper functions
//

func assertResultsEqualExpectedEvents(t *testing.T, expected []models.Event, actual interface{}) {
	assert.Len(t, actual, len(expected))

	results, ok := actual.([]interface{})
	assert.True(t, ok)
	for index, result := range results {
		event, ok := result.(models.Event)
		assert.True(t, ok)

		assert.Equal(t, expected[index], event)
	}

}

func marshalAndUnmarshalEventsToEnsureRelatedMarshallingCodeIsExecuted(
	t *testing.T,
	events []models.Event) (result []models.Event) {

	for _, event := range events {
		jsonEvent, err := json.Marshal(event)
		if err == nil {
			var unmarshalledEvent models.Event
			err = json.Unmarshal(jsonEvent, &unmarshalledEvent)
			if err == nil {
				result = append(result, unmarshalledEvent)
				continue
			}
		}
		assert.Fail(t, "marshalAndUnmarshalEventsToEnsureRelatedMarshallingCodeIsExecuted failed.")
	}
	return
}

func unmarshalEventsForAssert(t *testing.T, eventsJson []stub.SentInstance) (events []interface{}) {
	for _, eventJson := range eventsJson {
		var event models.Event
		err := json.Unmarshal(eventJson.Data, &event)
		assert.Nil(t, err)
		events = append(events, event)
	}
	return
}

func factorySendResultFalseOnceThenTrueFromThenOn() stub.SendResultFunc {
	callCount := 0
	return func() bool {
		callCount++
		return callCount != 1
	}
}

//
//  unit tests
//

func TestCallTransportWithNoParametersReturnsTrueForContinuePipeline(t *testing.T) {
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)

	continuePipeline, _ := sut.run(newEdgeXContextImpl())
	sut.waitForCompletion()

	assert.True(t, continuePipeline)
}

func TestCallTransportWithEventParameterReturnsPassedEventParameter(t *testing.T) {
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)
	event := stub.NewEvent()

	_, results := sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assertResultsEqualExpectedEvents(t, []models.Event{event}, results)
}

func TestCallTransportWithEventParametersReturnsPassedEventParameters(t *testing.T) {
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)
	event1 := stub.NewEvent()
	event2 := stub.NewEvent()

	_, results := sut.run(newEdgeXContextImpl(), event1, event2)
	sut.waitForCompletion()

	assertResultsEqualExpectedEvents(t, []models.Event{event1, event2}, results)
}

func TestCallWithEventParameterCallsSenderOnce(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)

	sut.run(newEdgeXContextImpl(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 1, sender.SendCalledCount)
}

func TestCallWithEventParameterCallsMarkedAsPushedOnce(t *testing.T) {
	edgeXContext := newEdgeXContextImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)

	sut.run(edgeXContext, stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 1, edgeXContext.MarkAsPushedCalledCount)
}

func TestCallWithEventParameterPassesParameterToSender(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assertResultsEqualExpectedEvents(
		t,
		marshalAndUnmarshalEventsToEnsureRelatedMarshallingCodeIsExecuted(t, []models.Event{event}),
		unmarshalEventsForAssert(t, sender.Sent))
}

func TestCallWithEventParametersCallsSenderOnceForEachParameter(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)

	sut.run(newEdgeXContextImpl(), stub.NewEvent(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 2, sender.SendCalledCount)
}

func TestCallWithEventParametersCallsMarkedAsPushedOnceForEachParameter(t *testing.T) {
	edgeXContext := newEdgeXContextImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)

	sut.run(edgeXContext, stub.NewEvent(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 2, edgeXContext.MarkAsPushedCalledCount)
}

func TestCallWithEventParametersPassesParameterToSenderForEachParameter(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)
	event1 := stub.NewEvent()
	event2 := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event1, event2)
	sut.waitForCompletion()

	assertResultsEqualExpectedEvents(
		t,
		marshalAndUnmarshalEventsToEnsureRelatedMarshallingCodeIsExecuted(t, []models.Event{event1, event2}),
		unmarshalEventsForAssert(t, sender.Sent))
}

func TestMarkAsPushedFailureLogsError(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newSUT(loggingClient, stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)
	errorMessage := uuid.New().String()
	edgeXContext := newEdgeXContextImplWithSpecificResult(errors.New(errorMessage))

	sut.run(edgeXContext, stub.NewEvent())
	sut.waitForCompletion()

	assert.True(t, loggingClient.SpecificErrorOccurred(errorMessage))
}

func TestMarshalFailureLogsWarning(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newSUT(loggingClient, stub.NewSenderImpl().Send, newNotifierImpl().notify, helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall())
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.True(t, loggingClient.SpecificWarningOccurred(marshalFailedLogMessage(event.ID, helper.JsonMarshalFuncFailureMessage)))
}

func TestMarshalSuccessDoesNotLogWarning(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newSUT(loggingClient, stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.False(t, loggingClient.SpecificWarningOccurred(marshalFailedLogMessage(event.ID, helper.JsonMarshalFuncFailureMessage)))
}

func TestMarshalFailureDoesNotCallSender(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall())

	sut.run(newEdgeXContextImpl(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 0, sender.SendCalledCount)
}

func TestMarshalFailureOnlyAffectsFailedParameter(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall())

	sut.run(newEdgeXContextImpl(), stub.NewEvent(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 1, sender.SendCalledCount)
}

func TestCallWithEventParametersDoesNotPassParameterToSenderWhenMarshalFails(t *testing.T) {
	sender := stub.NewSenderImpl()
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, helper.FactoryJsonMarshalFuncReturnsFailureOnFirstCall())
	event2 := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), stub.NewEvent(), event2)
	sut.waitForCompletion()

	assertResultsEqualExpectedEvents(
		t,
		marshalAndUnmarshalEventsToEnsureRelatedMarshallingCodeIsExecuted(t, []models.Event{event2}),
		unmarshalEventsForAssert(t, sender.Sent))
}

func TestSenderFailureResultsInRetry(t *testing.T) {
	sender := stub.NewSenderImplWithResultFunc(factorySendResultFalseOnceThenTrueFromThenOn())
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)

	sut.run(newEdgeXContextImpl(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 2, sender.SendCalledCount)
}

func TestSenderFailureRetryOnlyCallsMarkAsPushedOnce(t *testing.T) {
	sender := stub.NewSenderImplWithResultFunc(factorySendResultFalseOnceThenTrueFromThenOn())
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)
	edgeXContext := newEdgeXContextImpl()

	sut.run(edgeXContext, stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 1, edgeXContext.MarkAsPushedCalledCount)
}

func TestSenderFailureResultsInDelayBeforeRetry(t *testing.T) {
	sender := stub.NewSenderImplWithResultFunc(factorySendResultFalseOnceThenTrueFromThenOn())
	sut := newSUT(stub.NewLoggerStub(), sender.Send, newNotifierImpl().notify, json.Marshal)

	sut.run(newEdgeXContextImpl(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Len(t, sender.Sent, 2)
	assert.True(t, sender.Sent[1].When.UnixNano()-sender.Sent[0].When.UnixNano() > sendFailureWaitInNanosecondsForTesting)
}

func TestIfSenderCalledThenDebugLogged(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sender := stub.NewSenderImpl()
	sut := newSUT(loggingClient, sender.Send, newNotifierImpl().notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.Equal(t, 1, sender.SendCalledCount)
	assert.True(t, loggingClient.SpecificDebugOccurred(sentLogMessage(event.ID)))
}

func TestCallWithEventForNewDeviceCallsNotifier(t *testing.T) {
	notifier := newNotifierImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.Equal(t, 1, notifier.NotifyCalledCount)
}

func TestCallWithEventForNewDeviceCallsNotifierWithPassedEventParameter(t *testing.T) {
	notifier := newNotifierImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.Len(t, notifier.Notified, 1)
	assert.Equal(t, event, notifier.Notified[0])
}

func TestCallWithEventsForNewDeviceCallsNotifier(t *testing.T) {
	notifier := newNotifierImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)

	sut.run(newEdgeXContextImpl(), stub.NewEvent(), stub.NewEvent())
	sut.waitForCompletion()

	assert.Equal(t, 1, notifier.NotifyCalledCount)
}

func TestCallWithEventsForNewDevicesCallsNotifier(t *testing.T) {
	notifier := newNotifierImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)

	sut.run(newEdgeXContextImpl(), stub.NewEventForDevice("device1"), stub.NewEventForDevice("device2"))
	sut.waitForCompletion()

	assert.Equal(t, 2, notifier.NotifyCalledCount)
}

func TestCallWithEventsForNewDevicesCallsNotifierWithPassedEventParameters(t *testing.T) {
	notifier := newNotifierImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)
	event1 := stub.NewEventForDevice("device1")
	event2 := stub.NewEventForDevice("device2")

	sut.run(newEdgeXContextImpl(), event1, event2)
	sut.waitForCompletion()

	assert.Len(t, notifier.Notified, 2)
	assert.Equal(t, event1, notifier.Notified[0])
	assert.Equal(t, event2, notifier.Notified[1])
}

func TestNotifierSuccessLogsDebug(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	sut := newSUT(loggingClient, stub.NewSenderImpl().Send, newNotifierImpl().notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.True(t, loggingClient.SpecificDebugOccurred(detectedNewDeviceLogMessage(event.Device)))
}

func TestNotifierFailureDoesNotCauseLoggedDebug(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	notifier := newNotifierImplWithSpecificResult(false)
	sut := newSUT(loggingClient, stub.NewSenderImpl().Send, notifier.notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event)
	sut.waitForCompletion()

	assert.False(t, loggingClient.SpecificDebugOccurred(detectedNewDeviceLogMessage(event.Device)))
}

func TestNotifierSuccessDoesNotCallNotifierAgainForSameDevice(t *testing.T) {
	notifier := newNotifierImpl()
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event, event)
	sut.waitForCompletion()

	assert.Len(t, notifier.Notified, 1)
	assert.Equal(t, event, notifier.Notified[0])
}

func TestNotifierFailureCallsNotifierAgainForSameDevice(t *testing.T) {
	notifier := newNotifierImplWithSpecificResult(false)
	sut := newSUT(stub.NewLoggerStub(), stub.NewSenderImpl().Send, notifier.notify, json.Marshal)
	event := stub.NewEvent()

	sut.run(newEdgeXContextImpl(), event, event)
	sut.waitForCompletion()

	assert.Equal(t, 2, notifier.NotifyCalledCount)
}
