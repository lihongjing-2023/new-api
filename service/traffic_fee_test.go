package service

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func setTrafficFeeForTest(t *testing.T, fee float64) {
	t.Helper()
	old := operation_setting.GetTrafficFee()
	operation_setting.SetTrafficFeeForTest(fee)
	t.Cleanup(func() { operation_setting.SetTrafficFeeForTest(old) })
}

func TestTrafficFeeQuotaDisabledWhenFeeNonPositive(t *testing.T) {
	for _, fee := range []float64{0, -0.001} {
		setTrafficFeeForTest(t, fee)
		require.True(t, TrafficFeeQuotaDecimal(2.5).IsZero(), "fee=%v", fee)
	}
}

func TestTrafficFeeQuotaScalesWithFeeAndGroupRatio(t *testing.T) {
	setTrafficFeeForTest(t, 0.001)

	expected := decimal.NewFromFloat(0.001).
		Mul(decimal.NewFromFloat(common.QuotaPerUnit)).
		Mul(decimal.NewFromFloat(2.5))
	require.True(t, expected.Equal(TrafficFeeQuotaDecimal(2.5)), "got %s want %s", TrafficFeeQuotaDecimal(2.5), expected)
}

func TestTrafficFeeQuotaZeroForNonPositiveGroupRatio(t *testing.T) {
	setTrafficFeeForTest(t, 0.001)
	for _, gr := range []float64{0, -2.5} {
		require.True(t, TrafficFeeQuotaDecimal(gr).IsZero(), "groupRatio=%v", gr)
	}
}
