package validation

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func StrLenBetween(more, less int) validation.RuleFunc {
	message := "the value must be empty"

	if more == 0 && less > 0 {
		message = fmt.Sprintf("the length must be less than %v", less)
	} else if more >= 0 && less == 0 {
		message = fmt.Sprintf("the length must be more than %v", more)
	} else if more >= 0 && less > 0 {
		message = fmt.Sprintf("the length must be between (%v, %v) exclusive", more, less)
	}

	return func(value interface{}) error {
		value, isNil := validation.Indirect(value)
		if isNil {
			return nil
		}

		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("value not string")
		}

		l := len(s)

		if more >= 0 && l <= more || less > 0 && l >= less {
			return errors.New(message)
		}

		return nil
	}
}
func StrEquals(str string) validation.RuleFunc {
	return func(value interface{}) error {
		value, isNil := validation.Indirect(value)
		if isNil {
			return nil
		}
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("value not string")
		}
		if s != str {
			return errors.New("unexpected string")
		}
		return nil
	}
}
