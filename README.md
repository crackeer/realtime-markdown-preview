# Markdown实时预览工具

> **AI生成说明**：本项目从设计到实现全程由AI生成，包括代码、配置和文档。

一个基于Go语言开发的实时Markdown预览工具，支持文件变化监听、实时HTML转换和WebSocket推送更新。

## 功能特性

- 📝 **实时预览**：监听Markdown文件变化，自动刷新预览
- 🌐 **WebSocket推送**：使用WebSocket实时向浏览器推送更新内容
- 🎨 **GitHub风格**：支持GitHub Flavored Markdown（表格、任务列表、删除线等）
- 🚀 **自动打开浏览器**：启动服务后自动打开默认浏览器
- 🔧 **自定义端口**：支持通过命令行参数指定服务端口
- 💻 **跨平台支持**：支持Windows、macOS和Linux系统

## 技术栈

- **后端**：Go 1.16+
- **Web框架**：Gin
- **WebSocket**：Gorilla WebSocket
- **Markdown解析**：Goldmark
- **文件监听**：fsnotify
- **前端**：原生HTML/CSS/JavaScript

## 安装

### 方法1：使用Go命令直接运行

```bash
go run main.go your-markdown-file.md
```

### 方法2：编译后运行

```bash
# 编译
go build -o markdown-preview

# 运行
./markdown-preview your-markdown-file.md
```

## 使用说明

### 基本使用

```bash
go run main.go example.md
```

这将启动服务并自动打开浏览器，实时预览`example.md`文件的内容。

### 自定义端口

```bash
go run main.go -port 8081 example.md
```

### 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-port` | 指定服务端口 | 8080 |
| `-client` | 启动WebSocket客户端（用于测试） | false |
| `file_path` | Markdown文件路径（必填） | - |

## 项目结构

```
.
├── main.go          # 主程序入口
├── markdown.go      # Markdown转换功能
├── watcher.go       # 文件监听功能
├── websocket.go     # WebSocket管理功能
├── static/          # 静态文件目录
│   └── index.html   # 前端HTML页面
├── go.mod           # Go模块文件
├── go.sum           # Go依赖校验文件
└── README.md        # 项目说明文档
```

## 核心功能实现

### 1. 文件监听

使用`fsnotify`库监听Markdown文件的变化，当文件被修改时触发回调函数。

### 2. Markdown转换

使用`goldmark`库将Markdown内容转换为HTML，支持GitHub Flavored Markdown扩展。

### 3. WebSocket通信

- 服务器端管理多个WebSocket连接
- 客户端连接后立即获取初始HTML内容
- 文件变化时，服务器向所有客户端广播更新后的HTML

### 4. 自动浏览器打开

根据不同操作系统调用相应的命令打开浏览器：
- macOS: `open`命令
- Windows: `cmd /c start`命令
- Linux: `xdg-open`命令

## 浏览器兼容性

支持所有现代浏览器，包括：
- Chrome 57+
- Firefox 52+
- Safari 10.1+
- Edge 79+

## 示例

### 启动服务

```bash
go run main.go test.md
```

### 输出日志

```
2024/01/01 12:00:00 启动Markdown实时预览服务，监听文件: test.md，端口: 8080
2024/01/01 12:00:00 WebSocket服务器已启动
```

## 开发说明

### 依赖管理

```bash
go mod tidy
```

### 测试文件

项目包含两个测试Markdown文件：
- `test.md`：基本Markdown语法测试
- `test_table.md`：表格功能测试

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 联系方式

如有问题或建议，欢迎通过GitHub Issues反馈。