package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
    "github.com/pquerna/otp"
)

type dashboardStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type DashboardModel struct {
    styles dashboardStyles
    vault tlockvault.TLockVault
    current_index int
    token_current_index int
}

// Initialize root model
func InitializeDashboardModel(vault tlockvault.TLockVault) DashboardModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

    return DashboardModel {
        styles: dashboardStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(30).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(30).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            dimmedCenter: dimmed.Width(30).Copy().Align(lipgloss.Center),
        },
        vault: vault,
        current_index: 0,
    }
}

// Init
func (m DashboardModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "J":
            m.current_index = (m.current_index + 1) % len(m.vault.Data.Folders)
            m.token_current_index = 0
        case "K":
            if m.current_index == 0 {
                m.current_index = len(m.vault.Data.Folders) - 1
            } else {
                m.current_index -= 1
            }
            m.token_current_index = 0
        case "j":
            m.token_current_index = (m.token_current_index + 1) % len(m.vault.Data.Folders[m.current_index].Uris)
        case "k":
            if m.token_current_index == 0 {
                m.token_current_index = len(m.vault.Data.Folders[m.current_index].Uris) - 1
            } else {
                m.token_current_index -= 1
            }
        }
    }

	return m, nil
}

// View
func (m DashboardModel) View() string {
    width, height, _ := term.GetSize(0)

    style := lipgloss.NewStyle().Height(height).Width(30).Padding(1, 3)
    folder_style := lipgloss.NewStyle().Width(30).Padding(1, 3)

    // Folders
    folders := make([]string, 0)

    for index, folder := range m.vault.Data.Folders {
        render_fn := folder_style.Render

        ui := lipgloss.JoinVertical(
            lipgloss.Left,
            m.styles.title.Render(folder.Name),
            m.styles.dimmed.Render(fmt.Sprintf("%d tokens", len(folder.Uris))),
        )

        if index == m.current_index {
            render_fn = folder_style.Copy().Background(lipgloss.Color("#1E1E2E")).
                Width(23).
                Padding(1, 2).
                BorderBackground(lipgloss.Color("#1E1E2E")).
                Border(lipgloss.ThickBorder(), false, false, false, true).Render
        }

        folders = append(folders, render_fn(ui))
    }

    // Tokens
    tokens := make([]string, 0)

    for index, uri := range m.vault.Data.Folders[m.current_index].Uris {
        style := lipgloss.NewStyle().
            Width(width - 30 - 2).
            Padding(1, 3).
            MarginBottom(1)
        title := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
        issuer := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

        if index == m.token_current_index {
            style = style.Background(lipgloss.Color("#1E1E2E"))
            title = title.Background(lipgloss.Color("#1E1E2E")).Bold(true)
            issuer = issuer.Background(lipgloss.Color("#1E1E2E")).Bold(true)
        }

        totp, _ := otp.NewKeyFromURL(uri)

        tokens = append(tokens, style.Render(fmt.Sprintf("%s • %s", title.Render(totp.AccountName()), issuer.Render(totp.Issuer()))))
    }

    ui := []string {
        style.Render(lipgloss.JoinVertical(lipgloss.Left, folders...)),
    }

    ui = append(ui, lipgloss.JoinVertical(lipgloss.Left, tokens...))

    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        ui...
    )
}

