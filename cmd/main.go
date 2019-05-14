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

package main

import (
	"fmt"
	"github.com/dell-iot/cloudmqtt/internal/cloudmqtt"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
)

// main leverages the EdgeX Applications Functions SDK to create one-way export of device metadata and readings to
// Milky Way via MQTTS
func main() {
	sdk := &appsdk.AppFunctionsSDK{ServiceKey: "MilkyWay"}
	if err := sdk.Initialize(); err != nil {
		fmt.Printf("SDK initialization failed: %v", err)
		return
	}

	if err := sdk.SetFunctionsPipeline(cloudmqtt.FactoryTransport(sdk)); err != nil {
		sdk.LoggingClient.Error(fmt.Sprintf("main sdk.SetPipeline failed: %v", err))
		return
	}

	sdk.MakeItRun()
}
