# zoau - A Go module for Z Open Automation Utilities (ZOAU)
[![Go Reference](https://pkg.go.dev/badge/github.com/Stolkerve/zoau-go.svg)](https://pkg.go.dev/github.com/Stolkerve/zoau-go)

## Table of Contents

- [zoau - A Go module for Z Open Automation Utilities (ZOAU)](#zoau---a-nodejs-module-for-z-open-automation-utilities-zoau)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [System Requirements](#system-requirements)
  - [Build and Install](#build-and-install)
  - [Setup](#setup)
  - [Quick Start](#quick-start)
  - [API Documentation](#api-documentation)
  - [Contributing](#contributing)
  - [Legalities](#legalities)

## Overview

zoau - a Go module that exposes the Z Open Automation Utilities (ZOAU)
APIs in Go!

## System Requirements

- IBM® Open Enterprise SDK for Go
- z/OS® V2.5 or higher
- ZOAU v1.1.0 or higher is required on the system.
  - For more details, [see the zoau documentation.](https://www.ibm.com/docs/en/zoau/latest?topic=installing-configuring-zoa-utilities)

## Build and Install

- Before installing, [download and install IBM® Open Enterprise SDK for Go](https://www.ibm.com/products/open-enterprise-sdk-go-zos).

## Setup

- The PATH and LIBPATH environment variables need to include the location of the ZOAU
binaries and dlls, respectively.

```shell
export PATH=<path_to_zoau>/bin:$PATH
export LIBPATH=<path_to_zoau>/lib:$LIBPATH
```

For more details on setting up ZOAU, [see the documentation.](https://www.ibm.com/docs/en/zoau/latest?topic=installing-configuring-zoa-utilities)

## Quick Start

TODO

## More Examples

TODO

## API Documentation

TODO

## Contributing

See the zoau [CONTRIBUTING.md file](CONTRIBUTING.md) for details.

## Legalities

The zoau Go module is available under the Apache 2.0 license. See the [LICENSE
file](LICENSE) file for details