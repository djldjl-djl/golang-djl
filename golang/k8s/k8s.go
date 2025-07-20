package k8s

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type K8spod struct {
	Client *kubernetes.Clientset
}

func NewK8spod(client *kubernetes.Clientset) *K8spod {
	return &K8spod{
		Client: client,
	}
}

func (p *K8spod) Deletepod(ctx context.Context, namespace string, name string) error {
	err := p.Client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}
func (p *K8spod) Selectpod(namespace string, ctx context.Context) (error, []corev1.Pod) {
	// pod, err := p.Client.CoreV1().Pods("dinjialin").Get(context.Background(), "book-69774bd8bf-psz74", metav1.GetOptions{})
	pod, err := p.Client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err, nil
	}
	return nil, pod.Items
}

func (p *K8spod) GetLogPod(ctx context.Context, namespace, pname string) string {
	log := p.Client.CoreV1().Pods(namespace).GetLogs(pname, &corev1.PodLogOptions{})
	data, err := log.Stream(ctx)
	if err != nil {
		return "获取日志流失败"
	}
	defer data.Close()
	logs, err := io.ReadAll(data)
	if err != nil {
		return "读取日志失败"
	}
	return string(logs)
}
func (p *K8spod) WSGetLogPod(Conn *websocket.Conn, ctx context.Context, namespace, pname string) {
	log := p.Client.CoreV1().Pods(namespace).GetLogs(pname, &corev1.PodLogOptions{
		Follow: true,
	})
	data, err := log.Stream(ctx)
	if err != nil {
		Conn.WriteMessage(websocket.TextMessage, []byte("连接日志失败: "+err.Error()))
		return
	}
	defer data.Close()
	logs := bufio.NewReader(data)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := logs.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// 日志流结束，正常退出
					Conn.WriteMessage(websocket.TextMessage, []byte("日志读取结束"))
				} else {
					fmt.Println("读取日志出错:", err)
					Conn.WriteMessage(websocket.TextMessage, []byte("日志读取失败: "+err.Error()))
				}
				return
			}
			Conn.WriteMessage(websocket.TextMessage, []byte(line))
		}
	}
}

type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// 交互的结构体，接管输入和输出
type wsStream struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

func (t *wsStream) Read(p []byte) (int, error) {
	//从ws中读取消息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("读取消息错误: %v", err)
		return 0, err
	}
	//反序列化
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("读取消息语法错误: %v", err)
		return 0, err
	}
	//逻辑判断
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("消息类型错误'%s'", msg.Operation)
		return 0, fmt.Errorf("消息类型错误'%s'", msg.Operation)
	}
}

// 写数据的方法,拿到apiserver的返回内容，向web端输出
func (t *wsStream) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("写消息语法错误: %v", err)
		return 0, err
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("写消息错误: %v", err)
		return 0, err
	}

	return len(p), nil
}

func (p *K8spod) WebShell(config *rest.Config, Conn *websocket.Conn, namespace, pname, cname string) error {
	req := p.Client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pname).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: cname,
			Command:   []string{"/bin/sh"},
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	wsStream := &wsStream{
		wsConn:   Conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  wsStream,
		Stdout: wsStream,
		Stderr: wsStream,
		Tty:    true,
	})
}
