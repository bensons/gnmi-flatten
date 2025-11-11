# gnmi-flatten

A simple CLI tool to convert gNMI subscribe messages from JSON format to PathElem-style labels with their values.

## Overview

This tool reads gNMI (gRPC Network Management Interface) subscribe messages in JSON format and outputs the paths in a flattened, human-readable PathElem format along with their values.

## Features

- Parses gNMI Notification messages in JSON format
- Supports both single notifications and arrays of notifications
- Handles PathElem format with keys (e.g., `interface[name=Ethernet1]`)
- Processes prefix paths and combines them with update paths
- Supports all gNMI value types (string, int, uint, bool, float, bytes, JSON, etc.)
- Handles both update and delete operations

## Building

```bash
go build -o gnmi-flatten
```

## Usage

```bash
./gnmi-flatten -file <path-to-json-file>
```

### Output Format

The tool outputs paths in the format:

```
[origin:]path/element[key=value]/element = value
```

For example:
```
openconfig:/interfaces/interface[name=Ethernet1]/state/admin-status = UP
openconfig:/interfaces/interface[name=Ethernet1]/state/oper-status = UP
openconfig:/interfaces/interface[name=Ethernet1]/state/counters/in-octets = 1.23456789e+09
```

## Input Format

The tool accepts gNMI Notification messages in JSON format. The input can be either:

1. A single notification object
2. An array of notification objects

## Supported Value Types

The tool supports all gNMI value types:
- `string_val`
- `int_val`
- `uint_val`
- `bool_val`
- `float_val`
- `bytes_val`
- `decimal_val`
- `json_val`
- `leaflist_val`
- `any_val`
- `ascii_val`

## Dependencies

- [github.com/openconfig/gnmi](https://github.com/openconfig/gnmi) - gNMI protocol definitions
- [github.com/openconfig/ygot](https://github.com/openconfig/ygot) - YANG Go tools

