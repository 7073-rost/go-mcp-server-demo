package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// readFile 是 MCP Tool "read_file" 的处理函数。
// 它从 CallToolRequest 中提取 path 参数，读取本地文件并将内容以 TextContent 返回。
// 当文件不存在或读取失败时，通过 CallToolResult.IsError 告知调用方。
func readFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 校验 path 参数（使用 mcp-go 提供的便捷方法）
	path, err := request.RequireString("path")
	if err != nil || path == "" {
		return mcp.NewToolResultError("参数错误：缺少必填参数 'path'"), nil
	}

	// 2. 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return mcp.NewToolResultError(fmt.Sprintf("文件不存在：%s", path)), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("读取文件失败：%v", err)), nil
	}

	// 3. 包装成功结果
	return mcp.NewToolResultText(string(data)), nil
}

func main() {
	// 注册 read_file Tool 到 MCP Server
	s := server.NewMCPServer(
		"go-mcp-demo",
		"0.1.0",
	)

	tool := mcp.NewTool("read_file",
		mcp.WithDescription("读取本地文件并返回其内容"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("要读取的文件的绝对或相对路径"),
		),
	)

	s.AddTool(tool, readFile)

	// 以 stdio 方式启动（由 MCP 客户端驱动）
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP Server 运行错误: %v\n", err)
		os.Exit(1)
	}
}
