package dashboard

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"
	"github.com/eklairs/tlock/tlock-models/dashboard/tokens"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
)

// Bottom bar
func BottomBar(width int, currentUser string, styles tlockstyles.Styles) string {
    style := styles.Base.Copy().
        Width(width).
        Height(1).
        Margin(1)

    left := lipgloss.JoinHorizontal(
        lipgloss.Left,
        styles.AccentTitle.Render("TLock"),
        styles.DimmedTitle.Render("v1.0.0"),
    )

    right := lipgloss.JoinHorizontal(
        lipgloss.Right,
        styles.DimmedTitle.Render(time.Now().Format("2 January, 2006")),
        styles.AccentTitle.Render(currentUser),
    )

    bar := lipgloss.JoinHorizontal(
        lipgloss.Center,
        left,
        strings.Repeat(" ", width - lipgloss.Width(left) - lipgloss.Width(right) - 2),
        right,
    )

    return style.Render(bar)
}

// Dashboard key map
type dashboardKeyMap struct {
	Help        key.Binding
	Add         key.Binding
	ChangeTheme key.Binding
}

// ShortHelp()
func (k dashboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Add, k.ChangeTheme}
}

// FullHelp()
func (k dashboardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help},
		{k.Add},
        {k.ChangeTheme},
	}
}

// Keys
var dashboardKeys = dashboardKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help menu"),
	),
	Add: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add folder"),
	),
	ChangeTheme: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("Ctrl + T", "change theme"),
	),
}

// Dashboard Model
type DashboardModel struct {
	// Vault
	vault *tlockvault.TLockVault

	// Folders
	folders folders.Folders

	// Tokens
	tokens tokens.Tokens

	// Help
	help help.Model

	// Styles
	styles tlockstyles.Styles

	// Context
	context context.Context
}

func InitializeDashboardModel(vault tlockvault.TLockVault, context context.Context) DashboardModel {
	width, _, _ := term.GetSize(0)
	styles := tlockstyles.InitializeStyle(width, context.Theme)

	return DashboardModel{
		vault:   &vault,
		styles:  styles,
		context: context,
		help:    buildhelp.BuildHelp(styles),
		tokens:  tokens.InitializeTokens(&vault, context),
		folders: folders.InitializeFolders(&vault, context),
	}
}

// Init
func (m DashboardModel) Init() tea.Cmd {
	var cmd tea.Cmd

	if len(m.vault.Data.Folders) != 0 {
		cmd = func() tea.Msg {
			return folders.FolderChangedMsg{
				Folder: m.vault.Data.Folders[0].Name,
			}
		}
	}

	return cmd
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, dashboardKeys.Help):
			manager.PushScreen(InitializeHelpModel(m.context))
		}
	}

	return m, tea.Batch(m.folders.Update(msg, manager), m.tokens.Update(msg, manager))
}

// View
func (m DashboardModel) View() string {
    width, height, _ := term.GetSize(0)

	if len(m.vault.Data.Folders) == 0 {
		style := m.styles.Base.Copy().
			Height(height).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
		    m.styles.Title.Copy().UnsetWidth().Render(tokens.EmptyAsciiArt),
			m.styles.Title.Copy().UnsetBold().UnsetWidth().Render("So empty! How about adding a new folder?"), "",
			m.styles.Center.Copy().UnsetWidth().Render(m.help.View(dashboardKeys)),
		)

		return style.Render(ui)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
        lipgloss.JoinHorizontal(
            lipgloss.Left,
            m.folders.View(), "   ",
            m.tokens.View(),
        ),
        BottomBar(width, "Eklairs", m.styles),
	)
}
