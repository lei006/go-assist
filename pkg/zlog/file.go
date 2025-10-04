package zlog

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/ztime"
)

var LogMaxDurationDate = 15

// 提取共用的清理函数：删除指定目录下的过期文件
// 提取共用的清理函数：删除指定目录下的过期文件
func cleanExpiredFiles(archiveDir string, fileExt string, saveDay int) error {
	now := ztime.UnixMicro(ztime.Clock())
	return filepath.Walk(archiveDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// 解析文件名中的日期并判断是否过期
		date, err := ztime.Parse(strings.TrimSuffix(filepath.Base(path), fileExt), "Y-m-d")
		if err == nil && date.AddDate(0, 0, saveDay).Before(now) {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func openFile(filepa string, archive bool) (file *zfile.MemoryFile, fileName, fileDir string, err error) {
	fullPath := zfile.RealPath(filepa)
	fileDir, fileName = filepath.Split(fullPath)
	opt := []zfile.MemoryFileOption{zfile.MemoryFileAutoFlush(1)}
	if archive {
		ext := filepath.Ext(fileName)
		base := strings.TrimSuffix(fileName, ext)
		fileDir = zfile.RealPathMkdir(fileDir+base, true)
		fullPath = fileDir + fileName
		lastArchiveName := ztime.Now("Y-m-d") // 初始化为当前日期

		// 程序启动时立即执行一次清理（调用共用函数）
		//go cleanExpiredFiles(fileDir, ext, LogMaxDurationDate)

		opt = append(opt, zfile.MemoryFileFlushBefore(func(f *zfile.MemoryFile) error {
			archiveName := ztime.Now("Y-m-d")
			if lastArchiveName != archiveName {
				// 日期切换时再次执行清理（调用共用函数）
				go cleanExpiredFiles(fileDir, ext, LogMaxDurationDate)
				lastArchiveName = archiveName
			}
			fileName = archiveName + ext
			f.SetName(fileDir + fileName)
			return nil
		}))
	}
	f := zfile.NewMemoryFile(fullPath, opt...)
	return f, fileName, fileDir, nil
}

// SetFile Setting log file output
func (log *Logger) SetFile(filepath string, archive ...bool) {

	logArchive := len(archive) > 0 && archive[0]

	log.setLogfile(filepath, logArchive)
}

func (log *Logger) setLogfile(filepath string, archive bool) {

	fileObj, fileName, fileDir, _ := openFile(filepath, archive)
	log.mu.Lock()
	log.CloseFile()
	log.file = fileObj
	log.fileDir = fileDir
	log.fileName = fileName
	fileWriter := NewColorFilterWriter(fileObj)
	if log.fileAndStdout {
		log.out = io.MultiWriter(log.file, log.out)
	} else {
		log.out = fileWriter
	}
	log.mu.Unlock()
}

func (log *Logger) Discard() {
	log.mu.Lock()
	log.out = ioutil.Discard
	if log.file != nil {
		_ = log.file.Close()
	}
	log.level = LogNot
	log.mu.Unlock()
}

func (log *Logger) SetSaveFile(filepath string, archive ...bool) {

	log.SetFile(filepath, archive...)
	log.mu.Lock()
	defer log.mu.Unlock()

	log.fileAndStdout = true
	//log.out = io.MultiWriter(log.file, os.Stdout)
	// 控制台输出保留颜色，文件输出过滤颜色
	log.out = io.MultiWriter(os.Stdout, NewColorFilterWriter(log.file))

}

func (log *Logger) CloseFile() {
	if log.file != nil {
		_ = log.file.Close()
		log.file = nil
		log.out = os.Stdout
	}
}

// 为Logger添加清理方法，方便通过日志实例清理其关联的日志文件
func (log *Logger) CleanExpiredLogs(saveDay int) error {
	if log.fileDir == "" {
		return os.ErrInvalid
	}
	ext := filepath.Ext(log.fileName)
	if saveDay <= 0 {
		saveDay = LogMaxDurationDate
	}
	return cleanExpiredFiles(log.fileDir, ext, saveDay)
}
