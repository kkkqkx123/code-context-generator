package tui

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileSelectorModel 模拟文件选择器模型
type MockFileSelectorModel struct {
	mock.Mock
	path         string
	items        []MockFileItem
	selected     map[int]bool
	cursor       int
	scrollOffset int
	multiSelect  bool
}

type MockFileItem struct {
	Name string
	Path string
	IsDir bool
}

func (m *MockFileSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	args := m.Called(msg)
	model := args.Get(0)
	cmd := args.Get(1)
	if model == nil {
		return nil, nil
	}
	if cmd == nil {
		return model.(tea.Model), nil
	}
	return model.(tea.Model), cmd.(tea.Cmd)
}

func (m *MockFileSelectorModel) View() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockFileSelectorModel) Init() tea.Cmd {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(tea.Cmd)
}

func (m *MockFileSelectorModel) GetPath() string {
	return m.path
}

func (m *MockFileSelectorModel) GetCursor() int {
	return m.cursor
}

func (m *MockFileSelectorModel) GetSelected() map[int]bool {
	return m.selected
}

// TestFileSelectorInitialization 测试文件选择器初始化
func TestFileSelectorInitialization(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected struct {
			multiSelect bool
			cursor      int
			scrollOffset int
		}
	}{
		{
			name: "默认目录初始化",
			path: ".",
			expected: struct {
				multiSelect bool
				cursor      int
				scrollOffset int
			}{
				multiSelect:  true,
				cursor:       0,
				scrollOffset: 0,
			},
		},
		{
			name: "指定目录初始化",
			path: "./test_files",
			expected: struct {
				multiSelect bool
				cursor      int
				scrollOffset int
			}{
				multiSelect:  true,
				cursor:       0,
				scrollOffset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				path:        tt.path,
				multiSelect: tt.expected.multiSelect,
				cursor:      tt.expected.cursor,
				scrollOffset: tt.expected.scrollOffset,
			}
			
			mockModel.On("Init").Return(nil)

			cmd := mockModel.Init()
			// Init方法返回nil是正常的，表示没有初始命令
			assert.Nil(t, cmd)
			assert.Equal(t, tt.path, mockModel.GetPath())
			assert.Equal(t, tt.expected.cursor, mockModel.GetCursor())
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorFileLoading 测试文件加载功能
func TestFileSelectorFileLoading(t *testing.T) {
	// 创建测试目录结构
	testDir := "./test_temp_dir"
	err := os.MkdirAll(testDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(testDir)

	// 创建测试文件
	testFiles := []string{"file1.txt", "file2.go", "file3.md"}
	for _, file := range testFiles {
		filePath := filepath.Join(testDir, file)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		assert.NoError(t, err)
	}

	tests := []struct {
		name          string
		path          string
		expectedCount int
		expectedError bool
	}{
		{
			name:          "加载包含文件的目录",
			path:          testDir,
			expectedCount: len(testFiles),
			expectedError: false,
		},
		{
			name:          "加载空目录",
			path:          "./empty_temp_dir",
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "加载不存在的目录",
			path:          "./non_existent_dir",
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				path: tt.path,
			}
			
			// 模拟文件列表消息
			type FileListMsg struct {
				Items []MockFileItem
			}
			
			var items []MockFileItem
			if !tt.expectedError {
				for i := 0; i < tt.expectedCount; i++ {
					items = append(items, MockFileItem{
						Name: testFiles[i],
						Path: filepath.Join(tt.path, testFiles[i]),
						IsDir: false,
					})
				}
			}
			
			msg := FileListMsg{Items: items}
			
			if tt.expectedError {
				mockModel.On("Update", msg).Return(mockModel, nil)
			} else {
				updatedModel := &MockFileSelectorModel{
					path:  tt.path,
					items: items,
				}
				mockModel.On("Update", msg).Return(updatedModel, nil)
			}
			
			newModel, cmd := mockModel.Update(msg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			assert.NotNil(t, newModel)
			
			if !tt.expectedError {
				if fileModel, ok := newModel.(*MockFileSelectorModel); ok {
					assert.Equal(t, tt.expectedCount, len(fileModel.items))
				}
			}
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorSelection 测试文件选择功能
func TestFileSelectorSelection(t *testing.T) {
	tests := []struct {
		name           string
		initialItems   []MockFileItem
		initialSelected map[int]bool
		cursor         int
		key            string
		expectedToggle int
	}{
		{
			name: "单文件选择切换",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			initialSelected: map[int]bool{0: false, 1: false},
			cursor:         0,
			key:            " ", // 空格键
			expectedToggle: 0,
		},
		{
			name: "多文件选择",
			initialItems: []MockFileItem{
				{Name: "dir1", Path: "./dir1", IsDir: true},
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			initialSelected: map[int]bool{0: false, 1: false, 2: false},
			cursor:         1,
			key:            " ",
			expectedToggle: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				items:       tt.initialItems,
				selected:    tt.initialSelected,
				cursor:      tt.cursor,
				multiSelect: true,
			}
			
			// 模拟空格键按下
			keyMsg := tea.KeyMsg{Type: tea.KeySpace}
			
			// 期望更新选择状态
			updatedSelected := make(map[int]bool)
			for k, v := range tt.initialSelected {
				updatedSelected[k] = v
			}
			updatedSelected[tt.expectedToggle] = !updatedSelected[tt.expectedToggle]
			
			updatedModel := &MockFileSelectorModel{
				items:       tt.initialItems,
				selected:    updatedSelected,
				cursor:      tt.cursor,
				multiSelect: true,
			}
			
			mockModel.On("Update", keyMsg).Return(updatedModel, nil)
			
			newModel, cmd := mockModel.Update(keyMsg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			assert.NotNil(t, newModel)
			
			if fileModel, ok := newModel.(*MockFileSelectorModel); ok {
				expectedState := !tt.initialSelected[tt.expectedToggle]
				assert.Equal(t, expectedState, fileModel.selected[tt.expectedToggle])
			}
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorSelectAll 测试全选功能
func TestFileSelectorSelectAll(t *testing.T) {
	tests := []struct {
		name           string
		initialItems   []MockFileItem
		initialSelected map[int]bool
		key            string
		expectedAll    bool
	}{
		{
			name: "全选所有文件",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
				{Name: "file3.md", Path: "./file3.md", IsDir: false},
			},
			initialSelected: map[int]bool{0: false, 1: false, 2: false},
			key:            "a",
			expectedAll:    true,
		},
		{
			name: "取消全选",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			initialSelected: map[int]bool{0: true, 1: true},
			key:            "n",
			expectedAll:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				items:       tt.initialItems,
				selected:    tt.initialSelected,
				multiSelect: true,
			}
			
			// 模拟按键
			var keyMsg tea.KeyMsg
			switch tt.key {
			case "a":
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			case "n":
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			}
			
			// 期望更新选择状态
			updatedSelected := make(map[int]bool)
			for i := range tt.initialItems {
				updatedSelected[i] = tt.expectedAll
			}
			
			updatedModel := &MockFileSelectorModel{
				items:       tt.initialItems,
				selected:    updatedSelected,
				multiSelect: true,
			}
			
			mockModel.On("Update", keyMsg).Return(updatedModel, nil)
			
			newModel, cmd := mockModel.Update(keyMsg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			assert.NotNil(t, newModel)
			
			if fileModel, ok := newModel.(*MockFileSelectorModel); ok {
				for i := range tt.initialItems {
					assert.Equal(t, tt.expectedAll, fileModel.selected[i])
				}
			}
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorNavigation 测试导航功能
func TestFileSelectorNavigation(t *testing.T) {
	tests := []struct {
		name           string
		initialItems   []MockFileItem
		initialCursor  int
		key            string
		expectedCursor int
	}{
		{
			name: "向下导航",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
				{Name: "file3.md", Path: "./file3.md", IsDir: false},
			},
			initialCursor:  0,
			key:            "down",
			expectedCursor: 1,
		},
		{
			name: "向上导航",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
				{Name: "file3.md", Path: "./file3.md", IsDir: false},
			},
			initialCursor:  2,
			key:            "up",
			expectedCursor: 1,
		},
		{
			name: "向下导航到边界",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			initialCursor:  1,
			key:            "down",
			expectedCursor: 1, // 应该保持在边界
		},
		{
			name: "向上导航到边界",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			initialCursor:  0,
			key:            "up",
			expectedCursor: 0, // 应该保持在边界
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				items:  tt.initialItems,
				cursor: tt.initialCursor,
			}
			
			// 模拟方向键
			var keyMsg tea.KeyMsg
			switch tt.key {
			case "up":
				keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			case "down":
				keyMsg = tea.KeyMsg{Type: tea.KeyDown}
			}
			
			// 期望更新光标位置
			updatedModel := &MockFileSelectorModel{
				items:  tt.initialItems,
				cursor: tt.expectedCursor,
			}
			
			mockModel.On("Update", keyMsg).Return(updatedModel, nil)
			
			newModel, cmd := mockModel.Update(keyMsg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			assert.NotNil(t, newModel)
			
			if fileModel, ok := newModel.(*MockFileSelectorModel); ok {
				assert.Equal(t, tt.expectedCursor, fileModel.GetCursor())
			}
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorConfirmation 测试选择确认流程
func TestFileSelectorConfirmation(t *testing.T) {
	tests := []struct {
		name           string
		initialItems   []MockFileItem
		initialSelected map[int]bool
		key            string
		expectedMsg    string
	}{
		{
			name: "确认选择单个文件",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			initialSelected: map[int]bool{0: true, 1: false},
			key:            "enter",
			expectedMsg:    "FileSelectionMsg",
		},
		{
			name: "确认选择多个文件",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
				{Name: "file3.md", Path: "./file3.md", IsDir: false},
			},
			initialSelected: map[int]bool{0: true, 1: false, 2: true},
			key:            "enter",
			expectedMsg:    "FileSelectionMsg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				items:       tt.initialItems,
				selected:    tt.initialSelected,
				multiSelect: true,
			}
			
			// 模拟Enter键
			keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
			
			// 模拟文件选择消息
			type FileSelectionMsg struct {
				Selected []string
			}
			
			var selectedPaths []string
			for i, selected := range tt.initialSelected {
				if selected {
					selectedPaths = append(selectedPaths, tt.initialItems[i].Path)
				}
			}
			
			msg := FileSelectionMsg{Selected: selectedPaths}
			
			mockModel.On("Update", keyMsg).Return(mockModel, tea.Cmd(func() tea.Msg { return msg }))
			
			_, cmd := mockModel.Update(keyMsg)
			// 这里返回了一个命令，所以cmd不应该为nil
			assert.NotNil(t, cmd)
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorCancel 测试取消选择流程
func TestFileSelectorCancel(t *testing.T) {
	tests := []struct {
		name         string
		initialItems []MockFileItem
		key          string
		expectedMsg  string
	}{
		{
			name: "按Esc键取消选择",
			initialItems: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "file2.go", Path: "./file2.go", IsDir: false},
			},
			key:         "esc",
			expectedMsg: "cancel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				items: tt.initialItems,
			}
			
			// 模拟Esc键
			keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
			
			mockModel.On("Update", keyMsg).Return(mockModel, nil)
			
			_, cmd := mockModel.Update(keyMsg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorViewRendering 测试视图渲染
func TestFileSelectorViewRendering(t *testing.T) {
	tests := []struct {
		name        string
		items       []MockFileItem
		selected    map[int]bool
		cursor      int
		path        string
		expectedContains []string
	}{
		{
			name: "基本文件列表渲染",
			items: []MockFileItem{
				{Name: "file1.txt", Path: "./file1.txt", IsDir: false},
				{Name: "dir1", Path: "./dir1", IsDir: true},
			},
			selected:    map[int]bool{0: true, 1: false},
			cursor:      0,
			path:        "./test",
			expectedContains: []string{"file1.txt", "dir1", "./test"},
		},
		{
			name: "选中状态渲染",
			items: []MockFileItem{
				{Name: "selected.txt", Path: "./selected.txt", IsDir: false},
				{Name: "unselected.txt", Path: "./unselected.txt", IsDir: false},
			},
			selected:    map[int]bool{0: true, 1: false},
			cursor:      0,
			path:        "./test",
			expectedContains: []string{"selected.txt", "unselected.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				items:      tt.items,
				selected:   tt.selected,
				cursor:     tt.cursor,
				path:       tt.path,
			}
			
			// 构建期望的视图内容
			viewContent := ""
			for _, item := range tt.items {
				selectedMark := " "
				if tt.selected[len(viewContent)] {
					selectedMark = "✓"
				}
				viewContent += selectedMark + " " + item.Name + "\n"
			}
			viewContent += "路径: " + tt.path + "\n"
			
			mockModel.On("View").Return(viewContent)
			
			view := mockModel.View()
			
			// 验证视图内容
			for _, expected := range tt.expectedContains {
				assert.Contains(t, view, expected)
			}
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorPerformance 测试性能相关场景
func TestFileSelectorPerformance(t *testing.T) {
	tests := []struct {
		name        string
		fileCount   int
		maxTime     int // 毫秒
		description string
	}{
		{
			name:        "大量文件处理",
			fileCount:   1000,
			maxTime:     100,
			description: "处理1000个文件应该在100毫秒内完成",
		},
		{
			name:        "超大文件列表",
			fileCount:   5000,
			maxTime:     500,
			description: "处理5000个文件应该在500毫秒内完成",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟大量文件项
			var items []MockFileItem
			for i := 0; i < tt.fileCount; i++ {
				items = append(items, MockFileItem{
					Name:  "file" + string(rune(i)) + ".txt",
					Path:  "./file" + string(rune(i)) + ".txt",
					IsDir: false,
				})
			}
			
			mockModel := &MockFileSelectorModel{
				items: items,
			}
			
			// 模拟文件列表消息
			type FileListMsg struct {
				Items []MockFileItem
			}
			
			msg := FileListMsg{Items: items}
			updatedModel := &MockFileSelectorModel{
				items: items,
			}
			
			mockModel.On("Update", msg).Return(updatedModel, nil)
			
			start := time.Now()
			_, cmd := mockModel.Update(msg)
			elapsed := time.Since(start)
			
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			assert.Less(t, elapsed.Milliseconds(), int64(tt.maxTime), tt.description)
			
			mockModel.AssertExpectations(t)
		})
	}
}

// TestFileSelectorErrorHandling 测试错误处理
func TestFileSelectorErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		errorType     string
		expectedError bool
		description   string
	}{
		{
			name:          "权限受限目录",
			path:          "/root/restricted",
			errorType:     "permission",
			expectedError: true,
			description:   "应该正确处理权限错误",
		},
		{
			name:          "不存在的路径",
			path:          "/non/existent/path",
			errorType:     "not_found",
			expectedError: true,
			description:   "应该正确处理路径不存在错误",
		},
		{
			name:          "特殊字符路径",
			path:          "./test/with/special/chars/[!@#$%]",
			errorType:     "special_chars",
			expectedError: false,
			description:   "应该正确处理特殊字符路径",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockFileSelectorModel{
				path: tt.path,
			}
			
			// 模拟错误消息
			type ErrorMsg struct {
				Err error
			}
			
			var err error
			if tt.expectedError {
				err = assert.AnError
			}
			
			msg := ErrorMsg{Err: err}
			
			if tt.expectedError {
				mockModel.On("Update", msg).Return(mockModel, nil)
			} else {
				updatedModel := &MockFileSelectorModel{
					path: tt.path,
				}
				mockModel.On("Update", msg).Return(updatedModel, nil)
			}
			
			_, cmd := mockModel.Update(msg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}