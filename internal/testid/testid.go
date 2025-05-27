package testid

import "github.com/phrkdll/strongoid/pkg/strongoid"

type TestId strongoid.Id[string]

//go:generate go run github.com/phrkdll/strongoid/cmd/gen --modules=json,gorm
