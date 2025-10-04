// Package git Git集成功能实现
package git

import (
	"fmt"
	"time"

	"code-context-generator/pkg/types"

	"github.com/go-git/go-git/v5"
)

// Integration Git集成管理器
type Integration struct {
	detector *GitDetector
	history  *GitHistory
	diff     *GitDiff
	stats    *GitStats
	config   *types.GitIntegrationConfig
}

// NewIntegration 创建新的Git集成管理器
func NewIntegration(repoPath string, config *types.GitIntegrationConfig) (*Integration, error) {
	// 检测Git仓库
	detector := NewGitDetector(repoPath)
	if err := detector.Detect(); err != nil {
		return nil, fmt.Errorf("Git仓库检测失败: %w", err)
	}

	// 获取Git信息
	gitInfo, err := detector.GetGitInfo()
	if err != nil {
		return nil, fmt.Errorf("获取Git信息失败: %w", err)
	}

	if !gitInfo.IsGitRepo {
		return nil, fmt.Errorf("不是Git仓库")
	}

	// 打开仓库
	repo, err := git.PlainOpen(detector.repoPath)
	if err != nil {
		return nil, fmt.Errorf("打开Git仓库失败: %w", err)
	}

	return &Integration{
		detector: detector,
		history:  NewGitHistory(repo, detector.repoPath),
		diff:     NewGitDiff(repo, detector.repoPath),
		stats:    NewGitStats(repo, detector.repoPath),
		config:   config,
	}, nil
}

// GetGitIntegrationData 获取Git集成数据
func (i *Integration) GetGitIntegrationData() (*types.GitIntegrationData, error) {
	data := &types.GitIntegrationData{}

	// 获取Git基本信息
	gitInfo, err := i.detector.GetGitInfo()
	if err != nil {
		return nil, fmt.Errorf("获取Git信息失败: %w", err)
	}
	data.GitInfo = gitInfo

	// 如果启用了日志包含
	if i.config.IncludeLogs {
		// 解析时间过滤
		var since, until *time.Time
		if i.config.Filters.Since != "" {
		sinceTime, err := time.Parse("2006-01-02", i.config.Filters.Since)
			if err != nil {
				return nil, fmt.Errorf("解析开始时间失败: %w", err)
			}
			since = &sinceTime
		}
		if i.config.Filters.Until != "" {
			untilTime, err := time.Parse("2006-01-02", i.config.Filters.Until)
			if err != nil {
				return nil, fmt.Errorf("解析结束时间失败: %w", err)
			}
			until = &untilTime
		}

		// 获取提交历史
		gitHistory, err := i.history.GetCommitHistory(i.config.LogCount, since, until, i.config.Filters.Authors)
		if err != nil {
			return nil, fmt.Errorf("获取提交历史失败: %w", err)
		}
		data.GitHistory = gitHistory

		// 如果启用了差异包含
		if i.config.IncludeDiffs {
			var diffs []types.CommitDiff
			for _, commit := range gitHistory.Commits {
				commitDiff, err := i.diff.GetCommitDiff(commit.Hash, i.config.DiffFormat)
				if err != nil {
					continue // 跳过失败的差异获取
				}
				diffs = append(diffs, *commitDiff)
			}
			data.GitDiffs = diffs
		}
	}

	// 如果启用了统计
	if i.config.Stats.Enabled {
		gitStats, err := i.stats.GenerateStats(i.config.Stats.TimePeriod, i.config.Stats.AuthorsTop, i.config.Stats.FilesTop)
		if err != nil {
			return nil, fmt.Errorf("生成统计信息失败: %w", err)
		}
		data.GitStats = gitStats
	}

	return data, nil
}

// GetGitInfo 获取Git基本信息
func (i *Integration) GetGitInfo() (*types.GitInfo, error) {
	return i.detector.GetGitInfo()
}

// GetCommitHistory 获取提交历史
func (i *Integration) GetCommitHistory(count int, since, until *time.Time, authors []string) (*types.GitHistory, error) {
	return i.history.GetCommitHistory(count, since, until, authors)
}

// GetCommitDiff 获取提交差异
func (i *Integration) GetCommitDiff(commitHash string) (*types.CommitDiff, error) {
	return i.diff.GetCommitDiff(commitHash, i.config.DiffFormat)
}

// GetGitStats 获取Git统计
func (i *Integration) GetGitStats() (*types.GitStats, error) {
	return i.stats.GenerateStats(i.config.Stats.TimePeriod, i.config.Stats.AuthorsTop, i.config.Stats.FilesTop)
}

// IsGitRepository 检查是否为Git仓库
func (i *Integration) IsGitRepository() bool {
	return i.detector.isGitRepo
}

// GetRepositoryPath 获取仓库路径
func (i *Integration) GetRepositoryPath() string {
	return i.detector.repoPath
}