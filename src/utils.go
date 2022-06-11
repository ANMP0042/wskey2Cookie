// Package src /**
package src

import "os"

func Write2File(body, filePath string) {
	file, _ := os.Create(filePath)

	// 写入数据到文件中
	file.WriteString(body)
	defer file.Close()
	return
}
