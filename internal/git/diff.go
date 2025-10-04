// Package git Git集成功能实现
package git

import (
	"bytes"
	"fmt"
	"strings"

	"code-context-generator/pkg/types"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitDiff Git差异管理器
type GitDiff struct {
	repo     *git.Repository
	repoPath string
}

// NewGitDiff 创建新的Git差异管理器
func NewGitDiff(repo *git.Repository, repoPath string) *GitDiff {
	return &GitDiff{
		repo:     repo,
		repoPath: repoPath,
	}
}

// GetCommitDiff 获取提交的差异
func (gd *GitDiff) GetCommitDiff(commitHash string, format string) (*types.CommitDiff, error) {
	// 解析提交哈希
	hash, err := gd.repo.ResolveRevision(plumbing.Revision(commitHash))
	if err != nil {
		return nil, fmt.Errorf("解析提交哈希失败: %w", err)
	}

	commit, err := gd.repo.CommitObject(*hash)
	if err != nil {
		return nil, fmt.Errorf("获取提交对象失败: %w", err)
	}

	return gd.getCommitDiffFromCommit(commit, format)
}

// GetLatestCommitDiff 获取最新提交的差异
func (gd *GitDiff) GetLatestCommitDiff(format string) (*types.CommitDiff, error) {
	head, err := gd.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("获取HEAD失败: %w", err)
	}

	commit, err := gd.repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("获取提交对象失败: %w", err)
	}

	return gd.getCommitDiffFromCommit(commit, format)
}

// getCommitDiffFromCommit 从提交对象获取差异
func (gd *GitDiff) getCommitDiffFromCommit(commit *object.Commit, format string) (*types.CommitDiff, error) {
	diff := &types.CommitDiff{
		CommitHash: commit.Hash.String(),
		Files:      []types.FileDiff{},
	}

	// 获取父提交
	var parent *object.Commit
	if len(commit.ParentHashes) > 0 {
		var err error
		parent, err = gd.repo.CommitObject(commit.ParentHashes[0])
		if err != nil {
			return nil, fmt.Errorf("获取父提交失败: %w", err)
		}
	}

	// 获取文件变更
	var patch *object.Patch
	var err error
	
	if parent != nil {
		patch, err = parent.Patch(commit)
		if err != nil {
			return nil, fmt.Errorf("获取补丁失败: %w", err)
		}
	} else {
		// 初始提交，获取所有文件
		return gd.getInitialCommitDiff(commit, format)
	}

	// 获取文件变更
	filePatches := patch.FilePatches()

	// 处理每个变更
	for _, change := range filePatches {
		fileDiff, err := gd.processFilePatch(change, format)
		if err != nil {
			continue // 跳过处理失败的变更
		}
		
		diff.Files = append(diff.Files, *fileDiff)
		diff.TotalChanges++
	}

	return diff, nil
}

// getInitialCommitDiff 获取初始提交的差异
func (gd *GitDiff) getInitialCommitDiff(commit *object.Commit, format string) (*types.CommitDiff, error) {
	diff := &types.CommitDiff{
		CommitHash: commit.Hash.String(),
		Files:      []types.FileDiff{},
	}

	// 获取所有文件
	fileIter, err := commit.Files()
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %w", err)
	}

	err = fileIter.ForEach(func(file *object.File) error {
		content, err := file.Contents()
		if err != nil {
			return nil // 跳过无法读取的文件
		}

		lines := strings.Count(content, "\n")
		if !strings.HasSuffix(content, "\n") {
			lines++
		}

		fileDiff := types.FileDiff{
			FilePath:   file.Name,
			Status:     "added",
			Insertions: lines,
			Deletions:  0,
		}

		if format != "" {
			fileDiff.Diff = gd.formatDiff("", content, format)
		}

		diff.Files = append(diff.Files, fileDiff)
		diff.TotalChanges++
		return nil
	})

	return diff, err
}

// processFilePatch 处理单个文件补丁
func (gd *GitDiff) processFilePatch(patch diff.FilePatch, format string) (*types.FileDiff, error) {
	from, to := patch.Files()
	
	fileDiff := &types.FileDiff{}

	// 确定文件状态
	if from == nil && to != nil {
		fileDiff.Status = "added"
		fileDiff.FilePath = to.Path()
	} else if from != nil && to == nil {
		fileDiff.Status = "deleted"
		fileDiff.FilePath = from.Path()
	} else if from != nil && to != nil {
		fileDiff.Status = "modified"
		fileDiff.FilePath = to.Path()
	}

	// 使用chunks来获取变更信息
	if format != "" {
		// 从chunks构建差异内容
		var fromContent, toContent strings.Builder
		chunks := patch.Chunks()
		
		for _, chunk := range chunks {
			switch chunk.Type() {
			case diff.Add:
				toContent.WriteString(chunk.Content())
			case diff.Delete:
				fromContent.WriteString(chunk.Content())
			case diff.Equal:
				fromContent.WriteString(chunk.Content())
				toContent.WriteString(chunk.Content())
			}
		}
		
		diffContent := gd.formatDiff(fromContent.String(), toContent.String(), format)
		fileDiff.Diff = diffContent
		fileDiff.Insertions, fileDiff.Deletions = gd.countChanges(diffContent)
	} else {
		// 简单的行数统计
		chunks := patch.Chunks()
		for _, chunk := range chunks {
			switch chunk.Type() {
			case diff.Add:
				fileDiff.Insertions += strings.Count(chunk.Content(), "\n")
			case diff.Delete:
				fileDiff.Deletions += strings.Count(chunk.Content(), "\n")
			}
		}
	}

	return fileDiff, nil
}

// formatDiff 格式化差异
func (gd *GitDiff) formatDiff(fromContent, toContent, format string) string {
	fromLines := strings.Split(fromContent, "\n")
	toLines := strings.Split(toContent, "\n")

	var buffer bytes.Buffer

	switch format {
	case "unified":
		buffer.WriteString("@@ -1,")
		buffer.WriteString(fmt.Sprintf("%d", len(fromLines)))
		buffer.WriteString(" +1,")
		buffer.WriteString(fmt.Sprintf("%d", len(toLines)))
		buffer.WriteString(" @@\n")

		// 简化的差异显示
		for i, line := range toLines {
			if i < len(fromLines) && line != fromLines[i] {
				buffer.WriteString("-")
				buffer.WriteString(fromLines[i])
				buffer.WriteString("\n")
				buffer.WriteString("+")
				buffer.WriteString(line)
				buffer.WriteString("\n")
			} else if i >= len(fromLines) {
				buffer.WriteString("+")
				buffer.WriteString(line)
				buffer.WriteString("\n")
			}
		}

	case "context":
		buffer.WriteString("--- a/\n")
		buffer.WriteString("+++ b/\n")
		buffer.WriteString("***************\n")
		
		// 简化的上下文格式
		for _, line := range toLines {
			buffer.WriteString("  ")
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}

	case "raw":
		// 原始格式，只显示内容
		buffer.WriteString(toContent)

	default:
		// 默认统一格式
		return gd.formatDiff(fromContent, toContent, "unified")
	}

	return buffer.String()
}

// countChanges 计算差异中的变更行数
func (gd *GitDiff) countChanges(diff string) (insertions, deletions int) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			insertions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletions++
		}
	}
	return
}

// countContentChanges 计算内容变更的行数
func (gd *GitDiff) countContentChanges(fromContent, toContent string) (insertions, deletions int) {
	fromLines := strings.Split(fromContent, "\n")
	toLines := strings.Split(toContent, "\n")

	maxLen := len(fromLines)
	if len(toLines) > maxLen {
		maxLen = len(toLines)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(fromLines) {
			insertions++
		} else if i >= len(toLines) {
			deletions++
		} else if fromLines[i] != toLines[i] {
			insertions++
			deletions++
		}
	}

	return
}