package cli

import (
	"flag"
	"fmt"
	"github.com/JonatasFreireDev/spec-framework/internal/dispatch"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
	command := fs.String("command", "", "supervised executable")
	enable := fs.Bool("enable", false, "explicitly enable supervised execution")
	ids := fs.String("ids", "", "comma-separated assigned dispatch ids")
	max := fs.Int("max-parallel", 1, "maximum concurrent dispatches")
	diffHash := fs.String("diff-hash", "", "immutable working-tree diff hash")
	parent := fs.String("parent", "", "returned code-runner dispatch id")
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
		if !*yes || !*enable || *work == "" || *ids == "" || *command == "" {
			fmt.Fprintln(errout, "dispatch wave requires --work --ids --command --enable --yes")
			return 2
		}
		for _, r := range dispatch.RunWave(p, *work, splitCSV(*ids), *max, *enable, *command, fs.Args()) {
			if r.Error != "" {
				fmt.Fprintln(errout, r.ID, r.Error)
			} else {
				fmt.Fprintln(out, "RAN", r.ID, r.Transcript.OutputHash)
			}
		}
		return 0
	}
	fmt.Fprintln(errout, "unknown dispatch operation")
	return 2
}
