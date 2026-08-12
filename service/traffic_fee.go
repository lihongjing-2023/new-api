package service

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/shopspring/decimal"
)

// TrafficFeeQuotaDecimal 返回每次请求固定流量费换算的内部额度点数（decimal）。
// 预扣阶段不收取，按量计费结算阶段统一叠加；流量费非正数或分组倍率非正时返回 0，不影响现有行为。
func TrafficFeeQuotaDecimal(groupRatio float64) decimal.Decimal {
	fee := operation_setting.GetTrafficFee()
	if fee <= 0 || groupRatio <= 0 {
		return decimal.Zero
	}
	return decimal.NewFromFloat(fee).
		Mul(decimal.NewFromFloat(common.QuotaPerUnit)).
		Mul(decimal.NewFromFloat(groupRatio))
}
