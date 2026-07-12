package wizard

import (
	"fmt"
	"io"
	"strings"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

type migrationChoice struct {
	ID, Type, Scope string
	Apply           bool
}
type migrationModel struct{ form *huh.Form }

func (m migrationModel) Init() tea.Cmd { return m.form.Init() }
func (m migrationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.form.Update(msg)
	if f, ok := model.(*huh.Form); ok {
		m.form = f
	}
	if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
		return m, tea.Quit
	}
	return m, cmd
}
func (m migrationModel) View() tea.View {
	if m.form.State != huh.StateNormal {
		return tea.NewView("")
	}
	return tea.NewView(m.form.View())
}

func RunDecisionMigration(input io.Reader, output io.Writer, plan workflow.DecisionMigrationPlan) ([]workflow.DecisionMigrationItem, bool, error) {
	choices := make([]migrationChoice, len(plan.Items))
	groups := make([]*huh.Group, 0, len(plan.Items)+1)
	for i, item := range plan.Items {
		choices[i] = migrationChoice{ID: item.ID, Type: item.InferredType, Scope: strings.Join(item.Scope, ", "), Apply: true}
		c := &choices[i]
		groups = append(groups, huh.NewGroup(huh.NewConfirm().Title("Migrate "+item.ID+"?").Affirmative("Yes").Negative("Skip").Value(&c.Apply), huh.NewSelect[string]().Title(item.ID+" type").Options(huh.NewOption("Product", "product"), huh.NewOption("Architecture", "architecture"), huh.NewOption("Security", "security"), huh.NewOption("Data", "data"), huh.NewOption("Delivery", "delivery")).Value(&c.Type), huh.NewInput().Title(item.ID+" scope paths (comma-separated)").Value(&c.Scope)))
	}
	confirmed := true
	groups = append(groups, huh.NewGroup(huh.NewConfirm().Title("Apply migration with backup?").Affirmative("Apply").Negative("Cancel").Value(&confirmed)))
	form := huh.NewForm(groups...)
	final, err := tea.NewProgram(migrationModel{form: form}, tea.WithInput(input), tea.WithOutput(output)).Run()
	if err != nil {
		return nil, false, fmt.Errorf("decision migration wizard: %w", err)
	}
	fm, ok := final.(migrationModel)
	if !ok {
		return nil, false, fmt.Errorf("unexpected wizard model %T", final)
	}
	if fm.form.State == huh.StateAborted || !confirmed {
		return nil, false, nil
	}
	var out []workflow.DecisionMigrationItem
	for _, c := range choices {
		if !c.Apply {
			continue
		}
		var scope []string
		for _, x := range strings.Split(c.Scope, ",") {
			if x = strings.TrimSpace(x); x != "" {
				scope = append(scope, x)
			}
		}
		out = append(out, workflow.DecisionMigrationItem{ID: c.ID, InferredType: c.Type, Scope: scope})
	}
	return out, true, nil
}
