# Git集成功能扩展计划

## 功能概述

为项目添加Git集成功能，支持Git仓库分析、提交历史集成、变更追踪等特性，增强代码上下文生成的版本控制维度。

## 当前状态分析

### 现有功能
- ❌ 无Git集成功能
- ❌ 无法识别Git仓库
- ❌ 缺少版本历史信息

### 与repomix对比
- ❌ 缺少Git日志集成
- ❌ 缺少变更差异显示
- ❌ 缺少提交统计

## 扩展目标

### 核心功能
1. **Git仓库检测** - 自动识别Git仓库
2. **提交历史集成** - 包含Git日志信息
3. **变更差异显示** - 显示文件变更历史
4. **提交统计功能** - 代码变更分析

### 高级功能
- 分支信息集成
- 标签版本支持
- 代码作者统计
- 变更热点分析

## 技术实现方案

### 1. Git仓库检测系统

#### 检测逻辑
```go
type GitDetector struct {
    repoPath string
    isGitRepo bool
    gitDir string
}

func (gd *GitDetector) Detect() error
func (gd *GitDetector) GetGitInfo() (*GitInfo, error)
```

#### Git信息结构
```go
type GitInfo struct {
    IsGitRepo    bool       `json:"is_git_repo"`
    RepositoryPath string   `json:"repository_path"`
    CurrentBranch string    `json:"current_branch"`
    RemoteURL     string    `json:"remote_url"`
    CommitCount   int       `json:"commit_count"`
    LastCommit    *CommitInfo `json:"last_commit"`
}
```

### 2. 提交历史集成

#### 提交信息获取
```go
type CommitInfo struct {
    Hash        string    `json:"hash"`
    Author      string    `json:"author"`
    Email       string    `json:"email"`
    Date        time.Time `json:"date"`
    Message     string    `json:"message"`
    Files       []string  `json:"files"`
    Insertions  int       `json:"insertions"`
    Deletions   int       `json:"deletions"`
}

type GitHistory struct {
    Commits []CommitInfo `json:"commits"`
    TotalCommits int     `json:"total_commits"`
    TimeRange    TimeRange `json:"time_range"`
    Contributors []string `json:"contributors"`
}
```

#### 命令行集成
```bash
# 包含Git日志（默认50条）
c-gen --include-git-logs

# 指定提交数量
c-gen --include-git-logs --git-log-count 20

# 包含变更差异
c-gen --include-git-logs --include-diffs

# 指定时间范围
c-gen --include-git-logs --git-since "2024-01-01" --git-until "2024-12-31"
```

### 3. 变更差异显示

#### 差异信息结构
```go
type FileDiff struct {
    FilePath    string `json:"file_path"`
    Status      string `json:"status"`  // added, modified, deleted
    Insertions  int    `json:"insertions"`
    Deletions   int    `json:"deletions"`
    Diff        string `json:"diff"`
}

type CommitDiff struct {
    CommitHash string     `json:"commit_hash"`
    Files      []FileDiff `json:"files"`
    TotalChanges int      `json:"total_changes"`
}
```

#### 输出格式集成
```xml
<git_info>
  <repository>
    <path>/path/to/repo</path>
    <branch>main</branch>
    <remote>https://github.com/user/repo.git</remote>
  </repository>
  <commits>
    <commit hash="abc123" date="2024-01-01T10:00:00Z">
      <author>John Doe</author>
      <message>Initial commit</message>
      <files>
        <file path="src/main.go" status="added" insertions="50" deletions="0"/>
      </files>
    </commit>
  </commits>
</git_info>
```

### 4. 提交统计功能

#### 统计维度
- **提交频率分析** - 时间分布统计
- **代码变更量** - 行数变化统计
- **作者贡献度** - 开发者统计
- **文件变更热度** - 文件修改频率

#### 统计数据结构
```go
type GitStats struct {
    TimePeriod    TimeRange          `json:"time_period"`
    CommitStats   CommitStats        `json:"commit_stats"`
    AuthorStats   []AuthorStat       `json:"author_stats"`
    FileStats     []FileStat         `json:"file_stats"`
    ActivityHeatmap map[string]int   `json:"activity_heatmap"`
}

type CommitStats struct {
    TotalCommits int     `json:"total_commits"`
    AvgCommitsPerDay float64 `json:"avg_commits_per_day"`
    BusiestDay   string  `json:"busiest_day"`
}

type AuthorStat struct {
    Name     string `json:"name"`
    Commits  int    `json:"commits"`
    Changes  int    `json:"changes"`
    Percentage float64 `json:"percentage"`
}
```

## 实施步骤

### 第一阶段：基础集成（2周）
1. **Git仓库检测**
   - 实现Git检测逻辑
   - 集成go-git库
   - 基础信息获取

2. **提交历史获取**
   - 实现提交日志读取
   - 支持数量限制
   - 时间范围过滤

### 第二阶段：高级功能（2周）
3. **变更差异集成**
   - 实现差异计算
   - 支持多种diff格式
   - 性能优化

4. **统计功能**
   - 提交统计实现
   - 作者分析
   - 变更热点识别

### 第三阶段：优化集成（1周）
5. **输出格式优化**
   - 各格式Git信息集成
   - AI优化输出
   - 性能测试

6. **用户体验优化**
   - 命令行界面完善
   - 错误处理
   - 文档更新

## 代码架构

### 新增包结构
```
internal/git/
├── detector.go          # Git仓库检测
├── history.go           # 提交历史获取
├── diff.go              # 差异计算
├── stats.go             # 统计功能
└── integration.go       # 集成管理器

pkg/types/
├── git_info.go          # Git信息结构
├── commit_info.go       # 提交信息
└── git_stats.go         # 统计结构
```

### 依赖管理
- **主要依赖**: `github.com/go-git/go-git/v5`
- **可选依赖**: 系统git命令（备用）
- **兼容性**: 支持Git 2.0+

## 配置变更

### 新增配置项
```yaml
git:
  enabled: true
  include_logs: false
  log_count: 50
  include_diffs: false
  diff_format: "unified"  # unified, context, raw
  stats:
    enabled: false
    time_period: "1y"     # 1y, 6m, 30d
    authors_top: 10
    files_top: 20
  filters:
    authors: []
    paths: []
    since: ""
    until: ""
```

## 命令行参数

### 新增参数
```bash
# Git集成相关
--include-git-logs              # 包含Git日志
--git-log-count <number>        # 提交数量限制
--include-diffs                 # 包含变更差异
--git-since <date>              # 开始时间
--git-until <date>              # 结束时间

# 统计功能
--git-stats                     # 生成Git统计
--git-authors-top <number>      # 显示前N名作者
--git-files-top <number>        # 显示前N个变更文件

# 过滤选项
--git-author <name>             # 按作者过滤
--git-path <pattern>            # 按路径过滤
```

## 输出格式集成

### XML格式示例
```xml
<git_integration>
  <repository_info>
    <path>/project/path</path>
    <branch>main</branch>
    <commit_count>150</commit_count>
  </repository_info>
  <commit_history>
    <commit>
      <hash>abc123</hash>
      <author>John Doe</author>
      <date>2024-01-01T10:00:00Z</date>
      <message>Initial commit</message>
      <files_changed>3</files_changed>
    </commit>
  </commit_history>
  <git_stats>
    <total_commits>150</total_commits>
    <top_authors>
      <author name="John Doe" commits="100" percentage="66.7%"/>
    </top_authors>
  </git_stats>
</git_integration>
```

### JSON格式示例
```json
{
  "git_info": {
    "is_git_repo": true,
    "repository_path": "/project/path",
    "current_branch": "main",
    "commit_count": 150
  },
  "commit_history": [
    {
      "hash": "abc123",
      "author": "John Doe",
      "message": "Initial commit",
      "files": ["src/main.go"]
    }
  ]
}
```

## 测试策略

### 单元测试
- Git检测逻辑测试
- 提交历史解析测试
- 差异计算准确性测试

### 集成测试
- 真实Git仓库测试
- 多种Git操作场景测试
- 性能基准测试

### 验收测试
- 与git命令行输出对比
- 大型仓库性能测试
- 边缘情况处理测试

## 性能考虑

### 优化策略
- **懒加载**: Git操作按需执行
- **缓存**: 提交信息缓存
- **并行**: 多提交并行处理
- **增量**: 只处理变更部分

### 大仓库处理
- 支持分页获取提交
- 内存使用监控
- 超时处理机制

## 风险评估

### 技术风险
- Git库兼容性问题
- 大仓库性能问题
- 复杂Git历史处理

### 缓解措施
- 备用实现方案
- 性能监控和限制
- 渐进式功能启用

## 成功指标

### 功能指标
- 支持主流Git操作
- 性能影响<20%
- 准确率>99%

### 用户体验
- 响应时间合理
- 错误信息清晰
- 配置灵活

## 后续扩展

### 短期扩展
- GitHub/GitLab API集成
- 代码审查集成
- 自动化工作流

### 长期扩展
- 代码质量分析集成
- 团队协作分析
- 预测性分析