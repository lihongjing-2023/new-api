package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 本组用例源于真实用户配置（reasoning_effort 相关参数覆盖），锁定两个关键语义：
//  1. pass_missing_key 只在"字段缺失"时放行；字段存在但值不匹配时仍会被拦截；
//  2. 多条操作串行执行、作用于同一份请求体，后一条操作的条件能看到前一条操作的产物。
func userScenarioOverride() map[string]interface{} {
	cond := func(path, mode string, value interface{}, passMissingKey bool) map[string]interface{} {
		return map[string]interface{}{
			"path":             path,
			"mode":             mode,
			"value":            value,
			"pass_missing_key": passMissingKey,
		}
	}
	return map[string]interface{}{
		"operations": []interface{}{
			map[string]interface{}{
				"mode": "replace",
				"path": "reasoning_effort",
				"from": "max",
				"to":   "xhigh",
				"conditions": []interface{}{
					cond("reasoning_effort", "full", "max", false),
				},
				"logic": "OR",
			},
			map[string]interface{}{
				"mode": "replace",
				"path": "reasoning_effort",
				"from": "xhigh",
				"to":   "xhigh",
				"conditions": []interface{}{
					cond("reasoning_effort", "full", "xhigh", false),
				},
				"logic": "OR",
			},
			map[string]interface{}{
				"mode":  "set",
				"path":  "reasoning_effort",
				"value": "high",
				"conditions": []interface{}{
					cond("reasoning_effort", "full", "high", true),
					cond("thinking.type", "full", "enabled", true),
				},
				"logic": "AND",
			},
			map[string]interface{}{
				"mode": "replace",
				"path": "messages.0.role",
				"from": "system",
				"to":   "user",
			},
		},
	}
}

func applyUserScenario(t *testing.T, input []byte) (string, []string) {
	t.Helper()
	recorder := &paramOverrideAuditRecorder{}
	ctx := map[string]interface{}{
		paramOverrideContextAuditRecorder: recorder,
	}
	out, err := ApplyParamOverride(input, userScenarioOverride(), ctx)
	require.NoError(t, err)
	return string(out), recorder.lines
}

func auditContains(audit []string, substr string) bool {
	for _, l := range audit {
		if strings.Contains(l, substr) {
			return true
		}
	}
	return false
}

func TestPassMissingKeyTriggersSetWhenKeyAbsent(t *testing.T) {
	// 请求体里既没有 reasoning_effort 也没有 thinking：前两条 replace 因条件字段
	// 缺失被跳过；第三条 set 的两个条件都配了 pass_missing_key=true，缺失即放行，
	// AND 后条件满足，set 必定触发——这正是"无脑塞入 high"的根源。
	const msg = `"messages":[{"role":"system","content":"hi"}]`
	out, audit := applyUserScenario(t, []byte(`{"model":"gpt-5",`+msg+`}`))

	assert.False(t, auditContains(audit, "from max to xhigh"))
	assert.False(t, auditContains(audit, "from xhigh to xhigh"))
	assert.True(t, auditContains(audit, "set reasoning_effort = high"))
	assert.Contains(t, out, `"reasoning_effort":"high"`)
}

func TestOperationsSerialChainSeesPreviousResult(t *testing.T) {
	// reasoning_effort=max：第一条 max->xhigh 触发并把请求体改为 xhigh；
	// 因为操作串行作用于同一份请求体，第二条 xhigh->xhigh 的条件此刻也能看到
	// 改后的值，于是跟着"触发"（from==to，无实际变化，仅出现在审计中）。
	const msg = `"messages":[{"role":"system","content":"hi"}]`
	out, audit := applyUserScenario(t, []byte(`{"model":"gpt-5","reasoning_effort":"max",`+msg+`}`))

	assert.True(t, auditContains(audit, "from max to xhigh"))
	assert.True(t, auditContains(audit, "from xhigh to xhigh"))
	assert.Contains(t, out, `"reasoning_effort":"xhigh"`)
}

func TestPassMissingKeyDoesNotBypassExplicitMismatch(t *testing.T) {
	// thinking.type 显式存在但为 disabled：该条件 exists=true 且值不匹配 => false，
	// pass_missing_key 只在"字段缺失"时放行，不能绕过"存在但不匹配"的场景，
	// 因此 set 被拦截、请求体不会凭空多出 reasoning_effort。
	const msg = `"messages":[{"role":"system","content":"hi"}]`
	out, audit := applyUserScenario(t, []byte(`{"model":"gpt-5","thinking":{"type":"disabled"},`+msg+`}`))

	assert.False(t, auditContains(audit, "set reasoning_effort = high"))
	assert.NotContains(t, out, `"reasoning_effort"`)
}
