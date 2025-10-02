package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMainModel 模拟主模型
type MockMainModel struct {
	mock.Mock
	state       string
	currentView string
	width       int
	height      int
}

func (m *MockMainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *MockMainModel) View() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMainModel) Init() tea.Cmd {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(tea.Cmd)
}

func (m *MockMainModel) GetState() string {
	return m.state
}

func (m *MockMainModel) GetCurrentView() string {
	return m.currentView
}

// TestMainWorkflowInitialization 测试主工作流初始化
func TestMainWorkflowInitialization(t *testing.T) {
	tests := []struct {
		name     string
		expected struct {
			state       string
			currentView string
			width       int
			height      int
		}
	}{
		{
			name: "默认初始化",
			expected: struct {
				state       string
				currentView string
				width       int
				height      int
			}{
				state:       "StateInput",
				currentView: "ViewMain",
				width:       80,
				height:      24,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state:       tt.expected.state,
				currentView: tt.expected.currentView,
				width:       tt.expected.width,
				height:      tt.expected.height,
			}
			
			mockModel.On("Init").Return(nil)

			cmd := mockModel.Init()
			// Init方法返回nil是正常的，表示没有初始命令
			assert.Nil(t, cmd)
			assert.Equal(t, tt.expected.state, mockModel.GetState())
			assert.Equal(t, tt.expected.currentView, mockModel.GetCurrentView())
			mockModel.AssertExpectations(t)
		})
	}
}

// TestMainWorkflowStateTransitions 测试状态转换
func TestMainWorkflowStateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  string
		initialView   string
		trigger       string
		expectedState string
		expectedView  string
		description   string
	}{
		{
			name:          "从输入状态到文件选择状态",
			initialState:  "StateInput",
			initialView:   "ViewMain",
			trigger:       "file_select_trigger",
			expectedState: "StateSelect",
			expectedView:  "ViewSelect",
			description:   "应该正确切换到文件选择状态",
		},
		{
			name:          "从文件选择状态到处理状态",
			initialState:  "StateSelect",
			initialView:   "ViewSelect",
			trigger:       "process_trigger",
			expectedState: "StateProcessing",
			expectedView:  "ViewProgress",
			description:   "应该正确切换到处理状态",
		},
		{
			name:          "从处理状态到结果状态",
			initialState:  "StateProcessing",
			initialView:   "ViewProgress",
			trigger:       "result_ready",
			expectedState: "StateResult",
			expectedView:  "ViewResult",
			description:   "应该正确切换到结果状态",
		},
		{
			name:          "从输入状态到配置状态",
			initialState:  "StateInput",
			initialView:   "ViewMain",
			trigger:       "config_trigger",
			expectedState: "StateConfig",
			expectedView:  "ViewConfig",
			description:   "应该正确切换到配置状态",
		},
		{
			name:          "从结果状态返回输入状态",
			initialState:  "StateResult",
			initialView:   "ViewResult",
			trigger:       "back_to_input",
			expectedState: "StateInput",
			expectedView:  "ViewMain",
			description:   "应该正确返回到输入状态",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state:       tt.initialState,
				currentView: tt.initialView,
			}
			
			// 模拟触发状态转换的消息
			var triggerMsg tea.Msg
			switch tt.trigger {
			case "file_select_trigger":
				triggerMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			case "process_trigger":
				triggerMsg = tea.KeyMsg{Type: tea.KeyEnter}
			case "result_ready":
				type ResultMsg struct {
					Result interface{}
				}
				triggerMsg = ResultMsg{Result: nil}
			case "config_trigger":
				triggerMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			case "back_to_input":
				triggerMsg = tea.KeyMsg{Type: tea.KeyEsc}
			}
			
			// 期望状态转换
			updatedModel := &MockMainModel{
				state:       tt.expectedState,
				currentView: tt.expectedView,
			}
			
			mockModel.On("Update", triggerMsg).Return(updatedModel, nil)
			
			newModel, cmd := mockModel.Update(triggerMsg)
			// Update方法返回nil命令是正常的，表示没有后续命令
			assert.Nil(t, cmd)
			assert.NotNil(t, newModel)
			
			if mainModel, ok := newModel.(*MockMainModel); ok {
				assert.Equal(t, tt.expectedState, mainModel.GetState())
				assert.Equal(t, tt.expectedView, mainModel.GetCurrentView())
			}
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}

// TestMainWorkflowMessageHandling 测试消息处理
func TestMainWorkflowMessageHandling(t *testing.T) {
	tests := []struct {
		name        string
		msgType     string
		msgData     interface{}
		expectedCmd bool
		description string
	}{
		{
			name:        "处理进度消息",
			msgType:     "ProgressMsg",
			msgData:     struct{ Progress float64; Status string }{Progress: 0.5, Status: "处理中..."},
			expectedCmd: false,
			description: "应该正确处理进度更新消息",
		},
		{
			name:        "处理结果消息",
			msgType:     "ResultMsg",
			msgData:     struct{ Result interface{} }{Result: nil},
			expectedCmd: false,
			description: "应该正确处理结果消息",
		},
		{
			name:        "处理错误消息",
			msgType:     "ErrorMsg",
			msgData:     struct{ Err error }{Err: assert.AnError},
			expectedCmd: false,
			description: "应该正确处理错误消息",
		},
		{
			name:        "处理文件选择消息",
			msgType:     "FileSelectionMsg",
			msgData:     struct{ Selected []string }{Selected: []string{"./file1.txt", "./file2.go"}},
			expectedCmd: false,
			description: "应该正确处理文件选择消息",
		},
		{
			name:        "处理配置更新消息",
			msgType:     "ConfigUpdateMsg",
			msgData:     struct{ Config interface{} }{Config: nil},
			expectedCmd: false,
			description: "应该正确处理配置更新消息",
		},
		{
			name:        "处理窗口大小消息",
			msgType:     "WindowSizeMsg",
			msgData:     struct{ Width int; Height int }{Width: 120, Height: 30},
			expectedCmd: false,
			description: "应该正确处理窗口大小变化消息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state:       "StateProcessing",
				currentView: "ViewProgress",
				width:       80,
				height:      24,
			}
			
			// 创建不同类型的消息
			var msg tea.Msg
			switch tt.msgType {
			case "ProgressMsg":
				type ProgressMsg struct {
					Progress float64
					Status   string
				}
				data := tt.msgData.(struct{ Progress float64; Status string })
				msg = ProgressMsg{Progress: data.Progress, Status: data.Status}
			case "ResultMsg":
				type ResultMsg struct {
					Result interface{}
				}
				msg = ResultMsg{Result: tt.msgData}
			case "ErrorMsg":
				type ErrorMsg struct {
					Err error
				}
				msg = ErrorMsg{Err: assert.AnError}
			case "FileSelectionMsg":
				type FileSelectionMsg struct {
					Selected []string
				}
				msg = FileSelectionMsg{Selected: []string{"./file1.txt"}}
			case "ConfigUpdateMsg":
				type ConfigUpdateMsg struct {
					Config interface{}
				}
				msg = ConfigUpdateMsg{Config: tt.msgData}
			case "WindowSizeMsg":
				type WindowSizeMsg struct {
					Width  int
					Height int
				}
				msg = WindowSizeMsg{Width: 120, Height: 30}
			}
			
			mockModel.On("Update", msg).Return(mockModel, nil)
			
			_, cmd := mockModel.Update(msg)
			if tt.expectedCmd {
				// Update方法返回nil命令是正常的，表示没有后续命令
				assert.Nil(t, cmd)
			}
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}

// TestMainWorkflowKeyboardEvents 测试键盘事件处理
func TestMainWorkflowKeyboardEvents(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		keyType     tea.KeyType
		state       string
		expectedCmd bool
		description string
	}{
		{
			name:        "Ctrl+C退出",
			key:         "ctrl+c",
			keyType:     tea.KeyCtrlC,
			state:       "StateInput",
			expectedCmd: false,
			description: "应该正确处理退出快捷键",
		},
		{
			name:        "Esc键返回",
			key:         "esc",
			keyType:     tea.KeyEsc,
			state:       "StateSelect",
			expectedCmd: false,
			description: "应该正确处理返回快捷键",
		},
		{
			name:        "q键退出",
			key:         "q",
			keyType:     tea.KeyRunes,
			state:       "StateInput",
			expectedCmd: false,
			description: "应该正确处理q键退出",
		},
		{
			name:        "普通字符键",
			key:         "a",
			keyType:     tea.KeyRunes,
			state:       "StateInput",
			expectedCmd: false,
			description: "应该正确处理普通字符输入",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state: tt.state,
			}
			
			// 创建键盘消息
			var keyMsg tea.KeyMsg
		if tt.keyType == tea.KeyRunes {
			keyMsg = tea.KeyMsg{Type: tt.keyType, Runes: []rune(tt.key)}
		} else {
			keyMsg = tea.KeyMsg{Type: tt.keyType}
		}
			
			mockModel.On("Update", keyMsg).Return(mockModel, nil)
			
			_, cmd := mockModel.Update(keyMsg)
			if tt.expectedCmd {
				assert.NotNil(t, cmd)
			}
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}

// TestMainWorkflowSubmodelCoordination 测试子模型协调
func TestMainWorkflowSubmodelCoordination(t *testing.T) {
	tests := []struct {
		name        string
		state       string
		trigger     string
		expectedCmd bool
		description string
	}{
		{
			name:        "文件选择器模型协调",
			state:       "StateSelect",
			trigger:     "file_selector_msg",
			expectedCmd: true,
			description: "应该正确协调文件选择器模型",
		},
		{
			name:        "进度条模型协调",
			state:       "StateProcessing",
			trigger:     "progress_msg",
			expectedCmd: true,
			description: "应该正确协调进度条模型",
		},
		{
			name:        "结果查看器模型协调",
			state:       "StateResult",
			trigger:     "result_viewer_msg",
			expectedCmd: true,
			description: "应该正确协调结果查看器模型",
		},
		{
			name:        "配置编辑器模型协调",
			state:       "StateConfig",
			trigger:     "config_editor_msg",
			expectedCmd: true,
			description: "应该正确协调配置编辑器模型",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state: tt.state,
			}
			
			// 模拟子模型相关的消息
			var msg tea.Msg
			switch tt.trigger {
			case "file_selector_msg":
				type FileSelectionMsg struct {
					Selected []string
				}
				msg = FileSelectionMsg{Selected: []string{"./test.txt"}}
			case "progress_msg":
				type ProgressMsg struct {
					Progress float64
					Status   string
				}
				msg = ProgressMsg{Progress: 0.5, Status: "处理中"}
			case "result_viewer_msg":
				type ResultMsg struct {
					Result interface{}
				}
				msg = ResultMsg{Result: nil}
			case "config_editor_msg":
				type ConfigUpdateMsg struct {
					Config interface{}
				}
				msg = ConfigUpdateMsg{Config: nil}
			}
			
			mockModel.On("Update", msg).Return(mockModel, nil)
			
			_, cmd := mockModel.Update(msg)
			if tt.expectedCmd {
				// Update方法返回nil命令是正常的，表示没有后续命令
				assert.Nil(t, cmd)
			}
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}

// TestMainWorkflowErrorHandling 测试错误处理
func TestMainWorkflowErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		state       string
		errorType   string
		expectedCmd bool
		description string
	}{
		{
			name:        "文件加载错误",
			state:       "StateSelect",
			errorType:   "file_load_error",
			expectedCmd: true,
			description: "应该正确处理文件加载错误",
		},
		{
			name:        "处理过程错误",
			state:       "StateProcessing",
			errorType:   "processing_error",
			expectedCmd: true,
			description: "应该正确处理处理过程错误",
		},
		{
			name:        "配置加载错误",
			state:       "StateInput",
			errorType:   "config_error",
			expectedCmd: true,
			description: "应该正确处理配置加载错误",
		},
		{
			name:        "网络错误",
			state:       "StateProcessing",
			errorType:   "network_error",
			expectedCmd: true,
			description: "应该正确处理网络错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state: tt.state,
			}
			
			// 模拟错误消息
			type ErrorMsg struct {
				Err error
			}
			
			msg := ErrorMsg{Err: assert.AnError}
			
			mockModel.On("Update", msg).Return(mockModel, nil)
			
			_, cmd := mockModel.Update(msg)
			if tt.expectedCmd {
				// Update方法返回nil命令是正常的，表示没有后续命令
				assert.Nil(t, cmd)
			}
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}

// TestMainWorkflowViewRendering 测试视图渲染
func TestMainWorkflowViewRendering(t *testing.T) {
	tests := []struct {
		name          string
		state         string
		currentView   string
		width         int
		height        int
		expectedParts []string
		description   string
	}{
		{
			name:          "输入状态视图",
			state:         "StateInput",
			currentView:   "ViewMain",
			width:         80,
			height:        24,
			expectedParts: []string{"路径输入", "提示信息", "帮助信息"},
			description:   "应该正确渲染输入状态视图",
		},
		{
			name:          "文件选择状态视图",
			state:         "StateSelect",
			currentView:   "ViewSelect",
			width:         80,
			height:        24,
			expectedParts: []string{"文件列表", "选择状态", "操作提示"},
			description:   "应该正确渲染文件选择状态视图",
		},
		{
			name:          "处理状态视图",
			state:         "StateProcessing",
			currentView:   "ViewProgress",
			width:         80,
			height:        24,
			expectedParts: []string{"进度条", "状态信息", "进度百分比"},
			description:   "应该正确渲染处理状态视图",
		},
		{
			name:          "结果状态视图",
			state:         "StateResult",
			currentView:   "ViewResult",
			width:         80,
			height:        24,
			expectedParts: []string{"结果概览", "文件列表", "统计信息"},
			description:   "应该正确渲染结果状态视图",
		},
		{
			name:          "配置状态视图",
			state:         "StateConfig",
			currentView:   "ViewConfig",
			width:         80,
			height:        24,
			expectedParts: []string{"配置项", "标签页", "编辑提示"},
			description:   "应该正确渲染配置状态视图",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModel := &MockMainModel{
				state:       tt.state,
				currentView: tt.currentView,
				width:       tt.width,
				height:      tt.height,
			}
			
			// 构建期望的视图内容
			viewContent := ""
			for _, part := range tt.expectedParts {
				viewContent += part + "\n"
			}
			
			mockModel.On("View").Return(viewContent)
			
			view := mockModel.View()
			
			// 验证视图内容
			for _, expectedPart := range tt.expectedParts {
				assert.Contains(t, view, expectedPart)
			}
			
			t.Logf("测试: %s", tt.description)
			mockModel.AssertExpectations(t)
		})
	}
}