package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// TestReadFileTool_Success 测试 read_file 工具成功读取存在的本地文件。
// 运行方式：go test -v -run TestReadFileTool_Success
func TestReadFileTool_Success(t *testing.T) {
	// 创建临时测试文件
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "hello.txt")
	expectedContent := "Hello, MCP World!\n这是一个测试文件。"

	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("准备临时文件失败: %v", err)
	}

	// 构造 CallToolRequest，模拟客户端传入 path 参数
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "read_file",
			Arguments: map[string]any{"path": tmpFile},
		},
	}

	// 调用待测的 Tool Handler
	result, err := readFile(context.Background(), req)
	if err != nil {
		t.Fatalf("readFile 返回了意外的错误: %v", err)
	}
	if result == nil {
		t.Fatal("readFile 返回了 nil 结果")
	}
	if result.IsError {
		t.Fatalf("预期成功但返回了 IsError=true, content=%v", result.Content)
	}

	// 验证返回内容
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("预期 TextContent, 实际类型: %T", result.Content[0])
	}
	if textContent.Text != expectedContent {
		t.Errorf("文件内容不匹配\n预期: %q\n实际: %q", expectedContent, textContent.Text)
	}
}

// TestReadFileTool_FileNotFound 测试 read_file 工具在文件不存在时返回错误信息。
// 运行方式：go test -v -run TestReadFileTool_FileNotFound
func TestReadFileTool_FileNotFound(t *testing.T) {
	// 使用一个不存在的路径
	nonExistentPath := filepath.Join(t.TempDir(), "definitely_not_exists_12345.txt")

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "read_file",
			Arguments: map[string]any{"path": nonExistentPath},
		},
	}

	result, err := readFile(context.Background(), req)
	if err != nil {
		// 即使底层返回 error，也允许 handler 自行包装在 result 中
		t.Logf("readFile 返回 error（允许）: %v", err)
	}

	if result == nil {
		t.Fatal("readFile 返回了 nil 结果，预期包含错误信息")
	}

	// 核心断言：IsError 必须为 true，且 content 里必须有明确的错误描述
	if !result.IsError {
		t.Fatalf("预期 IsError=true, 实际为 false")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("预期 TextContent, 实际类型: %T", result.Content[0])
	}
	if !strings.Contains(textContent.Text, "不存在") && !strings.Contains(textContent.Text, "not exist") && !strings.Contains(textContent.Text, "no such") {
		t.Errorf("错误信息缺少关键提示词，实际内容: %q", textContent.Text)
	}
}

// TestReadFileTool_PathMissing 测试调用时未提供 path 参数的场景。
// 运行方式：go test -v -run TestReadFileTool_PathMissing
func TestReadFileTool_PathMissing(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "read_file",
			Arguments: map[string]any{}, // 空参数
		},
	}

	result, err := readFile(context.Background(), req)
	if err != nil {
		t.Logf("readFile 返回 error（允许）: %v", err)
	}
	if result == nil {
		t.Fatal("readFile 返回了 nil 结果")
	}
	if !result.IsError {
		t.Fatal("预期因缺少 path 参数而返回 IsError=true")
	}
}
