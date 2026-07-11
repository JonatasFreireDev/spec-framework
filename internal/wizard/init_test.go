package wizard

import (
	tea "charm.land/bubbletea/v2"
	"testing"
)

func key(code rune, text string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code, Text: text})
}
func TestInitModelBuildsConfirmedMultiAgentPlan(t *testing.T) {
	m := NewInitModel()
	next, _ := m.Update(key(tea.KeyDown, ""))
	m = next.(InitModel)
	next, _ = m.Update(key(tea.KeySpace, " "))
	m = next.(InitModel)
	next, _ = m.Update(key(tea.KeyEnter, ""))
	m = next.(InitModel)
	for _, r := range "../product" {
		next, _ = m.Update(key(r, string(r)))
		m = next.(InitModel)
	}
	next, _ = m.Update(key(tea.KeyEnter, ""))
	m = next.(InitModel)
	next, _ = m.Update(key('y', "y"))
	m = next.(InitModel)
	if !m.Confirmed || m.Target != "../product" || len(m.Agents()) != 2 {
		t.Fatalf("model=%+v agents=%v", m, m.Agents())
	}
}
