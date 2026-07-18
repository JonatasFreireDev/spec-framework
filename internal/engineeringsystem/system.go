package engineeringsystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

var allowedTriggers = map[string]bool{
	"new_dependency":               true,
	"external_integration":         true,
	"data_ownership_change":        true,
	"migration":                    true,
	"architecture_boundary_change": true,
	"deployment_change":            true,
	"security_boundary_change":     true,
	"operational_change":           true,
}

var allowedMaturity = map[string]bool{
	"baseline": true,
	"mapped":   true,
	"governed": true,
	"verified": true,
	"operated": true,
}

var technicalEntityPrefixes = map[string]string{
	"systems": "SYS-", "applications": "APP-", "components": "CMP-",
	"repositories": "REPO-", "data_stores": "DATA-", "interfaces": "IFACE-", "deployments": "DEPLOY-",
}

var technicalEntityTypes = map[string]string{
	"systems": "system", "applications": "application", "components": "component",
	"repositories": "repository", "data_stores": "data-store", "interfaces": "interface", "deployments": "deployment",
}

var allowedStandardCategories = map[string]bool{
	"architecture": true, "code": true, "api": true, "events": true, "data": true,
	"dependencies": true, "security": true, "observability": true, "testing": true, "delivery": true,
}

var allowedStandardLevels = map[string]bool{
	"required": true, "recommended": true, "experimental": true, "deprecated": true,
}

type Area struct {
	Name     string `json:"name"`
	Contract string `json:"contract"`
	Maturity string `json:"maturity"`
	Evidence int    `json:"evidence"`
}

type Inspection struct {
	ID                     string            `json:"id"`
	Status                 string            `json:"status"`
	Version                string            `json:"version"`
	OriginMode             string            `json:"originMode"`
	Scope                  string            `json:"scope"`
	Areas                  []Area            `json:"areas"`
	Decisions              int               `json:"decisions"`
	Standards              int               `json:"standards"`
	FitnessFunctions       int               `json:"fitnessFunctions"`
	QualitySystem          bool              `json:"qualitySystem"`
	QualityExceptions      []string          `json:"qualityExceptions,omitempty"`
	QualityExceptionScopes map[string]string `json:"qualityExceptionScopes,omitempty"`
	QualityEnvironments    []string          `json:"qualityEnvironments,omitempty"`
	QualityTestDataClasses []string          `json:"qualityTestDataClasses,omitempty"`
	QualityPlatforms       []string          `json:"qualityPlatforms,omitempty"`
	Blockers               []string          `json:"blockers,omitempty"`
}

type contextDocument struct {
	ID                  string   `yaml:"id"`
	Status              string   `yaml:"status"`
	Version             string   `yaml:"version"`
	OriginMode          string   `yaml:"origin_mode"`
	EngineeringTriggers []string `yaml:"engineering_triggers"`
}

type catalogDocument struct {
	SchemaVersion    int                    `yaml:"schema_version"`
	ID               string                 `yaml:"id"`
	Status           string                 `yaml:"status"`
	Version          string                 `yaml:"version"`
	OriginMode       string                 `yaml:"origin_mode"`
	Scope            string                 `yaml:"scope"`
	Areas            map[string]catalogArea `yaml:"areas"`
	Decisions        []any                  `yaml:"decisions"`
	Standards        []any                  `yaml:"standards"`
	FitnessFunctions []any                  `yaml:"fitness_functions"`
}

type catalogArea struct {
	Contract string   `yaml:"contract"`
	Maturity string   `yaml:"maturity"`
	Evidence []string `yaml:"evidence"`
}

type qualityCatalogDocument struct {
	SchemaVersion     int                           `yaml:"schema_version"`
	EngineeringSystem string                        `yaml:"engineering_system"`
	Version           string                        `yaml:"version"`
	Status            string                        `yaml:"status"`
	Areas             map[string]qualityCatalogArea `yaml:"areas"`
	GateSource        string                        `yaml:"gate_source"`
	Exceptions        qualityExceptionPolicy        `yaml:"exceptions"`
	Environments      []string                      `yaml:"environments"`
	TestDataClasses   []string                      `yaml:"test_data_classes"`
	Platforms         []string                      `yaml:"platforms"`
}

type qualityCatalogArea struct {
	Maturity         string   `yaml:"maturity"`
	Policy           string   `yaml:"policy"`
	DelegatedGate    string   `yaml:"delegated_gate"`
	RequiredEvidence []string `yaml:"required_evidence"`
}

type qualityExceptionPolicy struct {
	RequireOwner          bool                     `yaml:"require_owner"`
	RequireResidualRisk   bool                     `yaml:"require_residual_risk"`
	RequireExpiryOrReview bool                     `yaml:"require_expiry_or_review"`
	Records               []qualityExceptionRecord `yaml:"records"`
}

type qualityExceptionRecord struct {
	ID             string `yaml:"id"`
	Scope          string `yaml:"scope"`
	Owner          string `yaml:"owner"`
	Rationale      string `yaml:"rationale"`
	ResidualRisk   string `yaml:"residual_risk"`
	Mitigation     string `yaml:"mitigation"`
	ExpiryOrReview string `yaml:"expiry_or_review"`
	ReentryGate    string `yaml:"reentry_gate"`
	Status         string `yaml:"status"`
}

type technicalCatalogDocument struct {
	SchemaVersion int                          `yaml:"schema_version"`
	Entities      map[string]map[string]string `yaml:"entities"`
	Relations     []technicalRelationDocument  `yaml:"relations"`
}

type technicalEntityDocument struct {
	SchemaVersion int    `yaml:"schema_version"`
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
	Status        string `yaml:"status"`
}

type technicalRelationDocument struct {
	ID       string   `yaml:"id"`
	Type     string   `yaml:"type"`
	Source   string   `yaml:"source"`
	Target   string   `yaml:"target"`
	Evidence []string `yaml:"evidence"`
}

type standardsCatalogDocument struct {
	SchemaVersion int               `yaml:"schema_version"`
	Profiles      map[string]string `yaml:"profiles"`
	Standards     map[string]string `yaml:"standards"`
	Exceptions    map[string]string `yaml:"exceptions"`
}

type standardProfileDocument struct {
	SchemaVersion int      `yaml:"schema_version"`
	ID            string   `yaml:"id"`
	Version       string   `yaml:"version"`
	Status        string   `yaml:"status"`
	Extends       []string `yaml:"extends"`
	Standards     []string `yaml:"standards"`
}

type standardDocument struct {
	SchemaVersion int                    `yaml:"schema_version"`
	ID            string                 `yaml:"id"`
	Version       string                 `yaml:"version"`
	Status        string                 `yaml:"status"`
	Category      string                 `yaml:"category"`
	Level         string                 `yaml:"level"`
	Rules         []standardRuleDocument `yaml:"rules"`
}

type standardRuleDocument struct {
	ID           string   `yaml:"id"`
	Requirement  string   `yaml:"requirement"`
	Verification []string `yaml:"verification"`
}

type standardExceptionDocument struct {
	SchemaVersion int      `yaml:"schema_version"`
	ID            string   `yaml:"id"`
	Standard      string   `yaml:"standard"`
	Scope         []string `yaml:"scope"`
	Owner         string   `yaml:"owner"`
	Rationale     string   `yaml:"rationale"`
	ResidualRisk  string   `yaml:"residual_risk"`
	Mitigation    string   `yaml:"mitigation"`
	ExpiresOn     string   `yaml:"expires_on"`
	ReentryGate   string   `yaml:"reentry_gate"`
	Status        string   `yaml:"status"`
}

type operationsCatalogDocument struct {
	SchemaVersion int               `yaml:"schema_version"`
	Environments  map[string]string `yaml:"environments"`
	Deployments   map[string]string `yaml:"deployments"`
	Runbooks      map[string]string `yaml:"runbooks"`
}

func Inspect(root string) (Inspection, error) {
	dir := filepath.Join(root, "engineering")
	contextData, err := os.ReadFile(filepath.Join(dir, "context.md"))
	if err != nil {
		return Inspection{}, err
	}
	var context contextDocument
	if err := yaml.Unmarshal([]byte(yamlPayload(string(contextData))), &context); err != nil {
		return Inspection{}, fmt.Errorf("engineering/context.md has invalid YAML metadata: %w", err)
	}
	catalogData, err := os.ReadFile(filepath.Join(dir, "engineering-system.yaml"))
	if err != nil {
		return Inspection{}, err
	}
	var catalog catalogDocument
	if err := yaml.Unmarshal(catalogData, &catalog); err != nil {
		return Inspection{}, fmt.Errorf("engineering-system.yaml is invalid YAML: %w", err)
	}
	canonicalData, err := os.ReadFile(filepath.Join(dir, "engineering-system.md"))
	if err != nil {
		return Inspection{}, err
	}
	canonical := markdownSnapshot(string(canonicalData))
	i := Inspection{
		ID:               context.ID,
		Status:           context.Status,
		Version:          context.Version,
		OriginMode:       context.OriginMode,
		Scope:            catalog.Scope,
		Decisions:        len(catalog.Decisions),
		Standards:        len(catalog.Standards),
		FitnessFunctions: len(catalog.FitnessFunctions),
	}
	if catalog.SchemaVersion != 1 {
		i.Blockers = append(i.Blockers, "catalog schema_version must be 1")
	}
	if !regexp.MustCompile(`^ENGSYS-[A-Z0-9-]+$`).MatchString(i.ID) {
		i.Blockers = append(i.Blockers, "context engineering system id is invalid")
	}
	if !regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`).MatchString(i.Version) {
		i.Blockers = append(i.Blockers, "context semantic version is invalid")
	}
	if !oneOf(i.OriginMode, "generate", "evolve", "adopt") {
		i.Blockers = append(i.Blockers, "context origin mode is invalid")
	}
	if i.Scope == "" {
		i.Blockers = append(i.Blockers, "catalog scope is missing")
	}
	for field, values := range map[string][2]string{
		"id":          {context.ID, catalog.ID},
		"status":      {context.Status, catalog.Status},
		"version":     {context.Version, catalog.Version},
		"origin_mode": {context.OriginMode, catalog.OriginMode},
	} {
		if values[1] == "" || values[0] != values[1] {
			i.Blockers = append(i.Blockers, fmt.Sprintf("context and catalog %s do not match", field))
		}
	}
	for field, values := range map[string][2]string{
		"id":      {context.ID, canonical["id"]},
		"status":  {context.Status, canonical["status"]},
		"version": {context.Version, canonical["version"]},
	} {
		if values[1] == "" || values[0] != values[1] {
			i.Blockers = append(i.Blockers, fmt.Sprintf("context and canonical %s do not match", field))
		}
	}
	if len(catalog.Areas) == 0 {
		i.Blockers = append(i.Blockers, "catalog areas are missing")
	}
	for name, source := range catalog.Areas {
		area := Area{Name: name, Contract: source.Contract, Maturity: source.Maturity, Evidence: len(source.Evidence)}
		i.Areas = append(i.Areas, area)
		if area.Contract == "" || area.Maturity == "" {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s is missing contract or maturity", name))
			continue
		}
		if !allowedMaturity[area.Maturity] {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s has invalid maturity %s", name, area.Maturity))
		}
		contractPath := filepath.Clean(filepath.Join(dir, filepath.FromSlash(area.Contract)))
		relative, relErr := filepath.Rel(dir, contractPath)
		if relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(filepath.FromSlash(area.Contract)) {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s contract %s escapes engineering", name, area.Contract))
		} else if _, err := os.Stat(contractPath); err != nil {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s contract %s is missing", name, area.Contract))
		}
		if area.Maturity != "baseline" && area.Evidence == 0 {
			i.Blockers = append(i.Blockers, fmt.Sprintf("area %s maturity %s requires evidence", name, area.Maturity))
		}
	}
	if quality, exists := catalog.Areas["quality"]; exists {
		switch filepath.ToSlash(quality.Contract) {
		case "quality/quality-system.md":
			i.QualitySystem = true
			i.Blockers = append(i.Blockers, validateQualitySystem(dir, i)...)
			i.QualityExceptions, i.QualityExceptionScopes, i.Blockers = inspectQualityExceptions(dir, i.Blockers)
			i.QualityEnvironments, i.QualityTestDataClasses, i.QualityPlatforms = qualityDimensions(dir)
		case "quality/quality-model.md":
			// Legacy contract remains valid until explicit migration.
		default:
			i.Blockers = append(i.Blockers, "quality area contract must be quality/quality-system.md or legacy quality/quality-model.md")
		}
	}
	if source, exists := catalog.Areas["technical_catalog"]; exists {
		i.Blockers = append(i.Blockers, validateTechnicalCatalog(dir, source.Contract)...)
	}
	if source, exists := catalog.Areas["standards"]; exists {
		i.Blockers = append(i.Blockers, validateStandardsCatalog(dir, source.Contract)...)
	}
	if source, exists := catalog.Areas["operations"]; exists {
		i.Blockers = append(i.Blockers, validateOperationsCatalog(dir, source.Contract)...)
	}
	sort.Slice(i.Areas, func(left, right int) bool { return i.Areas[left].Name < i.Areas[right].Name })
	i.Blockers = unique(i.Blockers)
	sort.Strings(i.Blockers)
	return i, nil
}

func validateTechnicalCatalog(engineeringDir, contract string) []string {
	var catalog technicalCatalogDocument
	path, blockers := loadEngineeringYAML(engineeringDir, contract, "technical catalog", &catalog)
	if len(blockers) != 0 {
		return blockers
	}
	if catalog.SchemaVersion != 1 {
		blockers = append(blockers, "technical catalog schema_version must be 1")
	}
	for category := range catalog.Entities {
		if technicalEntityPrefixes[category] == "" {
			blockers = append(blockers, "technical catalog has unknown entity category "+category)
		}
	}
	entityIDs := map[string]bool{}
	for category, prefix := range technicalEntityPrefixes {
		for id, reference := range catalog.Entities[category] {
			if !regexp.MustCompile(`^[A-Z][A-Z0-9-]+$`).MatchString(id) || !strings.HasPrefix(id, prefix) {
				blockers = append(blockers, fmt.Sprintf("technical catalog %s id %s must start with %s", category, id, prefix))
			}
			entityIDs[id] = true
			entityPath, refBlockers := resolvedCatalogReference(engineeringDir, filepath.Dir(path), reference, "technical entity "+id)
			blockers = append(blockers, refBlockers...)
			if len(refBlockers) != 0 {
				continue
			}
			var entity technicalEntityDocument
			if err := readYAML(entityPath, &entity); err != nil {
				blockers = append(blockers, "technical entity "+id+" is invalid YAML")
				continue
			}
			if entity.SchemaVersion != 1 || entity.ID != id || entity.Type != technicalEntityTypes[category] || strings.TrimSpace(entity.Status) == "" {
				blockers = append(blockers, "technical entity "+id+" has invalid schema, identity, type, or status")
			}
		}
	}
	seenRelations := map[string]bool{}
	for _, relation := range catalog.Relations {
		if !regexp.MustCompile(`^REL-[A-Z0-9-]+$`).MatchString(relation.ID) || strings.TrimSpace(relation.Type) == "" {
			blockers = append(blockers, "technical relation "+relation.ID+" has invalid identity or type")
		}
		if seenRelations[relation.ID] {
			blockers = append(blockers, "technical relation "+relation.ID+" is duplicated")
		}
		seenRelations[relation.ID] = true
		if !entityIDs[relation.Source] || !entityIDs[relation.Target] {
			blockers = append(blockers, "technical relation "+relation.ID+" references an unknown source or target")
		}
	}
	return blockers
}

func validateStandardsCatalog(engineeringDir, contract string) []string {
	var catalog standardsCatalogDocument
	path, blockers := loadEngineeringYAML(engineeringDir, contract, "standards catalog", &catalog)
	if len(blockers) != 0 {
		return blockers
	}
	if catalog.SchemaVersion != 1 {
		blockers = append(blockers, "standards catalog schema_version must be 1")
	}
	base := filepath.Dir(path)
	profiles := map[string]standardProfileDocument{}
	for id, reference := range catalog.Profiles {
		if !strings.HasPrefix(id, "PROFILE-") {
			blockers = append(blockers, "standards profile id "+id+" must start with PROFILE-")
		}
		profilePath, refBlockers := resolvedCatalogReference(engineeringDir, base, reference, "standards profile "+id)
		blockers = append(blockers, refBlockers...)
		if len(refBlockers) != 0 {
			continue
		}
		var profile standardProfileDocument
		if err := readYAML(profilePath, &profile); err != nil {
			blockers = append(blockers, "standards profile "+id+" is invalid YAML")
			continue
		}
		if profile.SchemaVersion != 1 || profile.ID != id || !semanticVersion(profile.Version) {
			blockers = append(blockers, "standards profile "+id+" has invalid schema, identity, or semantic version")
		}
		profiles[id] = profile
	}
	for id, reference := range catalog.Standards {
		if !strings.HasPrefix(id, "STD-") {
			blockers = append(blockers, "standard id "+id+" must start with STD-")
		}
		standardPath, refBlockers := resolvedCatalogReference(engineeringDir, base, reference, "standard "+id)
		blockers = append(blockers, refBlockers...)
		if len(refBlockers) != 0 {
			continue
		}
		var standard standardDocument
		if err := readYAML(standardPath, &standard); err != nil {
			blockers = append(blockers, "standard "+id+" is invalid YAML")
			continue
		}
		if standard.SchemaVersion != 1 || standard.ID != id || !semanticVersion(standard.Version) {
			blockers = append(blockers, "standard "+id+" has invalid schema, identity, or semantic version")
		}
		if !allowedStandardCategories[standard.Category] {
			blockers = append(blockers, "standard "+id+" has invalid category "+standard.Category)
		}
		if !allowedStandardLevels[standard.Level] {
			blockers = append(blockers, "standard "+id+" has invalid obligation level "+standard.Level)
		}
		if len(standard.Rules) == 0 {
			blockers = append(blockers, "standard "+id+" must declare at least one verifiable rule")
		}
		seenRules := map[string]bool{}
		for _, rule := range standard.Rules {
			if !strings.HasPrefix(rule.ID, id+"-R") || strings.TrimSpace(rule.Requirement) == "" || len(rule.Verification) == 0 {
				blockers = append(blockers, "standard "+id+" has an invalid rule identity, requirement, or verification")
			}
			if seenRules[rule.ID] {
				blockers = append(blockers, "standard "+id+" rule "+rule.ID+" is duplicated")
			}
			seenRules[rule.ID] = true
		}
	}
	for id, reference := range catalog.Exceptions {
		if !strings.HasPrefix(id, "STDEX-") {
			blockers = append(blockers, "standard exception id "+id+" must start with STDEX-")
		}
		exceptionPath, refBlockers := resolvedCatalogReference(engineeringDir, base, reference, "standard exception "+id)
		blockers = append(blockers, refBlockers...)
		if len(refBlockers) != 0 {
			continue
		}
		var exception standardExceptionDocument
		if err := readYAML(exceptionPath, &exception); err != nil {
			blockers = append(blockers, "standard exception "+id+" is invalid YAML")
			continue
		}
		if exception.SchemaVersion != 1 || exception.ID != id || catalog.Standards[exception.Standard] == "" {
			blockers = append(blockers, "standard exception "+id+" has invalid schema, identity, or standard reference")
		}
		if len(exception.Scope) == 0 || strings.TrimSpace(exception.Owner) == "" || strings.TrimSpace(exception.Rationale) == "" || strings.TrimSpace(exception.ResidualRisk) == "" || strings.TrimSpace(exception.Mitigation) == "" || strings.TrimSpace(exception.ReentryGate) == "" {
			blockers = append(blockers, "standard exception "+id+" lacks scope, owner, rationale, residual risk, mitigation, or re-entry gate")
		}
		if !oneOf(exception.Status, "open", "closed") {
			blockers = append(blockers, "standard exception "+id+" has invalid status "+exception.Status)
		}
		expires, dateErr := time.Parse("2006-01-02", exception.ExpiresOn)
		if dateErr != nil {
			blockers = append(blockers, "standard exception "+id+" expires_on must use YYYY-MM-DD")
		} else if exception.Status == "open" && !expires.After(time.Now().UTC().Truncate(24*time.Hour)) {
			blockers = append(blockers, "standard exception "+id+" is open but expired")
		}
	}
	for id, profile := range profiles {
		for _, parent := range profile.Extends {
			if _, exists := catalog.Profiles[parent]; !exists {
				blockers = append(blockers, "standards profile "+id+" extends unknown profile "+parent)
			}
		}
		for _, standard := range profile.Standards {
			if _, exists := catalog.Standards[standard]; !exists {
				blockers = append(blockers, "standards profile "+id+" references unknown standard "+standard)
			}
		}
	}
	blockers = append(blockers, validateProfileCycles(profiles)...)
	return blockers
}

func validateOperationsCatalog(engineeringDir, contract string) []string {
	var catalog operationsCatalogDocument
	path, blockers := loadEngineeringYAML(engineeringDir, contract, "operations catalog", &catalog)
	if len(blockers) != 0 {
		return blockers
	}
	if catalog.SchemaVersion != 1 {
		blockers = append(blockers, "operations catalog schema_version must be 1")
	}
	for category, entries := range map[string]map[string]string{"environment": catalog.Environments, "deployment": catalog.Deployments, "runbook": catalog.Runbooks} {
		for id, reference := range entries {
			blockers = append(blockers, validateCatalogReference(engineeringDir, filepath.Dir(path), reference, category+" "+id)...)
		}
	}
	return blockers
}

func loadEngineeringYAML(engineeringDir, contract, label string, target any) (string, []string) {
	path, blockers := resolvedCatalogReference(engineeringDir, engineeringDir, contract, label)
	if len(blockers) != 0 {
		return path, blockers
	}
	if err := readYAML(path, target); err != nil {
		return path, []string{label + " is invalid YAML"}
	}
	return path, nil
}

func validateCatalogReference(engineeringDir, baseDir, reference, label string) []string {
	_, blockers := resolvedCatalogReference(engineeringDir, baseDir, reference, label)
	return blockers
}

func resolvedCatalogReference(engineeringDir, baseDir, reference, label string) (string, []string) {
	reference = strings.TrimSpace(reference)
	path := filepath.Clean(filepath.Join(baseDir, filepath.FromSlash(reference)))
	relative, err := filepath.Rel(engineeringDir, path)
	if reference == "" || filepath.IsAbs(filepath.FromSlash(reference)) || err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return path, []string{label + " reference escapes engineering"}
	}
	info, statErr := os.Stat(path)
	if statErr != nil || info.IsDir() {
		return path, []string{label + " reference " + reference + " is missing"}
	}
	return path, nil
}

func readYAML(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, target)
}

func semanticVersion(version string) bool {
	return regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:[-+][0-9A-Za-z.-]+)?$`).MatchString(version)
}

func validateProfileCycles(profiles map[string]standardProfileDocument) []string {
	state := map[string]int{}
	var blockers []string
	var visit func(string)
	visit = func(id string) {
		if state[id] == 1 {
			blockers = append(blockers, "standards profiles contain an inheritance cycle at "+id)
			return
		}
		if state[id] == 2 {
			return
		}
		state[id] = 1
		for _, parent := range profiles[id].Extends {
			if _, exists := profiles[parent]; exists {
				visit(parent)
			}
		}
		state[id] = 2
	}
	for id := range profiles {
		visit(id)
	}
	return blockers
}

func inspectQualityExceptions(engineeringDir string, blockers []string) ([]string, map[string]string, []string) {
	data, err := os.ReadFile(filepath.Join(engineeringDir, "quality", "quality-system.yaml"))
	if err != nil {
		return nil, nil, blockers
	}
	var catalog qualityCatalogDocument
	if yaml.Unmarshal(data, &catalog) != nil {
		return nil, nil, blockers
	}
	seen := map[string]bool{}
	scopes := map[string]string{}
	var ids []string
	today := time.Now().UTC().Truncate(24 * time.Hour)
	for _, record := range catalog.Exceptions.Records {
		if !regexp.MustCompile(`^QEX-[A-Z0-9-]+$`).MatchString(record.ID) {
			blockers = append(blockers, "quality exception id "+record.ID+" is invalid")
			continue
		}
		if seen[record.ID] {
			blockers = append(blockers, "quality exception "+record.ID+" is duplicated")
			continue
		}
		seen[record.ID] = true
		if strings.TrimSpace(record.Scope) == "" || strings.TrimSpace(record.Owner) == "" || strings.TrimSpace(record.Rationale) == "" || strings.TrimSpace(record.ResidualRisk) == "" || strings.TrimSpace(record.Mitigation) == "" || strings.TrimSpace(record.ExpiryOrReview) == "" || strings.TrimSpace(record.ReentryGate) == "" {
			blockers = append(blockers, "quality exception "+record.ID+" lacks scope, owner, rationale, residual risk, mitigation, expiry/review, or re-entry gate")
		}
		if !validExceptionScope(record.Scope) {
			blockers = append(blockers, "quality exception "+record.ID+" scope must be product or a safe domains/ path")
		}
		if !oneOf(record.Status, "open", "closed") {
			blockers = append(blockers, "quality exception "+record.ID+" has invalid status "+record.Status)
		}
		expires, dateErr := time.Parse("2006-01-02", record.ExpiryOrReview)
		if dateErr != nil {
			blockers = append(blockers, "quality exception "+record.ID+" expiry_or_review must use YYYY-MM-DD")
		}
		if record.Status == "open" && dateErr == nil && !expires.After(today) {
			blockers = append(blockers, "quality exception "+record.ID+" is open but expired or due for review")
		}
		if record.Status == "open" && dateErr == nil && expires.After(today) && validExceptionScope(record.Scope) {
			ids = append(ids, record.ID)
			scopes[record.ID] = filepath.ToSlash(strings.TrimSpace(record.Scope))
		}
	}
	sort.Strings(ids)
	return ids, scopes, blockers
}

func validExceptionScope(scope string) bool {
	scope = filepath.ToSlash(strings.TrimSpace(scope))
	if scope == "product" {
		return true
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(scope)))
	return strings.HasPrefix(clean, "domains/") && clean != "domains" && !strings.Contains(clean, "../") && !filepath.IsAbs(filepath.FromSlash(scope))
}

func qualityDimensions(engineeringDir string) ([]string, []string, []string) {
	data, err := os.ReadFile(filepath.Join(engineeringDir, "quality", "quality-system.yaml"))
	if err != nil {
		return nil, nil, nil
	}
	var catalog qualityCatalogDocument
	if yaml.Unmarshal(data, &catalog) != nil {
		return nil, nil, nil
	}
	return normalizedValues(catalog.Environments), normalizedValues(catalog.TestDataClasses), normalizedValues(catalog.Platforms)
}

func normalizedValues(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func validateQualitySystem(engineeringDir string, inspection Inspection) []string {
	humanData, err := os.ReadFile(filepath.Join(engineeringDir, "quality", "quality-system.md"))
	if err != nil {
		return []string{"quality system human contract quality/quality-system.md is missing"}
	}
	human := markdownSnapshot(string(humanData))
	humanMaturity := qualityHumanMaturity(string(humanData))
	var blockers []string
	if human["status"] != inspection.Status {
		blockers = append(blockers, "quality system human contract status must match the Engineering System status")
	}
	if human["engineering system"] != inspection.ID+" @ "+inspection.Version {
		blockers = append(blockers, "quality system human contract must pin the Engineering System id and version")
	}
	data, err := os.ReadFile(filepath.Join(engineeringDir, "quality", "quality-system.yaml"))
	if err != nil {
		return append(blockers, "quality system mechanical catalog quality/quality-system.yaml is missing")
	}
	var catalog qualityCatalogDocument
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return []string{"quality system mechanical catalog is invalid YAML"}
	}
	if catalog.SchemaVersion != 1 {
		blockers = append(blockers, "quality system schema_version must be 1")
	}
	if catalog.EngineeringSystem != inspection.ID {
		blockers = append(blockers, "quality system engineering_system must match the Engineering System id")
	}
	if catalog.Version != inspection.Version {
		blockers = append(blockers, "quality system version must match the Engineering System version")
	}
	if catalog.Status != inspection.Status {
		blockers = append(blockers, "quality system status must match the Engineering System status")
	}
	if filepath.ToSlash(catalog.GateSource) != "knowledge/conventions/gates.md" {
		blockers = append(blockers, "quality system gate_source must be knowledge/conventions/gates.md")
	}
	for _, required := range []string{"behavioral", "accessibility", "security_privacy", "performance_reliability", "observability"} {
		area, exists := catalog.Areas[required]
		if !exists {
			blockers = append(blockers, "quality system area "+required+" is missing")
			continue
		}
		if !allowedMaturity[area.Maturity] {
			blockers = append(blockers, "quality system area "+required+" has invalid maturity "+area.Maturity)
		}
		if humanMaturity[required] == "" {
			blockers = append(blockers, "quality system human contract has no "+required+" maturity")
		} else if humanMaturity[required] != area.Maturity {
			blockers = append(blockers, "quality system human and mechanical maturity differ for "+required)
		}
		if area.Policy == "" {
			blockers = append(blockers, "quality system area "+required+" has no policy")
		} else {
			qualityDir := filepath.Join(engineeringDir, "quality")
			policyPath := filepath.Clean(filepath.Join(qualityDir, filepath.FromSlash(area.Policy)))
			relative, relErr := filepath.Rel(qualityDir, policyPath)
			if relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(area.Policy) {
				blockers = append(blockers, "quality system policy "+area.Policy+" escapes engineering/quality")
			} else if _, err := os.Stat(policyPath); err != nil {
				blockers = append(blockers, "quality system policy "+area.Policy+" is missing")
			}
		}
		if required == "security_privacy" && area.DelegatedGate != "security-review" {
			blockers = append(blockers, "quality system security_privacy delegated_gate must be security-review")
		}
		if required != "security_privacy" && area.DelegatedGate != "" {
			blockers = append(blockers, "quality system area "+required+" has unsupported delegated_gate "+area.DelegatedGate)
		}
		if area.Maturity != "baseline" && len(area.RequiredEvidence) == 0 {
			blockers = append(blockers, "quality system area "+required+" maturity "+area.Maturity+" requires evidence")
		}
		for _, evidence := range area.RequiredEvidence {
			if !validEvidenceReference(filepath.Dir(engineeringDir), evidence) {
				blockers = append(blockers, "quality system area "+required+" has invalid or missing evidence "+evidence)
			}
		}
	}
	for name, values := range map[string][]string{"environments": catalog.Environments, "test_data_classes": catalog.TestDataClasses, "platforms": catalog.Platforms} {
		seen := map[string]bool{}
		for _, value := range values {
			value = strings.ToLower(strings.TrimSpace(value))
			if value == "" || seen[value] {
				blockers = append(blockers, "quality system "+name+" contains an empty or duplicate value")
			}
			seen[value] = true
		}
	}
	if !catalog.Exceptions.RequireOwner || !catalog.Exceptions.RequireResidualRisk || !catalog.Exceptions.RequireExpiryOrReview {
		blockers = append(blockers, "quality system exceptions must require owner, residual risk, and expiry or review")
	}
	return blockers
}

func qualityHumanMaturity(text string) map[string]string {
	aliases := map[string]string{
		"behavioral": "behavioral", "accessibility": "accessibility", "security and privacy": "security_privacy",
		"performance and reliability": "performance_reliability", "observability": "observability",
	}
	out := map[string]string{}
	for _, line := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		if len(cells) < 4 {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(cells[0]))
		if key := aliases[name]; key != "" {
			out[key] = strings.ToLower(strings.Trim(strings.TrimSpace(cells[len(cells)-1]), "`"))
		}
	}
	return out
}

func validEvidenceReference(productRoot, evidence string) bool {
	evidence = strings.TrimSpace(evidence)
	if evidence == "" {
		return false
	}
	lower := strings.ToLower(evidence)
	if strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "gate:") || strings.HasPrefix(lower, "ci:") || strings.HasPrefix(lower, "command:") {
		return true
	}
	if filepath.IsAbs(filepath.FromSlash(evidence)) {
		return false
	}
	path := filepath.Clean(filepath.Join(productRoot, filepath.FromSlash(evidence)))
	relative, err := filepath.Rel(productRoot, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func CompositeHash(root string, overrides map[string][]byte) (string, error) {
	dir := filepath.Join(root, "engineering")
	var paths []string
	if err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return "", err
	}
	sort.Strings(paths)
	var content strings.Builder
	for _, path := range paths {
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		data, exists := overrides[rel]
		if !exists {
			var err error
			data, err = os.ReadFile(path)
			if err != nil {
				return "", err
			}
		}
		content.WriteString(rel)
		content.WriteByte('\n')
		content.WriteString(normalize(string(data)))
		content.WriteByte('\n')
	}
	sum := sha256.Sum256([]byte(content.String()))
	return hex.EncodeToString(sum[:]), nil
}

func SynchronizeStatus(contextText string, catalogData []byte, status string) (string, []byte, error) {
	pattern := regexp.MustCompile(`(?m)^(\s*status:\s*)[^\r\n#]+`)
	if !pattern.MatchString(contextText) {
		return "", nil, fmt.Errorf("engineering context has no status field")
	}
	updatedContext := pattern.ReplaceAllString(contextText, "${1}"+status)
	var document yaml.Node
	if err := yaml.Unmarshal(catalogData, &document); err != nil {
		return "", nil, err
	}
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return "", nil, fmt.Errorf("engineering catalog must be a YAML mapping")
	}
	mapping := document.Content[0]
	found := false
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == "status" {
			mapping.Content[index+1].Value = status
			found = true
		}
	}
	if !found {
		return "", nil, fmt.Errorf("engineering catalog has no status field")
	}
	updatedCatalog, err := yaml.Marshal(&document)
	return updatedContext, updatedCatalog, err
}

func SynchronizeQualityStatus(catalogData []byte, status string) ([]byte, error) {
	var document yaml.Node
	if err := yaml.Unmarshal(catalogData, &document); err != nil {
		return nil, err
	}
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("quality system catalog must be a YAML mapping")
	}
	mapping := document.Content[0]
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == "status" {
			mapping.Content[index+1].Value = status
			return yaml.Marshal(&document)
		}
	}
	return nil, fmt.Errorf("quality system catalog has no status field")
}

func Validate(root string) (Inspection, error) { return Inspect(root) }

func Migrate(root string, dryRun bool) ([]string, error) {
	path := filepath.Join(root, "engineering", "engineering-system.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("engineering-system.yaml is invalid YAML: %w", err)
	}
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("engineering-system.yaml must be a YAML mapping")
	}
	mapping := document.Content[0]
	var changes []string
	hasSchema := false
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == "schema_version" {
			hasSchema = true
			if mapping.Content[index+1].Value != "1" {
				return nil, fmt.Errorf("unsupported schema_version %s", mapping.Content[index+1].Value)
			}
		}
	}
	if !hasSchema {
		key := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "schema_version"}
		value := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: "1"}
		mapping.Content = append([]*yaml.Node{key, value}, mapping.Content...)
		changes = append(changes, "ADD engineering/engineering-system.yaml schema_version: 1")
	}
	qualityMigration := false
	if areas := mappingNodeValue(mapping, "areas"); areas != nil {
		if quality := mappingNodeValue(areas, "quality"); quality != nil {
			if contract := mappingNodeValue(quality, "contract"); contract != nil && filepath.ToSlash(contract.Value) == "quality/quality-model.md" {
				contract.Value = "quality/quality-system.md"
				qualityMigration = true
				changes = append(changes,
					"UPDATE engineering/engineering-system.yaml quality contract",
					"ADD engineering/quality/quality-system.md",
					"ADD engineering/quality/quality-system.yaml",
					"ADD engineering/quality/test-strategy.md",
				)
			}
		}
	}
	if len(changes) == 0 {
		return []string{"Engineering System catalog already uses schema_version 1 and requires no quality migration"}, nil
	}
	if dryRun {
		return changes, nil
	}
	updated, err := yaml.Marshal(&document)
	if err != nil {
		return nil, err
	}
	var generated map[string][]byte
	if qualityMigration {
		var catalog catalogDocument
		if err := yaml.Unmarshal(updated, &catalog); err != nil {
			return nil, err
		}
		generated = qualityMigrationFiles(root, catalog)
	}
	if err := applyMigration(path, data, updated, generated); err != nil {
		return nil, err
	}
	return changes, nil
}

func mappingNodeValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			return mapping.Content[index+1]
		}
	}
	return nil
}

func qualityMigrationFiles(root string, catalog catalogDocument) map[string][]byte {
	dir := filepath.Join(root, "engineering", "quality")
	files := map[string]string{
		"quality-system.md":   fmt.Sprintf("# Engineering Quality System\n\n## Snapshot\n\n| Field | Value |\n| --- | --- |\n| Engineering System | `%s @ %s` |\n| Status | `%s` |\n| Mechanical catalog | [quality-system.yaml](quality-system.yaml) |\n| Quality model | [quality-model.md](quality-model.md) |\n| Test strategy | [test-strategy.md](test-strategy.md) |\n\n## Scope\n\nMigrated from the legacy Engineering Quality Model. Inspect real tests, environments, data, platforms, gates, and evidence before advancing maturity.\n\n## Capability Model\n\n| Area | Policy | Evidence | Maturity |\n| --- | --- | --- | --- |\n| Behavioral | [test-strategy.md](test-strategy.md) | Not configured | `baseline` |\n| Accessibility | [test-strategy.md](test-strategy.md) | Not configured | `baseline` |\n| Security and privacy | [test-strategy.md](test-strategy.md) | Not configured | `baseline` |\n| Performance and reliability | [quality-model.md](quality-model.md) | Not configured | `baseline` |\n| Observability | [quality-model.md](quality-model.md) | Not configured | `baseline` |\n\n## Exceptions\n\nNo exceptions were migrated. New exceptions require owner, residual risk, mitigation, and expiry or review date.\n", catalog.ID, catalog.Version, catalog.Status),
		"quality-system.yaml": fmt.Sprintf("schema_version: 1\nengineering_system: %s\nversion: %s\nstatus: %s\nareas:\n  behavioral: {maturity: baseline, policy: test-strategy.md, required_evidence: []}\n  accessibility: {maturity: baseline, policy: test-strategy.md, required_evidence: []}\n  security_privacy: {maturity: baseline, policy: test-strategy.md, delegated_gate: security-review, required_evidence: []}\n  performance_reliability: {maturity: baseline, policy: quality-model.md, required_evidence: []}\n  observability: {maturity: baseline, policy: quality-model.md, required_evidence: []}\ngate_source: knowledge/conventions/gates.md\nenvironments: []\ntest_data_classes: []\nplatforms: []\nexceptions:\n  require_owner: true\n  require_residual_risk: true\n  require_expiry_or_review: true\n  records: []\n", catalog.ID, catalog.Version, catalog.Status),
		"test-strategy.md":    "# Engineering Test Strategy\n\n## Scope\n\nMigrated baseline. Define shared test levels, risk coverage, environments, data, platforms, flaky-test handling, and delivery application from real evidence.\n\n## Delivery Application\n\nEach `tests.md` pins the Engineering System id/version, maps every `AC-*`, and declares deviations or exceptions.\n",
	}
	out := map[string][]byte{}
	for name, body := range files {
		out[filepath.Join(dir, name)] = []byte(body)
	}
	return out
}

func applyMigration(catalogPath string, originalCatalog, updatedCatalog []byte, generated map[string][]byte) error {
	var created []string
	rollback := func() {
		_ = os.WriteFile(catalogPath, originalCatalog, 0o644)
		for _, path := range created {
			_ = os.Remove(path)
		}
	}
	for path, body := range generated {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			rollback()
			return fmt.Errorf("migration target is a directory: %s", path)
		} else if err == nil {
			continue
		} else if !os.IsNotExist(err) {
			rollback()
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			rollback()
			return err
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			rollback()
			return err
		}
		_, writeErr := file.Write(body)
		if writeErr == nil {
			writeErr = file.Sync()
		}
		closeErr := file.Close()
		created = append(created, path)
		if writeErr != nil {
			rollback()
			return writeErr
		}
		if closeErr != nil {
			rollback()
			return closeErr
		}
	}
	tmp := catalogPath + ".quality-migration.tmp"
	if err := os.WriteFile(tmp, updatedCatalog, 0o644); err != nil {
		rollback()
		return err
	}
	if err := os.Rename(tmp, catalogPath); err != nil {
		_ = os.Remove(tmp)
		rollback()
		return err
	}
	return nil
}

func Triggers(text string) (valid, invalid []string) {
	var context contextDocument
	if err := yaml.Unmarshal([]byte(yamlPayload(text)), &context); err != nil {
		return nil, []string{"invalid_yaml"}
	}
	seen := map[string]bool{}
	for _, value := range context.EngineeringTriggers {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		if allowedTriggers[value] {
			valid = append(valid, value)
		} else {
			invalid = append(invalid, value)
		}
	}
	sort.Strings(valid)
	sort.Strings(invalid)
	return valid, invalid
}

func AllowedTriggers() []string {
	var out []string
	for trigger := range allowedTriggers {
		out = append(out, trigger)
	}
	sort.Strings(out)
	return out
}

func yamlPayload(text string) string {
	text = strings.TrimPrefix(text, "\ufeff")
	trimmed := strings.TrimSpace(text)
	if strings.HasPrefix(trimmed, "---") {
		lines := strings.Split(trimmed, "\n")
		for index := 1; index < len(lines); index++ {
			if strings.TrimSpace(lines[index]) == "---" {
				return strings.Join(lines[1:index], "\n")
			}
		}
	}
	lower := strings.ToLower(text)
	if start := strings.Index(lower, "```yaml"); start >= 0 {
		body := text[start+len("```yaml"):]
		if end := strings.Index(body, "```"); end >= 0 {
			return body[:end]
		}
	}
	return text
}

func markdownSnapshot(text string) map[string]string {
	out := map[string]string{}
	pattern := regexp.MustCompile(`(?m)^\|\s*([^|]+?)\s*\|\s*` + "`?" + `([^|` + "`" + `]+?)` + "`?" + `\s*\|$`)
	for _, match := range pattern.FindAllStringSubmatch(text, -1) {
		key := strings.ToLower(strings.TrimSpace(match[1]))
		if key != "field" {
			out[key] = strings.TrimSpace(match[2])
		}
	}
	return out
}

func normalize(text string) string {
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

func oneOf(value string, options ...string) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

func unique(items []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, item := range items {
		if item != "" && !seen[item] {
			seen[item] = true
			out = append(out, item)
		}
	}
	return out
}
