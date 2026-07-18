package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JonatasFreireDev/spec-framework/internal/dispatch"
	"github.com/JonatasFreireDev/spec-framework/internal/sourceimport"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func mustRecommendations(items []dispatch.Recommendation, err error) []dispatch.Recommendation {
	if err != nil {
		return []dispatch.Recommendation{{Kind: "unavailable", Detail: err.Error(), RequiresConfirmation: true}}
	}
	return items
}

func runDispatch(args []string, out, errout io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(errout, "dispatch requires plan, assign, return, observe, or reconcile")
		return 2
	}
	fs := flag.NewFlagSet("dispatch", flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	work := fs.String("work", "", "workspace")
	graph := fs.String("graph", "", "execution graph")
	task := fs.String("task", "", "task")
	agent := fs.String("agent", "", "agent")
	role := fs.String("role", "code-runner", "role")
	id := fs.String("id", "", "dispatch id")
	summary := fs.String("summary", "", "return summary")
	evidence := fs.String("evidence", "", "comma-separated evidence")
	dependencies := fs.String("depends-on", "", "comma-separated dispatch dependencies")
	outputHashes := fs.String("output-hashes", "", "comma-separated product-relative path=sha256 outputs")
	blockers := fs.String("blockers", "", "comma-separated delegated blockers")
	decisionCandidates := fs.String("decision-candidates", "", "comma-separated delegated decision candidates")
	command := fs.String("command", "", "supervised executable")
	enable := fs.Bool("enable", false, "explicitly enable supervised execution")
	wave := fs.String("wave", "", "persisted scheduler wave id")
	max := fs.Int("max-parallel", 1, "maximum concurrent dispatches")
	diffHash := fs.String("diff-hash", "", "immutable working-tree diff hash")
	parent := fs.String("parent", "", "returned code-runner dispatch id")
	runID := fs.String("run", "", "scalable import run")
	chunk := fs.String("chunk", "", "scalable import chunk")
	reviewInput := fs.String("review-input", "", "structured import review JSON")
	harnesses := fs.String("harnesses", "", "comma-separated allowed harness basenames")
	enabled := fs.Bool("enabled", false, "enable dispatch capability")
	retention := fs.Int("transcript-retention", 100, "transcripts to retain per workspace")
	yes := fs.Bool("yes", false, "confirm mutation")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	wd, _ := os.Getwd()
	p := *root
	if !filepath.IsAbs(p) {
		p = filepath.Join(wd, p)
	}
	g := *graph
	if g != "" && !filepath.IsAbs(g) {
		g = filepath.Join(p, filepath.FromSlash(g))
	}
	switch args[0] {
	case "configure":
		if !*yes {
			fmt.Fprintln(errout, "dispatch configure requires --yes")
			return 2
		}
		if err := dispatch.SaveConfig(p, dispatch.Config{Version: 1, Enabled: *enabled, Harnesses: splitCSV(*harnesses), MaxParallel: *max, TranscriptRetention: *retention}); err != nil {
			fmt.Fprintln(errout, err)
			return 1
		}
		fmt.Fprintln(out, "DISPATCH CONFIGURED")
		return 0
	case "plan":
		items, e := dispatch.Plan(p, g)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		for _, x := range items {
			fmt.Fprintf(out, "%s %s ready=%t %s\n", x.TaskID, x.Role, x.Ready, strings.Join(x.Blockers, "; "))
		}
		return 0
	case "assign":
		if isEngineeringDispatchRole(*role) {
			if !*yes || *work == "" || *task == "" || *agent == "" {
				fmt.Fprintln(errout, "engineering assignment requires --work --task <handoff-path> --role --agent --yes")
				return 2
			}
			path := *task
			if !filepath.IsAbs(path) {
				path = filepath.Join(p, filepath.FromSlash(path))
			}
			x, e := dispatch.AssignEngineering(p, *work, path, *role, *agent, splitCSV(*dependencies))
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "ASSIGNED", x.ID)
			return 0
		}
		if *role == "technical-discovery" {
			if !*yes || *work == "" || *task == "" || *agent == "" {
				fmt.Fprintln(errout, "technical-discovery assignment requires --work --task <question-path> --agent --yes")
				return 2
			}
			path := *task
			if !filepath.IsAbs(path) {
				path = filepath.Join(p, filepath.FromSlash(path))
			}
			x, e := dispatch.AssignResearch(p, *work, path, *agent)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "ASSIGNED", x.ID)
			return 0
		}
		if *role == "threat-modeler" {
			if !*yes || *work == "" || *task == "" || *agent == "" {
				fmt.Fprintln(errout, "threat-modeler assignment requires --work --task <boundary-path> --agent --yes")
				return 2
			}
			path := *task
			if !filepath.IsAbs(path) {
				path = filepath.Join(p, filepath.FromSlash(path))
			}
			x, e := dispatch.AssignThreatModel(p, *work, path, *agent)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "ASSIGNED", x.ID)
			return 0
		}
		if *role == "artifact-importer" {
			if !*yes || *work == "" || *runID == "" || *agent == "" {
				fmt.Fprintln(errout, "artifact-importer assignment requires --work --run --agent --yes")
				return 2
			}
			x, e := dispatch.AssignImportChunk(p, *work, *runID, *chunk, *agent)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "ASSIGNED", x.ID)
			return 0
		}
		if !*yes || *work == "" || *agent == "" || ((*role != "qa" && *role != "code-review" && *role != "security-review") && *task == "") || ((*role == "qa" || *role == "code-review" || *role == "security-review") && *parent == "") {
			fmt.Fprintln(errout, "dispatch assign requires --work --task --agent --yes")
			return 2
		}
		if *role == "qa" || *role == "code-review" || *role == "security-review" {
			x, e := dispatch.AssignReview(p, *work, *parent, *role, *agent)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "ASSIGNED", x.ID)
			return 0
		}
		x, e := dispatch.Assign(p, *work, g, *task, *role, *agent)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintln(out, "ASSIGNED", x.ID)
		return 0
	case "return":
		if !*yes || *work == "" || *id == "" || *agent == "" {
			fmt.Fprintln(errout, "dispatch return requires --work --id --agent --yes")
			return 2
		}
		if len(splitCSV(*outputHashes)) > 0 {
			x, e := dispatch.ReturnEngineering(p, *work, *id, *agent, *summary, splitCSV(*evidence), splitCSV(*outputHashes), splitCSV(*blockers), splitCSV(*decisionCandidates))
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "RETURNED", x.ID)
			return 0
		}
		if *reviewInput != "" {
			data, e := os.ReadFile(*reviewInput)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			var review sourceimport.ChunkReview
			if e = json.Unmarshal(data, &review); e != nil {
				fmt.Fprintln(errout, e)
				return 2
			}
			x, e := dispatch.ReturnImport(p, *work, *id, *agent, *summary, review)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "RETURNED", x.ID)
			return 0
		}
		x, e := dispatch.Return(p, *work, *id, *agent, *summary, *diffHash, splitCSV(*evidence))
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintln(out, "RETURNED", x.ID)
		return 0
	case "observe":
		if *work == "" {
			fmt.Fprintln(errout, "dispatch observe requires --work")
			return 2
		}
		xs, e := dispatch.Observe(p, *work)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		for _, x := range xs {
			fmt.Fprintf(out, "%s %s %s %s\n", x.ID, x.TaskID, x.Role, x.Status)
		}
		return 0
	case "reconcile":
		if *work == "" {
			fmt.Fprintln(errout, "dispatch reconcile requires --work")
			return 2
		}
		xs, e := dispatch.Reconcile(p, *work)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		for _, x := range xs {
			fmt.Fprintf(out, "%s %s -> %s\n", x.Kind, x.DispatchID, x.Owner)
		}
		return 0
	case "run":
		if !*yes || !*enable || *work == "" || *id == "" || *command == "" {
			fmt.Fprintln(errout, "dispatch run requires --work --id --command --enable --yes")
			return 2
		}
		t, e := dispatch.Run(p, *work, *id, *enable, *command, fs.Args())
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintln(out, "RAN", t.DispatchID, t.OutputHash)
		return 0
	case "wave":
		if !*yes || !*enable || *work == "" || *wave == "" || *command == "" {
			fmt.Fprintln(errout, "dispatch wave requires --work --wave --command --enable --yes")
			return 2
		}
		waveIDs, e := dispatch.WaveIDs(p, *work, *wave)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		for _, r := range dispatch.RunWave(p, *work, waveIDs, *max, *enable, *command, fs.Args()) {
			if r.Error != "" {
				fmt.Fprintln(errout, r.ID, r.Error)
			} else {
				fmt.Fprintln(out, "RAN", r.ID, r.Transcript.OutputHash)
			}
		}
		return 0
	case "recommend":
		if *work == "" {
			fmt.Fprintln(errout, "dispatch recommend requires --work")
			return 2
		}
		for _, r := range mustRecommendations(dispatch.Recommend(p, *work, *max)) {
			fmt.Fprintln(out, r.Kind, r.Detail)
		}
		return 0
	}
	fmt.Fprintln(errout, "unknown dispatch operation")
	return 2
}

func isEngineeringDispatchRole(role string) bool {
	switch role {
	case "technical-landscape", "engineering-standards", "operations-baseline", "engineering-evidence", "engineering-system":
		return true
	default:
		return false
	}
}
