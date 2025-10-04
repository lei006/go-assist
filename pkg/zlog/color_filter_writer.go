package zlog

import (
	"io"
	"regexp"
)

// 定义ANSI颜色转义符的正则表达式
var ansiColorRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// ColorFilterWriter 用于过滤ANSI颜色转义符的Writer
type ColorFilterWriter struct {
	writer io.Writer // 底层真实的Writer（如文件）
}

// Write 实现io.Writer接口，写入时过滤颜色转义符
func (w *ColorFilterWriter) Write(p []byte) (n int, err error) {
	// 过滤ANSI颜色转义符
	filtered := ansiColorRegex.ReplaceAll(p, []byte(""))
	// 写入过滤后的内容到底层Writer
	return w.writer.Write(filtered)
}

// NewColorFilterWriter 创建一个带颜色过滤功能的Writer
func NewColorFilterWriter(w io.Writer) *ColorFilterWriter {
	return &ColorFilterWriter{writer: w}
}
