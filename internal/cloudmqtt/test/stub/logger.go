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

type loggingClient struct {
	errors   []string
	warnings []string
	debugs   []string
}

func NewLoggerStub() *loggingClient {
	return &loggingClient{}
}

func (l *loggingClient) SetLogLevel(logLevel string) error {
	return nil
}

func (l *loggingClient) Debug(msg string, args ...interface{}) {
	l.debugs = append(l.debugs, msg)
}

func (l *loggingClient) Error(msg string, args ...interface{}) {
	l.errors = append(l.errors, msg)
}

func (l *loggingClient) Info(msg string, args ...interface{})  {}
func (l *loggingClient) Trace(msg string, args ...interface{}) {}

func (l *loggingClient) Warn(msg string, args ...interface{}) {
	l.warnings = append(l.warnings, msg)
}

func (l *loggingClient) occurred(msgs []string, expectedMessage string) bool {
	for _, msg := range msgs {
		if msg == expectedMessage {
			return true
		}
	}
	return false
}

func (l *loggingClient) SpecificDebugOccurred(expectedMessage string) bool {
	return l.occurred(l.debugs, expectedMessage)
}

func (l *loggingClient) ErrorsOccurred() bool {
	return len(l.errors) > 0
}

func (l *loggingClient) SpecificErrorOccurred(expectedMessage string) bool {
	return l.occurred(l.errors, expectedMessage)
}

func (l *loggingClient) SpecificWarningOccurred(expectedMessage string) bool {
	return l.occurred(l.warnings, expectedMessage)
}
