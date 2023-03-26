package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// 文件或文件夹是否存在
func FileOrDirIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 获取文件的SHA1值(字母小写)
func GetFileSha1(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetSha1 获取[]byte的SHA1值(字母小写)
func GetSha1(data []byte) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 创建文件夹
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func CreateCoverFromLocal(videoPath string, coverPath string) error {
	// capture first video frame as jpg
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vframes", "1", coverPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 本地获取视频时长
func GetVideoDuration(videoPath string) (time.Duration, error) {
	cmd := exec.Command("ffprobe", "-i", videoPath, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	//cmd := exec.Command("cmd", "/C", cmdStr)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	duration, err := time.ParseDuration(strings.TrimSpace(string(out)) + "s")
	if err != nil {
		return 0, err
	}
	return duration, nil
}

// 替换文件名中的非法字符为下划线
func MakeFileNameLegal(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "*", "_")
	s = strings.ReplaceAll(s, "?", "_")
	s = strings.ReplaceAll(s, "\"", "_")
	s = strings.ReplaceAll(s, "<", "_")
	s = strings.ReplaceAll(s, ">", "_")
	s = strings.ReplaceAll(s, "|", "_")
	return s
}

// 获取绝对路径
// filePath为相对路径或绝对路径(可以是文件或文件夹)
func GetAbsPath(filePath string) (string, error) {
	return filepath.Abs(filePath)
}

// 仅在有写入时才创建文件
type LazyFileWriter struct {
	filePath string
	file     *os.File
}

func (w *LazyFileWriter) Write(p []byte) (n int, err error) {
	if w.file == nil {
		w.file, err = os.OpenFile(w.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func GetNewLazyFileWriter(filePath string) *LazyFileWriter {
	return &LazyFileWriter{filePath: filePath}
}
