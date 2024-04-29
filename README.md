# mikrotik-cabletest

mikrotik-cabletest is a utility for testing cables using MikroTik RouterOS devices.

## Introduction

mikrotik-cabletest provides a convenient solution for testing the integrity of cables using MikroTik RouterOS devices with cable testing capabilities. It streamlines the cable testing process, making it easier to identify and troubleshoot cable-related issues.

## Features

- **Cable wiring test** 
- **Speed Measurement**

## Requirements

- 2 MikroTik RouterOS device with cable testing functionality.
- Go programming language installed on the local system.

## Installation

To install mikrotik-cabletest, use the following `go get` command:

```bash
go get github.com/ui-kreinhard/mikrotik-cabletest
```

This command will download and install the mikrotik-cabletest utility along with its dependencies.

## Usage

After installation, you can use the mikrotik-cabletest command to initiate cable tests. Here's a basic example:

```bash
env SWITCH_IP=192.168.88.2 SWITCH_USERNAME=admin SWITCH_PASSWORD=admin SSH_PORT=22 PORT_TO_TEST=ether3 required=true=22 mikrotik-cabletest
```
