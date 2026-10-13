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
	"github.com/muesli/termenv"
)

type viewStyles struct {
	banner  lipgloss.Style
	heading lipgloss.Style
	panel   lipgloss.Style
	bullet  lipgloss.Style
	key     lipgloss.Style
	quiet   lipgloss.Style
}

func newViewStyles(renderer *lipgloss.Renderer) viewStyles {
	return viewStyles{
		banner: renderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Padding(0, 2),
		heading: renderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")),
		panel: renderer.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69")).
			Padding(1, 2),
		bullet: renderer.NewStyle().
			Foreground(lipgloss.Color("212")),
		key: renderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("205")).
			Padding(0, 1),
		quiet: renderer.NewStyle().
			Foreground(lipgloss.Color("245")),
	}
}

type model struct {
	visitor string
	width   int
	height  int
	styles  viewStyles
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
	connection := fmt.Sprintf("Hello %s. You connected to an application, not a shell.", m.visitor)
	capabilities := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.heading.Render("WHAT WISH ADDS"),
		"",
		m.styles.bullet.Render("◆")+" An SSH server embedded in this Go program",
		m.styles.bullet.Render("◆")+" A separate Bubble Tea app for each connection",
		m.styles.bullet.Render("◆")+" Terminal-size events sent to the app",
	)
	meaning := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.heading.Render("WHY THAT MATTERS"),
		"",
		"Plain SSH gives you a remote command prompt. Wish lets the server",
		"own the whole terminal experience. It does not replace SSH routing,",
		"accounts, access control, or a collaboration relay.",
	)

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n\n%s\n\n%s quits this session.\n",
		m.styles.banner.Render("  WISH  ·  SSH-HOSTED TERMINAL APP  "),
		connection,
		m.styles.panel.Render(capabilities),
		m.styles.panel.Render(meaning),
		m.styles.key.Render("q"),
	) + m.styles.quiet.Render("  Local demonstration only · no Vesta or internet connection") + "\n"
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
			bubbletea.MiddlewareWithColorProfile(
				func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
					return model{
						visitor: sess.User(),
						styles:  newViewStyles(bubbletea.MakeRenderer(sess)),
					}, []tea.ProgramOption{tea.WithAltScreen()}
				},
				termenv.TrueColor,
			),
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
