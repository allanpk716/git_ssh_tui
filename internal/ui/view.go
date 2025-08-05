package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

// View 渲染主视图
func (m Model) View() string {
	switch m.state {
	case ListView:
		return m.listView()
	case AddView:
		return m.addView()
	case EditView:
		return m.editView()
	case DeleteConfirmView:
		return m.deleteConfirmView()
	default:
		return "未知状态"
	}
}

// listView 渲染列表视图
func (m Model) listView() string {
	var content strings.Builder

	// 主列表
	content.WriteString(m.list.View())
	content.WriteString("\n")

	// 错误信息
	if m.err != nil {
		content.WriteString(errorStyle.Render(fmt.Sprintf("错误: %s", m.err.Error())))
		content.WriteString("\n")
	}

	// 帮助信息
	helpText := []string{
		"a/n: 添加新配置",
		"e: 编辑配置",
		"d/x: 删除配置",
		"q: 退出",
	}
	content.WriteString(helpStyle.Render(strings.Join(helpText, " • ")))

	return content.String()
}

// addView 渲染添加视图
func (m Model) addView() string {
	var content strings.Builder

	// 标题
	content.WriteString(titleStyle.Render("添加新的 SSH 配置"))
	content.WriteString("\n\n")

	// 表单
formContent := m.renderForm()
content.WriteString(GetFormStyle(m.width).Render(formContent))

	// 警告信息
	if m.warning != "" {
		content.WriteString("\n")
		content.WriteString(warningStyle.Render(fmt.Sprintf("⚠️  %s", m.warning)))
	}

	// 错误信息
	if m.err != nil {
		content.WriteString("\n")
		content.WriteString(errorStyle.Render(fmt.Sprintf("错误: %s", m.err.Error())))
	}

	// 帮助信息
	content.WriteString("\n\n")
	helpText := []string{
		"Tab: 下一个字段",
		"Shift+Tab: 上一个字段",
		"Enter: 提交",
		"Esc: 取消",
	}
	content.WriteString(helpStyle.Render(strings.Join(helpText, " • ")))

	return content.String()
}

// renderForm 渲染表单
func (m Model) renderForm() string {
	var form strings.Builder

	// Host 字段
	form.WriteString(m.renderFormField("Host:", m.form.inputs[0], 0))
	form.WriteString("\n")

	// HostName 字段
	form.WriteString(m.renderFormField("HostName:", m.form.inputs[1], 1))
	form.WriteString("\n")

	// User 字段
	form.WriteString(m.renderFormField("User:", m.form.inputs[2], 2))
	form.WriteString("\n")

	// Port 字段
	form.WriteString(m.renderFormField("Port:", m.form.inputs[3], 3))
	form.WriteString("\n")

	// IdentityFile 字段
	form.WriteString(m.renderFormField("IdentityFile:", m.form.inputs[4], 4))
	form.WriteString("\n\n")

	// 提交按钮
	submitButton := "[ 提交 ]"
	if m.form.focusIndex == len(m.form.inputs) {
		submitButton = buttonStyle.Render("[ 提交 ]")
	} else {
		submitButton = cancelButtonStyle.Render("[ 提交 ]")
	}
	form.WriteString(submitButton)

	return form.String()
}

// renderFormField 渲染表单字段
func (m Model) renderFormField(label string, input interface{}, index int) string {
	var field strings.Builder

	// 标签
	if m.form.focusIndex == index {
		field.WriteString(focusedStyle.Render(label))
	} else {
		field.WriteString(labelStyle.Render(label))
	}
	field.WriteString(" ")

	// 输入框
	if textInput, ok := input.(textinput.Model); ok {
		field.WriteString(textInput.View())
	}

	return field.String()
}

// editView 渲染编辑视图
func (m Model) editView() string {
	var content strings.Builder

	// 标题
	content.WriteString(titleStyle.Render("编辑 SSH 配置"))
	content.WriteString("\n\n")

	// 表单
formContent := m.renderForm()
content.WriteString(GetFormStyle(m.width).Render(formContent))

	// 警告信息
	if m.warning != "" {
		content.WriteString("\n")
		content.WriteString(warningStyle.Render(fmt.Sprintf("⚠️  %s", m.warning)))
	}

	// 错误信息
	if m.err != nil {
		content.WriteString("\n")
		content.WriteString(errorStyle.Render(fmt.Sprintf("错误: %s", m.err.Error())))
	}

	// 帮助信息
	content.WriteString("\n\n")
	helpText := []string{
		"Tab: 下一个字段",
		"Shift+Tab: 上一个字段",
		"Enter: 保存修改",
		"Esc: 取消",
	}
	content.WriteString(helpStyle.Render(strings.Join(helpText, " • ")))

	return content.String()
}

// deleteConfirmView 渲染删除确认视图
func (m Model) deleteConfirmView() string {
	var content strings.Builder

	// 获取要删除的主机信息
	hosts := m.sshConfig.GetHosts()
	if m.deleteIndex >= 0 && m.deleteIndex < len(hosts) {
		host := hosts[m.deleteIndex]

		// 确认对话框内容
		portInfo := ""
		if host.Port != "" {
			portInfo = fmt.Sprintf("Port: %s\n", host.Port)
		}
		dialogContent := fmt.Sprintf(
			"确定要删除以下 SSH 配置吗？\n\n"+
				"Host: %s\n"+
				"HostName: %s\n"+
				"User: %s\n"+
				"%s"+
				"IdentityFile: %s\n\n"+
				"此操作无法撤销！\n\n"+
				"[Y] 确认删除    [N] 取消",
			host.Host,
			host.HostName,
			host.User,
			portInfo,
			host.IdentityFile,
		)

		content.WriteString(confirmDialogStyle.Render(dialogContent))
	} else {
		content.WriteString(errorStyle.Render("错误: 无效的选择"))
	}

	return content.String()
}