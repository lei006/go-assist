package multicast

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// MultiInterfaceReceiver 多网卡组播接收器管理器
type MultiInterfaceReceiver struct {
	receivers     []*MulticastReceiver
	isRunning     bool
	mu            sync.RWMutex
	onMessage     func(ReceiverMessage)
	onError       func(error)
	multicastIP   string
	multicastPort uint16

	MulticastAddr string
}

// NewMultiReceiver 创建多网卡组播接收器
func NewMultiReceiver(multicastIP string, multicastPort uint16) (*MultiInterfaceReceiver, error) {
	// 验证组播地址
	ip := net.ParseIP(multicastIP)
	if ip == nil || !ip.IsMulticast() {
		return nil, errors.New("无效的组播地址")
	}

	// 验证端口
	if multicastPort < 1 || multicastPort > 65535 {
		return nil, errors.New("端口必须在1-65535范围内")
	}

	return &MultiInterfaceReceiver{
		multicastIP:   multicastIP,
		multicastPort: multicastPort,
		onMessage:     func(ReceiverMessage) {},
		onError:       func(error) {},
		MulticastAddr: fmt.Sprintf("%s:%d", multicastIP, multicastPort),
	}, nil
}

// Start 启动多网卡接收（所有可用网卡）
func (mr *MultiInterfaceReceiver) Start() error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if mr.isRunning {
		return errors.New("多网卡接收器已在运行")
	}

	// 获取所有可用网卡
	ifaces, err := getAvailableInterfaces()
	if err != nil {
		return fmt.Errorf("获取网卡列表失败: %w", err)
	}

	if len(ifaces) == 0 {
		return errors.New("未找到可用的网络接口")
	}

	multicastAddr := fmt.Sprintf("%s:%d", mr.multicastIP, mr.multicastPort)
	
	// 为每个网卡创建接收器
	for _, iface := range ifaces {
		receiver := NewSingleReceiver(iface.Name, multicastAddr)
		
		// 设置回调（转发到多网卡接收器的回调）
		receiver.SetOnMessage(func(msg ReceiverMessage) {
			msg.InterfaceName = iface.Name
			mr.mu.RLock()
			callback := mr.onMessage
			mr.mu.RUnlock()
			callback(msg)
		})

		receiver.SetOnError(func(err error) {
			wrappedErr := fmt.Errorf("网卡[%s]接收错误: %w", iface.Name, err)
			mr.mu.RLock()
			callback := mr.onError
			mr.mu.RUnlock()
			callback(wrappedErr)
		})

		// 启动单网卡接收器
		if err := receiver.Start(); err != nil {
			mr.stop() // 启动失败时清理已创建的接收器
			return fmt.Errorf("启动网卡[%s]接收器失败: %w", iface.Name, err)
		}

		mr.receivers = append(mr.receivers, receiver)
	}

	mr.isRunning = true
	return nil
}

// 多网卡接收器停止方法
func (mr *MultiInterfaceReceiver) Stop() {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.stop()
}

// 内部停止方法（无锁版本）
func (mr *MultiInterfaceReceiver) stop() {
	if !mr.isRunning {
		return
	}

	// 停止所有单网卡接收器
	for _, receiver := range mr.receivers {
		receiver.Stop()
	}
	mr.receivers = nil
	mr.isRunning = false
}

// 设置消息回调
func (mr *MultiInterfaceReceiver) OnMessage(callback func(ReceiverMessage)) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.onMessage = callback
}

// 设置错误回调
func (mr *MultiInterfaceReceiver) OnError(callback func(error)) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.onError = callback
}

// 检查是否运行中
func (mr *MultiInterfaceReceiver) IsRunning() bool {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	return mr.isRunning
}
    