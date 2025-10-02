package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"code-context-generator/pkg/types"
)

// TestShouldIncludeFileDebug 测试文件包含/排除逻辑（带调试输出）
func TestShouldIncludeFileDebug(t *testing.T) {
	// 创建临时测试目录结构
	tempDir := t.TempDir()
	
	// 创建测试文件结构
	testFiles := []string{
		"node_modules/package/index.js",
		"vendor/lib/file.go",
	}

	// 创建目录和文件
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("创建目录失败 %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("创建文件失败 %s: %v", file, err)
		}
	}

	walker := &FileSystemWalker{}

	// 测试 node_modules 排除
	nodeModulesPath := filepath.Join(tempDir, "node_modules/package/index.js")
	options := &types.WalkOptions{
		ExcludePatterns: []string{"node_modules/"},
	}
	
	fmt.Printf("\n=== 测试 node_modules 排除 ===\n")
	fmt.Printf("文件路径: %s\n", nodeModulesPath)
	fmt.Printf("根路径: %s\n", tempDir)
	fmt.Printf("排除模式: %v\n", options.ExcludePatterns)
	
	result := walker.shouldIncludeFile(nodeModulesPath, tempDir, options)
	fmt.Printf("结果: %v (期望: false)\n", result)
	
	if result != false {
		t.Errorf("node_modules 排除失败: got %v, want false", result)
	}

	// 测试 vendor 排除
	vendorPath := filepath.Join(tempDir, "vendor/lib/file.go")
	options2 := &types.WalkOptions{
		ExcludePatterns: []string{"vendor/"},
	}
	
	fmt.Printf("\n=== 测试 vendor 排除 ===\n")
	fmt.Printf("文件路径: %s\n", vendorPath)
	fmt.Printf("根路径: %s\n", tempDir)
	fmt.Printf("排除模式: %v\n", options2.ExcludePatterns)
	
	result2 := walker.shouldIncludeFile(vendorPath, tempDir, options2)
	fmt.Printf("结果: %v (期望: false)\n", result2)
	
	if result2 != false {
		t.Errorf("vendor 排除失败: got %v, want false", result2)
	}
}