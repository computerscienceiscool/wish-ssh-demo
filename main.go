package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	keyStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
)

type model struct {
	visitor string
	width   int
	height  int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"\n%s\n\nHello %s. You connected to an application, not a shell.\n\nWish provides:\n  • an SSH server embedded in the Go program\n  • a separate Bubble Tea application for each connection\n  • terminal-size events passed to the application\n\nPlain SSH gives you a remote command prompt. Wish gives the server control\nover the full terminal experience. It does not replace SSH routing, accounts,\nor access control.\n\n%s quits this session.\n",
		titleStyle.Render("Wish SSH-hosted TUI demo"),
		m.visitor,
		keyStyle.Render("q"),
	)
}

func main() {
	address := flag.String("address", "127.0.0.1:2222", "SSH listen address")
	hostKey := flag.String("host-key", ".wish-data/host-key", "host key path; created if absent")
	flag.Parse()
	if err := os.MkdirAll(filepath.Dir(*hostKey), 0o700); err != nil {
		log.Fatalf("create host-key directory: %v", err)
	}

	server, err := wish.NewServer(
		wish.WithAddress(*address),
		wish.WithHostKeyPath(*hostKey),
		wish.WithMiddleware(
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return model{visitor: sess.User()}, []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		log.Fatalf("create SSH server: %v", err)
	}

	log.Printf("Wish demo listening on ssh://%s", *address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("serve Wish demo: %v", err)
	}
}
