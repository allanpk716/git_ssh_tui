# SSH 配置管理工具 (git_ssh_tui)

一个基于 Bubbletea 构建的跨平台 SSH 配置管理 TUI 工具，帮助用户更方便地管理 `~/.ssh/config` 文件。

## 功能特性

- 📋 **列出现有配置**: 以友好的列表形式展示所有 SSH 配置
- ➕ **添加新配置**: 通过表单界面轻松添加新的 SSH 主机配置
- ✏️ **编辑配置**: 修改现有的 SSH 配置，支持所有字段的编辑
- 🗑️ **删除配置**: 安全删除不需要的 SSH 配置（带确认提示）
- ⚠️ **智能警告**: 自动检测 PuTTY 格式密钥并给出转换提示
- 🔒 **安全默认**: 自动为所有配置添加 `IdentitiesOnly yes` 设置
- 🌍 **跨平台支持**: 支持 Windows、macOS 和 Linux
- 🎨 **美观界面**: 使用 Lipgloss 打造的现代化 TUI 界面

## 安装和运行

### 前置要求

- Go 1.23.4 或更高版本

### 构建和运行

```bash
# 克隆项目
git clone https://github.com/allanpk716/git_ssh_tui.git
cd git_ssh_tui

# 下载依赖
go mod tidy

# 构建
go build -o ssh-config-manager

# 运行
./ssh-config-manager
```

### Windows 用户

```powershell
# 构建
go build -o ssh-config-manager.exe

# 运行
.\ssh-config-manager.exe
```

## 使用说明

### 主界面操作

- `a` 或 `n`: 添加新的 SSH 配置
- `e`: 编辑选中的配置
- `d` 或 `x`: 删除选中的配置
- `↑`/`↓`: 在列表中导航
- `q`: 退出程序

### 添加/编辑配置界面

- `Tab`: 切换到下一个输入字段
- `Shift+Tab`: 切换到上一个输入字段
- `Enter`: 提交表单（在最后一个字段时）
- `Esc`: 取消并返回主界面

### 删除确认界面

- `Y`: 确认删除
- `N` 或 `Esc`: 取消删除

## 项目结构

```
git_ssh_tui/
├── main.go                    # 主程序入口
├── go.mod                     # Go 模块文件
├── README.md                  # 项目说明
└── internal/
    ├── config/
    │   └── ssh_config.go      # SSH 配置文件处理
    └── ui/
        ├── model.go           # Bubbletea 模型和状态管理
        ├── view.go            # 界面渲染逻辑
        └── styles.go          # UI 样式定义
```

## 技术栈

- **语言**: Go 1.23.4
- **TUI 框架**: [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **样式库**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **组件库**: [Bubbles](https://github.com/charmbracelet/bubbles)

## 特殊功能

### PuTTY 密钥检测

当用户在 `IdentityFile` 字段中输入 `.ppk` 格式的文件时，程序会自动显示警告信息：

> ⚠️ 这是一个 PuTTY 格式的密钥，标准的 OpenSSH 可能无法使用。请使用 puttygen.exe 将其转换为 OpenSSH 格式后再使用。

### 自动目录创建

如果 `~/.ssh` 目录或 `~/.ssh/config` 文件不存在，程序会自动创建它们，确保用户可以立即开始使用。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证。