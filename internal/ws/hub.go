package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client 代表一个WebSocket连接
type Client struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	TeamID int    `json:"team_id"`
	Conn   *websocket.Conn
	Send   chan []byte
}

// Message WebSocket消息结构
type Message struct {
	Type      string          `json:"type"`
	UserID    int             `json:"user_id,omitempty"`
	Username  string          `json:"username,omitempty"`
	TeamID    int             `json:"team_id,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Timestamp int64           `json:"timestamp"`
}

// LocationUpdate 位置更新消息
type LocationUpdate struct {
	UserID    int     `json:"user_id"`
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Distance  float64 `json:"distance"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Type     string `json:"type"`
}

// Hub 管理所有WebSocket连接
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.RWMutex
}

var GlobalHub *Hub

func init() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[WS] 客户端已连接: user=%d team=%d", client.UserID, client.TeamID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("[WS] 客户端已断开: user=%d team=%d", client.UserID, client.TeamID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToTeam 向车队内所有成员广播消息
func (h *Hub) BroadcastToTeam(teamID int, msg Message) {
	msg.Timestamp = time.Now().Unix()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[WS] 序列化消息失败: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.TeamID == teamID {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// BroadcastToUser 向指定用户发送消息
func (h *Hub) BroadcastToUser(userID int, msg Message) {
	msg.Timestamp = time.Now().Unix()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[WS] 序列化消息失败: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// GetOnlineCount 获取在线用户数
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetTeamOnlineCount 获取车队在线人数
func (h *Hub) GetTeamOnlineCount(teamID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for client := range h.clients {
		if client.TeamID == teamID {
			count++
		}
	}
	return count
}

// RegisterClient 注册客户端
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient 注销客户端
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}
