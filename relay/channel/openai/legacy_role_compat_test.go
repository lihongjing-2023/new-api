package openai

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyLegacyRoleCompat(t *testing.T) {
	request := &dto.GeneralOpenAIRequest{
		Model: "gpt-4",
		Messages: []dto.Message{
			{Role: "system"},
			{Role: "developer"},
			{Role: "user"},
			{Role: "system"},
			{Role: "assistant"},
			{Role: "tool"},
		},
	}

	applyLegacyRoleCompat(request)

	require.Len(t, request.Messages, 6)
	// system 一律改写为 user，首条也不例外
	assert.Equal(t, "user", request.Messages[0].Role)
	// developer 一律改写为 user
	assert.Equal(t, "user", request.Messages[1].Role)
	assert.Equal(t, "user", request.Messages[2].Role)
	assert.Equal(t, "user", request.Messages[3].Role)
	// assistant/tool 不受影响
	assert.Equal(t, "assistant", request.Messages[4].Role)
	assert.Equal(t, "tool", request.Messages[5].Role)
}

func TestApplyLegacyRoleCompatEmptyAndNil(t *testing.T) {
	applyLegacyRoleCompat(nil)
	applyLegacyRoleCompat(&dto.GeneralOpenAIRequest{})
}

func newRoleCompatRelayInfo(enabled bool) *relaycommon.RelayInfo {
	return &relaycommon.RelayInfo{
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:       constant.ChannelTypeOpenAI,
			UpstreamModelName: "gpt-4",
			ChannelSetting:    dto.ChannelSettings{LegacyRoleCompat: enabled},
		},
	}
}

// 开关打开时，ConvertOpenAIRequest 必须在所有渠道级适配之后统一改写角色
func TestConvertOpenAIRequestLegacyRoleCompatEnabled(t *testing.T) {
	adaptor := &Adaptor{}
	request := &dto.GeneralOpenAIRequest{
		Model: "gpt-4",
		Messages: []dto.Message{
			{Role: "system"},
			{Role: "developer"},
			{Role: "user"},
			{Role: "system"},
		},
	}

	result, err := adaptor.ConvertOpenAIRequest(nil, newRoleCompatRelayInfo(true), request)
	require.NoError(t, err)

	converted, ok := result.(*dto.GeneralOpenAIRequest)
	require.True(t, ok)
	roles := make([]string, 0, len(converted.Messages))
	for _, message := range converted.Messages {
		roles = append(roles, message.Role)
	}
	assert.Equal(t, []string{"user", "user", "user", "user"}, roles)
}

// 开关关闭时不得改写任何角色
func TestConvertOpenAIRequestLegacyRoleCompatDisabled(t *testing.T) {
	adaptor := &Adaptor{}
	request := &dto.GeneralOpenAIRequest{
		Model: "gpt-4",
		Messages: []dto.Message{
			{Role: "system"},
			{Role: "developer"},
			{Role: "user"},
			{Role: "system"},
		},
	}

	result, err := adaptor.ConvertOpenAIRequest(nil, newRoleCompatRelayInfo(false), request)
	require.NoError(t, err)

	converted, ok := result.(*dto.GeneralOpenAIRequest)
	require.True(t, ok)
	roles := make([]string, 0, len(converted.Messages))
	for _, message := range converted.Messages {
		roles = append(roles, message.Role)
	}
	assert.Equal(t, []string{"system", "developer", "user", "system"}, roles)
}

// 开关打开时必须压过 o 系列的 system->developer 适配（最终强制）
func TestConvertOpenAIRequestLegacyRoleCompatOverridesOModelAdaptation(t *testing.T) {
	adaptor := &Adaptor{}
	info := newRoleCompatRelayInfo(true)
	info.UpstreamModelName = "o3-mini"
	request := &dto.GeneralOpenAIRequest{
		Model: "o3-mini",
		Messages: []dto.Message{
			{Role: "system"},
			{Role: "user"},
		},
	}

	result, err := adaptor.ConvertOpenAIRequest(nil, info, request)
	require.NoError(t, err)

	converted, ok := result.(*dto.GeneralOpenAIRequest)
	require.True(t, ok)
	assert.Equal(t, "user", converted.Messages[0].Role)
	assert.Equal(t, "user", converted.Messages[1].Role)
}
