package testid

import "github.com/phrkdll/strongoid/pkg/strongoid"

//go:generate go run github.com/phrkdll/strongoid/cmd/gen --modules=json,gorm

type TestStringId strongoid.Id[string]
type TestIntId strongoid.Id[int64]
