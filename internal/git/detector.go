// Package git Git集成功能实现
package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code-context-generator/pkg/types"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitDetector Git仓库检测器
type GitDetector struct {
	repoPath string
	isGitRepo bool
	gitDir   string
	repo     *git.Repository
}

// NewGitDetector 创建新的Git检测器
func NewGitDetector(repoPath string) *GitDetector {
	return &GitDetector{
		repoPath: repoPath,
	}
}

// Detect 检测是否为Git仓库
func (gd *GitDetector) Detect() error {
	// 检查当前目录是否为Git仓库
	repo, err := git.PlainOpen(gd.repoPath)
	if err == nil {
		gd.repo = repo
		gd.isGitRepo = true
		gd.gitDir = filepath.Join(gd.repoPath, ".git")
		return nil
	}

	// 检查父目录是否为Git仓库
	currentPath := gd.repoPath
	for {
		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			break // 到达根目录
		}
		
		repo, err := git.PlainOpen(parent)
		if err == nil {
			gd.repo = repo
			gd.isGitRepo = true
			gd.gitDir = filepath.Join(parent, ".git")
			gd.repoPath = parent
			return nil
		}
		currentPath = parent
	}

	return fmt.Errorf("未找到Git仓库")
}

// GetGitInfo 获取Git信息
func (gd *GitDetector) GetGitInfo() (*types.GitInfo, error) {
	if !gd.isGitRepo {
		return nil, fmt.Errorf("不是Git仓库")
	}

	info := &types.GitInfo{
		IsGitRepo:      gd.isGitRepo,
		RepositoryPath: gd.repoPath,
	}

	// 获取当前分支
	head, err := gd.repo.Head()
	if err == nil {
		info.CurrentBranch = head.Name().Short()
	}

	// 获取远程URL
	remotes, err := gd.repo.Remotes()
	if err == nil && len(remotes) > 0 {
		config := remotes[0].Config()
		if len(config.URLs) > 0 {
			info.RemoteURL = config.URLs[0]
		}
	}

	// 获取提交数量
	commitIter, err := gd.repo.CommitObjects()
	if err == nil {
		count := 0
		err = commitIter.ForEach(func(c *object.Commit) error {
			count++
			return nil
		})
		if err == nil {
			info.CommitCount = count
		}
	}

	// 获取最新提交
	if head != nil {
		commit, err := gd.repo.CommitObject(head.Hash())
		if err == nil {
			info.LastCommit = &types.CommitInfo{
				Hash:    commit.Hash.String(),
				Author:  commit.Author.Name,
				Email:   commit.Author.Email,
				Date:    commit.Author.When,
				Message: strings.TrimSpace(commit.Message),
			}
		}
	}

	return info, nil
}

// IsGitRepository 检查是否为Git仓库
func IsGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return true
	}
	return false
}

// FindGitRepository 查找Git仓库根目录
func FindGitRepository(startPath string) (string, error) {
	currentPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		if IsGitRepository(currentPath) {
			return currentPath, nil
		}

		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			break // 到达根目录
		}
		currentPath = parent
	}

	return "", fmt.Errorf("未找到Git仓库")
}