package openai

import (
	"github.com/QuantumNous/new-api/relaykit/dto"
)

// applyLegacyRoleCompat 把新版 OpenAI 角色约定改写为旧版上游可接受的形式：
//   - developer -> user（旧版兼容上游不识别 developer 角色）
//   - system -> user（部分兼容上游只接受 user/assistant，不接受任何 system）
//
// 不改写 assistant/tool，避免破坏工具调用链。
// 只在已解析的消息结构上原地修改 Role 字符串，不触发任何 JSON 重排，
// 用于替代 channels 参数覆盖里 messages.*.role 的通配符改写。
func applyLegacyRoleCompat(request *dto.GeneralOpenAIRequest) {
	if request == nil {
		return
	}
	for i := range request.Messages {
		switch request.Messages[i].Role {
		case "developer", "system":
			request.Messages[i].Role = "user"
		}
	}
}
