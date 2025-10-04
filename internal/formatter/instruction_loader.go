package formatter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"code-context-generator/pkg/types"
)

// InstructionLoader AI指令加载器
type InstructionLoader struct {
	config *types.Config
}

// NewInstructionLoader 创建指令加载器
func NewInstructionLoader(config *types.Config) *InstructionLoader {
	return &InstructionLoader{
		config: config,
	}
}

// LoadInstructions 加载AI指令
func (l *InstructionLoader) LoadInstructions() (string, error) {
	if l.config == nil || !l.config.Output.AIInstructions.Enabled {
		return "", nil
	}

	// 优先从文件加载
	if l.config.Output.AIInstructions.FilePath != "" {
		content, err := l.loadFromFile(l.config.Output.AIInstructions.FilePath)
		if err != nil {
			return "", fmt.Errorf("加载指令文件失败: %w", err)
		}
		return content, nil
	}

	// 其次使用内联内容
	if l.config.Output.AIInstructions.Content != "" {
		return l.processTemplate(l.config.Output.AIInstructions.Content), nil
	}

	// 返回默认指令
	return l.getDefaultInstructions(), nil
}

// loadFromFile 从文件加载指令
func (l *InstructionLoader) loadFromFile(filePath string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 尝试相对路径
		if !filepath.IsAbs(filePath) {
			cwd, err := os.Getwd()
			if err != nil {
				return "", err
			}
			filePath = filepath.Join(cwd, filePath)
		}
		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return "", fmt.Errorf("指令文件不存在: %s", filePath)
		}
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取指令文件失败: %s", err)
	}

	return l.processTemplate(string(content)), nil
}

// processTemplate 处理模板变量替换
func (l *InstructionLoader) processTemplate(content string) string {
	// 简单的模板变量替换
	// 可以扩展为更复杂的模板系统
	
	replacements := map[string]string{
		"{{TOOL_NAME}}":     "code-context-generator",
		"{{CURRENT_DATE}}":  fmt.Sprintf("%v", os.Getenv("DATE")),
		"{{REPO_NAME}}":     l.getRepoName(),
	}

	result := content
	for key, value := range replacements {
		result = strings.ReplaceAll(result, key, value)
	}

	return result
}

// getRepoName 获取仓库名称
func (l *InstructionLoader) getRepoName() string {
	// 从当前目录名获取
	cwd, err := os.Getwd()
	if err != nil {
		return "unknown-project"
	}
	
	return filepath.Base(cwd)
}

// getDefaultInstructions 获取默认指令
func (l *InstructionLoader) getDefaultInstructions() string {
	return `## AI Analysis Instructions

Please analyze this codebase for the following aspects:

### Code Quality
- Identify potential code quality issues
- Suggest improvements for readability and maintainability
- Check for code duplication

### Security
- Look for potential security vulnerabilities
- Identify hardcoded credentials or sensitive data
- Check for unsafe coding practices

### Performance
- Suggest performance optimization opportunities
- Identify potential bottlenecks
- Recommend efficient algorithms or data structures

### Architecture
- Analyze the overall architecture and design patterns
- Suggest improvements for modularity and extensibility
- Identify coupling and cohesion issues

### Documentation
- Check for missing or outdated documentation
- Suggest improvements for code comments
- Identify areas that need better documentation

### Testing
- Identify areas that lack test coverage
- Suggest test cases for critical functionality
- Check for proper error handling

Please provide specific, actionable recommendations for each area.`
}

// GetPresetInstructions 获取预设指令
func (l *InstructionLoader) GetPresetInstructions(preset string) string {
	switch preset {
	case "security":
		return `## Security Analysis Instructions

Focus on security-related aspects of this codebase:

1. **Input Validation**
   - Check for SQL injection vulnerabilities
   - Identify XSS (Cross-Site Scripting) risks
   - Look for buffer overflow possibilities
   - Validate all user inputs

2. **Authentication & Authorization**
   - Review authentication mechanisms
   - Check authorization implementations
   - Identify privilege escalation risks
   - Validate session management

3. **Data Protection**
   - Look for sensitive data exposure
   - Check encryption implementations
   - Identify insecure data storage
   - Validate data transmission security

4. **Code Injection**
   - Check for command injection vulnerabilities
   - Look for LDAP injection risks
   - Identify XML injection possibilities
   - Validate dynamic code execution

5. **Error Handling**
   - Review error message disclosure
   - Check for information leakage
   - Validate exception handling
   - Identify debugging information exposure

Provide specific security recommendations and risk assessments.`

	case "performance":
		return `## Performance Analysis Instructions

Analyze this codebase for performance optimization opportunities:

1. **Algorithm Efficiency**
   - Identify inefficient algorithms (O(n²) when O(n log n) possible)
   - Look for unnecessary nested loops
   - Check for optimal data structure usage
   - Suggest better sorting/searching algorithms

2. **Memory Usage**
   - Identify memory leaks
   - Check for excessive object creation
   - Look for inefficient memory patterns
   - Suggest memory pooling where applicable

3. **I/O Operations**
   - Minimize file system operations
   - Optimize database queries
   - Reduce network calls
   - Implement proper caching strategies

4. **Concurrency**
   - Identify blocking operations
   - Suggest async/await patterns
   - Look for race conditions
   - Recommend parallel processing

5. **Resource Management**
   - Check resource cleanup
   - Identify resource contention
   - Suggest connection pooling
   - Optimize resource allocation

Provide specific performance metrics and improvement suggestions.`

	case "documentation":
		return `## Documentation Analysis Instructions

Review this codebase for documentation quality:

1. **Code Comments**
   - Check for missing function/method comments
   - Identify unclear or outdated comments
   - Look for TODO/FIXME items
   - Validate comment accuracy

2. **API Documentation**
   - Review public API documentation
   - Check parameter descriptions
   - Validate return value documentation
   - Identify missing examples

3. **Architecture Documentation**
   - Check for high-level architecture docs
   - Identify missing design decisions
   - Look for outdated system documentation
   - Validate component interactions

4. **README and Guides**
   - Review installation instructions
   - Check usage examples
   - Validate configuration documentation
   - Identify missing setup guides

5. **Inline Documentation**
   - Check for complex algorithm explanations
   - Identify business logic documentation
   - Look for data structure documentation
   - Validate code examples

Suggest specific documentation improvements and additions.`

	default:
		return l.getDefaultInstructions()
	}
}