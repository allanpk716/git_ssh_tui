package main

import (
	"fmt"
	"log"
	"os"

	"github.com/allanpk716/git_ssh_tui/internal/ui"
	"github.com/charmbracelet/bubbletea"
)

func main() {
	// 创建模型
	m, err := ui.NewModel()
	if err != nil {
		log.Fatalf("无法初始化应用: %v", err)
	}

	// 创建程序
	p := tea.NewProgram(m, tea.WithAltScreen())

	// 运行程序
	if _, err := p.Run(); err != nil {
		fmt.Printf("程序运行出错: %v", err)
		os.Exit(1)
	}
}