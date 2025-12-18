package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketManager WebSocket管理器
type WebSocketManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	upgrader   websocket.Upgrader
}

// NewWebSocketManager 创建新的WebSocket管理器
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源
			},
		},
	}
}

// HandleConnection 处理WebSocket连接请求
func (wm *WebSocketManager) HandleConnection(c *gin.Context) {
	conn, err := wm.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// 将新客户端添加到客户端列表
	wm.clients[conn] = true
	log.Printf("New WebSocket client connected, total clients: %d", len(wm.clients))

	// 发送初始HTML内容 - 从上下文获取文件路径
	if filePath, exists := c.Get("filePath"); exists {
		if fp, ok := filePath.(string); ok {
			html, err := GetMarkdownHTML(fp)
			if err == nil {
				conn.WriteMessage(websocket.TextMessage, []byte(html))
			} else {
				log.Printf("转换Markdown失败: %v", err)
			}
		}
	}

	// 读取客户端消息（保持连接）
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			delete(wm.clients, conn)
			log.Printf("WebSocket client disconnected, total clients: %d", len(wm.clients))
			break
		}
	}
}

// Start 启动WebSocket管理器
func (wm *WebSocketManager) Start() {
	go func() {
		for {
			// 从广播通道接收消息
			message := <-wm.broadcast
			// 向所有客户端发送消息
			for client := range wm.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("WebSocket write error: %v", err)
					client.Close()
					delete(wm.clients, client)
				}
			}
		}
	}()
}

// Broadcast 向所有客户端广播消息
func (wm *WebSocketManager) Broadcast(message []byte) {
	wm.broadcast <- message
}
