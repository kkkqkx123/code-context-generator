package models

import (
	"fmt"
	"strings"
)

// getFileIcon 根据文件扩展名返回对应的图标
func getFileIcon(filename string, isDir bool) string {
	if isDir {
		return "📂" // 目录使用打开的文件夹图标
	}

	// 获取文件扩展名
	ext := strings.ToLower(strings.TrimPrefix(filename, "."))
	if dotIndex := strings.LastIndex(filename, "."); dotIndex != -1 && dotIndex < len(filename)-1 {
		ext = strings.ToLower(filename[dotIndex+1:])
	}

	// 文档类文件使用📝图标
	switch ext {
	case "md", "txt", "csv", "doc", "docx", "pdf", "rtf":
		return "📝"
	// 配置文件使用⚙️图标
	case "json", "xml", "toml", "yaml", "yml", "ini", "conf", "config", "properties":
		return "⚙️"
	// 代码文件使用💻图标
	case "go", "py", "js", "ts", "java", "cpp", "c", "h", "cs", "php", "rb", "swift", "kt", "rs":
		return "💻"
	// 样式文件使用🎨图标
	case "css", "scss", "sass", "less", "html", "htm":
		return "🎨"
	// 脚本文件使用📜图标
	case "sh", "bat", "cmd", "ps1", "bash", "zsh":
		return "📜"
	// 压缩文件使用📦图标
	case "zip", "rar", "7z", "tar", "gz", "bz2":
		return "📦"
	// 图片文件使用🖼️图标
	case "jpg", "jpeg", "png", "gif", "bmp", "svg", "ico":
		return "🖼️"
	// 音频文件使用🎵图标
	case "mp3", "wav", "flac", "aac", "ogg":
		return "🎵"
	// 视频文件使用🎬图标
	case "mp4", "avi", "mkv", "mov", "wmv", "flv":
		return "🎬"
	// 数据库文件使用🗄️图标
	case "db", "sqlite", "mdb", "accdb":
		return "🗄️"
	// 日志文件使用📋图标
	case "log":
		return "📋"
	// 默认文件图标
	default:
		return "📄"
	}
}

// formatFileSize 格式化文件大小显示
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}