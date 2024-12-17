#!/bin/bash
export GOOS=linux
export GOARCH=${1:-amd64}
go build -o x-main

# 设置变量
TARGET_DIR="x-ui"
ARCHIVE_NAME="x-ui-${GOOS}-${GOARCH}.tar.gz"

# 清理并创建目标文件夹
if [ -d "$TARGET_DIR" ]; then
    echo "Cleaning up existing $TARGET_DIR folder..."
    rm -rf "$TARGET_DIR"
fi

echo "Creating folder: $TARGET_DIR"
mkdir -p "$TARGET_DIR"

# 复制文件到目标文件夹
cp x-main "$TARGET_DIR/x-ui"
cp x-ui.sh "$TARGET_DIR/"
cp x-ui.service "$TARGET_DIR/"
cp -rf  bin "$TARGET_DIR/"

# 打包目标文件夹
echo "Packing folder: $ARCHIVE_NAME"
tar -czf "$ARCHIVE_NAME" "$TARGET_DIR"

if [ $? -eq 0 ]; then
    echo "Successfully packing: $ARCHIVE_NAME"
else
    echo "Error creating archive!"
    exit 1
fi

# 清理文件
echo "Cleaning up folder: $TARGET_DIR"
rm -rf x-main
rm -rf "$TARGET_DIR"

echo "completed!"