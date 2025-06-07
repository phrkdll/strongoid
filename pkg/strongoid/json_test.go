package strongoid_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/phrkdll/strongoid/pkg/strongoid"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	for _, tc := range stringIdTestCases {
		t.Run(tc.in.Inner, func(t *testing.T) {
			jsonBytes, err := strongoid.Id[string](tc.in).MarshalJSON()
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
			jsonBytes, err := strongoid.Id[int64](tc.in).MarshalJSON()
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
			b := []byte("\"" + tc.in.Inner + "\"")

			err := (*strongoid.Id[string])(&id).UnmarshalJSON(b)
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
