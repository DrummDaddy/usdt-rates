package rates

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func TopN(values []decimal.Decimal, n int) (decimal.Decimal, error) {
	if n < 1 {
		return decimal.Zero, fmt.Errorf(`n must be greater than zero`)

	}
	if n > len(values) {
		return decimal.Zero, fmt.Errorf(`n=%d must be greater or equal to length=%d of values`, n, len(values))

	}
	return values[n-1], nil

}

func AvgNM(values []decimal.Decimal, n, m int) (decimal.Decimal, error) {
	if n < 1 || m < 1 {
		return decimal.Zero, fmt.Errorf(`n must be greater or equal to zero`)
	}
	if n > m {
		return decimal.Zero, fmt.Errorf("n=%d must be greater or equal to m=%d", n, m)

	}

	if m > len(values) {
		return decimal.Zero, fmt.Errorf("m=%d must be greater or equal to length=%d", m, len(values))
	}

	sub := values[n-1 : m]

	sum := decimal.Zero
	for _, v := range sub {
		sum = sum.Add(v)
	}
	avg := sum.Div(decimal.NewFromInt(int64(len(sub))))
	return avg, nil

}
