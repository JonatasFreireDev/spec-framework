package design

type Version struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Screen struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Asset     string   `json:"asset,omitempty"`
	Viewports []string `json:"viewports,omitempty"`
	States    []string `json:"states,omitempty"`
}

type SourceManifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	ID            string            `json:"id"`
	Type          string            `json:"type"`
	Authority     string            `json:"authority"`
	Location      string            `json:"location"`
	Version       Version           `json:"version"`
	Adapter       string            `json:"adapter,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Screens       []Screen          `json:"screens"`
}

type Mapping struct {
	Requirement string `json:"requirement"`
	Criterion   string `json:"criterion,omitempty"`
	Screen      string `json:"screen,omitempty"`
	State       string `json:"state,omitempty"`
	Coverage    string `json:"coverage"`
	Notes       string `json:"notes,omitempty"`
}

type UseCaseManifest struct {
	SchemaVersion  int       `json:"schemaVersion"`
	UseCase        string    `json:"useCase"`
	OriginMode     string    `json:"originMode"`
	Maturity       string    `json:"maturity"`
	FidelityPolicy string    `json:"fidelityPolicy"`
	Sources        []string  `json:"sources"`
	Mappings       []Mapping `json:"mappings,omitempty"`
	NonProduction  bool      `json:"nonProduction,omitempty"`
}

type Inspection struct {
	UseCase        string   `json:"useCase"`
	OriginMode     string   `json:"originMode"`
	Maturity       string   `json:"maturity"`
	FidelityPolicy string   `json:"fidelityPolicy"`
	Sources        []string `json:"sources"`
	Screens        int      `json:"screens"`
	Mappings       int      `json:"mappings"`
	Blockers       []string `json:"blockers,omitempty"`
}
