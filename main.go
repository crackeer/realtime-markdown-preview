package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFS embed.FS

// openBrowser 自动打开浏览器
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	case "windows": // Windows
		err = exec.Command("cmd", "/c", "start", url).Start()
	default: // Linux and others
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		log.Printf("自动打开浏览器失败: %v", err)
	}
}

var (
	channelMap = make(map[string]map[string]chan string)
	port       = "8080"
	dirPath    string
	locker     *sync.Mutex
)

func init() {
	channelMap = make(map[string]map[string]chan string)
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	locker = &sync.Mutex{}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("错误：必须指定Markdown文件夹路径作为第一个参数")
	}
	dirPath = os.Args[1]

	log.Printf("启动Markdown实时预览服务，监听文件夹: %s，端口: %s", dirPath, port)

	// 初始化gin路由
	r := gin.Default()

	r.GET("/html/*filepath", getHTML)

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost:" + port)
	}()

	// 404路由 - 直接返回index.html
	r.NoRoute(func(c *gin.Context) {

		// 直接返回index.html文件内容 - 使用嵌入的资源
		content, err := staticFS.ReadFile("static/index.html")
		if err != nil {
			c.String(500, "无法读取index.html文件: %v", err)
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func getHTML(c *gin.Context) {

	// 获取文件路径
	relativePath := c.Param("filepath")
	fmt.Printf("请求路径: %s\n", relativePath)
	if relativePath == "/" {
		c.String(http.StatusBadRequest, "缺少文件路径")
		return
	}

	fullPath := filepath.Join(dirPath, relativePath)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "文件不存在")
		return
	}

	// 检查是否为.md文件
	if filepath.Ext(fullPath) != ".md" {
		c.String(http.StatusBadRequest, "只支持.md文件")
		return
	}
	html, err := GetMarkdownHTML(fullPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "转换Markdown失败: %v", err)
		return
	}
	c.String(http.StatusOK, html)
	return
}
