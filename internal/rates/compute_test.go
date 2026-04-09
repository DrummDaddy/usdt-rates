package rates

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTopN(t *testing.T) {
	values := []decimal.Decimal{
		decimal.NewFromInt(10),
		decimal.NewFromInt(20),
		decimal.NewFromInt(30),
	}
	v, err := TopN(values, 2)
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(20), v)
}

func TestAvgNM(t *testing.T) {
	values := []decimal.Decimal{
		decimal.NewFromInt(10),
		decimal.NewFromInt(20),
		decimal.NewFromInt(30),
		decimal.NewFromInt(40),
	}

	// avg of [2;4] => (20 + 30 + 40) / 3 = 30
	avg, err := AvgNM(values, 2, 4)
	require.NoError(t, err)

	expected := decimal.NewFromInt(30)
	require.True(t, avg.Equal(expected), "expected %s, got %s", expected, avg)
}
