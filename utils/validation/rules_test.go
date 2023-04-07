package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrLenBetween(t *testing.T) {
	type test struct {
		val                any
		moreThan, lessThan int
		valid              bool
	}

	tests := []test{
		{val: "123456", moreThan: 5, lessThan: 10, valid: true},
		{val: "12345", moreThan: 5, lessThan: 10, valid: false},
		{val: "1234567890", moreThan: 5, lessThan: 10, valid: false},
		{val: "123456", moreThan: 5, lessThan: 0, valid: true},
		{val: "1234", moreThan: 0, lessThan: 5, valid: true},
		{val: 12, moreThan: 0, lessThan: 5, valid: false},
		{val: "", moreThan: -1, lessThan: 1, valid: true},
	}
	for _, tc := range tests {
		rule := StrLenBetween(tc.moreThan, tc.lessThan)
		err := rule(tc.val)
		if tc.valid {
			assert.NoError(t, err)
			continue
		}
		assert.Error(t, err)
	}
}

func TestStrEquals(t *testing.T) {
	type test struct {
		val      any
		expected string
		valid    bool
	}

	tests := []test{
		{val: "value", expected: "value", valid: true},
		{val: "", expected: "value", valid: false},
		{val: 12, expected: "value", valid: false},
	}
	for _, tc := range tests {
		rule := StrEquals(tc.expected)
		err := rule(tc.val)
		if tc.valid {
			assert.NoError(t, err)
			continue
		}
		assert.Error(t, err)
	}
}
