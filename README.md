# sprite-cuter 项目

sprite-cuter 是一个用于处理精灵图（Sprite Sheet）的工具，它包含一个 Go 语言编写的后端服务和一个基于 React 的前端界面。该项目旨在提供一个用户友好的界面，用于上传、处理和下载精灵图。

## 项目结构

项目主要分为 `backend` 和 `frontend` 两个部分。

```
sprite-cuter/
├── backend/             # 后端服务 (Go)
│   ├── controller/      # HTTP 请求控制器
│   ├── core/            # 核心业务逻辑 (例如精灵图处理)
│   ├── export/          # 导出文件存储目录
│   ├── go.mod
│   ├── go.sum
│   ├── main.go          # 后端入口文件
│   ├── router/          # 路由配置
│   ├── uploads/         # 上传文件存储目录
│   └── utils/           # 工具函数
├── frontend/            # 前端应用 (React)
│   ├── public/
│   ├── src/             # React 源代码
│   ├── package.json     # 前端依赖管理
│   └── vite.config.js   # Vite 配置
├── design_and_tasks.md  # 设计和任务文档
├── start.bat            # 启动脚本 (Windows)
└── README.md            # 项目说明文件
```

## 功能特性

- **文件上传**: 用户可以通过前端界面上传图片文件。
- **精灵图处理**: 后端服务能够处理上传的图片，可能包括切割、合并或生成精灵图。
- **文件下载**: 处理后的文件可以从后端下载。

## 安装与运行

### 前提条件

- Go 1.16+ (或更高版本)
- Node.js 14+ (或更高版本)
- npm 或 yarn

### 步骤

1. **克隆仓库**

   ```bash
   git clone <仓库地址>
   cd sprite-cuter
   ```

2. **启动后端服务**

   进入 `backend` 目录并运行：

   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

   后端服务默认运行在 `http://localhost:8080`。

3. **启动前端应用**

   进入 `frontend` 目录并安装依赖，然后启动开发服务器：

   ```bash
   cd frontend
   npm install  # 或者 yarn install
   npm run dev  # 或者 yarn dev
   ```

   前端应用默认运行在 `http://localhost:5173` (Vite 默认端口)。

4. **访问应用**

   在浏览器中打开 `http://localhost:5173` 即可访问 sprite-cuter 应用。

## 使用方法

1. 在前端界面上选择并上传您的图片文件。
2. 等待后端处理完成。
3. 下载处理后的精灵图或相关文件。

## 贡献

欢迎贡献！如果您有任何建议或发现 Bug，请随时提交 Issue 或 Pull Request。

## 许可证

本项目采用 [MIT 许可证](LICENSE) 发布。