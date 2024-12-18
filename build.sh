#!/bin/bash

# =====================
# 环境变量与参数
# =====================
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=${1:-amd64}    # 通过外部参数指定 GOARCH，默认 amd64

# 检查 GOARCH 参数是否有效
if [[ "$GOARCH" != "amd64" && "$GOARCH" != "arm64" ]]; then
    echo "Error: Unsupported GOARCH. Only 'amd64' or 'arm64' are accepted."
    exit 1
fi

# 项目变量
DST_DIR="x-ui"
EXE_NAME="x-main"
ZIP_NAME="x-ui-${GOOS}-${GOARCH}.tar.gz"
# 工具链变量
HOST_ARCH=$(uname -m)
TOOLCHAIN_DIR="gcc-linaro-7.5.0-2019.12-x86_64_aarch64-linux-gnu"
TOOLCHAIN_URL="https://releases.linaro.org/components/toolchain/binaries/latest-7/aarch64-linux-gnu/${TOOLCHAIN_DIR}.tar.xz"
TOOLCHAIN_ARCHIVE="${TOOLCHAIN_DIR}.tar.xz"

# 错误检查
check_error() {
    if [ $? -ne 0 ]; then
        echo "Error: $1"
        cleanup
        exit 1
    fi
}

# 下载并设置工具链
setup_toolchain() {
    if [[ "$HOST_ARCH" = "x86_64" && "$GOARCH" = "arm64" ]]; then
        wget -c -O "$TOOLCHAIN_ARCHIVE" "$TOOLCHAIN_URL"
        check_error "Failed to download toolchain"
        tar -xf "$TOOLCHAIN_ARCHIVE"
        check_error "Failed to extract toolchain"
        export PATH=$(pwd)/$TOOLCHAIN_DIR/bin:$PATH
        export CC=aarch64-linux-gnu-gcc
        export CXX=aarch64-linux-gnu-g++
    fi
}

# 清理目标目录
clean_dst() {
    [ -d "$DST_DIR" ] && rm -rf "$DST_DIR"
    mkdir -p "$DST_DIR/bin/"
}

# 构建项目
build_project() {
    echo "Building for CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH..."
    go build -ldflags="-s -w" -o "$EXE_NAME"
    check_error "Build failed"
}

# 复制文件
copy_files() {
    cp "$EXE_NAME" "$DST_DIR/x-ui"
    cp "x-ui.sh" "x-ui.service" "$DST_DIR/"

    # 复制 bin 目录内容
    [ -f "bin/xray-$GOOS-$GOARCH" ] && cp "bin/xray-$GOOS-$GOARCH" "$DST_DIR/bin/"
    compgen -G "bin/*.dat" > /dev/null && cp bin/*.dat "$DST_DIR/bin/"
}

# 打包目标文件夹
pack_files() {
    echo "Packing folder to: $ZIP_NAME"
    tar -czf "$ZIP_NAME" "$DST_DIR"
    check_error "Packing failed"
}

# 清理临时文件
cleanup() {
    rm -f "$EXE_NAME"
    rm -rf "$DST_DIR"
    if [ "$GOARCH" = "arm64" ]; then
        rm -rf "$TOOLCHAIN_DIR" "$TOOLCHAIN_ARCHIVE"
    fi
}

# 主逻辑
main() {
    setup_toolchain
    build_project
    clean_dst
    copy_files
    pack_files
    cleanup
    echo "Completed successfully! Archive: $ZIP_NAME"
}

main
