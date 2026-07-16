package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/reviewfinding"
	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func runRuntime(command string, args []string, out, errout io.Writer) int {
	op := ""
	if (command == "lease" || command == "commands" || command == "integrate" || command == "runtime" || command == "reviews") && len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		op, args = args[0], args[1:]
	}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(errout)
	root := fs.String("product-root", "product", "product root")
	work := fs.String("work", "", "workspace id")
	task := fs.String("task", "", "task id")
	agent := fs.String("agent", "", "agent id")
	graph := fs.String("graph", "", "graph path")
	from := fs.String("from", "", "handoff source")
	to := fs.String("to", "", "handoff target")
	summary := fs.String("summary", "", "summary")
	step := fs.String("step", "", "checkpoint step")
	input := fs.String("input-hash", "", "input hash")
	output := fs.String("output-hash", "", "output hash")
	base := fs.String("base-commit", "", "base commit")
	risk := fs.String("risk", "R0", "R0 or R1")
	source := fs.String("source", "human", "command source")
	cwd := fs.String("cwd", ".", "repository-relative cwd")
	timeout := fs.Int("timeout", 300, "timeout seconds")
	max := fs.Int("max-parallel", 4, "maximum parallel tasks")
	interval := fs.Duration("interval", 2*time.Second, "watch polling interval")
	count := fs.Int("count", 0, "watch observations before stopping (0 keeps watching)")
	dry := fs.Bool("dry-run", false, "preview migration")
	jsonOutput := fs.Bool("json", false, "print structured output")
	yes := fs.Bool("yes", false, "confirm mutation")
	commits := fs.String("commits", "", "comma-separated commits")
	isolate := fs.Bool("isolate", false, "create a task Git worktree")
	importFile := fs.String("input", "", "JSON array of normalized review findings")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	p := *root
	if !filepath.IsAbs(p) {
		wd, _ := os.Getwd()
		p = filepath.Join(wd, p)
	}
	g := *graph
	if g != "" && !filepath.IsAbs(g) {
		g = filepath.Join(p, filepath.FromSlash(g))
	}
	rest := fs.Args()
	if op != "" {
		rest = append([]string{op}, rest...)
	}
	switch command {
	case "reviews":
		if len(rest) != 1 || rest[0] != "import" || strings.TrimSpace(*source) == "" || strings.TrimSpace(*importFile) == "" {
			fmt.Fprintln(errout, "reviews import requires --source and --input")
			return 2
		}
		data, e := os.ReadFile(*importFile)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		var findings []reviewfinding.Finding
		if e = json.Unmarshal(data, &findings); e != nil {
			fmt.Fprintln(errout, "review input must be a JSON array:", e)
			return 2
		}
		imported, e := reviewfinding.Import(p, *source, findings)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		for _, finding := range imported {
			fmt.Fprintf(out, "IMPORTED %s ROUTE %s\n", finding.ID, finding.Route())
		}
		return 0
	case "resume":
		if *work == "" {
			fmt.Fprintln(errout, "resume requires --work")
			return 2
		}
		s, e := workflow.Resume(p, *work)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintf(out, "Workspace %s runtime v%d: %s (%s)\n", s.WorkspaceID, s.Version, s.Phase, s.Status)
		for _, b := range s.Blockers {
			fmt.Fprintln(out, "BLOCKED:", b)
		}
		return 0
	case "handoff":
		if *work == "" || *from == "" || *to == "" {
			fmt.Fprintln(errout, "handoff requires --work, --from, and --to")
			return 2
		}
		h, e := workflow.WriteHandoff(p, *work, *from, *to, *summary)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintln(out, "Created", h.ID)
		return 0
	case "checkpoint":
		c, e := workflow.WriteCheckpoint(p, *work, *step, *base, *input, *output)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintln(out, "Created", c.ID)
		return 0
	case "lease":
		if len(rest) == 0 {
			fmt.Fprintln(errout, "lease requires claim, heartbeat, or recover")
			return 2
		}
		switch rest[0] {
		case "cleanup":
			if !*yes || !*isolate || *work == "" || *task == "" {
				fmt.Fprintln(errout, "lease cleanup requires --work, --task, --isolate, and --yes")
				return 2
			}
			if e := workflow.RemoveTaskWorktree(filepath.Dir(p), *work, *task); e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "WORKTREE CLEANED", *task)
			return 0
		case "claim":
			if *isolate && strings.TrimSpace(*work) == "" {
				fmt.Fprintln(errout, "lease claim --isolate requires --work")
				return 2
			}
			l, e := workflow.ClaimLease(p, g, *task, *agent, 30*time.Minute)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintf(out, "LEASED %s to %s until %s\n", l.TaskID, l.Agent, l.ExpiresAt)
			if *isolate {
				repo := filepath.Dir(p)
				path, x := workflow.CreateTaskWorktree(repo, *work, *task)
				if x != nil {
					_ = workflow.ReleaseLease(p, *task, *agent)
					fmt.Fprintln(errout, x)
					return 1
				}
				fmt.Fprintln(out, "WORKTREE", path)
			}
			return 0
		case "heartbeat":
			l, e := workflow.Heartbeat(p, *task, *agent, 30*time.Minute)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "HEARTBEAT", l.TaskID, l.ExpiresAt)
			return 0
		case "recover":
			xs, e := workflow.RecoverLeases(p)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			for _, x := range xs {
				fmt.Fprintln(out, "RECOVERED", x)
			}
			return 0
		}
	case "commands":
		if len(rest) == 0 {
			fmt.Fprintln(errout, "commands requires plan or execute")
			return 2
		}
		if rest[0] == "plan" {
			argv := rest[1:]
			pl, e := workflow.CreateCommandPlan(p, *work, *task, *cwd, *source, *risk, argv, *timeout)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "PLANNED", pl.ID)
			return 0
		}
		if rest[0] == "execute" {
			if !*yes {
				fmt.Fprintln(errout, "commands execute requires --yes")
				return 2
			}
			if len(rest) < 2 {
				fmt.Fprintln(errout, "commands execute requires plan id")
				return 2
			}
			ev, e := workflow.ExecuteCommandPlan(p, *work, rest[1], *agent)
			fmt.Fprint(out, ev.Output)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			return 0
		}
	case "schedule":
		if len(rest) == 2 && rest[0] == "activate" {
			if !*yes || !*isolate || *agent == "" {
				fmt.Fprintln(errout, "schedule activate requires wave id, --agent, --isolate, and --yes")
				return 2
			}
			paths, e := workflow.ActivateScheduledWave(p, *work, g, rest[1], *agent)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			for _, path := range paths {
				fmt.Fprintln(out, "WORKTREE", path)
			}
			return 0
		}
		s, e := workflow.BuildSchedule(p, *work, g, *max)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		for _, w := range s.Waves {
			fmt.Fprintf(out, "%s %s\n", w.ID, strings.Join(w.Tasks, ","))
		}
		return 0
	case "integrate":
		if len(rest) == 0 {
			fmt.Fprintln(errout, "integrate requires plan or apply")
			return 2
		}
		if rest[0] == "plan" {
			i, e := workflow.CreateIntegration(p, *work, *base, splitCSV(*commits))
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, "PLANNED", i.ID)
			return 0
		}
		if rest[0] == "apply" {
			if !*yes || len(rest) < 2 {
				fmt.Fprintln(errout, "integrate apply requires id and --yes")
				return 2
			}
			i, e := workflow.ApplyIntegration(p, rest[1])
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			fmt.Fprintln(out, i.Status, i.IntegratedDiffHash)
			return 0
		}
	case "runtime":
		if *work == "" {
			fmt.Fprintln(errout, "runtime requires --work")
			return 2
		}
		if len(rest) > 0 && rest[0] == "memory" {
			if len(rest) < 2 || (rest[1] != "inspect" && rest[1] != "compact") {
				fmt.Fprintln(errout, "runtime memory requires inspect or compact")
				return 2
			}
			if rest[1] == "compact" && !*yes {
				fmt.Fprintln(errout, "runtime memory compact requires --yes")
				return 2
			}
			var assessment workflow.MemoryAssessment
			var e error
			if rest[1] == "compact" {
				assessment, e = workflow.CompactRuntimeMemory(p, *work, *task)
			} else {
				assessment, e = workflow.AssessRuntimeMemory(p, *work, *task)
			}
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			if *jsonOutput {
				data, _ := json.Marshal(assessment)
				fmt.Fprintln(out, string(data))
			} else {
				fmt.Fprintf(out, "Memory %s: %d source(s), %d active risk(s), %d contradiction(s)\n", assessment.Path, len(assessment.Sources), len(assessment.ActiveRisks), len(assessment.Contradictions))
			}
			return 0
		}
		if len(rest) > 0 && rest[0] == "reconcile" {
			findings, e := workflow.ReconcileRuntime(p, *work)
			if e != nil {
				fmt.Fprintln(errout, e)
				return 1
			}
			if *jsonOutput {
				data, _ := json.Marshal(findings)
				fmt.Fprintln(out, string(data))
			} else {
				for _, finding := range findings {
					fmt.Fprintf(out, "%s %s -> %s\n", finding.Kind, finding.Detail, finding.Owner)
				}
			}
			return 0
		}
		if len(rest) > 0 && (rest[0] == "status" || rest[0] == "watch") {
			watch := rest[0] == "watch"
			if *interval <= 0 {
				fmt.Fprintln(errout, "watch interval must be positive")
				return 2
			}
			iterations := 1
			if watch && *count > 0 {
				iterations = *count
			}
			for index := 0; index < iterations || (watch && *count == 0); index++ {
				observation, e := workflow.ObserveRuntime(p, *work)
				if e != nil {
					fmt.Fprintln(errout, e)
					return 1
				}
				if *jsonOutput {
					data, _ := json.Marshal(observation)
					fmt.Fprintln(out, string(data))
				} else {
					fmt.Fprintf(out, "Workspace %s: %s (%s); %d event(s), %d checkpoint(s), %d evidence record(s), %d active lease(s)\n", observation.State.WorkspaceID, observation.State.Phase, observation.State.Status, observation.Events, observation.Checkpoints, observation.Evidence, len(observation.ActiveLeases))
				}
				if !watch || (*count > 0 && index+1 >= *count) {
					break
				}
				time.Sleep(*interval)
			}
			return 0
		}
		msg, e := workflow.MigrateWorkspace(p, *work, *dry)
		if e != nil {
			fmt.Fprintln(errout, e)
			return 1
		}
		fmt.Fprintln(out, msg)
		return 0
	}
	fmt.Fprintln(errout, "unsupported runtime operation")
	return 2
}
