package tokens

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/eklairs/tlock/tlock-internal/boundedinteger"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"

	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Edit folder key map
type tokenKeyMap struct {
	Manual key.Binding
	Screen key.Binding
}

// ShortHelp()
func (k tokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Manual, k.Screen}
}

// FullHelp()
func (k tokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Manual},
		{k.Screen},
	}
}

// Keys
var tokenKeys = tokenKeyMap{
	Manual: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add token"),
	),
	Screen: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "add from screen"),
	),
}

var EmptyAsciiArt = `
\    /\
 )  ( ')
(  /  )
 \(__)|
`

// Tokens
type Tokens struct {
	// Context
	context context.Context

	// Vault
	vault tlockvault.TLockVault

	// Focused index
	focused_index boundedinteger.BoundedInteger

	// Styles
	styles tlockstyles.Styles

	// Folder
	folder string

	// Help
	help help.Model
}

// Initializes a new instance of folders
func InitializeTokens(vault tlockvault.TLockVault, context context.Context, folder string) Tokens {
	// Terminal size
	width, _, _ := term.GetSize(0)

	// Styles
	styles := tlockstyles.InitializeStyle(width-folders.FOLDERS_WIDTH, context.Theme)

	return Tokens{
		vault:   vault,
		styles:  styles,
		context: context,
		folder:  folder,
		help:    buildhelp.BuildHelp(styles),
	}
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "j":
			tokens.focused_index.Increase()
		case "k":
			tokens.focused_index.Decrease()
        case "s":
            manager.PushScreen(InitializeTokenFromScreen(tokens.context))
		}
	}

	return nil
}

// View
func (tokens Tokens) View() string {
	// Get term size
	_, height, _ := term.GetSize(0)

	// Get URIs
	uris := tokens.vault.GetTokens(tokens.folder)

	if len(uris) == 0 {
		style := tokens.styles.Base.Copy().
			Height(height).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Left,
			tokens.styles.Center.Render(tokens.styles.Title.Copy().UnsetWidth().Render(EmptyAsciiArt)),
			tokens.styles.Center.Render(tokens.styles.Base.Copy().UnsetWidth().Render("So empty! How about adding a new token?")),
			tokens.styles.Center.Copy().UnsetWidth().Render(tokens.help.View(tokenKeys)),
		)

		return style.Render(ui)
	}

	// List of items
	items := make([]string, 0)

	// Render
	return lipgloss.JoinVertical(lipgloss.Center, items...)
}