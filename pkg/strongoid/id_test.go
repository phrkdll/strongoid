package strongoid_test

import "github.com/phrkdll/strongoid/pkg/strongoid"

type StringId strongoid.Id[string]
type IntegerId strongoid.Id[int64]

var stringIdTestCases = []struct {
	in          StringId
	expectError bool
}{
	{StringId{Inner: "person/1"}, false},
	{StringId{Inner: "vehicle/22"}, false},
}

var integerIdTestCases = []struct {
	in          IntegerId
	expectError bool
}{
	{IntegerId{Inner: 1}, false},
	{IntegerId{Inner: 2}, false},
}
