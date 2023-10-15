package copyutil

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
)

func Copy(from interface{}, to interface{}) {
	copier.CopyWithOption(to, from, copier.Option{
		Converters: []copier.TypeConverter{
			DecimalToFloat64Converter(),
		},
	})
}

func DecimalToFloat64Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: decimal.Decimal{},
		DstType: copier.Float64,
		Fn: func(src interface{}) (interface{}, error) {
			d, ok := src.(decimal.Decimal)
			if !ok {
				return nil, errors.New("src type is not decimal.Decimal")
			}
			return d.InexactFloat64(), nil
		},
	}
}
