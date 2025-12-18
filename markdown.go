package main

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// ReadMarkdownFile 读取markdown文件内容
func ReadMarkdownFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ConvertMarkdownToHTML 将markdown转换为HTML
func ConvertMarkdownToHTML(markdown string) (string, error) {
	// 创建带有扩展的goldmark实例
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown，包含表格、任务列表、删除线等
		),
	)
	
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetMarkdownHTML 读取markdown文件并转换为HTML
func GetMarkdownHTML(filePath string) (string, error) {
	content, err := ReadMarkdownFile(filePath)
	if err != nil {
		return "", err
	}
	return ConvertMarkdownToHTML(content)
}
