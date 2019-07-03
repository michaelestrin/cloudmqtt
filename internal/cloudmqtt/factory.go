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
	"fmt"
	"github.com/dell-iot/cloudmqtt/internal/cloudmqtt/impl"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/types"
	"os"
	"time"
)

// setting function translates setting's key to value (or logs and exits if the requested key does not exist).
func setting(loggingClient logger.LoggingClient, settings map[string]string, key string) (value string) {
	value, ok := settings[key]
	if !ok {
		loggingClient.Error(fmt.Sprintf("main.setting missing setting: %s", key))
		os.Exit(-1)
	}
	return
}

// FactoryTransport returns a function that can be called by the EdgeX Applications Functions SDK.
func FactoryTransport(sdk *appsdk.AppFunctionsSDK) func(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	settings := sdk.ApplicationSettings()

	mqtt := impl.NewMqttInstanceForCloud(
		sdk.LoggingClient,
		setting(sdk.LoggingClient, settings, "certFile"),
		setting(sdk.LoggingClient, settings, "keyFile"),
		setting(sdk.LoggingClient, settings, "clientId"),
		setting(sdk.LoggingClient, settings, "userName"),
		setting(sdk.LoggingClient, settings, "password"),
		setting(sdk.LoggingClient, settings, "server"),
		setting(sdk.LoggingClient, settings, "eventTopic"),
		setting(sdk.LoggingClient, settings, "newDeviceTopic"))

	marshaller := json.Marshal

	metadataClient := metadata.NewDeviceClient(
		types.EndpointParams{
			ServiceKey:  clients.CoreMetaDataServiceKey,
			Path:        clients.ApiDeviceRoute,
			UseRegistry: false,
			Url:         setting(sdk.LoggingClient, settings, "edgeXMetaDataUri") + clients.ApiDeviceRoute,
			Interval:    clients.ClientMonitorDefault,
		},
		nil)

	notifier := impl.NewNotifier(sdk.LoggingClient, mqtt.NewDeviceSender, marshaller, metadataClient)

	return NewTransport(sdk.LoggingClient, 1*time.Second, mqtt.EventSender, notifier.Notify, marshaller).Run
}
