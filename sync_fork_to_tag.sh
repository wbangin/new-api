#!/bin/bash

# 遇到错误立即退出
set -e

# 定义变量
UPSTREAM_URL="https://github.com/QuantumNous/new-api.git"
BRANCH="main"

echo "=========================================================="
echo "开始同步 Fork 仓库 (保留本地 Commits，基于上游最新 Release Tag)"
echo "=========================================================="

# 1. 检查是否已经配置了 upstream 远程仓库
echo "[1/6] 检查 upstream 远程仓库配置..."
if ! git remote | grep -q "^upstream$"; then
    echo "未找到 upstream，正在添加: $UPSTREAM_URL"
    git remote add upstream "$UPSTREAM_URL"
else
    echo "upstream 已配置。"
fi

# 2. 从上游仓库拉取最新代码和所有的 Tags
echo "[2/6] 从 upstream 获取最新数据和 Tags..."
git fetch upstream --tags

# 3. 自动获取最新的 Release Tag
echo "[3/6] 正在寻找最新的 Release Tag..."
# 使用 -v:refname 按版本号倒序排列所有 tag，并取第一行作为最新 tag
LATEST_TAG=$(git tag --list --sort=-v:refname | head -n 1)

if [ -z "$LATEST_TAG" ]; then
    echo "❌ 错误: 在仓库中没有找到任何 Tag。请确认该项目是否发布过 Tag。"
    exit 1
fi
echo "💡 成功获取到最新 Tag: [ $LATEST_TAG ]"

# 4. 确保当前处于目标分支
echo "[4/6] 切换到目标分支: $BRANCH..."
git checkout $BRANCH

# 5. 执行变基 (Rebase)
# 这会将你的 2 个 commit 放到最新 Tag 对应的代码状态之上
echo "[5/6] 正在将本地 commits 变基到 Tag [ $LATEST_TAG ] 之上..."
git rebase "$LATEST_TAG"

# 6. 强制推送到你的个人远程仓库 (origin)
echo "[6/6] 强制推送到你的远程 Fork 仓库 (origin)..."
git push -f origin $BRANCH

echo "=========================================================="
echo "✅ 同步成功！"
echo "你的分支现在已基于上游的稳定版本 [ $LATEST_TAG ]，且保留了你的个人修改。"
echo "=========================================================="
