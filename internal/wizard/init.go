package wizard

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	Target            string
	Agents            []install.Agent
	StartingPoint     string
	Sources           []string
	InstallImpeccable bool
	ImpeccableVersion string
	Confirmed         bool
	Cancelled         bool
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
	form              *huh.Form
	selected          *[]install.Agent
	target            *string
	startingPoint     *string
	sources           *string
	installImpeccable *bool
	impeccableVersion *string
	confirmed         *bool
	width             int
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
	b.WriteString("\n\n" + label.Render("Starting point") + "\n  " + value(*m.startingPoint))
	adapter := "skip"
	if *m.installImpeccable {
		adapter = "install @ " + *m.impeccableVersion
	}
	b.WriteString("\n\n" + label.Render("Impeccable") + "\n  " + value(adapter))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(0, 2).
		Width(30)
	return box.Render(b.String())
}

// RunInit drives the interactive installer wizard and returns the plan.
func RunInit(input io.Reader, output io.Writer) (Result, error) {
	selected := []install.Agent{install.Codex, install.Cursor, install.Claude}
	var target string
	startingPoint := "new-product"
	var sources string
	installImpeccable := false
	impeccableVersion := "latest"
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
			huh.NewSelect[string]().
				Title("Where is this project starting?").
				Options(
					huh.NewOption("New product", "new-product"),
					huh.NewOption("Existing product", "existing-product"),
					huh.NewOption("Existing documents / epics / PRDs", "existing-documents"),
					huh.NewOption("Existing feature", "existing-feature"),
					huh.NewOption("Existing implementation", "existing-implementation"),
					huh.NewOption("Audit only", "audit-only"),
				).
				Value(&startingPoint),
		),
		huh.NewGroup(
			huh.NewInput().Title("Source paths (comma-separated; required for existing documents)").Value(&sources),
		).WithHideFunc(func() bool {
			return !showSourcePaths(startingPoint)
		}),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Install the optional Impeccable design adapter?").
				Description("Runs the official version-pinned npx installer after product initialization.").
				Affirmative("Install").
				Negative("Skip").
				Value(&installImpeccable),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Impeccable CLI version").
				Description("Use latest to resolve the current recommended version, or enter an exact semantic version.").
				Placeholder("latest").
				Value(&impeccableVersion).
				Validate(func(s string) error {
					if installImpeccable && strings.TrimSpace(s) == "" {
						return errors.New("an explicit Impeccable version is required")
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
				Value(&confirmed).
				Validate(func(bool) error {
					return validateSources(startingPoint, sources)
				}),
		),
	)

	model := initModel{form: form, selected: &selected, target: &target, startingPoint: &startingPoint, sources: &sources, installImpeccable: &installImpeccable, impeccableVersion: &impeccableVersion, confirmed: &confirmed}
	final, err := tea.NewProgram(model, programOptions(input, output)...).Run()
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
		Target:            strings.TrimSpace(target),
		Agents:            selected,
		StartingPoint:     startingPoint,
		Sources:           splitValues(sources),
		InstallImpeccable: installImpeccable,
		ImpeccableVersion: strings.TrimSpace(impeccableVersion),
		Confirmed:         confirmed,
		Cancelled:         !confirmed,
	}, nil
}

func validateSources(startingPoint, sources string) error {
	if showSourcePaths(startingPoint) && strings.TrimSpace(sources) == "" {
		return errors.New("at least one source path is required for existing documents")
	}
	return nil
}

func showSourcePaths(startingPoint string) bool {
	return startingPoint == "existing-documents"
}

func programOptions(input io.Reader, output io.Writer) []tea.ProgramOption {
	options := []tea.ProgramOption{tea.WithOutput(output)}
	// Let Bubble Tea manage the process stdin so it can fall back to /dev/tty
	// when a Unix launcher redirects stdin. Injected readers still support tests.
	if input != os.Stdin {
		options = append(options, tea.WithInput(input))
	}
	return options
}

func splitValues(value string) []string {
	var out []string
	for _, part := range strings.Split(value, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
