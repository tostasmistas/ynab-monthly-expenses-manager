package backend

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

func StringToDecimalHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() != reflect.String {
			return data, nil
		}

		if to != reflect.TypeOf(decimal.Decimal{}) {
			return data, nil
		}

		return decimal.NewFromString(data.(string))
	}
}
