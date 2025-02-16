package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func CreateTempFile(dir, filename string) (*os.File, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		// 如果有其他类型的错误，则打印错误信息
		return nil, err
	}

	// 创建临时文件
	file, err := os.CreateTemp(dir, filename)

	if err != nil {
		return nil, err
	}
	return file, nil
	//defer func() {
	//	// 关闭文件并删除它
	//	if err := file.Close(); err != nil {
	//		log.Printf("Error closing file: %v", err)
	//	}
	//	if err := os.Remove(file.Name()); err != nil {
	//		log.Printf("Error removing temp file: %v", err)
	//	}
	//}()
}

func CreateFile(dir, filename string) (*os.File, error) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		// 如果有其他类型的错误，则打印错误信息
		return nil, err
	}

	// 创建临时文件
	file, err := os.Create(path.Join(dir, filename))

	if err != nil {
		return nil, err
	}
	return file, nil
}

// DownloadFile 从给定的 URL 下载文件并保存到本地
func DownloadFile(url string, localFilePath string) (*os.File, error) {
	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: HTTP status code %d", resp.StatusCode)
	}

	// 创建本地文件
	file, err := os.Create(localFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 将响应内容复制到本地文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}
