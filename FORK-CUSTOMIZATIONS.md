# Fork 个性化改造清单

> 本文件用于跟踪本仓库（fork）相对上游官方仓库的个性化改动，方便后期维护、同步上游代码时参考。

## 仓库信息

| 项 | 值 |
|----|----|
| 上游仓库 | `https://github.com/QuantumNous/new-api.git` |
| 本仓库（fork） | `https://github.com/lihongjing-2023/new-api.git` |
| 上游基线 commit | `ccd535ef`（`fix: harden concurrent quota and status updates`） |
| 当前 HEAD | `a8aab7a0`（`ci: add build-artifacts workflow ...`） |
| 同步策略 | 以 `git merge origin/main` / `git pull` 方式跟踪上游，改造尽量以独立提交保持，便于 cherry-pick |

## 改造总览

| # | 改造项 | 类型 | 状态 |
|---|--------|------|------|
| 1 | 流量计费功能 | 功能新增 | 已合并 |
| 2 | CNB 云开发环境初始化 | 开发环境 | 已合并 |
| 3 | CI 产物构建 workflow | CI/CD | 已合并 |

---

## 改造明细

### 1. 流量计费功能

- **commit**: `83177481 feat: 完成流量计费功能`
- **说明**: 按流量对任务计费，管理端可配置流量费率；用量日志展示流量费用明细。
- **涉及文件**:
  - `service/traffic_fee.go`（新增，流量费计算）
  - `service/traffic_fee_test.go`（新增，单元测试）
  - `setting/operation_setting/traffic_fee_setting.go`（新增，费率配置项）
  - `service/task_billing.go` / `service/task_polling.go` / `service/text_quota.go`（改动，计费链路接入）
  - `web/src/features/system-settings/billing/traffic-fee-settings.tsx`（新增，管理端配置界面）
  - `web/src/features/system-settings/billing/index.tsx`、`section-registry.tsx`、`types.ts`（配置页注册）
  - `web/src/features/usage-logs/...`（用量日志展示）
  - `web/src/i18n/locales/en.json`、`zh.json`（前端文案）
  - 对应测试文件：`service/task_billing_test.go`、`task_polling.go` 等
- **同步上游注意事项**: 与上游冲突点主要在 `service/task_billing.go`、`service/text_quota.go`、`web/src/features/system-settings/billing/*` 和 i18n 词条文件，合入上游新代码时需重点检查这几处。

### 2. CNB 云开发环境初始化

- **commit**: `b8623929 chore: add CNB cloud dev environment initialization config`
- **说明**: 增加云开发环境（CNB / CloudStudio）的初始化配置。
- **涉及文件**:
  - `.cnb.yml`（新增）
  - `.ide/Dockerfile`（新增）
- **同步上游注意事项**: 与上游冲突可能性低，基本无侵入。

### 3. CI 产物构建 workflow

- **commit**: `a8aab7a0 ci: add build-artifacts workflow (frontend + docker image as GitHub artifacts)`
- **说明**: 在 fork 仓库新增 `build-artifacts` workflow，手动触发或 push 到 `main` / 打 `v*` tag 时：
  - 构建前端静态产物（`web/dist`），上传 Artifact `frontend-dist`
  - 构建完整 Docker 镜像，导出为 tar.gz，上传 Artifact `new-api-docker-image`
  - 产物下载后自行使用：`docker load -i new-api-image.tar.gz`
- **涉及文件**:
  - `.github/workflows/build-artifacts.yml`（新增）
- **同步上游注意事项**: 上游官方仓库同名目录下无此 workflow，无冲突。

---

## 历史改动记录（已回退 / 未保留）

以下改动曾在本地工作区存在，但因回退已丢弃，**不在当前分支中**。记录于此以防后期需要重新实现：

| 项 | 说明 | 状态 |
|----|------|------|
| Dockerfile 国内镜像加速 | 增加 `NPM_CONFIG_REGISTRY`（npmmirror）与 `GOPROXY`（goproxy.cn） | 已回退，未保留 |
| rsbuild 部署配置 | `serverUrl` 指向 `https://buddybackend.cloud`、`assetPrefix: '/'` | 已回退，未保留 |
| i18n 兜底语言 | `fallbackLng` 改为 `zhCN`，检测顺序去掉 `navigator` | 已回退，未保留 |
| 充值表单 i18n | `Pay {{amount}}` / `Save {{amount}}` 文案国际化 | 已回退，未保留 |
| `lightningcss-win32-x64-msvc` 依赖 | `web/package.json` 增加 devDependency 以支持本地构建 | 已回退，未保留 |
| 废弃文件 `nul` | Windows 下误创建的空文件 | 已删除 |

---

## 上游同步操作手册

### 拉取上游最新代码并合并

```bash
git fetch origin
git merge origin/main
# 若发生冲突，优先保留本 fork 的功能性改动（流量计费、CI workflow），
# 参考上面各改造项的"同步上游注意事项"
```

### 查看本 fork 相对上游的差异

```bash
git log --oneline origin/main..HEAD      # 本 fork 独有提交
git diff origin/main..HEAD --stat        # 本 fork 改动文件总览
```

### 更新本清单

每次新增/回退个性化改动后，同步更新本文件中的"改造明细"或"历史改动记录"。
