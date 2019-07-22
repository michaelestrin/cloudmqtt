# michaelestrin/cloudmqtt -- Example EdgeX-Cloud MQTTS Service

## Overview

This service is build upon the EdgeX Applications Functions SDK.  

It receives northbound events and forwards them to a user-defined topic on an upstream MQTTS server.  When an event 
    references a device this service hasn't seen before, the service will separately query the EdgeX core-metadata 
    service for the device's metadata and forward the result to a separate user-defined topic on an upstream MQTTS 
    server.
    
It receives southbound commands and logs them.  As there is no defined standard for translating command content received
    via MQTT message into specific core-command api calls, the provided `impl.commandHandler.Receiver` must be modified
    to translate the incoming command string into call to appropriate core-command endpoint.

## Table of Contents

- [Basic Setup and Usage](#basic-setup-and-usage)
    - [External Service Dependencies](#external-service-dependencies)
    - [Configuration](#configuration)
- [Code](#code)
    - [Architecture](#architecture)
    - [Project Layout](#project-layout)
    - [Assumptions](#assumptions)
- [Connecting to the Cloud](#connecting-to-the-cloud)
    

## Basic Setup And Usage

### External Service Dependencies

This service requires a running instance of EdgeX.

During development, I used the tip of master (updating as changes were made) of 
    [edgexfoundry/edgex-go](https://github.com/edgexfoundry/edgex-go) and 
    [edgexfoundry/device-random](https://github.com/edgexfoundry/device-random) to run the supporting services 
    locally. Since this sample service is build upon the EdgeX Applications Functions SDK, it required a post-Delhi 
    release version of the source.  
    
### Configuration

This service is driven by a `configuration.toml` [TOML-based](https://en.wikipedia.org/wiki/TOML) file that resides in
    the directory the application is executed from.  The configuration file extends the configuration file defined by 
    the EdgeX Application Functions SDK.  Configuration specific to this sample service is constrained to the 
    `[ApplicationSettings]` section and has several required/optional key/value pairs:
    
- `certFile` - a string, this defines the path and name of a file containing the public key to use for the MQTTS 
    connection.  Required if `keyfile` is provided, otherwise optional.
- `keyFile` - a string, this defines the path and name of a file containing the private key to use for the MQTTS 
    connection.  Required if `certfile` is provided, otherwise optional.
- `clientId` - a string, this defines the value passed to the MQTTS instance to uniquely identify the adapter.
- `userName` - a string, this defines the value passed to the MQTTS instance to uniquely identify the user.
- `password` - a string, this defines the value passed to the MQTTS instance to uniquely identify the password.
- `server` - a string, this defines the address for a running MQTTS instance that will receive events and metadata.
- `edgeXMetaDataUri` - a string, this defines the address for a running instance of the EdgeX core-metadata service.  
- `dataTopic` - a string, this defines the MQTT topic that will receive device events/readings.
- `commandTopic` - a string, this defines the MQTT topic that will receive device metadata.

A sample configuration file can be found at 
    [`configs/configuration.toml`](https://github.com/michaelestrin/cloudmqtt/blob/master/configs/configuration.toml).
    
    
## Code

### Architecture

This service implementation leverages constructor-based 
    [dependency injection](https://en.wikipedia.org/wiki/Dependency_injection) to facilitate decoupling and testability. 
    A factory function is called to return an Applications Functions SDK-compatible transform function that is 
    subsequently called by the SDK. 
    
### Project Layout

The project loosely follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout). 
    Specifically, the project is structured as follows:

```
cmd/                     - main project application
    main.go              - service's execution entry point
configs/                 - configuration file templates or default configs
    configuration.toml   - sample configuration file
internal/                - private application and library code
    app/                 - private application code
        cloudmqtt/       - service's application code
            contract/    - service's contracts
            impl/        - service's implementation (MQTT, notification) created by application code's factory function 
            test/        - support for unit testing
                helper/  - unit test-specific helper functions used by more than one unit test package
                stub/    - test stub implementations of system objects (used instead of mocks)
```

### Assumptions

- Service was developed and tested with EdgeX's latest version (i.e. tip of master) after the Edinburgh. It may or may 
    not function properly with other EdgeX versions.
- MQTTS is used for transport.
- Events/readings and metadata will be pushed onto configured MQTT topics transformed to JSON but otherwise 
    with content as received.
- A specific device's metadata is sent once per instance execution.    
- There is no shared knowledge of existing devices across service instances. Each service instance tracks its own 
    devices and forwards metadata for any device for which it has not seen a reading from before.
- Knowledge of existing devices is not persisted across instance executions. A newly restarted service will transmit 
    metadata for each device connected to it (even if a previously executed instance sent that same metadata). 
- Device metadata and the first reading for a device may be received by the northbound application in an unpredictable 
    order.  That is, the new device's first reading may show up before, at the same time as, or after the device's 
    metadata. 
    
## Connecting to the Cloud

[Connecting to AWS IoT](docs/aws/README.md)

[Connecting to Azure IoT](docs/azure/README.md)

[Connecting to Dell Boomi](docs/boomi/README.md)