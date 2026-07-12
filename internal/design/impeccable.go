package design

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type AdapterPlan struct {
	Adapter       string   `json:"adapter"`
	Mode          string   `json:"mode"`
	Maturity      string   `json:"maturity"`
	WorkingDir    string   `json:"workingDir"`
	Commands      []string `json:"commands"`
	OutputScope   string   `json:"outputScope"`
	NonProduction bool     `json:"nonProduction"`
	RequiresHuman bool     `json:"requiresHumanReview"`
}

func ImpeccablePlan(productRoot, useCase, maturity string) (AdapterPlan, error) {
	uc, _, err := resolveUseCase(productRoot, useCase)
	if err != nil {
		return AdapterPlan{}, err
	}
	if !oneOf(maturity, "wireframe", "mockup", "prototype") {
		return AdapterPlan{}, fmt.Errorf("impeccable target maturity must be wireframe, mockup, or prototype")
	}
	commands := map[string][]string{
		"wireframe": {"/impeccable shape"},
		"mockup":    {"/impeccable craft", "/impeccable critique"},
		"prototype": {"/impeccable craft", "/impeccable harden", "/impeccable adapt", "/impeccable audit", "/impeccable polish"},
	}
	return AdapterPlan{Adapter: "impeccable", Mode: "generate", Maturity: maturity, WorkingDir: filepath.ToSlash(uc), Commands: commands[maturity], OutputScope: filepath.ToSlash(filepath.Join(productRoot, "design", "use-cases", filepath.Base(uc))), NonProduction: true, RequiresHuman: true}, nil
}

func EncodeAdapterPlan(plan AdapterPlan) string {
	data, _ := json.MarshalIndent(plan, "", "  ")
	return string(data) + "\n"
}
