# 使用指南

## 快速开始

1. **构建项目**
   ```bash
   go build -o ssh-config-manager.exe
   ```

2. **运行程序**
   ```bash
   ./ssh-config-manager.exe
   ```

## 使用场景示例

### 场景 1: 添加 GitHub SSH 配置

1. 启动程序后，按 `a` 键进入添加模式
2. 填写以下信息：
   - **Host**: `github`
   - **HostName**: `github.com`
   - **User**: `git`
   - **IdentityFile**: `~/.ssh/id_rsa`
3. 按 `Tab` 键在字段间切换，最后按 `Enter` 提交

### 场景 2: 添加公司 GitLab 配置

1. 按 `a` 键进入添加模式
2. 填写以下信息：
   - **Host**: `gitlab-work`
   - **HostName**: `gitlab.company.com`
   - **User**: `git`
   - **IdentityFile**: `~/.ssh/work_key`
3. 提交后即可使用 `git clone git@gitlab-work:project/repo.git`

### 场景 3: 编辑现有配置

1. 在主界面选择要编辑的配置项
2. 按 `e` 键进入编辑模式
3. 表单会自动填充当前配置的所有信息
4. 修改需要更改的字段（如更换密钥文件路径）
5. 按 `Enter` 保存修改，或按 `Esc` 取消

### 场景 4: 处理 PuTTY 密钥警告

如果你在 Windows 上使用 TortoiseGit 的 .ppk 文件：

1. 在 **IdentityFile** 字段输入 `.ppk` 文件路径时
2. 程序会显示警告："这是一个 PuTTY 格式的密钥..."
3. 使用 `puttygen.exe` 转换密钥：
   - 打开 PuTTY Key Generator
   - 加载你的 .ppk 文件
   - 选择 "Conversions" → "Export OpenSSH key"
   - 保存为新的私钥文件（如 `id_rsa`）
   - 在程序中使用新的私钥路径

## 键盘快捷键总结

### 主界面
- `a` / `n`: 添加新配置
- `e`: 编辑选中配置
- `d` / `x`: 删除选中配置
- `↑` / `↓`: 上下导航
- `q`: 退出程序

### 添加/编辑界面
- `Tab`: 下一个字段
- `Shift+Tab`: 上一个字段
- `Enter`: 提交表单
- `Esc`: 取消并返回

### 删除确认界面
- `Y`: 确认删除
- `N` / `Esc`: 取消删除

## 生成的配置文件示例

程序会在 `~/.ssh/config` 中生成如下格式的配置：

```
Host github
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_rsa
    IdentitiesOnly yes

Host gitlab-work
    HostName gitlab.company.com
    User git
    IdentityFile ~/.ssh/work_key
    IdentitiesOnly yes
```

## 安全特性

### IdentitiesOnly yes
程序会自动为每个SSH配置添加 `IdentitiesOnly yes` 设置，这个设置的作用是：
- 只使用明确指定的身份文件（IdentityFile）进行认证
- 防止SSH客户端尝试使用SSH代理中的其他密钥
- 提高安全性，避免意外使用错误的密钥
- 确保连接使用预期的身份验证方式

## 故障排除

### 权限问题
如果遇到权限错误，确保：
- 有权限访问 `~/.ssh` 目录
- 私钥文件权限正确（通常是 600）

### 路径问题
- Windows 用户可以使用 `C:\Users\username\.ssh\id_rsa` 格式
- 或使用 `~/.ssh/id_rsa` 格式（程序会自动解析）

### 测试配置
添加配置后，可以使用以下命令测试：
```bash
ssh -T git@github  # 测试 GitHub 连接
ssh -T git@gitlab-work  # 测试 GitLab 连接
```