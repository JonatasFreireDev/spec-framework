package install

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	framework "github.com/JonatasFreireDev/spec-framework"
)

const initContractSchemaVersion = 1

type initAssetCatalog struct {
	SchemaVersion int                        `json:"schema_version"`
	Sets          map[string][]initAssetSpec `json:"sets"`
	Modules       map[string]initModuleSpec  `json:"modules"`
	Adapters      []string                   `json:"approval_adapters"`
}

type initModuleSpec struct {
	AssetSets     []string `json:"asset_sets"`
	ArtifactTypes []string `json:"artifact_types"`
}

type initContract struct {
	SchemaVersion    int              `json:"schema_version"`
	ID               string           `json:"id"`
	AssetSets        []string         `json:"asset_sets"`
	Directories      []string         `json:"directories"`
	Files            []initFileSpec   `json:"files"`
	Patches          []initPatchSpec  `json:"patches"`
	Registry         initRegistrySpec `json:"registry"`
	BootstrapProfile string           `json:"bootstrap_profile"`
	Actions          []string         `json:"actions"`
}

type initAssetSpec struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type initFileSpec struct {
	Source       string            `json:"source"`
	Target       string            `json:"target"`
	Replacements map[string]string `json:"replacements"`
}

type initPatchSpec struct {
	Target  string `json:"target"`
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

type initRegistrySpec struct {
	ExcludeTypes         []string              `json:"exclude_types"`
	PrependArtifacts     []map[string]any      `json:"prepend_artifacts"`
	AppendParentsByType  map[string][]string   `json:"append_parents_by_type"`
	ReplaceParentsByType map[string][]string   `json:"replace_parents_by_type"`
	Modules              []string              `json:"modules"`
	RequiredArtifacts    []artifactRequirement `json:"required_artifacts"`
	ApprovalAdapters     map[string]string     `json:"approval_adapters"`
	AllowedArtifactTypes []string              `json:"-"`
	AllowedAdapters      map[string]bool       `json:"-"`
}

type artifactRequirement struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type initContractLoader struct {
	assets fs.FS
}

type initializationPlan struct {
	Contract    initContract
	Directories []string
	Files       map[string]plannedFile
	Actions     []string
}

type plannedFile struct {
	Source string
	Data   []byte
}

func defaultInitContractLoader() initContractLoader {
	return initContractLoader{assets: framework.Assets}
}

func loadInitContract(point string) (initContract, initAssetCatalog, error) {
	return defaultInitContractLoader().load(point)
}

func (loader initContractLoader) load(point string) (initContract, initAssetCatalog, error) {
	var catalog initAssetCatalog
	if err := loader.readStrictJSON("framework/init/catalog.json", &catalog); err != nil {
		return initContract{}, initAssetCatalog{}, err
	}
	if catalog.SchemaVersion != initContractSchemaVersion {
		return initContract{}, initAssetCatalog{}, fmt.Errorf("unsupported init asset catalog schema %d", catalog.SchemaVersion)
	}
	if len(catalog.Sets) == 0 {
		return initContract{}, initAssetCatalog{}, errors.New("init asset catalog has no sets")
	}

	var contract initContract
	name := "framework/init/contracts/" + point + ".json"
	if err := loader.readStrictJSON(name, &contract); err != nil {
		return initContract{}, initAssetCatalog{}, err
	}
	if contract.SchemaVersion != initContractSchemaVersion {
		return initContract{}, initAssetCatalog{}, fmt.Errorf("unsupported init contract schema %d in %s", contract.SchemaVersion, name)
	}
	if contract.ID == "" || contract.ID != point {
		return initContract{}, initAssetCatalog{}, fmt.Errorf("init contract id %q does not match starting point %q", contract.ID, point)
	}
	if contract.BootstrapProfile != point {
		return initContract{}, initAssetCatalog{}, fmt.Errorf("init contract %q selects bootstrap profile %q", point, contract.BootstrapProfile)
	}
	selectedAssets := map[string]bool{}
	allowedTypes := map[string]bool{}
	allowedAdapters := map[string]bool{}
	for _, adapter := range catalog.Adapters {
		allowedAdapters[adapter] = true
	}
	for _, module := range contract.Registry.Modules {
		spec, ok := catalog.Modules[module]
		if !ok {
			return initContract{}, initAssetCatalog{}, fmt.Errorf("unknown artifact module %q", module)
		}
		for _, assetSet := range spec.AssetSets {
			selectedAssets[assetSet] = true
		}
		for _, kind := range spec.ArtifactTypes {
			allowedTypes[kind] = true
		}
	}
	if len(contract.Registry.Modules) > 0 {
		for _, assetSet := range contract.AssetSets {
			if !selectedAssets[assetSet] {
				return initContract{}, initAssetCatalog{}, fmt.Errorf("asset set %q is not provided by selected modules", assetSet)
			}
		}
		for assetSet := range selectedAssets {
			found := false
			for _, selected := range contract.AssetSets {
				if selected == assetSet {
					found = true
					break
				}
			}
			if !found {
				return initContract{}, initAssetCatalog{}, fmt.Errorf("selected module asset set %q is missing from contract", assetSet)
			}
		}
	}
	for kind, adapter := range contract.Registry.ApprovalAdapters {
		if !allowedAdapters[adapter] {
			return initContract{}, initAssetCatalog{}, fmt.Errorf("unknown approval adapter %q for artifact type %q", adapter, kind)
		}
	}
	contract.Registry.AllowedArtifactTypes = make([]string, 0, len(allowedTypes))
	for kind := range allowedTypes {
		contract.Registry.AllowedArtifactTypes = append(contract.Registry.AllowedArtifactTypes, kind)
	}
	contract.Registry.AllowedAdapters = allowedAdapters
	if err := loader.validateContract(contract, catalog); err != nil {
		return initContract{}, initAssetCatalog{}, fmt.Errorf("invalid init contract %q: %w", point, err)
	}
	return contract, catalog, nil
}

func (loader initContractLoader) readStrictJSON(name string, target any) error {
	data, err := fs.ReadFile(loader.assets, name)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", name, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("parse embedded %s: %w", name, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("parse embedded %s: trailing JSON value", name)
		}
		return fmt.Errorf("parse embedded %s: %w", name, err)
	}
	return nil
}

func (loader initContractLoader) validateContract(contract initContract, catalog initAssetCatalog) error {
	if len(contract.AssetSets) == 0 {
		return errors.New("contract has no asset sets")
	}
	for _, module := range contract.Registry.Modules {
		if _, ok := catalog.Modules[module]; !ok {
			return fmt.Errorf("unknown artifact module %q", module)
		}
	}
	seenSets := map[string]bool{}
	for _, name := range contract.AssetSets {
		if seenSets[name] {
			return fmt.Errorf("duplicate asset set %q", name)
		}
		seenSets[name] = true
		assets, ok := catalog.Sets[name]
		if !ok {
			return fmt.Errorf("unknown asset set %q", name)
		}
		if len(assets) == 0 {
			return fmt.Errorf("asset set %q is empty", name)
		}
		for _, asset := range assets {
			if err := loader.validateAssetPath(asset.Source, asset.Target); err != nil {
				return fmt.Errorf("asset set %q: %w", name, err)
			}
		}
	}
	for _, file := range contract.Files {
		if err := loader.validateAssetPath(file.Source, file.Target); err != nil {
			return err
		}
		for find := range file.Replacements {
			if find == "" {
				return fmt.Errorf("init file %q has an empty replacement anchor", file.Target)
			}
		}
	}
	seenDirectories := map[string]bool{}
	for _, directory := range contract.Directories {
		if directory == "" || directory == "." {
			return errors.New("explicit product directory must not be empty or root")
		}
		if _, err := safeProductPath("product", directory); err != nil {
			return err
		}
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(directory)))
		if seenDirectories[clean] {
			return fmt.Errorf("duplicate product directory %q", clean)
		}
		seenDirectories[clean] = true
	}
	for _, patch := range contract.Patches {
		if _, err := safeProductPath("product", patch.Target); err != nil {
			return err
		}
		if patch.Find == "" {
			return fmt.Errorf("patch target %q has an empty find value", patch.Target)
		}
	}
	knownActions := map[string]bool{"create-import-run": true}
	seenActions := map[string]bool{}
	for _, action := range contract.Actions {
		if !knownActions[action] {
			return fmt.Errorf("unknown init action %q", action)
		}
		if seenActions[action] {
			return fmt.Errorf("duplicate init action %q", action)
		}
		seenActions[action] = true
	}
	return nil
}

func (loader initContractLoader) validateAssetPath(source, target string) error {
	if source == "" || !fs.ValidPath(source) {
		return fmt.Errorf("unsafe embedded asset source %q", source)
	}
	if _, err := fs.Stat(loader.assets, source); err != nil {
		return fmt.Errorf("embedded asset source %q does not exist: %w", source, err)
	}
	if _, err := safeProductPath("product", target); err != nil {
		return err
	}
	return nil
}

func safeProductPath(root, relative string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == "." {
		return root, nil
	}
	if filepath.IsAbs(clean) || filepath.VolumeName(clean) != "" || strings.HasPrefix(clean, string(filepath.Separator)) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("unsafe product target %q", relative)
	}
	return filepath.Join(root, clean), nil
}

func buildInitializationPlan(point string) (initializationPlan, error) {
	return defaultInitContractLoader().buildPlan(point)
}

func (loader initContractLoader) buildPlan(point string) (initializationPlan, error) {
	contract, catalog, err := loader.load(point)
	if err != nil {
		return initializationPlan{}, err
	}
	plan := initializationPlan{Contract: contract, Directories: append([]string{}, contract.Directories...), Files: map[string]plannedFile{}, Actions: append([]string{}, contract.Actions...)}
	for _, setName := range contract.AssetSets {
		for _, asset := range catalog.Sets[setName] {
			if err := loader.expandAsset(&plan, setName, asset); err != nil {
				return initializationPlan{}, err
			}
		}
	}
	for _, file := range contract.Files {
		data, err := fs.ReadFile(loader.assets, file.Source)
		if err != nil {
			return initializationPlan{}, err
		}
		text := string(data)
		keys := make([]string, 0, len(file.Replacements))
		for find := range file.Replacements {
			keys = append(keys, find)
		}
		sort.Strings(keys)
		for _, find := range keys {
			count := strings.Count(text, find)
			if count != 1 {
				return initializationPlan{}, fmt.Errorf("init file replacement in %s expected one match for %q, found %d", file.Source, find, count)
			}
			text = strings.Replace(text, find, file.Replacements[find], 1)
		}
		if _, exists := plan.Files[file.Target]; exists {
			return initializationPlan{}, fmt.Errorf("entry file target %q collides with a selected asset", file.Target)
		}
		plan.Files[file.Target] = plannedFile{Source: file.Source, Data: []byte(text)}
	}
	for _, patch := range contract.Patches {
		file, ok := plan.Files[patch.Target]
		if !ok {
			return initializationPlan{}, fmt.Errorf("init patch target is not planned: %s", patch.Target)
		}
		count := strings.Count(string(file.Data), patch.Find)
		if count != 1 {
			return initializationPlan{}, fmt.Errorf("init patch %s expected one match, found %d", patch.Target, count)
		}
		file.Data = []byte(strings.Replace(string(file.Data), patch.Find, patch.Replace, 1))
		plan.Files[patch.Target] = file
	}
	if err := configurePlannedRegistry(&plan, contract.Registry); err != nil {
		return initializationPlan{}, err
	}
	if err := validatePlannedDirectories(plan.Directories, plan.Files); err != nil {
		return initializationPlan{}, err
	}
	for i, directory := range plan.Directories {
		plan.Directories[i] = filepath.ToSlash(filepath.Clean(filepath.FromSlash(directory)))
	}
	sort.Slice(plan.Directories, func(i, j int) bool {
		leftDepth := strings.Count(plan.Directories[i], "/")
		rightDepth := strings.Count(plan.Directories[j], "/")
		if leftDepth == rightDepth {
			return plan.Directories[i] < plan.Directories[j]
		}
		return leftDepth < rightDepth
	})
	return plan, nil
}

func validatePlannedDirectories(directories []string, files map[string]plannedFile) error {
	for _, directory := range directories {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(directory)))
		if _, collision := files[clean]; collision {
			return fmt.Errorf("planned directory %q collides with a file", clean)
		}
		for parent := clean; parent != "." && parent != ""; parent = filepath.ToSlash(filepath.Dir(filepath.FromSlash(parent))) {
			if _, collision := files[parent]; collision {
				return fmt.Errorf("planned directory %q is nested below file %q", clean, parent)
			}
		}
	}
	return nil
}

func (loader initContractLoader) expandAsset(plan *initializationPlan, setName string, asset initAssetSpec) error {
	info, err := fs.Stat(loader.assets, asset.Source)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		data, err := fs.ReadFile(loader.assets, asset.Source)
		if err != nil {
			return err
		}
		return addPlannedAsset(plan, asset.Target, asset.Source, setName, data)
	}
	return fs.WalkDir(loader.assets, asset.Source, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(filepath.FromSlash(asset.Source), filepath.FromSlash(name))
		if err != nil {
			return err
		}
		target := filepath.ToSlash(filepath.Join(filepath.FromSlash(asset.Target), rel))
		data, err := fs.ReadFile(loader.assets, name)
		if err != nil {
			return err
		}
		return addPlannedAsset(plan, target, name, setName, data)
	})
}

func addPlannedAsset(plan *initializationPlan, target, source, setName string, data []byte) error {
	if _, err := safeProductPath("product", target); err != nil {
		return err
	}
	if existing, ok := plan.Files[target]; ok {
		return fmt.Errorf("asset target collision %q between %s and %s in set %q", target, existing.Source, source, setName)
	}
	plan.Files[target] = plannedFile{Source: source, Data: append([]byte{}, data...)}
	return nil
}

func configurePlannedRegistry(plan *initializationPlan, spec initRegistrySpec) error {
	const registryTarget = ".product/artifacts.json"
	file, ok := plan.Files[registryTarget]
	if !ok {
		return fmt.Errorf("planned product is missing %s", registryTarget)
	}
	var registry struct {
		Artifacts         []map[string]any      `json:"artifacts"`
		Modules           []string              `json:"modules,omitempty"`
		RequiredArtifacts []artifactRequirement `json:"required_artifacts,omitempty"`
		ApprovalAdapters  map[string]string     `json:"approval_adapters,omitempty"`
	}
	if err := json.Unmarshal(file.Data, &registry); err != nil {
		return fmt.Errorf("parse planned registry: %w", err)
	}
	excluded := stringSet(spec.ExcludeTypes)
	foundTypes := map[string]bool{}
	filtered := make([]map[string]any, 0, len(registry.Artifacts)+len(spec.PrependArtifacts))
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		foundTypes[kind] = true
		if excluded[kind] {
			continue
		}
		if parents, ok := spec.ReplaceParentsByType[kind]; ok {
			artifact["parentIds"] = parents
		}
		if additions := spec.AppendParentsByType[kind]; len(additions) > 0 {
			artifact["parentIds"] = append(stringSlice(artifact["parentIds"]), additions...)
		}
		filtered = append(filtered, artifact)
	}
	for _, kind := range spec.ExcludeTypes {
		if !foundTypes[kind] {
			return fmt.Errorf("registry exclusion type %q matched no artifact", kind)
		}
	}
	for kind := range spec.ReplaceParentsByType {
		if !foundTypes[kind] || excluded[kind] {
			return fmt.Errorf("registry parent replacement type %q matched no active artifact", kind)
		}
	}
	for kind := range spec.AppendParentsByType {
		if !foundTypes[kind] || excluded[kind] {
			return fmt.Errorf("registry parent append type %q matched no active artifact", kind)
		}
	}
	registry.Artifacts = append(append([]map[string]any{}, spec.PrependArtifacts...), filtered...)
	registry.Modules = append([]string{}, spec.Modules...)
	registry.RequiredArtifacts = append([]artifactRequirement{}, spec.RequiredArtifacts...)
	registry.ApprovalAdapters = map[string]string{}
	for kind, adapter := range spec.ApprovalAdapters {
		registry.ApprovalAdapters[kind] = adapter
	}
	for _, artifact := range registry.Artifacts {
		kind, _ := artifact["type"].(string)
		if adapter := registry.ApprovalAdapters[kind]; adapter != "" {
			if _, exists := artifact["approval_adapter"]; !exists {
				artifact["approval_adapter"] = adapter
			}
		}
		if adapter, exists := artifact["approval_adapter"].(string); exists && len(spec.AllowedAdapters) > 0 && !spec.AllowedAdapters[adapter] {
			return fmt.Errorf("unknown approval adapter %q on artifact type %q", adapter, kind)
		}
	}
	allowedTypes := stringSet(spec.AllowedArtifactTypes)
	if len(allowedTypes) > 0 {
		for _, artifact := range registry.Artifacts {
			kind, _ := artifact["type"].(string)
			if !allowedTypes[kind] {
				return fmt.Errorf("artifact type %q is not provided by selected modules", kind)
			}
		}
	}
	if err := validatePlannedRegistry(registry.Artifacts, plan.Files); err != nil {
		return err
	}
	for _, required := range registry.RequiredArtifacts {
		found := false
		for _, artifact := range registry.Artifacts {
			kind, _ := artifact["type"].(string)
			path, _ := artifact["path"].(string)
			if kind == required.Type && filepath.ToSlash(path) == filepath.ToSlash(required.Path) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("required artifact %s at %s is not registered", required.Type, required.Path)
		}
	}
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return err
	}
	plan.Files[registryTarget] = plannedFile{Source: "generated:init-registry", Data: append(data, '\n')}
	return nil
}

func validatePlannedRegistry(artifacts []map[string]any, files map[string]plannedFile) error {
	ids := map[string]bool{}
	kinds := map[string]string{}
	paths := map[string]bool{}
	parentsByID := map[string][]string{}
	for index, artifact := range artifacts {
		id, _ := artifact["id"].(string)
		kind, _ := artifact["type"].(string)
		status, _ := artifact["status"].(string)
		path, _ := artifact["path"].(string)
		if id == "" || kind == "" || status == "" || path == "" {
			return fmt.Errorf("registry artifact %d requires id, type, status, and path", index)
		}
		if ids[id] {
			return fmt.Errorf("duplicate registry artifact id %q", id)
		}
		if paths[path] {
			return fmt.Errorf("duplicate registry artifact path %q", path)
		}
		if strings.TrimSpace(kind) == "" {
			return fmt.Errorf("registry artifact %q has empty type", id)
		}
		if status != "draft" {
			return fmt.Errorf("initial artifact %q has invalid status %q", id, status)
		}
		if _, ok := files[path]; !ok {
			return fmt.Errorf("registry artifact %q points to unplanned path %q", id, path)
		}
		parents, err := registryStringSlice(artifact, "parentIds")
		if err != nil {
			return fmt.Errorf("registry artifact %q: %w", id, err)
		}
		seenParents := map[string]bool{}
		for _, parent := range parents {
			if seenParents[parent] {
				return fmt.Errorf("registry artifact %q has duplicate parent %q", id, parent)
			}
			seenParents[parent] = true
		}
		if documents, exists := artifact["documents"]; exists {
			values, ok := documents.(map[string]any)
			if !ok {
				return fmt.Errorf("registry artifact %q documents must be an object", id)
			}
			for name, value := range values {
				documentPath, ok := value.(string)
				if !ok || documentPath == "" {
					return fmt.Errorf("registry artifact %q document %q must be a path", id, name)
				}
				if _, ok := files[documentPath]; !ok {
					return fmt.Errorf("registry artifact %q document %q points to unplanned path %q", id, name, documentPath)
				}
			}
		}
		ids[id], paths[path] = true, true
		kinds[id] = kind
		parentsByID[id] = parents
	}
	for _, artifact := range artifacts {
		id, _ := artifact["id"].(string)
		for _, parent := range parentsByID[id] {
			if !ids[parent] {
				return fmt.Errorf("registry artifact %q has unknown parent %q", id, parent)
			}
		}
		if rawTarget, exists := artifact["targetFeature"]; exists {
			target, ok := rawTarget.(string)
			if !ok || target == "" {
				return fmt.Errorf("registry artifact %q targetFeature must be a non-empty id", id)
			}
			if !ids[target] {
				return fmt.Errorf("registry artifact %q has unknown targetFeature %q", id, target)
			}
			if kinds[target] != "feature" {
				return fmt.Errorf("registry artifact %q targetFeature %q is not a feature", id, target)
			}
		}
	}
	return validateRegistryParentCycles(parentsByID)
}

func registryStringSlice(artifact map[string]any, field string) ([]string, error) {
	value, exists := artifact[field]
	if !exists {
		return nil, nil
	}
	raw, ok := value.([]any)
	if !ok {
		if typed, typedOK := value.([]string); typedOK {
			return append([]string{}, typed...), nil
		}
		return nil, fmt.Errorf("%s must be an array of strings", field)
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		text, ok := item.(string)
		if !ok || text == "" {
			return nil, fmt.Errorf("%s must contain non-empty strings", field)
		}
		out = append(out, text)
	}
	return out, nil
}

func validateRegistryParentCycles(parentsByID map[string][]string) error {
	state := map[string]uint8{}
	var visit func(string) error
	visit = func(id string) error {
		if state[id] == 1 {
			return fmt.Errorf("registry parent cycle includes %q", id)
		}
		if state[id] == 2 {
			return nil
		}
		state[id] = 1
		for _, parent := range parentsByID[id] {
			if _, known := parentsByID[parent]; known {
				if err := visit(parent); err != nil {
					return err
				}
			}
		}
		state[id] = 2
		return nil
	}
	for id := range parentsByID {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}

func writeInitializationPlan(productRoot string, plan initializationPlan) error {
	for _, directory := range plan.Directories {
		path, err := safeProductPath(productRoot, directory)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	targets := make([]string, 0, len(plan.Files))
	for target := range plan.Files {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	for _, target := range targets {
		path, err := safeProductPath(productRoot, target)
		if err != nil {
			return err
		}
		if err := writeFile(path, plan.Files[target].Data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func stringSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		out[value] = true
	}
	return out
}

func stringSlice(value any) []string {
	raw, _ := value.([]any)
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		if text, ok := item.(string); ok {
			out = append(out, text)
		}
	}
	if typed, ok := value.([]string); ok {
		return append([]string{}, typed...)
	}
	return out
}

func commitStagedProduct(stagedProduct, productRoot string) error {
	if _, err := os.Stat(productRoot); err == nil {
		return fmt.Errorf("target already contains product/: %s", filepath.Dir(productRoot))
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.Rename(stagedProduct, productRoot)
}
