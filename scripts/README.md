# Scripts

此目录包含了项目的维护脚本和 Git Hooks 配置，用于规范开发流程。

## 🎣 Git Hooks

位于 `githooks/templates` 目录下，提供了以下自动化检查：

- **pre-commit**: 分支名称检查。
  - 强制要求分支名符合 `feature/xxx`, `fix/xxx`, `hotfix/xxx` 或 `main`, `develop` 等规范。
- **commit-msg**: Commit Message 格式检查。
  - 遵循 Angular Commit Message 规范（如 `feat: add new feature`, `fix: resolve bug`）。

## ⚙️ 配置方法

你可以选择以下任一方式启用 Git Hooks：

### 方法一：配置 Git Hooks 路径 (推荐)

直接将 Git 的 hooks 路径指向此目录：

```bash
git config core.hooksPath scripts/githooks/templates
```

### 方法二：手动复制

将脚本复制到本地 `.git/hooks` 目录并添加执行权限：

```bash
cp scripts/githooks/templates/* .git/hooks/
chmod +x .git/hooks/*
```
