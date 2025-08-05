package ui

import (
	"fmt"

	"github.com/allanpk716/git_ssh_tui/internal/config"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

// ViewState 表示当前视图状态
type ViewState int

const (
	ListView ViewState = iota
	AddView
	EditView
	DeleteConfirmView
)

// Model 是应用的主要模型
type Model struct {
	sshConfig   *config.SSHConfig
	state       ViewState
	list        list.Model
	form        FormModel
	selected    int
	err         error
	warning     string
	deleteIndex int
	editIndex   int
	isEditing   bool
	width       int
	height      int
}

// FormModel 表示添加/编辑表单的模型
type FormModel struct {
	hostInput         textinput.Model
	hostnameInput     textinput.Model
	userInput         textinput.Model
	identityFileInput textinput.Model
	focusIndex        int
	inputs            []textinput.Model
}

// HostItem 实现 list.Item 接口
type HostItem struct {
	host config.SSHHost
}

func (h HostItem) FilterValue() string {
	return h.host.Host
}

func (h HostItem) Title() string {
	return h.host.Host
}

func (h HostItem) Description() string {
	return fmt.Sprintf("%s@%s", h.host.User, h.host.HostName)
}

// NewModel 创建新的模型
func NewModel() (*Model, error) {
	sshConfig, err := config.NewSSHConfig()
	if err != nil {
		return nil, err
	}

	// 创建列表项
	items := make([]list.Item, len(sshConfig.GetHosts()))
	for i, host := range sshConfig.GetHosts() {
		items[i] = HostItem{host: host}
	}

	// 创建列表模型
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH 配置管理"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// 创建表单模型
	form := NewFormModel()

	return &Model{
		sshConfig: sshConfig,
		state:     ListView,
		list:      l,
		form:      form,
	}, nil
}

// NewFormModel 创建新的表单模型
func NewFormModel() FormModel {
	// Host 输入框
	hostInput := textinput.New()
	hostInput.Placeholder = "例如: gitlab-work"
	hostInput.Focus()
	hostInput.CharLimit = 50
	hostInput.Width = 30

	// HostName 输入框
	hostnameInput := textinput.New()
	hostnameInput.Placeholder = "例如: gitlab.example.com"
	hostnameInput.CharLimit = 100
	hostnameInput.Width = 30

	// User 输入框
	userInput := textinput.New()
	userInput.Placeholder = "例如: git"
	userInput.CharLimit = 50
	userInput.Width = 30

	// IdentityFile 输入框
	identityFileInput := textinput.New()
	identityFileInput.Placeholder = "例如: ~/.ssh/id_rsa"
	identityFileInput.CharLimit = 200
	identityFileInput.Width = 50

	// 创建inputs数组，直接引用上面创建的输入框
	inputs := []textinput.Model{hostInput, hostnameInput, userInput, identityFileInput}

	return FormModel{
		hostInput:         hostInput,
		hostnameInput:     hostnameInput,
		userInput:         userInput,
		identityFileInput: identityFileInput,
		focusIndex:        0,
		inputs:            inputs,
	}
}

// Init 初始化模型
func (m Model) Init() tea.Cmd {
	return nil
}

// Update 处理消息更新
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 3)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case ListView:
			return m.updateListView(msg)
		case AddView:
			return m.updateAddView(msg)
		case EditView:
			return m.updateEditView(msg)
		case DeleteConfirmView:
			return m.updateDeleteConfirmView(msg)
		}
	}

	return m, nil
}

// updateListView 更新列表视图
func (m Model) updateListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "a", "n":
		m.state = AddView
		m.form = NewFormModel()
		m.isEditing = false
		m.warning = ""
		return m, nil
	case "e":
		if len(m.list.Items()) > 0 {
			m.editIndex = m.list.Index()
			m.state = EditView
			m.isEditing = true
			m.form = m.createFormWithData(m.editIndex)
			m.warning = ""
		}
		return m, nil
	case "d", "x":
		if len(m.list.Items()) > 0 {
			m.deleteIndex = m.list.Index()
			m.state = DeleteConfirmView
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// createFormWithData 创建预填充数据的表单
func (m Model) createFormWithData(index int) FormModel {
	hosts := m.sshConfig.GetHosts()
	if index < 0 || index >= len(hosts) {
		return NewFormModel()
	}

	host := hosts[index]
	form := NewFormModel()

	// 预填充数据
	form.hostInput.SetValue(host.Host)
	form.hostnameInput.SetValue(host.HostName)
	form.userInput.SetValue(host.User)
	form.identityFileInput.SetValue(host.IdentityFile)

	// 同时更新inputs数组
	form.inputs[0].SetValue(host.Host)
	form.inputs[1].SetValue(host.HostName)
	form.inputs[2].SetValue(host.User)
	form.inputs[3].SetValue(host.IdentityFile)

	return form
}

// updateAddView 更新添加视图
func (m Model) updateAddView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = ListView
		m.warning = ""
		return m, nil
	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		if s == "enter" && m.form.focusIndex == len(m.form.inputs) {
			// 提交表单
			return m.submitForm()
		}

		if s == "up" || s == "shift+tab" {
			m.form.focusIndex--
		} else {
			m.form.focusIndex++
		}

		if m.form.focusIndex > len(m.form.inputs) {
			m.form.focusIndex = 0
		} else if m.form.focusIndex < 0 {
			m.form.focusIndex = len(m.form.inputs)
		}

		cmds := make([]tea.Cmd, len(m.form.inputs))
		for i := 0; i <= len(m.form.inputs)-1; i++ {
			if i == m.form.focusIndex {
				cmds[i] = m.form.inputs[i].Focus()
				m.form.inputs[i].PromptStyle = focusedStyle
				m.form.inputs[i].TextStyle = focusedStyle
				continue
			}
			m.form.inputs[i].Blur()
			m.form.inputs[i].PromptStyle = noStyle
			m.form.inputs[i].TextStyle = noStyle
		}

		// 检查 IdentityFile 警告
		if m.form.focusIndex == 3 { // IdentityFile 输入框
			identityFile := m.form.inputs[3].Value()
			if identityFile != "" {
				if valid, warning := config.ValidateIdentityFile(identityFile); !valid {
					m.warning = warning
				} else {
					m.warning = ""
				}
			}
		}

		return m, tea.Batch(cmds...)
	}

	// 处理输入框更新
	cmd := m.updateInputs(msg)
	return m, cmd
}

// updateEditView 更新编辑视图
func (m Model) updateEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = ListView
		m.warning = ""
		m.isEditing = false
		return m, nil
	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		if s == "enter" && m.form.focusIndex == len(m.form.inputs) {
			// 提交表单
			return m.submitForm()
		}

		if s == "up" || s == "shift+tab" {
			m.form.focusIndex--
		} else {
			m.form.focusIndex++
		}

		if m.form.focusIndex > len(m.form.inputs) {
			m.form.focusIndex = 0
		} else if m.form.focusIndex < 0 {
			m.form.focusIndex = len(m.form.inputs)
		}

		cmds := make([]tea.Cmd, len(m.form.inputs))
		for i := 0; i <= len(m.form.inputs)-1; i++ {
			if i == m.form.focusIndex {
				cmds[i] = m.form.inputs[i].Focus()
				m.form.inputs[i].PromptStyle = focusedStyle
				m.form.inputs[i].TextStyle = focusedStyle
				continue
			}
			m.form.inputs[i].Blur()
			m.form.inputs[i].PromptStyle = noStyle
			m.form.inputs[i].TextStyle = noStyle
		}

		// 检查 IdentityFile 警告
		if m.form.focusIndex == 3 { // IdentityFile 输入框
			identityFile := m.form.inputs[3].Value()
			if identityFile != "" {
				if valid, warning := config.ValidateIdentityFile(identityFile); !valid {
					m.warning = warning
				} else {
					m.warning = ""
				}
			}
		}

		return m, tea.Batch(cmds...)
	}

	// 处理输入框更新
	cmd := m.updateInputs(msg)
	return m, cmd
}

// updateDeleteConfirmView 更新删除确认视图
func (m Model) updateDeleteConfirmView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "ctrl+c":
		return m, tea.Quit
	case "y", "Y":
		// 确认删除
		if err := m.sshConfig.RemoveHost(m.deleteIndex); err != nil {
			m.err = err
		} else {
			if err := m.sshConfig.Save(); err != nil {
				m.err = err
			} else {
				// 更新列表
				m.refreshList()
			}
		}
		m.state = ListView
		return m, nil
	case "n", "N", "esc":
		// 取消删除
		m.state = ListView
		return m, nil
	}
	return m, nil
}

// updateInputs 更新输入框
func (m Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.form.inputs))

	for i := range m.form.inputs {
		m.form.inputs[i], cmds[i] = m.form.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// submitForm 提交表单
func (m Model) submitForm() (tea.Model, tea.Cmd) {
	host := config.SSHHost{
		Host:         m.form.inputs[0].Value(),
		HostName:     m.form.inputs[1].Value(),
		User:         m.form.inputs[2].Value(),
		IdentityFile: m.form.inputs[3].Value(),
	}

	// 验证必填字段
	if host.Host == "" {
		m.err = fmt.Errorf("Host 字段不能为空")
		return m, nil
	}

	if m.isEditing {
		// 编辑模式：更新现有配置
		if err := m.sshConfig.UpdateHost(m.editIndex, host); err != nil {
			m.err = err
			return m, nil
		}
	} else {
		// 添加模式：添加新配置
		m.sshConfig.AddHost(host)
	}

	if err := m.sshConfig.Save(); err != nil {
		m.err = err
		return m, nil
	}

	// 刷新列表并返回列表视图
	m.refreshList()
	m.state = ListView
	m.warning = ""
	m.isEditing = false
	return m, nil
}

// refreshList 刷新列表
func (m *Model) refreshList() {
	items := make([]list.Item, len(m.sshConfig.GetHosts()))
	for i, host := range m.sshConfig.GetHosts() {
		items[i] = HostItem{host: host}
	}
	m.list.SetItems(items)
}