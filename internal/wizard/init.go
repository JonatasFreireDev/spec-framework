package wizard

import (
	"fmt"
	"io"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/JonatasFreireDev/spec-framework/internal/install"
)

type Stage uint8

const (
	SelectAgents Stage = iota
	EnterTarget
	Review
	Finished
)

type InitModel struct {
	Stage     Stage
	Cursor    int
	Selected  map[int]bool
	Target    string
	Confirmed bool
	Cancelled bool
}

var choices = []struct {
	Name  string
	Agent install.Agent
}{{"Codex", install.Codex}, {"Cursor", install.Cursor}, {"Claude Code", install.Claude}}

func NewInitModel() InitModel     { return InitModel{Stage: SelectAgents, Selected: map[int]bool{0: true}} }
func (m InitModel) Init() tea.Cmd { return nil }
func (m InitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "ctrl+c", "esc":
		m.Cancelled = true
		m.Stage = Finished
		return m, tea.Quit
	}
	if key.String() == "q" && m.Stage != EnterTarget {
		m.Cancelled = true
		m.Stage = Finished
		return m, tea.Quit
	}
	switch m.Stage {
	case SelectAgents:
		switch key.Code {
		case tea.KeyUp:
			if m.Cursor > 0 {
				m.Cursor--
			}
		case tea.KeyDown:
			if m.Cursor < len(choices)-1 {
				m.Cursor++
			}
		case tea.KeySpace:
			m.Selected[m.Cursor] = !m.Selected[m.Cursor]
		case tea.KeyEnter:
			if len(m.Agents()) > 0 {
				m.Stage = EnterTarget
			}
		}
	case EnterTarget:
		switch key.Code {
		case tea.KeyEnter:
			if strings.TrimSpace(m.Target) != "" {
				m.Target = strings.TrimSpace(m.Target)
				m.Stage = Review
			}
		case tea.KeyBackspace:
			if len(m.Target) > 0 {
				r := []rune(m.Target)
				m.Target = string(r[:len(r)-1])
			}
		default:
			if key.Text != "" {
				m.Target += key.Text
			}
		}
	case Review:
		switch strings.ToLower(key.String()) {
		case "y", "enter":
			m.Confirmed = true
			m.Stage = Finished
			return m, tea.Quit
		case "n":
			m.Cancelled = true
			m.Stage = Finished
			return m, tea.Quit
		case "b":
			m.Stage = SelectAgents
		}
	}
	return m, nil
}
func (m InitModel) View() tea.View {
	var b strings.Builder
	switch m.Stage {
	case SelectAgents:
		b.WriteString("Which agent skill formats should be installed?\n\n")
		for i, c := range choices {
			cursor := " "
			if i == m.Cursor {
				cursor = ">"
			}
			check := " "
			if m.Selected[i] {
				check = "x"
			}
			fmt.Fprintf(&b, "%s [%s] %s\n", cursor, check, c.Name)
		}
		b.WriteString("\nSpace selects, Enter continues, Esc cancels.\n")
	case EnterTarget:
		b.WriteString("Target directory:\n> " + m.Target + "\n\nEnter continues, Esc cancels.\n")
	case Review:
		fmt.Fprintf(&b, "Installation plan\n\nTarget: %s\nAgents: %s\n\nApply this plan? [Y/n/b]\n", m.Target, strings.Join(m.AgentNames(), ", "))
	case Finished:
		if m.Cancelled {
			b.WriteString("Installation cancelled; no files were written.\n")
		}
	}
	return tea.NewView(b.String())
}
func (m InitModel) Agents() []install.Agent {
	var out []install.Agent
	for i, c := range choices {
		if m.Selected[i] {
			out = append(out, c.Agent)
		}
	}
	return out
}
func (m InitModel) AgentNames() []string {
	var out []string
	for i, c := range choices {
		if m.Selected[i] {
			out = append(out, string(c.Agent))
		}
	}
	return out
}

func RunInit(input io.Reader, output io.Writer) (InitModel, error) {
	model, err := tea.NewProgram(NewInitModel(), tea.WithInput(input), tea.WithOutput(output)).Run()
	if err != nil {
		return InitModel{}, err
	}
	result, ok := model.(InitModel)
	if !ok {
		return InitModel{}, fmt.Errorf("unexpected wizard model %T", model)
	}
	return result, nil
}
