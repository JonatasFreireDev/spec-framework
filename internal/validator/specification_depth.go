package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

type specificationContractRegistry struct {
	Version int `yaml:"version"`
	Modules []struct {
		ID                  string   `yaml:"id"`
		Path                string   `yaml:"path"`
		Tiers               []string `yaml:"tiers"`
		AllowsNotApplicable bool     `yaml:"allows_not_applicable"`
		RequiredSections    []string `yaml:"required_sections"`
	} `yaml:"modules"`
}

func validateSpecificationDepth(s Snapshot) []Diagnostic {
	var out []Diagnostic
	var registry *specificationContractRegistry
	for rel, contextText := range s.Text {
		if !strings.HasSuffix(rel, "/context.md") || !strings.Contains(rel, "/use-cases/") {
			continue
		}
		meta := metadata(contextText)
		if strings.EqualFold(meta["rigor_tier"], "N/A") {
			continue
		}
		base := filepath.ToSlash(filepath.Dir(rel))
		specPath := base + "/specification.md"
		specification, exists := s.Text[specPath]
		if !exists {
			continue
		}
		version := strings.TrimSpace(meta["specification_contract_version"])
		if version == "" {
			out = append(out, Diagnostic{Note, "specification-depth-migration", specPath, "Specification uses the compatible legacy depth contract.", "Audit the bundle and explicitly add specification_contract_version: 2 when the product is ready to migrate; upgrade never adds it automatically."})
			continue
		}
		if registry == nil {
			loaded, err := loadSpecificationContractRegistry(s.FrameworkRoot)
			if err != nil {
				out = append(out, Diagnostic{Error, "specification-contract-registry", "framework/skills/specification/references/contracts.yaml", err.Error(), "Restore the shipped Specification contract registry or remove the v2 opt-in to remain on the legacy contract."})
				return out
			}
			registry = &loaded
		}
		if version != fmt.Sprint(registry.Version) {
			out = append(out, Diagnostic{Error, "specification-contract-version", rel, "Unsupported specification_contract_version " + version + ".", fmt.Sprintf("Use %d or remove the field to remain on the legacy contract.", registry.Version)})
			continue
		}

		status := strings.ToLower(markdownStatus(specification))
		severity := Warning
		if status == "proposed" || requiresApproval(status) {
			severity = Error
		}
		for _, heading := range []string{"Evidence And Boundary", "Cross-Contract Synthesis", "Traceability Summary", "Adversarial Review", "Open Questions And Decisions"} {
			if !hasMarkdownSection(specification, heading) {
				out = append(out, Diagnostic{severity, "specification-depth", specPath, "Specification v2 root is missing section " + heading + ".", "Use the v2 root template as an index and cross-contract synthesis."})
			}
		}
		if hasSpecificationPlaceholder(specification) || containsDeferredSpecificationLanguage(specification) {
			out = append(out, Diagnostic{severity, "specification-depth", specPath, "Specification v2 root contains placeholder or deferred content.", "Replace template values and future promises with evidence, a concrete non-blocking assumption, or a blocking question."})
		}
		if (status == "proposed" || requiresApproval(status)) && hasBlockingSpecificationQuestion(specification) {
			out = append(out, Diagnostic{Error, "specification-depth", specPath, "Proposed-or-later Specification has an unresolved blocking question.", "Resolve the question or return the Specification to draft."})
		}

		tier := strings.ToUpper(meta["rigor_tier"])
		requirementOwner := map[string]string{}
		contractBodies := map[string][]string{}
		for _, module := range registry.Modules {
			if !containsFold(module.Tiers, tier) {
				continue
			}
			path := base + "/" + filepath.ToSlash(module.Path)
			body, ok := s.Text[path]
			if !ok {
				out = append(out, Diagnostic{severity, "specification-depth", path, "Specification v2 required contract is missing.", "Create it from the concern-specific template declared in contracts.yaml."})
				continue
			}
			contractStatus := strings.ToLower(markdownStatus(body))
			if contractStatus == "not_applicable" {
				if !module.AllowsNotApplicable {
					out = append(out, Diagnostic{severity, "specification-depth", path, module.ID + " cannot be not_applicable for rigor " + tier + ".", "Complete the required contract."})
				}
				if !validatorMeaningful(tableFields(body)["rationale"]) {
					out = append(out, Diagnostic{severity, "specification-depth", path, "Not-applicable contract has no concrete rationale.", "Explain the evidence-backed reason this concern does not apply."})
				}
				continue
			}
			if (status == "proposed" || requiresApproval(status)) && contractStatus == "draft" {
				out = append(out, Diagnostic{Error, "specification-depth", path, "Draft contract cannot feed a proposed-or-later Specification v2.", "Complete the module and make it proposed before advancing the root Specification."})
			}
			for _, heading := range module.RequiredSections {
				section := markdownSection(body, heading)
				if section == "" || hasSpecificationPlaceholder(section) {
					out = append(out, Diagnostic{severity, "specification-depth", path, module.ID + " contract has missing, empty, or placeholder section " + heading + ".", "Complete the concern-specific section with evidence-backed content."})
				}
			}
			if containsDeferredSpecificationLanguage(body) {
				out = append(out, Diagnostic{severity, "specification-depth", path, module.ID + " contract defers required content to a future evolution.", "Keep the Specification draft and write the contract now or record a blocking question."})
			}
			rows := markdownTableRows(body, "Requirements")
			if len(rows) == 0 {
				out = append(out, Diagnostic{severity, "specification-traceability", path, module.ID + " contract has no requirement rows.", "Add sourced REQ-* rows mapped to observable AC-* criteria."})
			}
			for _, row := range rows {
				id := firstRequirementID(row["id"])
				if id == "" {
					out = append(out, Diagnostic{severity, "specification-traceability", path, "Requirement row has no stable REQ-* id.", "Assign a unique REQ-* id."})
					continue
				}
				if owner, duplicate := requirementOwner[id]; duplicate {
					out = append(out, Diagnostic{severity, "specification-traceability", path, "Requirement " + id + " duplicates " + owner + ".", "Keep one canonical owner and link reuse instead of duplicating the requirement."})
				} else {
					requirementOwner[id] = path
				}
				if !regexp.MustCompile(`\bAC-\d+\b`).MatchString(row["acceptance criteria"]) {
					out = append(out, Diagnostic{severity, "specification-traceability", path, "Requirement " + id + " has no observable AC-* mapping.", "Link at least one acceptance criterion."})
				}
				if !validatorMeaningful(row["source"]) || !strings.Contains(row["source"], "](") {
					out = append(out, Diagnostic{severity, "specification-traceability", path, "Requirement " + id + " has no navigable source.", "Link the approved parent, evidence, or decision that supports the requirement."})
				}
			}
			normalized := normalizedSpecificationContractBody(body)
			if normalized != "" {
				contractBodies[normalized] = append(contractBodies[normalized], path)
			}
		}
		for _, paths := range contractBodies {
			if len(paths) < 2 {
				continue
			}
			sort.Strings(paths)
			for _, path := range paths {
				out = append(out, Diagnostic{severity, "specification-depth", path, "Contract body duplicates another concern-specific module: " + strings.Join(paths, ", ") + ".", "Replace the generic shell with concern-specific content."})
			}
		}
	}
	return out
}

func loadSpecificationContractRegistry(frameworkRoot string) (specificationContractRegistry, error) {
	var registry specificationContractRegistry
	if strings.TrimSpace(frameworkRoot) == "" {
		return registry, fmt.Errorf("framework root is required to load Specification contract registry")
	}
	var data []byte
	var err error
	for _, rel := range []string{"framework/skills/specification/references/contracts.yaml", "skills/specification/references/contracts.yaml"} {
		data, err = os.ReadFile(filepath.Join(frameworkRoot, filepath.FromSlash(rel)))
		if err == nil {
			break
		}
	}
	if err != nil {
		return registry, fmt.Errorf("cannot read Specification contract registry: %w", err)
	}
	if err := yaml.Unmarshal(data, &registry); err != nil {
		return registry, fmt.Errorf("invalid Specification contract registry: %w", err)
	}
	if registry.Version < 1 || len(registry.Modules) == 0 {
		return registry, fmt.Errorf("Specification contract registry has no version or modules")
	}
	return registry, nil
}

func containsFold(values []string, expected string) bool {
	for _, value := range values {
		if strings.EqualFold(value, expected) {
			return true
		}
	}
	return false
}

func hasMarkdownSection(text, heading string) bool {
	return markdownSection(text, heading) != ""
}

func markdownSection(text, heading string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	active := false
	var content []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if active {
				break
			}
			active = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), heading)
			continue
		}
		if active {
			content = append(content, line)
		}
	}
	return strings.TrimSpace(strings.Join(content, "\n"))
}

func hasSpecificationPlaceholder(text string) bool {
	return regexp.MustCompile("`\\[[^]\\n]+\\]`").MatchString(text) || regexp.MustCompile(`(?i)\b(TBD|TODO|PLACEHOLDER)\b`).MatchString(text)
}

func containsDeferredSpecificationLanguage(text string) bool {
	lower := strings.ToLower(text)
	for _, phrase := range []string{"will receive stable", "next approved content evolution", "to be defined later", "will be detailed later", "future evolution"} {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

func firstRequirementID(value string) string {
	return regexp.MustCompile(`\bREQ-\d+\b`).FindString(value)
}

func normalizedSpecificationContractBody(text string) string {
	lines := strings.Split(normalizedText(text), "\n")
	if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "# ") {
		lines = lines[1:]
	}
	return strings.ToLower(strings.TrimSpace(strings.Join(lines, "\n")))
}

func hasBlockingSpecificationQuestion(text string) bool {
	for _, row := range markdownTableRows(text, "Open Questions And Decisions") {
		question := strings.TrimSpace(row["question/decision"])
		blocks := strings.TrimSpace(row["blocks"])
		if validatorMeaningful(question) && validatorMeaningful(blocks) {
			return true
		}
	}
	return false
}
