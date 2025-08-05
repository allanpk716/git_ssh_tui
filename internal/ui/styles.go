package ui

import "github.com/charmbracelet/lipgloss"

var (
	// 基础颜色
	primaryColor   = lipgloss.Color("205")
	secondaryColor = lipgloss.Color("99")
	successColor   = lipgloss.Color("46")
	errorColor     = lipgloss.Color("196")
	warningColor   = lipgloss.Color("214")
	mutedColor     = lipgloss.Color("241")

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// 分页样式
	paginationStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// 帮助样式
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// 聚焦样式
	focusedStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// 无样式
	noStyle = lipgloss.NewStyle()

	// 错误样式
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// 警告样式
	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	// 成功样式
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// 表单容器样式
formStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Margin(1, 0)

	// 输入框标签样式
	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Width(15).
			Align(lipgloss.Right)

	// 按钮样式
	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 2).
			Margin(0, 1)

	// 取消按钮样式
	cancelButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(mutedColor).
			Padding(0, 2).
			Margin(0, 1)

	// 确认对话框样式
	confirmDialogStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			Margin(2, 4).
			Align(lipgloss.Center)

	// 状态栏样式
statusBarStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1)
)

// GetFormStyle 根据终端宽度返回动态调整的表单样式
func GetFormStyle(width int) lipgloss.Style {
	formWidth := width - 8 // 减去边距和边框
	if formWidth < 40 {
		formWidth = 40
	}
	return formStyle.Copy().Width(formWidth)
}