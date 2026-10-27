#!/usr/bin/env bash
# 用法：bash scripts/git-commit-merge.sh "本次提交备注"
set -euo pipefail

if [ $# -eq 0 ] || [[ "$*" =~ ^[[:space:]]*$ ]]; then
  echo '用法: bash scripts/git-commit-merge.sh "本次提交备注"' >&2
  exit 1
fi

MESSAGE="$*"
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd -- "$SCRIPT_DIR/.."

DEV_BRANCH="dev"
MAIN_BRANCH="main"

# 在暂存或提交前检查分支与远端配置。
git show-ref --verify --quiet "refs/heads/$DEV_BRANCH"
git show-ref --verify --quiet "refs/heads/$MAIN_BRANCH"
git remote get-url origin > /dev/null

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ "$CURRENT_BRANCH" != "$DEV_BRANCH" ]; then
  echo "当前在 $CURRENT_BRANCH，切换到 $DEV_BRANCH"
  git checkout "$DEV_BRANCH"
fi

if [ -n "$(git status --porcelain)" ]; then
  git add -A
  git commit -m "$MESSAGE"
else
  echo "工作区干净，无需提交"
fi

git push origin "$DEV_BRANCH"

git checkout "$MAIN_BRANCH"
git merge --no-edit "$DEV_BRANCH"
git push origin "$MAIN_BRANCH"

git checkout "$DEV_BRANCH"
echo "完成：$DEV_BRANCH 已提交并合并到 ${MAIN_BRANCH}，当前在 $DEV_BRANCH 分支"
