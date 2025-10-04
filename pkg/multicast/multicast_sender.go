package multicast

import (
	"errors"
	"fmt"
	"net"
)

// MulticastSender 组播发送器，支持动态指定网卡发送
type MulticastSender struct {
	multicastAddr string // 组播组地址
	multicastPort uint16 // 组播端口
	ttl           int    // 固定TTL值
}

// NewMulticastSender 创建组播发送器实例
func NewMulticastSender(multicastAddr string, port uint16) (*MulticastSender, error) {
	// 验证组播地址有效性
	ip := net.ParseIP(multicastAddr)
	if ip == nil {
		return nil, errors.New("无效的组播地址格式")
	}
	if !ip.IsMulticast() {
		return nil, errors.New("不是有效的组播地址（需在224.0.0.0/4或ff00::/8范围内）")
	}

	// 验证端口范围
	if port == 0 {
		return nil, errors.New("端口必须在1-65535范围内")
	}

	return &MulticastSender{
		multicastAddr: multicastAddr,
		multicastPort: port,
		ttl:           255,
	}, nil
}

// Send 发送字符串消息到组播组
func (s *MulticastSender) Send(message string, localIP string) (map[string]error, error) {
	return s.SendBytes([]byte(message), localIP)
}

// SendBytes 发送字节数据到组播组
func (s *MulticastSender) SendBytes(data []byte, localIP string) (map[string]error, error) {
	targetAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.multicastAddr, s.multicastPort))
	if err != nil {
		return nil, fmt.Errorf("解析组播地址失败: %w", err)
	}

	var sendIPs []string
	if localIP != "" {
		ip := net.ParseIP(localIP)
		if ip == nil {
			return nil, errors.New("无效的本地IP地址")
		}
		// 验证指定的IP是否为IPv4
		if ip.To4() == nil {
			return nil, errors.New("只支持IPv4地址")
		}
		sendIPs = []string{localIP}
	} else {
		ifaces, err := GetLocalInterfaces()
		if err != nil {
			return nil, fmt.Errorf("获取网卡列表失败: %w", err)
		}
		if len(ifaces) == 0 {
			return nil, errors.New("未找到可用网卡")
		}

		// 收集所有网卡的第一个IPv4地址
		for _, ips := range ifaces {
			for _, ipStr := range ips {
				ip := net.ParseIP(ipStr)
				// 只保留IPv4地址
				if ip != nil && ip.To4() != nil {
					sendIPs = append(sendIPs, ipStr)
					break // 只取第一个IPv4
				}
			}
		}

		if len(sendIPs) == 0 {
			return nil, errors.New("未找到可用的IPv4地址")
		}
	}

	result := make(map[string]error)
	for _, ip := range sendIPs {
		iface, err := getInterfaceByIP(net.ParseIP(ip))
		ifaceName := "unknown"
		if iface != nil {
			ifaceName = iface.Name
		}

		err = s.sendByIP(data, targetAddr, ip)
		result[ifaceName] = err
	}

	return result, nil
}

// 内部方法：通过指定IP发送数据（修改参数，传入 localIP 供 Windows 绑定使用）
func (s *MulticastSender) sendByIP(data []byte, targetAddr *net.UDPAddr, localIP string) error {
	// 解析本地IP（转为 net.IP 类型，供后续绑定使用）
	localIPObj := net.ParseIP(localIP)
	if localIPObj == nil || localIPObj.To4() == nil {
		return errors.New("无效的本地IPv4地址")
	}

	// 解析本地地址（端口0表示随机分配）
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:0", localIP))
	if err != nil {
		return fmt.Errorf("解析本地地址失败: %w", err)
	}

	// 创建临时UDP连接
	conn, err := net.DialUDP("udp", localAddr, targetAddr)
	if err != nil {
		return fmt.Errorf("创建UDP连接失败: %w", err)
	}
	defer conn.Close() // 确保连接关闭

	// 查找IP对应的网卡（用于 Linux 绑定索引）
	iface, err := getInterfaceByIP(localIPObj)
	if err != nil {
		return fmt.Errorf("获取网卡信息失败: %w", err)
	}

	// 配置组播参数（传入 localIPObj 供 Windows 处理）
	if err := configureMulticastSocket(conn, iface, s.ttl, localIPObj); err != nil {
		return fmt.Errorf("配置组播参数失败: %w", err)
	}

	// 发送数据
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("发送数据失败: %w", err)
	}
	return nil
}
