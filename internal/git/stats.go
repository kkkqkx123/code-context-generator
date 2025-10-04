// Package git Git集成功能实现
package git

import (
	"fmt"
	"sort"
	"time"

	"code-context-generator/pkg/types"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitStats Git统计管理器
type GitStats struct {
	repo     *git.Repository
	repoPath string
}

// NewGitStats 创建新的Git统计管理器
func NewGitStats(repo *git.Repository, repoPath string) *GitStats {
	return &GitStats{
		repo:     repo,
		repoPath: repoPath,
	}
}

// GenerateStats 生成Git统计信息
func (gs *GitStats) GenerateStats(timePeriod string, authorsTop, filesTop int) (*types.GitStats, error) {
	// 解析时间周期
	var since *time.Time
	if timePeriod != "" {
		sinceTime, err := ParseTimePeriod(timePeriod)
		if err != nil {
			return nil, err
		}
		since = sinceTime
	}

	// 获取提交历史
	history, err := gs.getCommitHistoryForStats(since)
	if err != nil {
		return nil, fmt.Errorf("获取提交历史失败: %w", err)
	}

	stats := &types.GitStats{
		TimePeriod: types.TimeRange{
			Start: history.startTime,
			End:   history.endTime,
		},
		AuthorStats:     []types.AuthorStat{},
		FileStats:       []types.FileStat{},
		ActivityHeatmap: make(map[string]int),
	}

	// 计算提交统计
	stats.CommitStats = gs.calculateCommitStats(history.commits)

	// 计算作者统计
	stats.AuthorStats = gs.calculateAuthorStats(history.authorCommits, authorsTop)

	// 计算文件统计
	stats.FileStats = gs.calculateFileStats(history.fileChanges, filesTop)

	// 计算活动热力图
	stats.ActivityHeatmap = gs.calculateActivityHeatmap(history.commits)

	return stats, nil
}

// commitHistoryForStats 用于统计的提交历史
type commitHistoryForStats struct {
	commits       []types.CommitInfo
	authorCommits map[string]int
	fileChanges   map[string]*fileChangeInfo
	startTime     time.Time
	endTime       time.Time
}

type fileChangeInfo struct {
	changes    int
	insertions int
	deletions  int
}

// getCommitHistoryForStats 获取用于统计的提交历史
func (gs *GitStats) getCommitHistoryForStats(since *time.Time) (*commitHistoryForStats, error) {
	commitIter, err := gs.repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	history := &commitHistoryForStats{
		commits:       []types.CommitInfo{},
		authorCommits: make(map[string]int),
		fileChanges:   make(map[string]*fileChangeInfo),
	}

	var commits []types.CommitInfo

	err = commitIter.ForEach(func(commit *object.Commit) error {
		// 时间过滤
		if since != nil && commit.Author.When.Before(*since) {
			return nil
		}

		commitInfo := types.CommitInfo{
			Hash:    commit.Hash.String(),
			Author:  commit.Author.Name,
			Email:   commit.Author.Email,
			Date:    commit.Author.When,
			Message: commit.Message,
		}

		commits = append(commits, commitInfo)

		// 统计作者提交数
		history.authorCommits[commit.Author.Name]++

		// 获取文件变更
		files, err := gs.getCommitFilesForStats(commit)
		if err == nil {
			commitInfo.Files = files
			
			// 统计文件变更
			for _, file := range files {
				if history.fileChanges[file] == nil {
					history.fileChanges[file] = &fileChangeInfo{}
				}
				history.fileChanges[file].changes++
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按时间排序
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	history.commits = commits

	// 设置时间范围
	if len(commits) > 0 {
		history.startTime = commits[len(commits)-1].Date
		history.endTime = commits[0].Date
	}

	return history, nil
}

// getCommitFilesForStats 获取用于统计的提交文件列表
func (gs *GitStats) getCommitFilesForStats(commit *object.Commit) ([]string, error) {
	var files []string

	// 获取父提交
	var parent *object.Commit
	if len(commit.ParentHashes) > 0 {
		var err error
		parent, err = gs.repo.CommitObject(commit.ParentHashes[0])
		if err != nil {
			return files, nil // 没有父提交，可能是初始提交
		}
	}

	// 获取文件变更
	var patch *object.Patch
	var err error
	
	if parent != nil {
		patch, err = parent.Patch(commit)
		if err != nil {
			return files, err
		}
	} else {
		// 初始提交，获取所有文件
		iter, err := commit.Files()
		if err != nil {
			return files, err
		}
		
		err = iter.ForEach(func(file *object.File) error {
			files = append(files, file.Name)
			return nil
		})
		return files, err
	}

	// 从patch中提取文件路径
	filePatches := patch.FilePatches()
	for _, filePatch := range filePatches {
		from, to := filePatch.Files()
		
		if from != nil {
			files = append(files, from.Path())
		}
		if to != nil && (from == nil || from.Path() != to.Path()) {
			files = append(files, to.Path())
		}
	}
	
	// 使用diff包避免导入未使用
	_ = diff.Equal

	return files, nil
}

// calculateCommitStats 计算提交统计
func (gs *GitStats) calculateCommitStats(commits []types.CommitInfo) types.CommitStats {
	stats := types.CommitStats{
		TotalCommits: len(commits),
	}

	if len(commits) == 0 {
		return stats
	}

	// 计算平均每天提交数
	duration := commits[0].Date.Sub(commits[len(commits)-1].Date)
	days := duration.Hours() / 24
	if days > 0 {
		stats.AvgCommitsPerDay = float64(len(commits)) / days
	}

	// 找出最活跃的一天
	dailyCommits := make(map[string]int)
	for _, commit := range commits {
		day := commit.Date.Format("2006-01-02")
		dailyCommits[day]++
	}

	maxCommits := 0
	busiestDay := ""
	for day, count := range dailyCommits {
		if count > maxCommits {
			maxCommits = count
			busiestDay = day
		}
	}

	if busiestDay != "" {
		stats.BusiestDay = busiestDay
	}

	return stats
}

// calculateAuthorStats 计算作者统计
func (gs *GitStats) calculateAuthorStats(authorCommits map[string]int, top int) []types.AuthorStat {
	var stats []types.AuthorStat

	totalCommits := 0
	for _, count := range authorCommits {
		totalCommits += count
	}

	for author, commits := range authorCommits {
		percentage := 0.0
		if totalCommits > 0 {
			percentage = float64(commits) * 100.0 / float64(totalCommits)
		}

		stats = append(stats, types.AuthorStat{
			Name:       author,
			Commits:    commits,
			Percentage: percentage,
		})
	}

	// 按提交数排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Commits > stats[j].Commits
	})

	// 限制数量
	if top > 0 && len(stats) > top {
		stats = stats[:top]
	}

	return stats
}

// calculateFileStats 计算文件统计
func (gs *GitStats) calculateFileStats(fileChanges map[string]*fileChangeInfo, top int) []types.FileStat {
	var stats []types.FileStat

	for filePath, info := range fileChanges {
		stats = append(stats, types.FileStat{
			FilePath:   filePath,
			Changes:    info.changes,
			Insertions: info.insertions,
			Deletions:  info.deletions,
		})
	}

	// 按变更次数排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Changes > stats[j].Changes
	})

	// 限制数量
	if top > 0 && len(stats) > top {
		stats = stats[:top]
	}

	return stats
}

// calculateActivityHeatmap 计算活动热力图
func (gs *GitStats) calculateActivityHeatmap(commits []types.CommitInfo) map[string]int {
	heatmap := make(map[string]int)

	for _, commit := range commits {
		day := commit.Date.Format("2006-01-02")
		heatmap[day]++
	}

	return heatmap
}