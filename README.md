# Strongoid

A simple approach to strongly typed IDs in Go.

Usage:

Get the dependency:
```shell
go get github.com/phrkdll/strongoid
```

Then in your code, import and create a custome type based on the Id type provided by this package:
```golang
package main

import (
    "github.com/phrkdll/strongoid"
)

type UserId = strongoid.Id[int64]
```