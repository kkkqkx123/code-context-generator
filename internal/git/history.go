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

// GitHistory Git历史记录管理器
type GitHistory struct {
	repo     *git.Repository
	repoPath string
}

// NewGitHistory 创建新的Git历史记录管理器
func NewGitHistory(repo *git.Repository, repoPath string) *GitHistory {
	return &GitHistory{
		repo:     repo,
		repoPath: repoPath,
	}
}

// GetCommitHistory 获取提交历史
func (gh *GitHistory) GetCommitHistory(count int, since, until *time.Time, authorFilter []string) (*types.GitHistory, error) {
	commitIter, err := gh.repo.CommitObjects()
	if err != nil {
		return nil, fmt.Errorf("获取提交对象失败: %w", err)
	}

	var commits []types.CommitInfo
	contributors := make(map[string]bool)

	err = commitIter.ForEach(func(commit *object.Commit) error {
		// 时间过滤
		if since != nil && commit.Author.When.Before(*since) {
			return nil
		}
		if until != nil && commit.Author.When.After(*until) {
			return nil
		}

		// 作者过滤
		if len(authorFilter) > 0 {
			found := false
			for _, author := range authorFilter {
				if commit.Author.Name == author {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}

		commitInfo := types.CommitInfo{
			Hash:    commit.Hash.String(),
			Author:  commit.Author.Name,
			Email:   commit.Author.Email,
			Date:    commit.Author.When,
			Message: commit.Message,
		}

		// 获取文件变更列表
		files, err := gh.getCommitFiles(commit)
		if err == nil {
			commitInfo.Files = files
		}

		commits = append(commits, commitInfo)
		contributors[commit.Author.Name] = true

		// 限制数量
		if count > 0 && len(commits) >= count {
			return fmt.Errorf("达到数量限制")
		}

		return nil
	})

	if err != nil && err.Error() != "达到数量限制" {
		return nil, err
	}

	// 按时间排序（最新的在前）
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	// 获取贡献者列表
	var contributorList []string
	for contributor := range contributors {
		contributorList = append(contributorList, contributor)
	}
	sort.Strings(contributorList)

	// 计算时间范围
	var timeRange types.TimeRange
	if len(commits) > 0 {
		timeRange.Start = commits[len(commits)-1].Date
		timeRange.End = commits[0].Date
	}

	return &types.GitHistory{
		Commits:      commits,
		TotalCommits: len(commits),
		TimeRange:    timeRange,
		Contributors: contributorList,
	}, nil
}

// getCommitFiles 获取提交的文件列表
func (gh *GitHistory) getCommitFiles(commit *object.Commit) ([]string, error) {
	var files []string

	// 获取父提交
	var parent *object.Commit
	if len(commit.ParentHashes) > 0 {
		var err error
		parent, err = gh.repo.CommitObject(commit.ParentHashes[0])
		if err != nil {
			return files, nil // 没有父提交，可能是初始提交
		}
	}

	// 获取文件差异
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
	// 去重并保持顺序
	seen := make(map[string]bool)
	var unique []string
	for _, f := range files {
		if !seen[f] {
			seen[f] = true
			unique = append(unique, f)
		}
	}
	files = unique

	return files, nil
}

// GetLatestCommit 获取最新提交
func (gh *GitHistory) GetLatestCommit() (*types.CommitInfo, error) {
	head, err := gh.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("获取HEAD失败: %w", err)
	}

	commit, err := gh.repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("获取提交对象失败: %w", err)
	}

	files, err := gh.getCommitFiles(commit)
	if err != nil {
		files = []string{}
	}

	return &types.CommitInfo{
		Hash:    commit.Hash.String(),
		Author:  commit.Author.Name,
		Email:   commit.Author.Email,
		Date:    commit.Author.When,
		Message: commit.Message,
		Files:   files,
	}, nil
}

// ParseTimePeriod 解析时间周期字符串
func ParseTimePeriod(period string) (*time.Time, error) {
	now := time.Now()
	
	switch period {
	case "1y", "1year", "1 year":
		return &[]time.Time{now.AddDate(-1, 0, 0)}[0], nil
	case "6m", "6month", "6 months":
		return &[]time.Time{now.AddDate(0, -6, 0)}[0], nil
	case "3m", "3month", "3 months":
		return &[]time.Time{now.AddDate(0, -3, 0)}[0], nil
	case "1m", "1month", "1 month":
		return &[]time.Time{now.AddDate(0, -1, 0)}[0], nil
	case "30d", "30days", "30 days":
		return &[]time.Time{now.AddDate(0, 0, -30)}[0], nil
	case "7d", "7days", "7 days":
		return &[]time.Time{now.AddDate(0, 0, -7)}[0], nil
	case "1d", "1day", "1 day":
		return &[]time.Time{now.AddDate(0, 0, -1)}[0], nil
	default:
		return nil, fmt.Errorf("不支持的时间周期: %s", period)
	}
}