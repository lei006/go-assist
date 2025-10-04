package multicast

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// IsMulticastAddress 判断IP地址是否为多播地址
// 处理所有常见格式：标准IPv4、标准IPv6、IPv4映射到IPv6的地址
func IsMulticastAddress(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// 关键修复：先通过To4()提取纯IPv4地址（处理映射格式）
	if ipv4 := ip.To4(); ipv4 != nil {
		// IPv4多播地址范围：224.0.0.0 - 239.255.255.255
		// 直接检查IPv4地址的第一个字节
		return ipv4[0] >= 224 && ipv4[0] <= 239
	}

	// 处理IPv6地址
	if len(ip) == net.IPv6len {
		// IPv6多播地址以0xff开头（ff00::/8范围）
		return ip[0] == 0xff
	}

	// 非标准IP地址（如长度异常）
	return false
}

// setSocketOptionInt 设置socket选项的整数参数
// 这是一个内部辅助函数，封装了系统调用
func setSocketOptionInt(fd syscall.Handle, level, opt, value int) error {
	return syscall.SetsockoptInt(fd, level, opt, value)
}

// 跨平台设置多播环回（核心修复：避免使用conn.File()）
func setMulticastLoopbackCrossPlatform(conn *net.UDPConn, enable bool) error {
	// 使用SyscallConn获取底层套接字控制接口
	sysConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	flag := 0
	if enable {
		flag = 1
	}

	var controlErr error
	// 通过Control方法安全操作底层套接字
	err = sysConn.Control(func(fd uintptr) {
		// 根据操作系统类型处理
		if runtime.GOOS == "windows" {
			// Windows系统调用
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			if localAddr.IP.To4() != nil {
				controlErr = windows.SetsockoptInt(
					windows.Handle(fd),
					windows.IPPROTO_IP,
					windows.IP_MULTICAST_LOOP,
					flag,
				)
			} else {
				controlErr = windows.SetsockoptInt(
					windows.Handle(fd),
					windows.IPPROTO_IPV6,
					windows.IPV6_MULTICAST_LOOP,
					flag,
				)
			}
		} else {
			// Linux/Unix系统调用
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			if localAddr.IP.To4() != nil {
				controlErr = syscall.SetsockoptInt(
					syscall.Handle(fd),
					syscall.IPPROTO_IP,
					syscall.IP_MULTICAST_LOOP,
					flag,
				)
			} else {
				controlErr = syscall.SetsockoptInt(
					syscall.Handle(fd),
					syscall.IPPROTO_IPV6,
					syscall.IPV6_MULTICAST_LOOP,
					flag,
				)
			}
		}
	})

	if err != nil {
		return err
	}
	return controlErr
}

// 跨平台选择合适的网络接口
func selectSuitableInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}
		return &iface, nil
	}

	return nil, fmt.Errorf("未找到可用的网络接口（需非环回且已启动）")
}

// GetLocalInterfaces 获取所有可用的非环回网卡及其IP
func GetLocalInterfaces() (map[string][]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网卡列表失败: %w", err)
	}

	result := make(map[string][]string)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		var ips []string
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}
			// 只保留IPv4地址
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}

		if len(ips) > 0 {
			result[iface.Name] = ips
		}
	}

	return result, nil
}

// 根据IP查找对应的网卡
func getInterfaceByIP(ip net.IP) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.Equal(ip) {
				return &iface, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到IP为%s的网卡", ip.String())
}

// 配置组播套接字参数（修复跨系统 IP_MULTICAST_IF 格式问题）
func configureMulticastSocket(conn *net.UDPConn, iface *net.Interface, ttl int, localIP net.IP) error {
	// 获取底层套接字文件描述符
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return fmt.Errorf("获取原始连接失败: %w", err)
	}

	var setErr error
	err = rawConn.Control(func(fd uintptr) {
		// 1. 先设置 TTL（跨系统格式一致，无兼容性问题）
		setErr = syscall.SetsockoptInt(
			syscall.Handle(fd),
			syscall.IPPROTO_IP,
			syscall.IP_MULTICAST_TTL,
			ttl,
		)
		if setErr != nil {
			setErr = fmt.Errorf("设置TTL失败: %w", setErr)
			return
		}

		// 2. 根据操作系统适配 IP_MULTICAST_IF 参数格式
		if iface == nil || localIP == nil {
			setErr = errors.New("网卡或本地IP为空，无法绑定组播发送网卡")
			return
		}

		switch runtime.GOOS {
		case "windows":
			// Windows：需要传入 网卡IPv4地址的32位整数（网络字节序）
			// 步骤1：将 net.IP 转换为 IPv4 地址（确保是4字节）
			ipv4 := localIP.To4()
			if ipv4 == nil {
				setErr = errors.New("仅支持IPv4地址绑定（Windows）")
				return
			}

			// 步骤2：将4字节IPv4地址转换为32位整数（网络字节序 → 主机字节序）
			// 注意：syscall.SetsockoptInt 会自动转为网络字节序，此处直接用主机字节序
			ipInt := *(*uint32)(unsafe.Pointer(&ipv4[0]))

			// 步骤3：设置 IP_MULTICAST_IF（参数为32位IP整数）
			setErr = syscall.SetsockoptInt(
				syscall.Handle(fd),
				syscall.IPPROTO_IP,
				syscall.IP_MULTICAST_IF,
				int(ipInt),
			)

		case "linux":
			// Linux：传入 网卡索引（iface.Index）
			setErr = syscall.SetsockoptInt(
				syscall.Handle(fd),
				syscall.IPPROTO_IP,
				syscall.IP_MULTICAST_IF,
				iface.Index,
			)

		default:
			setErr = fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
			return
		}

		// 3. 处理绑定失败错误
		if setErr != nil {
			setErr = fmt.Errorf("绑定网卡失败（%s）: %w", runtime.GOOS, setErr)
			return
		}
	})

	if err != nil {
		return fmt.Errorf("执行系统调用失败: %w", err)
	}
	return setErr
}



// 获取所有可用的非环回网卡
func getAvailableInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []net.Interface
	for _, iface := range ifaces {
		// 跳过禁用或环回网卡
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 检查是否有IPv4地址
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil { // 只保留有IPv4地址的网卡
				result = append(result, iface)
				break
			}
		}
	}

	return result, nil
}
