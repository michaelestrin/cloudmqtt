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

import "time"

type SendResultFunc func() bool

type SentInstance struct {
	When time.Time
	Data []byte
}

type Sender struct {
	SendCalledCount int
	Sent            []SentInstance
	sendResultFunc  SendResultFunc
}

func NewSenderImplWithResultFunc(resultFunc SendResultFunc) *Sender {
	return &Sender{
		SendCalledCount: 0,
		sendResultFunc:  resultFunc,
	}
}

func NewSenderImpl() *Sender {
	return NewSenderImplWithResultFunc(
		func() bool {
			return true
		})
}

func (s *Sender) Send(data []byte) bool {
	s.SendCalledCount++
	s.Sent = append(s.Sent, SentInstance{When: time.Now(), Data: data})
	return s.sendResultFunc()
}
