package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

// runWebSocketClient 运行WebSocket客户端
func runWebSocketClient(filePath, port string) {
	// 构建WebSocket URL
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:" + port,
		Path:   "/ws",
	}
	q := u.Query()
	q.Add("file", filePath)
	u.RawQuery = q.Encode()

	log.Printf("连接到 WebSocket 服务器: %s", u.String())

	// 连接WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("WebSocket 连接失败: %v", err)
	}
	defer conn.Close()

	// 接收信号以退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 接收消息的goroutine
	go func() {
		defer conn.Close()
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("读取消息失败: %v", err)
				return
			}
			log.Printf("收到消息类型: %d, 长度: %d 字节", messageType, len(message))
			// 只打印前500个字符
			if len(message) > 500 {
				log.Printf("消息内容: %s...", string(message[:500]))
			} else {
				log.Printf("消息内容: %s", string(message))
			}
		}
	}()

	log.Println("WebSocket 客户端已连接，等待消息...")

	// 保持连接
	select {
	case <-interrupt:
		log.Println("收到中断信号，关闭连接...")
		return
	case <-time.After(30 * time.Second):
		log.Println("30秒超时，关闭连接...")
		return
	}
}

func main() {
	// 自定义命令行参数解析，支持文件路径作为第一个参数
	var port string = "8080"
	var runClient bool = false
	var filePath string

	// 解析命令行参数
	for i, arg := range os.Args[1:] {
		if arg == "-port" && i+1 < len(os.Args[1:]) {
			port = os.Args[1:][i+1]
		} else if arg == "-client" {
			runClient = true
		} else if filePath == "" && !strings.HasPrefix(arg, "-") {
			// 第一个非选项参数是文件路径
			filePath = arg
		}
	}

	if filePath == "" {
		log.Fatal("错误：必须指定Markdown文件路径作为第一个参数")
	}

	log.Printf("启动Markdown实时预览服务，监听文件: %s，端口: %s", filePath, port)

	// 初始化WebSocket管理器
	wsManager := NewWebSocketManager()
	wsManager.Start()

	// 初始化文件监听器
	watcher, err := NewFileWatcher(filePath, func() {
		// 文件变化时，转换为HTML并广播给所有客户端
		html, err := GetMarkdownHTML(filePath)
		if err != nil {
			log.Printf("转换Markdown失败: %v", err)
			return
		}
		wsManager.Broadcast([]byte(html))
	})
	if err != nil {
		log.Fatalf("创建文件监听器失败: %v", err)
	}
	watcher.Start()

	// 初始化gin路由
	r := gin.Default()

	// WebSocket路由 - 不再需要URL参数，直接使用服务器配置的文件路径
	r.GET("/ws", func(c *gin.Context) {
		// 将文件路径传递给WebSocket处理函数
		c.Set("filePath", filePath)
		wsManager.HandleConnection(c)
	})

	// 静态文件服务
	r.Static("/static", "./static")

	// 主页面路由 - 直接返回HTML，不再需要URL参数
	r.GET("/", func(c *gin.Context) {
		// 直接返回index.html文件内容
		content, err := os.ReadFile("./static/index.html")
		if err != nil {
			c.String(500, "无法读取index.html文件: %v", err)
			return
		}
		c.Data(200, "text/html; charset=utf-8", content)
	})

	// 后台打开浏览器
	go func() {
		// 等待服务器启动
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost:" + port)
	}()

	// 如果需要启动客户端，在后台运行
	if runClient {
		go runWebSocketClient(filePath, port)
	}

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
