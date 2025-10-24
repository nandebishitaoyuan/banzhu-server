package file

import (
	"errors"
	"os"
)

// DeleteFile 删除文件
func DeleteFile(path string) error {
	// 检查文件是否存在
	if !FileExists(path) {
		return errors.New("文件不存在")
	}
	return os.Remove(path)
}

// DeleteDir 删除文件夹及其内容
func DeleteDir(path string) error {
	if !FileExists(path) {
		return errors.New("文件夹不存在")
	}
	return os.RemoveAll(path)
}

// GetFile 读取文件并返回 *os.File
func GetFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ListDir 返回指定路径下的子文件和子文件夹
func ListDir(path string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []os.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // 忽略出错的条目
		}
		result = append(result, info)
	}
	return result, nil
}

func CountText(filePath string) (charCount int, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	content := string(data)

	// 按 Unicode 字符统计（中文、英文、符号都算1个）
	charCount = len([]rune(content))

	return charCount, nil
}
