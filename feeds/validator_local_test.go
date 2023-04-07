package feeds

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalValidator(t *testing.T) {
	type test struct {
		article Article
		valid   bool
		errMsgs []string
	}

	ctx := context.TODO()
	lp, err := NewLocalRulesProvider(ctx, "../configs/rules.yaml")
	assert.NoError(t, err)

	v := NewValidator(lp)

	tests := []test{
		{
			article: Article{
				Topic:       "A",
				Name:        "a",
				Description: "valid description",
			},
			valid: true,
		},
		{
			article: Article{
				Topic:       "A",
				Name:        "x",
				Description: "invalid",
			},
			valid: false,
			errMsgs: []string{
				"description: the length must be between (10, 100) exclusive",
				"name: unexpected string",
			},
		},
		{
			article: Article{
				Topic:       "B",
				Name:        "b",
				Description: "a very super long description that >> 40",
			},
			valid: false,
			errMsgs: []string{
				"description: the length must be less than 40",
			},
		},
		{
			article: Article{
				Topic:       "C",
				Name:        "c",
				Description: "valid description",
			},
			valid: true,
		},
		{
			article: Article{
				Topic:       "C",
				Name:        "c",
				Description: "invalid",
			},
			valid: false,
			errMsgs: []string{
				"description: the length must be more than 10",
			},
		},
	}

	for _, tc := range tests {
		err := v.Validate(ctx, tc.article)
		if tc.valid {
			assert.NoError(t, err)
			continue
		}
		assert.Error(t, err)
		for _, em := range tc.errMsgs {
			assert.ErrorContains(t, err, em)
		}
	}
}
