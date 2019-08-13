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
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/google/uuid"
	"github.com/michaelestrin/cloudmqtt/internal/cloudmqtt/test/stub"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//  SUT factory
//

func newCommandHandlerSUT(loggingClient logger.LoggingClient) *commandHandler {
	return NewCommandHandler(loggingClient)
}

//
//  unit tests
//

func TestHandlerCallLogsDebug(t *testing.T) {
	loggingClient := stub.NewLoggerStub()
	command := uuid.New().String()
	sut := newCommandHandlerSUT(loggingClient)

	sut.Receiver(command)

	assert.True(t, loggingClient.SpecificDebugOccurred(receivedCommandLogMessage(command)))
}
