#!/bin/bash

# =====================
# 环境变量与参数
# =====================
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=${1:-amd64}    # 通过外部参数指定 GOARCH，默认 amd64

# 项目变量
DST_DIR="x-ui"               
EXE_NAME="x-main"            
ZIP_NAME="x-ui-${GOOS}-${GOARCH}.tar.gz"

# =====================
# 错误检查函数
# =====================
check_error() {
    if [ $? -ne 0 ]; then
        echo "Error: $1"
        exit 1
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
}

# 主逻辑
main() {
    build_project
    clean_dst
    copy_files
    pack_files
    cleanup
    echo "Completed successfully! Archive: $ZIP_NAME"
}

main
