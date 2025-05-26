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
```

Add generator comment in your code (within the package that contains your IDs) to enable generated JSON un-/marshalling for all your strongly typed IDs.

```go
//go:generate go run github.com/phrkdll/strongoid/cmd/gen
```

Execute the following to generate code:

```shell
go generate ./...
```
