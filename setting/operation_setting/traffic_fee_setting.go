package operation_setting

import "github.com/QuantumNous/new-api/setting/config"

// TrafficFeeSetting 配置按量计费模型每次请求的固定附加费（流量费 / 底费）。
type TrafficFeeSetting struct {
	// TrafficFee 每次请求固定流量费（美元），如 0.001；设为 0 关闭，不影响现有行为。
	TrafficFee float64 `json:"traffic_fee"`
}

// 默认配置：流量费关闭
var trafficFeeSetting = TrafficFeeSetting{TrafficFee: 0}

func init() {
	// 注册到全局配置管理器，后台运营设置可持久化
	config.GlobalConfig.Register("traffic_fee_setting", &trafficFeeSetting)
}

// GetTrafficFee 返回每次请求的固定流量费（美元）。
func GetTrafficFee() float64 {
	return trafficFeeSetting.TrafficFee
}

// SetTrafficFeeForTest 注入流量费配置。Tests only.
func SetTrafficFeeForTest(fee float64) {
	trafficFeeSetting.TrafficFee = fee
}
