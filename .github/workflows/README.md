# GitHub Actions Workflows

本目录包含项目的 GitHub Actions 工作流配置文件。

## 📋 工作流列表

### 1. `ghcr.yml` - GHCR 镜像构建与推送

自动构建 Docker 镜像并推送到 GitHub Container Registry (GHCR)。

**触发条件：**
- 推送到 `main` 或 `master` 分支
- 创建新的版本标签 (如 `v1.0.0`)
- 手动触发（通过 GitHub Actions 页面）

**功能特性：**
- ✅ 多架构支持（amd64 和 arm64）
- ✅ 自动版本标签生成
- ✅ 构建缓存优化
- ✅ 自动推送到 GHCR

**生成的镜像标签：**
- `ghcr.io/用户名/仓库名:latest` - 默认分支的最新版本
- `ghcr.io/用户名/仓库名:main` - main 分支
- `ghcr.io/用户名/仓库名:v1.0.0` - 版本标签
- `ghcr.io/用户名/仓库名:1.0` - 主次版本
- `ghcr.io/用户名/仓库名:1` - 主版本
- `ghcr.io/用户名/仓库名:main-abc1234` - 分支+提交哈希

### 2. `docker-image-alpha.yml` - Alpha 版本构建

自动构建和发布 Alpha 测试版本到 Docker Hub 和 GHCR。

**触发条件：**
- 推送到 `alpha` 分支
- 手动触发

### 3. 其他工作流

- `docker-image-arm64.yml` - ARM64 架构专用构建
- `release.yml` - 正式版本发布
- `electron-build.yml` - Electron 应用构建
- `sync-to-gitee.yml` - 同步到 Gitee

## 🚀 使用说明

### 拉取 GHCR 镜像

```bash
# 拉取最新版本
docker pull ghcr.io/用户名/仓库名:latest

# 拉取指定版本
docker pull ghcr.io/用户名/仓库名:v1.0.0

# 拉取指定分支
docker pull ghcr.io/用户名/仓库名:main
```

### 手动触发构建

1. 进入 GitHub 仓库
2. 点击 `Actions` 标签
3. 选择 `Build and Push to GHCR` 工作流
4. 点击 `Run workflow` 按钮
5. 选择分支并点击运行

### 发布新版本

创建并推送版本标签即可自动触发构建：

```bash
# 创建版本标签
git tag v1.0.0

# 推送标签
git push origin v1.0.0
```

## 🔧 配置要求

### 必需的权限

工作流需要以下权限（已在配置中设置）：
- `contents: read` - 读取仓库内容
- `packages: write` - 写入 GitHub Packages

### GitHub Token

工作流使用 `${{ secrets.GITHUB_TOKEN }}` 自动认证，无需额外配置。此 token 由 GitHub Actions 自动提供。

### 可选配置

如果需要推送到 Docker Hub，需要设置以下 Secrets：
- `DOCKERHUB_USERNAME` - Docker Hub 用户名
- `DOCKERHUB_TOKEN` - Docker Hub 访问令牌

## 📦 镜像使用示例

### Docker Compose

```yaml
version: '3.8'

services:
  new-api:
    image: ghcr.io/用户名/仓库名:latest
    ports:
      - "3000:3000"
    volumes:
      - ./data:/data
    environment:
      - SESSION_SECRET=your-secret
    restart: unless-stopped
```

### Docker Run

```bash
docker run -d \
  --name new-api \
  -p 3000:3000 \
  -v $(pwd)/data:/data \
  -e SESSION_SECRET=your-secret \
  ghcr.io/用户名/仓库名:latest
```

## 🔍 查看镜像信息

访问 GitHub Packages 页面查看所有可用镜像：
```
https://github.com/用户名/仓库名/pkgs/container/仓库名
```

## 📝 注意事项

1. **公开访问**：默认情况下，GHCR 镜像是公开的，任何人都可以拉取
2. **私有仓库**：如果需要私有镜像，需要在仓库设置中配置
3. **存储限额**：GitHub 提供免费的存储空间，但有使用限制
4. **多架构支持**：镜像支持 amd64 和 arm64 架构，会自动选择合适的架构

## 🐛 故障排除

### 构建失败

1. 检查 Actions 日志查看详细错误信息
2. 确认 Dockerfile 语法正确
3. 验证所有依赖文件都已提交

### 推送失败

1. 确认仓库设置中启用了 GitHub Actions
2. 检查工作流权限设置
3. 验证 `GITHUB_TOKEN` 有足够的权限

### 拉取镜像失败

1. 确认镜像标签正确
2. 对于私有镜像，需要先登录：
   ```bash
   echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
   ```

## 📚 相关文档

- [GitHub Actions 文档](https://docs.github.com/actions)
- [GitHub Container Registry 文档](https://docs.github.com/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Docker Buildx 文档](https://docs.docker.com/buildx/working-with-buildx/)