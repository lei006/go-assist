package multicast

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"time"
)

// MulticastReceiver 单网卡组播接收器
type MulticastReceiver struct {
	conn           *net.UDPConn
	iface          *net.Interface // 绑定的网卡
	isRunning      bool
	mu             sync.RWMutex
	onMessage      func(ReceiverMessage)
	onError        func(error)
	enableLoopback bool
	readTimeout    time.Duration
	bufferSize     int
	multicastAddr  string
	ifaceName      string
}

// ReceiverMessage 多播消息结构
type ReceiverMessage struct {
	Data          []byte
	MulticastAddr string
	FromAddr      net.Addr
	Timestamp     time.Time
	Size          int
	InterfaceName string // 接收消息的网卡名
}

// NewSingleReceiver 创建单网卡接收器
func NewSingleReceiver(ifaceName, multicastAddr string) *MulticastReceiver {
	return &MulticastReceiver{
		ifaceName:      ifaceName,
		multicastAddr:  multicastAddr,
		bufferSize:     4096,
		readTimeout:    30 * time.Second,
		enableLoopback: true,
		onMessage:      func(ReceiverMessage) {},
		onError:        func(error) {},
	}
}

// 启动单网卡接收器
func (r *MulticastReceiver) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		return fmt.Errorf("接收器已在运行")
	}

	addr, err := net.ResolveUDPAddr("udp", r.multicastAddr)
	if err != nil {
		return fmt.Errorf("解析多播地址失败: %w", err)
	}

	var iface *net.Interface
	if r.ifaceName != "" {
		iface, err = net.InterfaceByName(r.ifaceName)
		if err != nil {
			return fmt.Errorf("获取网络接口 %s 失败: %w", r.ifaceName, err)
		}
	}

	// 创建多播连接
	conn, err := net.ListenMulticastUDP("udp", iface, addr)
	if err != nil {
		return fmt.Errorf("创建多播连接失败: %w", err)
	}

	// 设置多播环回（跨平台兼容）
	if err := setMulticastLoopbackCrossPlatform(conn, r.enableLoopback); err != nil {
		conn.Close()
		return fmt.Errorf("设置多播环回失败: %w", err)
	}

	// 设置缓冲区大小
	if r.bufferSize > 0 {
		if err := conn.SetReadBuffer(r.bufferSize); err != nil {
			conn.Close()
			return fmt.Errorf("设置接收缓冲区失败: %w", err)
		}
	}

	// 设置读取超时
	if r.readTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(r.readTimeout)); err != nil {
			conn.Close()
			return fmt.Errorf("设置读取超时失败: %w", err)
		}
	}

	r.conn = conn
	r.iface = iface
	r.isRunning = true

	// 启动接收循环
	go r.receiveLoop()

	return nil
}

// 接收消息循环
func (r *MulticastReceiver) receiveLoop() {
	buffer := make([]byte, r.bufferSize)

	for r.isRunning {
		n, fromAddr, err := r.conn.ReadFromUDP(buffer)
		if err != nil {
			if !r.isRunning {
				break
			}

			// 过滤超时错误：不触发回调
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 仅重置超时，不调用错误回调
				if r.readTimeout > 0 {
					r.conn.SetReadDeadline(time.Now().Add(r.readTimeout))
				}
				continue
			}

			r.mu.RLock()
			errCallback := r.onError
			r.mu.RUnlock()
			if errCallback != nil {
				errCallback(fmt.Errorf("读取消息失败: %w", err))
			}

			// 重置超时
			if r.readTimeout > 0 {
				r.conn.SetReadDeadline(time.Now().Add(r.readTimeout))
			}
			continue
		}

		if n <= 0 {
			continue
		}

		// 复制数据避免缓冲区覆盖
		data := make([]byte, n)
		copy(data, buffer[:n])

		// 构造消息
		msg := ReceiverMessage{
			Data:          data,
			MulticastAddr: r.multicastAddr,
			FromAddr:      fromAddr,
			Timestamp:     time.Now(),
			Size:          n,
		}

		// 调用回调
		r.mu.RLock()
		msgCallback := r.onMessage
		r.mu.RUnlock()
		if msgCallback != nil {
			msgCallback(msg)
		}

		// 重置超时
		if r.readTimeout > 0 {
			r.conn.SetReadDeadline(time.Now().Add(r.readTimeout))
		}
	}
}

// 单网卡接收器停止方法
func (r *MulticastReceiver) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRunning {
		return
	}

	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
	r.isRunning = false
}

// 设置消息回调
func (r *MulticastReceiver) SetOnMessage(callback func(ReceiverMessage)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onMessage = callback
}

// 设置错误回调
func (r *MulticastReceiver) SetOnError(callback func(error)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onError = callback
}

// 消息字符串表示
func (m *ReceiverMessage) String() string {
	return fmt.Sprintf(
		"[%s] 网卡[%s] 来自[%s] (%d字节): %s",
		m.Timestamp.Format("15:04:05"),
		m.InterfaceName,
		m.FromAddr,
		m.Size,
		string(m.Data),
	)
}

// 消息十六进制表示
func (m *ReceiverMessage) HexString() string {
	return hex.Dump(m.Data)
}
