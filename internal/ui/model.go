package ui

import (
	"fmt"
	"time"

	"github.com/mordvn/dotfiles-sync/internal/config"
	"github.com/mordvn/dotfiles-sync/internal/sync"

	tea "github.com/charmbracelet/bubbletea"
)

type ItemStatus struct {
	Path   config.PathConfig
	Status string // "TODO", "SYNCING", "OK", "ERROR"
	Error  string
}

type Model struct {
	config       config.Config
	items        []ItemStatus
	selectedIdx  int
	syncing      bool
	gitStatus    string
	lastUpdate   time.Time
	spinnerFrame int
	copier       *sync.Copier
	gitMgr       *sync.GitManager
}

const (
	spinnerSpeed = time.Millisecond * 100
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func NewModel(cfg config.Config) *Model {
	items := make([]ItemStatus, len(cfg.Paths))
	for i, p := range cfg.Paths {
		items[i] = ItemStatus{
			Path:   p,
			Status: "TODO",
		}
	}

	return &Model{
		config:    cfg,
		items:     items,
		gitStatus: "checking...",
		copier:    sync.NewCopier(cfg.DotfilesDir),
		gitMgr:    sync.NewGitManager(cfg.DotfilesDir, cfg.GitRepo),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "s":
			if !m.syncing {
				m.syncing = true
				m.selectedIdx = 0

				for i := range m.items {
					m.items[i].Status = "TODO"
				}
				return m, m.syncNext()
			}

		case "p":
			if err := m.gitMgr.CommitAndPush(""); err != nil {
				m.gitStatus = fmt.Sprintf("❌ Push failed: %v", err)
			} else {
				m.gitStatus = "✓ Pushed successfully"
			}
			return m, nil

		case "r":
			status, _ := m.gitMgr.GetStatus()
			m.gitStatus = status
			return m, nil

		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "j":
			if m.selectedIdx < len(m.items)-1 {
				m.selectedIdx++
			}
		}

	case spinnerMsg:
		m.spinnerFrame++
		if m.syncing && m.selectedIdx < len(m.items) {
			return m, m.syncCurrent()
		}

	case syncDoneMsg:
		if msg.err != nil {
			m.items[msg.idx].Status = "ERROR"
			m.items[msg.idx].Error = msg.err.Error()
		} else {
			m.items[msg.idx].Status = "OK"
			m.items[msg.idx].Error = ""
		}
		m.lastUpdate = time.Now()

		m.selectedIdx++
		if m.selectedIdx >= len(m.items) {
			m.syncing = false
			status, _ := m.gitMgr.GetStatus()
			m.gitStatus = status
			return m, nil
		}

		return m, m.syncNext()

	case tickMsg:
		if m.syncing {
			m.spinnerFrame = (m.spinnerFrame + 1) % len(spinnerFrames)
			return m, m.tick()
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.config.DotfilesDir == "" {
		return "Error: config not loaded\n"
	}

	header := StyleHeader.Render("📦 DOTFILES")

	var itemsStr string
	for i, item := range m.items {
		icon := GetStatusIcon(item.Status)
		statusStyle := GetStatusStyle(item.Status)

		status := item.Status
		if status == "SYNCING" {
			status = string(spinnerFrames[m.spinnerFrame%len(spinnerFrames)])
		}

		dest := fmt.Sprintf("→ %s/%s", m.config.DotfilesDir, item.Path.Name)

		statusStr := fmt.Sprintf("[%s]", status)
		if item.Error != "" {
			statusStr = statusStr + " " + item.Error
		}

		lineContent := fmt.Sprintf(" %s %s  %s  %s",
			icon,
			item.Path.Name,
			dest,
			statusStr,
		)

		if i == m.selectedIdx {
			lineContent = StyleHighlight.Render(lineContent)
		} else {
			lineContent = statusStyle.Render(lineContent)
		}

		itemsStr += lineContent + "\n"
	}

	gitStatusLine := fmt.Sprintf("Git: %s\n", m.gitStatus)
	if m.lastUpdate.Unix() > 0 {
		gitStatusLine = fmt.Sprintf("Last update: %s | Git: %s\n",
			m.lastUpdate.Format("15:04:05"),
			m.gitStatus,
		)
	}

	footer := StyleFooter.Render(
		fmt.Sprintf(
			"[%s] Sync | [%s] Push | [%s] Refresh | [↑↓] Navigate | [q] Quit",
			StyleShortcut.Render("s"),
			StyleShortcut.Render("p"),
			StyleShortcut.Render("r"),
		),
	)

	content := header + "\n" + itemsStr + "\n" + gitStatusLine + "\n" + footer

	return StyleBorder.Render(content)
}

type spinnerMsg struct{}

func (m Model) tick() tea.Cmd {
	return tea.Tick(spinnerSpeed, func(time.Time) tea.Msg {
		return spinnerMsg{}
	})
}

type syncDoneMsg struct {
	idx int
	err error
}

func (m Model) syncNext() tea.Cmd {
	m.items[m.selectedIdx].Status = "SYNCING"
	return m.syncCurrent()
}

func (m Model) syncCurrent() tea.Cmd {
	idx := m.selectedIdx
	item := m.items[idx]

	return func() tea.Msg {
		time.Sleep(time.Millisecond * 500)

		err := m.copier.Copy(item.Path)

		return syncDoneMsg{idx: idx, err: err}
	}
}

type tickMsg struct{}
