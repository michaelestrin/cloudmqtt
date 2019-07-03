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
	"crypto/tls"
	"fmt"
	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"os"
	"time"
)

const qosAtLeastOnce = 1

// mqtt is a receiver wrapping a one-way MQTTS implementation.
type mqtt struct {
	loggingClient  logger.LoggingClient
	client         mqttlib.Client
	eventTopic     string
	newDeviceTopic string
}

// NewMqttInstanceForCloud is a constructor that returns an mqtt receiver configured for cloud-based MQTTS.
func NewMqttInstanceForCloud(
	loggingClient logger.LoggingClient,
	certFile string,
	keyFile string,
	clientId string,
	userName string,
	password string,
	server string,
	eventTopic string,
	newDeviceTopic string) (q *mqtt) {

	q = &mqtt{
		loggingClient:  loggingClient,
		eventTopic:     eventTopic,
		newDeviceTopic: newDeviceTopic,
	}

	tlsConfig := &tls.Config{}
	if len(certFile) > 0 && len(keyFile) > 0 {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			loggingClient.Error(fmt.Sprintf("mqtt.mqttInstanceForCloud LoadX509KeyPair failed: %v", err))
			os.Exit(-1)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	} else {
		tlsConfig.ClientAuth = tls.NoClientCert
		tlsConfig.ClientCAs = nil
	}

	options := mqttlib.ClientOptions{
		ClientID:             clientId,
		Username:             userName,
		Password:             password,
		CleanSession:         true,
		AutoReconnect:        true,
		MaxReconnectInterval: 1 * time.Second,
		KeepAlive:            int64(30 * time.Second),
		TLSConfig:            tlsConfig,
	}
	options.AddBroker(server)
	q.client = mqttlib.NewClient(&options)

	if token := q.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return
}

// send function publishes content on designated northbound MQTT topic.
func send(q *mqtt, topicName string, content []byte) bool {
	if token := q.client.Publish(topicName, qosAtLeastOnce, false, content); token.Wait() && token.Error() != nil {
		q.loggingClient.Warn("mqtt send to " + topicName + " failed (" + token.Error().Error() + ")")
		return false
	}
	return true
}

// EventSender method transmits content to northbound MQTT event topic.
func (q *mqtt) EventSender(content []byte) bool {
	return send(q, q.eventTopic, content)
}

// NewDeviceSender method transmits content to northbound MQTT new device topic.
func (q *mqtt) NewDeviceSender(content []byte) bool {
	return send(q, q.newDeviceTopic, content)
}
