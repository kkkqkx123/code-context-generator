package autocomplete

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"code-context-generator/pkg/types"
)

func TestNewAutocompleter(t *testing.T) {
	// 测试默认配置
	ac := NewAutocompleter(nil)
	if ac == nil {
		t.Fatal("NewAutocompleter returned nil")
	}

	// 测试自定义配置
	config := &types.AutocompleteConfig{
		Enabled:        true,
		MinChars:       3,
		MaxSuggestions: 10,
	}
	ac = NewAutocompleter(config)
	if ac == nil {
		t.Fatal("NewAutocompleter with config returned nil")
	}
}

func TestFilePathAutocompleter_Complete(t *testing.T) {
	// 创建临时目录结构
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建子目录文件
	subFile := filepath.Join(subDir, "subtest.go")
	if err := os.WriteFile(subFile, []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	config := &types.AutocompleteConfig{
		Enabled:        true,
		MinChars:       1,
		MaxSuggestions: 5,
	}
	ac := NewAutocompleter(config).(*FilePathAutocompleter)

	tests := []struct {
		name        string
		input       string
		context     *types.CompleteContext
		wantResults bool
		minResults  int
	}{
		{
			name:        "complete file path",
			input:       filepath.Join(tempDir, "tes"),
			context:     &types.CompleteContext{Type: types.CompleteFilePath},
			wantResults: true,
			minResults:  1,
		},
		{
			name:        "complete directory",
			input:       tempDir,
			context:     &types.CompleteContext{Type: types.CompleteDirectory},
			wantResults: true,
			minResults:  1,
		},
		{
			name:        "complete extension",
			input:       ".g",
			context:     &types.CompleteContext{Type: types.CompleteExtension},
			wantResults: true,
			minResults:  1,
		},
		{
			name:        "complete pattern",
			input:       filepath.Join(tempDir, "*.txt"),
			context:     &types.CompleteContext{Type: types.CompletePattern},
			wantResults: true,
			minResults:  1,
		},
		{
			name:        "complete generic",
			input:       filepath.Join(tempDir, "sub"),
			context:     &types.CompleteContext{Type: types.CompleteGeneric},
			wantResults: true,
			minResults:  1,
		},
		{
			name:        "disabled autocompleter",
			input:       tempDir,
			context:     &types.CompleteContext{Type: types.CompleteGeneric},
			wantResults: false,
			minResults:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "disabled autocompleter" {
				ac.config.Enabled = false
				defer func() { ac.config.Enabled = true }()
			}

			results, err := ac.Complete(tt.input, tt.context)
			if err != nil && tt.wantResults {
				t.Errorf("Complete() error = %v, wantResults %v", err, tt.wantResults)
				return
			}

			if tt.wantResults && len(results) < tt.minResults {
				t.Errorf("Complete() got %d results, want at least %d", len(results), tt.minResults)
			}

			if !tt.wantResults && len(results) > 0 {
				t.Errorf("Complete() got %d results, want 0", len(results))
			}
		})
	}
}

func TestFilePathAutocompleter_GetSuggestions(t *testing.T) {
	config := &types.AutocompleteConfig{
		Enabled:        true,
		MinChars:       1,
		MaxSuggestions: 3,
	}
	ac := NewAutocompleter(config).(*FilePathAutocompleter)

	// 添加一些缓存数据
	ac.cache["test"] = []string{"test1", "test2", "testing"}

	tests := []struct {
		name           string
		input          string
		maxSuggestions int
		wantCount      int
	}{
		{
			name:           "get suggestions with max limit",
			input:          "test",
			maxSuggestions: 2,
			wantCount:      2,
		},
		{
			name:           "get suggestions without max limit",
			input:          "test",
			maxSuggestions: 0,
			wantCount:      3,
		},
		{
			name:           "no matching suggestions",
			input:          "nomatch",
			maxSuggestions: 5,
			wantCount:      0,
		},
		{
			name:           "disabled autocompleter",
			input:          "test",
			maxSuggestions: 5,
			wantCount:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "disabled autocompleter" {
				ac.config.Enabled = false
				defer func() { ac.config.Enabled = true }()
			}

			suggestions := ac.GetSuggestions(tt.input, tt.maxSuggestions)
			if len(suggestions) != tt.wantCount {
				t.Errorf("GetSuggestions() got %d suggestions, want %d", len(suggestions), tt.wantCount)
			}
		})
	}
}

func TestFilePathAutocompleter_CacheOperations(t *testing.T) {
	config := &types.AutocompleteConfig{
		Enabled:        true,
		MinChars:       1,
		MaxSuggestions: 5,
	}
	ac := NewAutocompleter(config).(*FilePathAutocompleter)

	// 测试更新缓存
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := ac.UpdateCache(tempDir); err != nil {
		t.Errorf("UpdateCache() error = %v", err)
	}

	// 验证缓存已更新
	if size := ac.GetCacheSize(); size != 1 {
		t.Errorf("GetCacheSize() = %d, want 1", size)
	}

	// 测试清除缓存
	ac.ClearCache()
	if size := ac.GetCacheSize(); size != 0 {
		t.Errorf("GetCacheSize() after ClearCache() = %d, want 0", size)
	}
}

func TestCommandAutocompleter(t *testing.T) {
	cmdAc := NewCommandAutocompleter()

	// 注册测试命令
	cmdInfo := &CommandInfo{
		Name:        "test",
		Description: "Test command",
		Aliases:     []string{"t", "tst"},
		Subcommands: []string{"sub1", "sub2"},
		Options:     []string{"--help", "--version"},
	}
	cmdAc.RegisterCommand(cmdInfo)

	tests := []struct {
		name     string
		input    string
		wantLen  int
		contains string
	}{
		{
			name:     "complete command name",
			input:    "te",
			wantLen:  1,
			contains: "test",
		},
		{
			name:     "complete command alias",
			input:    "t",
			wantLen:  3, // 会匹配 test, t, tst
			contains: "t",
		},
		{
			name:     "no match",
			input:    "nomatch",
			wantLen:  0,
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := cmdAc.Complete(tt.input)
			if len(results) != tt.wantLen {
				t.Errorf("Complete() = %d results, want %d", len(results), tt.wantLen)
			}
			if tt.contains != "" && len(results) > 0 {
				found := false
				for _, result := range results {
					if result == tt.contains {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Complete() results do not contain %s", tt.contains)
				}
			}
		})
	}

	// 测试获取命令信息
	if info, exists := cmdAc.GetCommandInfo("test"); !exists || info.Name != "test" {
		t.Errorf("GetCommandInfo() failed to retrieve command info")
	}
}

func TestCompositeSuggestionProvider(t *testing.T) {
	// 创建模拟的建议提供者
	mockProvider1 := &mockSuggestionProvider{
		suggestions: []Suggestion{
			{Text: "suggestion1", Description: "First suggestion"},
			{Text: "suggestion2", Description: "Second suggestion"},
		},
	}

	mockProvider2 := &mockSuggestionProvider{
		suggestions: []Suggestion{
			{Text: "suggestion2", Description: "Duplicate suggestion"},
			{Text: "suggestion3", Description: "Third suggestion"},
		},
	}

	composite := NewCompositeSuggestionProvider(mockProvider1, mockProvider2)
	context := &types.CompleteContext{Type: types.CompleteGeneric}

	suggestions, err := composite.GetSuggestions("test", context)
	if err != nil {
		t.Errorf("GetSuggestions() error = %v", err)
	}

	// 应该去重，所以期望3个建议
	if len(suggestions) != 3 {
		t.Errorf("GetSuggestions() = %d suggestions, want 3", len(suggestions))
	}
}

func TestAutocompleterOptions(t *testing.T) {
	opts := AutocompleterOptions{
		Enabled:        true,
		MinChars:       2,
		MaxSuggestions: 10,
		CacheSize:      100,
		Timeout:        5 * time.Second,
	}

	if !opts.Enabled {
		t.Error("AutocompleterOptions.Enabled should be true")
	}
	if opts.MinChars != 2 {
		t.Errorf("AutocompleterOptions.MinChars = %d, want 2", opts.MinChars)
	}
	if opts.MaxSuggestions != 10 {
		t.Errorf("AutocompleterOptions.MaxSuggestions = %d, want 10", opts.MaxSuggestions)
	}
	if opts.CacheSize != 100 {
		t.Errorf("AutocompleterOptions.CacheSize = %d, want 100", opts.CacheSize)
	}
	if opts.Timeout != 5*time.Second {
		t.Errorf("AutocompleterOptions.Timeout = %v, want 5s", opts.Timeout)
	}
}

// 模拟建议提供者用于测试
type mockSuggestionProvider struct {
	suggestions []Suggestion
	err         error
}

func (m *mockSuggestionProvider) GetSuggestions(input string, context *types.CompleteContext) ([]Suggestion, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.suggestions, nil
}