package main

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// FileWatcher 文件监听器
type FileWatcher struct {
	watcher  *fsnotify.Watcher
	filePath string
	OnChange func()
}

// NewFileWatcher 创建新的文件监听器
func NewFileWatcher(filePath string, onChange func()) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// 监听文件所在目录
	dir := filepath.Dir(filePath)
	if err := watcher.Add(dir); err != nil {
		return nil, err
	}

	return &FileWatcher{
		watcher:  watcher,
		filePath: filePath,
		OnChange: onChange,
	}, nil
}

// Start 开始监听文件变化
func (fw *FileWatcher) Start() {
	go func() {
		defer fw.watcher.Close()
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return
				}
				// 检查是否是目标文件的变化
				if event.Name == fw.filePath && (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
					log.Printf("File changed: %s", fw.filePath)
					if fw.OnChange != nil {
						fw.OnChange()
					}
				}
			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()
}
