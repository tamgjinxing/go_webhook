package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// 处理 /webhook 路径的 POST 请求
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为 POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 打印所有请求头信息
	for key, values := range r.Header {
		fmt.Printf("Header Key: %s, Header Values: %v\n", key, values)
	}

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 打印请求体内容（可以替换为你自己的处理逻辑）
	fmt.Printf("Received webhook data: %s\n", string(body))

	ParseWebhookDataAndCreated9735Event(string(body))

	// 在响应中返回一个确认消息
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	// 启动 HTTP 服务器并监听 3002 端口
	fmt.Println("Server starting on port 3002...")
	if err := http.ListenAndServe("0.0.0.0:3002", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	select {}
}
