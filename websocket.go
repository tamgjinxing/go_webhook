package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn *websocket.Conn
	url  string
	stop chan bool // 用于通知停止监听
}

// NewWebSocketClient 创建一个新的 WebSocket 客户端
func NewWebSocketClient(relayUrl string) (*WebSocketClient, error) {
	client := &WebSocketClient{
		url:  relayUrl,
		stop: make(chan bool), // 初始化 stop channel
	}
	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial(relayUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	return client, nil
}

// SendMessage 发送消息到 WebSocket 服务器
func (client *WebSocketClient) SendMessage(message string) error {
	err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	log.Printf("Message sent to %s: %s", client.url, message)
	return nil
}

// ReadMessage 从 WebSocket 服务器读取消息
func (client *WebSocketClient) ReadMessage() (string, error) {
	_, message, err := client.conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("failed to read message: %w", err)
	}
	return string(message), nil
}

// Close 关闭 WebSocket 连接
func (client *WebSocketClient) Close() error {
	err := client.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close WebSocket: %w", err)
	}
	log.Printf("Connection to %s closed.", client.url)
	return nil
}

// ListenAndServe 持续监听 WebSocket 消息
func (client *WebSocketClient) ListenAndServe(wg *sync.WaitGroup) {
	defer wg.Done() // 结束后通知 WaitGroup

	for {
		select {
		case <-client.stop:
			log.Printf("Stopping listener for %s.", client.url)
			return
		default:
			// 从 WebSocket 服务器读取消息
			message, err := client.ReadMessage()
			if err != nil {
				log.Printf("Error reading message from %s: %v", client.url, err)
				break
			}
			log.Printf("Received message from %s: %s", client.url, message)
		}
	}
}

// Stop 发送停止信号，关闭监听
func (client *WebSocketClient) Stop() {
	close(client.stop)
}

// StartRelayConnections 启动多个 WebSocket 连接并持续监听
func StartRelayConnections(relayUrls []interface{}, eventString string) {
	var wg sync.WaitGroup
	var clients []*WebSocketClient

	// 创建并启动每个 WebSocket 客户端
	for _, url := range relayUrls {
		wg.Add(1) // 每个连接需要一个 goroutine

		client, err := NewWebSocketClient(url.(string))
		if err != nil {
			log.Printf("Error creating WebSocket client for %s: %v", url, err)
			continue
		}

		clients = append(clients, client)

		// 启动监听消息的 goroutine
		go client.ListenAndServe(&wg)

		// 发送消息
		err = client.SendMessage(eventString)
		if err != nil {
			log.Printf("Error sending message to %s: %v", url, err)
		}

		time.Sleep(5 * time.Second)
		// 停止监听并关闭连接
		client.Stop()
		err = client.Close()
		if err != nil {
			log.Printf("Error closing WebSocket connection for %s: %v", url, err)
		}
	}

	// 等待所有的 WebSocket 连接完成
	wg.Wait()
}
