package feeds

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	valext "github.com/govinda-attal/winning11/utils/validation"
)

type validator struct {
	rulesPrv RulesProvider
}

func NewValidator(rp RulesProvider) Validator {
	return &validator{
		rulesPrv: rp,
	}
}

func (v *validator) Validate(ctx context.Context, a Article) error {
	rules, err := v.rulesPrv.Rules(ctx, a.Topic)
	if err != nil {
		// unable to fetch rules for a given topic
		return err
	}

	return validation.ValidateStruct(&a,
		validation.Field(&a.Name,
			validation.By(
				valext.StrEquals(rules.Name.Value),
			),
		),
		validation.Field(&a.Description,
			validation.By(
				valext.StrLenBetween(rules.Desc.LenMoreThan, rules.Desc.LenLessThan),
			),
		),
	)
}
