package strongoid_test

import (
	"fmt"
	"testing"

	"github.com/phrkdll/strongoid/pkg/strongoid"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	for _, tc := range stringIdTestCases {
		t.Run(tc.in.Inner, func(t *testing.T) {
			driverVal, err := strongoid.Id[string](tc.in).Value()
			assert.NotEmpty(t, driverVal)

			assert.Equal(t, tc.in.Inner, driverVal)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	for _, tc := range integerIdTestCases {
		t.Run(fmt.Sprintf("%v", tc.in.Inner), func(t *testing.T) {
			driverVal, err := strongoid.Id[int64](tc.in).Value()
			assert.NotEmpty(t, driverVal)

			assert.Equal(t, tc.in.Inner, driverVal)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestScan(t *testing.T) {
	for _, tc := range stringIdTestCases {
		t.Run(tc.in.Inner, func(t *testing.T) {
			id := StringId{}

			err := (*strongoid.Id[string])(&id).Scan(tc.in.Inner)
			assert.NotEmpty(t, id)

			assert.Equal(t, tc.in.Inner, id.Inner)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	for _, tc := range integerIdTestCases {
		t.Run(fmt.Sprint(tc.in.Inner), func(t *testing.T) {
			id := IntegerId{}

			err := (*strongoid.Id[int64])(&id).Scan(tc.in.Inner)
			assert.NotEmpty(t, id)

			assert.Equal(t, tc.in.Inner, id.Inner)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
