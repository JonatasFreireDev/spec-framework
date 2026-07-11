package wizard

import (
	"errors"
	"fmt"
	"io"
	"strings"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/JonatasFreireDev/spec-framework/internal/install"
)

// choice pairs a human-facing label with the agent it installs.
type choice struct {
	Name  string
	Agent install.Agent
}

// choices are the agent skill formats offered by the installer, in order.
var choices = []choice{
	{"Codex", install.Codex},
	{"Cursor", install.Cursor},
	{"Claude Code", install.Claude},
}

// Result captures the outcome of the interactive init wizard.
type Result struct {
	Target    string
	Agents    []install.Agent
	Confirmed bool
	Cancelled bool
}

// AgentNames returns the selected agent identifiers as strings.
func (r Result) AgentNames() []string {
	out := make([]string, 0, len(r.Agents))
	for _, a := range r.Agents {
		out = append(out, string(a))
	}
	return out
}

// agentOptions builds the multi-select options. None are preselected.
func agentOptions() []huh.Option[install.Agent] {
	opts := make([]huh.Option[install.Agent], 0, len(choices))
	for _, c := range choices {
		opts = append(opts, huh.NewOption(c.Name, c.Agent))
	}
	return opts
}

// summaryReserve is the horizontal space kept for the choices panel so the
// huh form on the left never grows into it.
const summaryReserve = 40

// minFormWidth keeps the question column usable on narrow terminals.
const minFormWidth = 32

// initModel wraps a huh form and, while the form is active, renders a live
// summary of the choices made so far in a panel to the right of the current
// question. Splitting the form into one group per field makes huh advance
// one question at a time.
type initModel struct {
	form      *huh.Form
	selected  *[]install.Agent
	target    *string
	confirmed *bool
	width     int
}

func (m initModel) Init() tea.Cmd { return m.form.Init() }

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		formWidth := ws.Width - summaryReserve
		if formWidth < minFormWidth {
			formWidth = minFormWidth
		}
		ws.Width = formWidth
		msg = ws
	}
	model, cmd := m.form.Update(msg)
	if f, ok := model.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
		return m, tea.Quit
	}
	return m, cmd
}

func (m initModel) View() tea.View {
	if m.form.State != huh.StateNormal {
		return tea.NewView("")
	}
	form := lipgloss.NewStyle().MarginRight(4).Render(m.form.View())
	joined := lipgloss.JoinHorizontal(lipgloss.Top, form, m.summaryView())
	return tea.NewView(joined)
}

// summaryView renders the choices-so-far panel shown beside the question.
func (m initModel) summaryView() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	label := lipgloss.NewStyle().Bold(true)
	filled := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	empty := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	value := func(s string) string {
		if strings.TrimSpace(s) == "" {
			return empty.Render("—")
		}
		return filled.Render(s)
	}

	agents := strings.Join(Result{Agents: *m.selected}.AgentNames(), ", ")

	var b strings.Builder
	b.WriteString(title.Render("Your choices") + "\n\n")
	b.WriteString(label.Render("Agents") + "\n  " + value(agents) + "\n\n")
	b.WriteString(label.Render("Target") + "\n  " + value(*m.target))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 2).
		Width(30)
	return box.Render(b.String())
}

// RunInit drives the interactive installer wizard and returns the plan.
func RunInit(input io.Reader, output io.Writer) (Result, error) {
	var selected []install.Agent
	var target string
	confirmed := true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[install.Agent]().
				Title("Which agent formats to install?").
				Options(agentOptions()...).
				// Height must cover every option plus the title line; huh
				// subtracts the title height from the viewport, so without
				// this the last option is pushed out of view and scrolls.
				Height(len(choices)+2).
				Filterable(false).
				Value(&selected).
				Validate(func(v []install.Agent) error {
					if len(v) == 0 {
						return errors.New("select at least one agent")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Target directory").
				Placeholder("../product").
				Value(&target).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("target directory is required")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Apply this installation plan?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirmed),
		),
	)

	model := initModel{form: form, selected: &selected, target: &target, confirmed: &confirmed}
	final, err := tea.NewProgram(model, tea.WithInput(input), tea.WithOutput(output)).Run()
	if err != nil {
		return Result{}, fmt.Errorf("init wizard: %w", err)
	}
	fm, ok := final.(initModel)
	if !ok {
		return Result{}, fmt.Errorf("unexpected wizard model %T", final)
	}
	if fm.form.State == huh.StateAborted {
		return Result{Cancelled: true}, nil
	}

	return Result{
		Target:    strings.TrimSpace(target),
		Agents:    selected,
		Confirmed: confirmed,
		Cancelled: !confirmed,
	}, nil
}
