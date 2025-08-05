package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SSHHost 表示一个 SSH 配置条目
type SSHHost struct {
	Host         string
	HostName     string
	User         string
	Port         string
	IdentityFile string
}

// SSHConfig 管理 SSH 配置文件
type SSHConfig struct {
	configPath string
	hosts      []SSHHost
}

// NewSSHConfig 创建新的 SSH 配置管理器
func NewSSHConfig() (*SSHConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	configPath := filepath.Join(sshDir, "config")

	// 确保 .ssh 目录存在
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return nil, fmt.Errorf("无法创建 .ssh 目录: %w", err)
	}

	// 如果配置文件不存在，创建空文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			return nil, fmt.Errorf("无法创建配置文件: %w", err)
		}
		file.Close()
	}

	config := &SSHConfig{
		configPath: configPath,
	}

	if err := config.Load(); err != nil {
		return nil, err
	}

	return config, nil
}

// Load 加载 SSH 配置文件
func (c *SSHConfig) Load() error {
	file, err := os.Open(c.configPath)
	if err != nil {
		return fmt.Errorf("无法打开配置文件: %w", err)
	}
	defer file.Close()

	c.hosts = []SSHHost{}
	scanner := bufio.NewScanner(file)
	var currentHost *SSHHost

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], " ")

		switch key {
		case "host":
			if currentHost != nil {
				c.hosts = append(c.hosts, *currentHost)
			}
			currentHost = &SSHHost{Host: value}
		case "hostname":
			if currentHost != nil {
				currentHost.HostName = value
			}
		case "user":
			if currentHost != nil {
				currentHost.User = value
			}
		case "port":
			if currentHost != nil {
				currentHost.Port = value
			}
		case "identityfile":
			if currentHost != nil {
				currentHost.IdentityFile = value
			}
		}
	}

	if currentHost != nil {
		c.hosts = append(c.hosts, *currentHost)
	}

	return scanner.Err()
}

// Save 保存 SSH 配置文件
func (c *SSHConfig) Save() error {
	file, err := os.Create(c.configPath)
	if err != nil {
		return fmt.Errorf("无法创建配置文件: %w", err)
	}
	defer file.Close()

	for _, host := range c.hosts {
		fmt.Fprintf(file, "Host %s\n", host.Host)
		if host.HostName != "" {
			fmt.Fprintf(file, "    HostName %s\n", host.HostName)
		}
		if host.User != "" {
			fmt.Fprintf(file, "    User %s\n", host.User)
		}
		if host.Port != "" {
			fmt.Fprintf(file, "    Port %s\n", host.Port)
		}
		if host.IdentityFile != "" {
			fmt.Fprintf(file, "    IdentityFile %s\n", host.IdentityFile)
		}
		// 默认添加 IdentitiesOnly yes 设置
		fmt.Fprintf(file, "    IdentitiesOnly yes\n")
		fmt.Fprintln(file)
	}

	return nil
}

// GetHosts 获取所有主机配置
func (c *SSHConfig) GetHosts() []SSHHost {
	return c.hosts
}

// AddHost 添加新的主机配置
func (c *SSHConfig) AddHost(host SSHHost) {
	c.hosts = append(c.hosts, host)
}

// RemoveHost 删除指定的主机配置
func (c *SSHConfig) RemoveHost(index int) error {
	if index < 0 || index >= len(c.hosts) {
		return fmt.Errorf("索引超出范围")
	}
	c.hosts = append(c.hosts[:index], c.hosts[index+1:]...)
	return nil
}

// UpdateHost 更新指定的主机配置
func (c *SSHConfig) UpdateHost(index int, host SSHHost) error {
	if index < 0 || index >= len(c.hosts) {
		return fmt.Errorf("索引超出范围")
	}
	c.hosts[index] = host
	return nil
}

// ValidateIdentityFile 验证身份文件路径并给出警告
func ValidateIdentityFile(path string) (bool, string) {
	if strings.HasSuffix(strings.ToLower(path), ".ppk") {
		return false, "这是一个 PuTTY 格式的密钥，标准的 OpenSSH 可能无法使用。请使用 puttygen.exe 将其转换为 OpenSSH 格式后再使用。"
	}
	return true, ""
}