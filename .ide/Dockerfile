# CNB 云原生开发环境镜像
# 基于 CNB 默认开发环境模板 cnbcool/default-dev-env（已内置：
#   - code-server WebIDE / VSCode server（单容器模式）
#   - Node.js 22、git、git-lfs、zsh、oh-my-zsh、CNB CLI、CNB skills 等
#   - zh_CN.UTF-8 中文 locale、openssh-server（支持桌面端 SSH 连接）
# ）
# 在此之上补充 new-api 项目开发所需：
#   - Go 1.26.1（后端，>= go.mod 声明的 go 1.25.1）
#   - Bun（前端包管理器与脚本运行器，对应 web/bun.lock）
#   - make（项目 makefile 构建入口依赖）
FROM cnbcool/default-dev-env:latest

# 安装构建工具链：build-essential 提供 gcc/g++/make 及 libc6-dev（C 头文件，
# Go 的 cgo 编译必需 stdint.h/stdlib.h 等，模板镜像未预装）；
# unzip 供 Bun 安装脚本解压使用；curl/ca-certificates 用于下载 Go 与 Bun
RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    build-essential curl ca-certificates unzip && \
    rm -rf /var/lib/apt/lists/*

# 安装 Go 1.26.1（官方预编译二进制，按架构下载解压到 /usr/local/go）
ENV GO_VERSION=1.26.1
RUN ARCH=$(dpkg --print-architecture 2>/dev/null || uname -m) && \
    case "$ARCH" in \
      amd64|x86_64) GO_ARCH=amd64 ;; \
      arm64|aarch64) GO_ARCH=arm64 ;; \
      *) echo "unsupported arch: $ARCH" && exit 1 ;; \
    esac && \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" -o /tmp/go.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz && \
    /usr/local/go/bin/go version
ENV PATH="/usr/local/go/bin:${PATH}"

# 安装 Bun（前端包管理器，1.x 与 web/bun.lock 兼容；安装脚本由 GitHub Releases 分发）
# 显式验证版本：管道中安装脚本失败时退出码可能仍为 0，须在构建期暴露问题
RUN curl -fsSL https://bun.sh/install | bash && \
    /root/.bun/bin/bun --version
ENV BUN_INSTALL="/root/.bun" \
    PATH="/root/.bun/bin:${PATH}"

# 字符集：支持命令行中文输入（官方推荐配置）
ENV LANG C.UTF-8
ENV LANGUAGE C.UTF-8

# CNB 会将仓库挂载到 /workspace，作为统一工作目录
WORKDIR /workspace
