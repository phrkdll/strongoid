# strongoid

A simple approach to strongly typed IDs in Go. Avoid confusing IDs because all of them have the same primitive type.

## Usage

Get the dependency:

```shell
go get github.com/phrkdll/strongoid
```

Then in your code, import and create a custom type based on the Id type provided by this package:

```go
package main

import (
    "github.com/phrkdll/strongoid/pkg/strongoid"
)

type UserId strongoid.Id[int64]

//go:generate go run github.com/phrkdll/strongoid/cmd/gen --modules=json,gorm
```

The generator comment enables generated code for integrating your strongly typed IDs.
This needs to be done once within each package that contains IDs that you want to generate code for.

Integrations are currently supported for:
- [JSON](https://pkg.go.dev/encoding/json) (MarshalJSON / UnmarshalJSON)
- [GORM](https://gorm.io/) (Scan / Value)

To generate code, execute the following command in your project's root folder:

```shell
go generate ./...
```
