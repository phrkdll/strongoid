package strongoid

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StringId Id[string]
type IntegerId Id[int64]

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

func TestMarshalJSON(t *testing.T) {
	for _, tc := range stringIdTestCases {
		t.Run(tc.in.Inner, func(t *testing.T) {
			jsonBytes, err := Id[string](tc.in).MarshalJSON()
			assert.NotEmpty(t, jsonBytes)

			jsonString := string(jsonBytes)

			assert.Equal(t, tc.in.Inner, strings.Trim(jsonString, "\""))

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	for _, tc := range integerIdTestCases {
		t.Run(fmt.Sprintf("%v", tc.in.Inner), func(t *testing.T) {
			jsonBytes, err := Id[int64](tc.in).MarshalJSON()
			assert.NotEmpty(t, jsonBytes)

			assert.Equal(t, fmt.Sprintf("%v", tc.in.Inner), string(jsonBytes))

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	for _, tc := range stringIdTestCases {
		t.Run(tc.in.Inner, func(t *testing.T) {
			id := StringId{}
			bytes := []byte("\"" + tc.in.Inner + "\"")

			err := (*Id[string])(&id).UnmarshalJSON(bytes)
			assert.NotEmpty(t, id)

			assert.Equal(t, tc.in, id)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
