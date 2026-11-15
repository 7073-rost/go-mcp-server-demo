# Go AI-Native MCP Server Demo

这是一个使用 Go 语言编写的极简 MCP (Model Context Protocol) Server。该项目展示了如何通过 **AI IDE (Trae)** 结合 **TDD（测试驱动开发）** 与 **Context Engineering（上下文工程）** 高效、规范地完成研发。

## ✨ 核心亮点
- **AI-Native 开发**: 全程使用 Trae 内置的 Seed-Code 模型辅助生成测试与实现。
- **TDD 驱动**: 严格遵循 Red-Green-Refactor 流程，覆盖 `read_file` Tool 的成功、文件不存在、参数缺失等边界场景。
- **上下文工程**: 配置了 `.trae/rules/go-style.md`，强制约束 AI 输出符合 Go 官方规范（如显式错误处理）的高质量代码。
- **工程护栏 (Harness)**: 集成 GitHub Actions，每次 push 自动运行 `go test -v ./...` 确保主分支代码安全。

## 🛠️ 快速开始
1. 拉取依赖：`go mod tidy`
2. 运行测试：`go test -v`
3. 构建：`go build -o mcp-server.exe`

## 🔌 接入 AI 客户端 (如 Claude Desktop)
在 `claude_desktop_config.json` 中添加：
```json
{
  "mcpServers": {
    "go-file-reader": {
      "command": "C:\\projects\\go-mcp-demo\\mcp-server.exe"
    }
  }
}