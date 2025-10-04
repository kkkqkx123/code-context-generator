// Package types Git相关类型定义
package types

import (
	"time"
)

// GitInfo Git仓库信息
type GitInfo struct {
	IsGitRepo      bool        `json:"is_git_repo" yaml:"is_git_repo" xml:"is_git_repo"`
	RepositoryPath string      `json:"repository_path" yaml:"repository_path" xml:"repository_path"`
	CurrentBranch  string      `json:"current_branch" yaml:"current_branch" xml:"current_branch"`
	RemoteURL      string      `json:"remote_url" yaml:"remote_url" xml:"remote_url"`
	CommitCount    int         `json:"commit_count" yaml:"commit_count" xml:"commit_count"`
	LastCommit     *CommitInfo `json:"last_commit" yaml:"last_commit" xml:"last_commit"`
}

// CommitInfo 提交信息
type CommitInfo struct {
	Hash       string    `json:"hash" yaml:"hash" xml:"hash"`
	Author     string    `json:"author" yaml:"author" xml:"author"`
	Email      string    `json:"email" yaml:"email" xml:"email"`
	Date       time.Time `json:"date" yaml:"date" xml:"date"`
	Message    string    `json:"message" yaml:"message" xml:"message"`
	Files      []string  `json:"files" yaml:"files" xml:"files"`
	Insertions int       `json:"insertions" yaml:"insertions" xml:"insertions"`
	Deletions  int       `json:"deletions" yaml:"deletions" xml:"deletions"`
}

// GitHistory Git历史记录
type GitHistory struct {
	Commits      []CommitInfo `json:"commits" yaml:"commits" xml:"commits"`
	TotalCommits int          `json:"total_commits" yaml:"total_commits" xml:"total_commits"`
	TimeRange    TimeRange    `json:"time_range" yaml:"time_range" xml:"time_range"`
	Contributors []string     `json:"contributors" yaml:"contributors" xml:"contributors"`
}

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time `json:"start" yaml:"start" xml:"start"`
	End   time.Time `json:"end" yaml:"end" xml:"end"`
}

// FileDiff 文件差异
type FileDiff struct {
	FilePath   string `json:"file_path" yaml:"file_path" xml:"file_path"`
	Status     string `json:"status" yaml:"status" xml:"status"` // added, modified, deleted
	Insertions int    `json:"insertions" yaml:"insertions" xml:"insertions"`
	Deletions  int    `json:"deletions" yaml:"deletions" xml:"deletions"`
	Diff       string `json:"diff" yaml:"diff" xml:"diff"`
}

// CommitDiff 提交差异
type CommitDiff struct {
	CommitHash   string     `json:"commit_hash" yaml:"commit_hash" xml:"commit_hash"`
	Files        []FileDiff `json:"files" yaml:"files" xml:"files"`
	TotalChanges int        `json:"total_changes" yaml:"total_changes" xml:"total_changes"`
}

// GitStats Git统计信息
type GitStats struct {
	TimePeriod      TimeRange          `json:"time_period" yaml:"time_period" xml:"time_period"`
	CommitStats     CommitStats        `json:"commit_stats" yaml:"commit_stats" xml:"commit_stats"`
	AuthorStats     []AuthorStat       `json:"author_stats" yaml:"author_stats" xml:"author_stats"`
	FileStats       []FileStat         `json:"file_stats" yaml:"file_stats" xml:"file_stats"`
	ActivityHeatmap map[string]int     `json:"activity_heatmap" yaml:"activity_heatmap" xml:"activity_heatmap"`
}

// CommitStats 提交统计
type CommitStats struct {
	TotalCommits     int     `json:"total_commits" yaml:"total_commits" xml:"total_commits"`
	AvgCommitsPerDay float64 `json:"avg_commits_per_day" yaml:"avg_commits_per_day" xml:"avg_commits_per_day"`
	BusiestDay       string  `json:"busiest_day" yaml:"busiest_day" xml:"busiest_day"`
}

// AuthorStat 作者统计
type AuthorStat struct {
	Name       string  `json:"name" yaml:"name" xml:"name"`
	Commits    int     `json:"commits" yaml:"commits" xml:"commits"`
	Changes    int     `json:"changes" yaml:"changes" xml:"changes"`
	Percentage float64 `json:"percentage" yaml:"percentage" xml:"percentage"`
}

// FileStat 文件统计
type FileStat struct {
	FilePath   string `json:"file_path" yaml:"file_path" xml:"file_path"`
	Changes    int    `json:"changes" yaml:"changes" xml:"changes"`
	Insertions int    `json:"insertions" yaml:"insertions" xml:"insertions"`
	Deletions  int    `json:"deletions" yaml:"deletions" xml:"deletions"`
}

// GitIntegrationConfig Git集成配置
type GitIntegrationConfig struct {
	Enabled      bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
	IncludeLogs  bool   `json:"include_logs" yaml:"include_logs" xml:"include_logs"`
	LogCount     int    `json:"log_count" yaml:"log_count" xml:"log_count"`
	IncludeDiffs bool   `json:"include_diffs" yaml:"include_diffs" xml:"include_diffs"`
	DiffFormat   string `json:"diff_format" yaml:"diff_format" xml:"diff_format"` // unified, context, raw
	Stats        struct {
		Enabled    bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
		TimePeriod string `json:"time_period" yaml:"time_period" xml:"time_period"` // 1y, 6m, 30d
		AuthorsTop int    `json:"authors_top" yaml:"authors_top" xml:"authors_top"`
		FilesTop   int    `json:"files_top" yaml:"files_top" xml:"files_top"`
	} `json:"stats" yaml:"stats" xml:"stats"`
	Filters struct {
		Authors []string `json:"authors" yaml:"authors" xml:"authors"`
		Paths   []string `json:"paths" yaml:"paths" xml:"paths"`
		Since   string   `json:"since" yaml:"since" xml:"since"`
		Until   string   `json:"until" yaml:"until" xml:"until"`
	} `json:"filters" yaml:"filters" xml:"filters"`
}

// GitIntegrationData Git集成数据
type GitIntegrationData struct {
	GitInfo      *GitInfo      `json:"git_info,omitempty" yaml:"git_info,omitempty" xml:"git_info,omitempty"`
	GitHistory   *GitHistory   `json:"git_history,omitempty" yaml:"git_history,omitempty" xml:"git_history,omitempty"`
	GitDiffs     []CommitDiff  `json:"git_diffs,omitempty" yaml:"git_diffs,omitempty" xml:"git_diffs,omitempty"`
	GitStats     *GitStats     `json:"git_stats,omitempty" yaml:"git_stats,omitempty" xml:"git_stats,omitempty"`
}

// ContextDataWithGit 带Git信息的上下文数据
type ContextDataWithGit struct {
	ContextData
	GitIntegrationData `json:"git_integration" yaml:"git_integration" xml:"git_integration"`
}