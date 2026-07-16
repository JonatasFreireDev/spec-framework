package cli

import (
	"bytes"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type commandExitError struct{ code int }

func (e commandExitError) Error() string { return fmt.Sprintf("command exited with status %d", e.code) }

// NewCommand builds the stable command tree. Leaf adapters deliberately keep
// their existing flag parsing during this migration so command syntax and exit
// behavior remain backward compatible while Cobra owns discovery and help.
func (app App) NewCommand(stdout, stderr io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           "spec-framework",
		Short:         "Specification Driven Development framework CLI",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "init" || cmd.Name() == "upgrade" || cmd.Name() == "update" || cmd.Name() == "uninstall" || cmd.Name() == "version" {
				return nil
			}
			if blocked, message := auditOnlyMutation(append([]string{cmd.Name()}, args...)); blocked {
				fmt.Fprintln(stderr, message)
				return commandExitError{code: 1}
			}
			return nil
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(
		app.versionCommand(stdout),
		app.updateCommand(stdout, stderr),
		app.uninstallCommand(stdout, stderr),
		app.legacyCommand("init", "Initialize product and the versioned runtime.", app.runInit, stdout, stderr),
		app.legacyCommand("upgrade", "Refresh the external runtime and pinned manifest.", app.runUpgrade, stdout, stderr),
		app.legacyCommand("migrate", "Preview or apply legacy runtime migration.", app.runMigrate, stdout, stderr),
		app.legacyCommand("move", "Move an artifact and update references.", func(args []string, out, errout io.Writer) int { return runMove(args, out, errout) }, stdout, stderr),
		app.legacyCommand("validate", "Validate a product repository.", func(args []string, out, errout io.Writer) int { return runValidate(args, out, errout) }, stdout, stderr),
		app.legacyCommand("template", "Audit canonical artifact template conformance.", func(args []string, out, errout io.Writer) int { return runTemplate(args, out, errout) }, stdout, stderr),
		app.legacyCommand("import", "Materialize approved source mappings as drafts.", func(args []string, out, errout io.Writer) int { return runImport(args, out, errout) }, stdout, stderr),
		app.legacyCommand("design", "Manage Design assets.", func(args []string, out, errout io.Writer) int { return runDesign(args, out, errout) }, stdout, stderr),
		app.legacyCommand("design-system", "Manage the product Design System.", func(args []string, out, errout io.Writer) int { return runDesignSystem(args, out, errout) }, stdout, stderr),
		app.legacyCommand("engineering-system", "Manage the product Engineering System.", func(args []string, out, errout io.Writer) int { return runEngineeringSystem(args, out, errout) }, stdout, stderr),
		app.legacyCommand("skill", "Resolve a pinned specialized skill path.", func(args []string, out, errout io.Writer) int { return runSkill(args, out, errout) }, stdout, stderr),
		app.legacyCommand("adapters", "List and manage optional external adapters.", func(args []string, out, errout io.Writer) int { return runAdapters(args, out, errout) }, stdout, stderr),
		app.legacyCommand("work", "Create a concurrent workspace for a feature.", func(args []string, out, errout io.Writer) int { return runWork(args, out, errout) }, stdout, stderr),
		app.legacyCommand("status", "Show workspace readiness and blockers.", func(args []string, out, errout io.Writer) int { return runWorkStatus("status", args, out, errout) }, stdout, stderr),
		app.legacyCommand("next", "Show the next skill for a workspace.", func(args []string, out, errout io.Writer) int { return runWorkStatus("next", args, out, errout) }, stdout, stderr),
		app.legacyCommand("approve", "Record an explicit artifact approval.", func(args []string, out, errout io.Writer) int { return runApprove(args, out, errout) }, stdout, stderr),
		app.legacyCommand("approve-batch", "Preview or approve multiple eligible artifacts.", func(args []string, out, errout io.Writer) int { return runApproveBatch(args, out, errout) }, stdout, stderr),
		app.legacyCommand("gates", "Check implementation gate readiness.", func(args []string, out, errout io.Writer) int { return runGates(args, out, errout) }, stdout, stderr),
		app.legacyCommand("graph", "Inspect and operate execution graphs.", func(args []string, out, errout io.Writer) int { return runGraph(args, out, errout) }, stdout, stderr),
		app.legacyCommand("task", "Inspect task readiness.", func(args []string, out, errout io.Writer) int { return runTask(args, out, errout) }, stdout, stderr),
		app.legacyCommand("guide", "Explain the current workspace gate and next action.", func(args []string, out, errout io.Writer) int { return runGuide(args, out, errout) }, stdout, stderr),
		app.legacyCommand("review", "Preview a workspace stage approval.", func(args []string, out, errout io.Writer) int { return runStage("review", args, out, errout) }, stdout, stderr),
		app.legacyCommand("approve-stage", "Approve eligible artifacts in a stage atomically.", func(args []string, out, errout io.Writer) int { return runStage("approve-stage", args, out, errout) }, stdout, stderr),
		app.legacyCommand("impact", "Inspect a decision's validity and propagation.", func(args []string, out, errout io.Writer) int { return runImpact(args, out, errout) }, stdout, stderr),
		app.legacyCommand("dashboard", "Show a consolidated workflow dashboard.", func(args []string, out, errout io.Writer) int { return runDashboard(args, out, errout) }, stdout, stderr),
		app.legacyCommand("server", "Run the local project-status server.", func(args []string, out, errout io.Writer) int { return runServer(args, out, errout) }, stdout, stderr),
		app.legacyCommand("decisions", "Check or migrate product decisions.", func(args []string, out, errout io.Writer) int { return runDecisions(args, out, errout) }, stdout, stderr),
		app.legacyCommand("dispatch", "Plan and supervise governed subagent assignments.", func(args []string, out, errout io.Writer) int { return runDispatch(args, out, errout) }, stdout, stderr),
	)
	for _, name := range []string{"resume", "handoff", "checkpoint", "lease", "commands", "schedule", "integrate", "runtime", "reviews"} {
		command := name
		root.AddCommand(app.legacyCommand(command, "Operate the resumable runtime.", func(args []string, out, errout io.Writer) int {
			return runRuntime(command, args, out, errout)
		}, stdout, stderr))
	}
	return root
}

func (app App) versionCommand(stdout io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(stdout, "spec-framework %s\n", app.version)
		},
	}
}

func (app App) legacyCommand(use, short string, run func([]string, io.Writer, io.Writer) int, stdout, stderr io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                use,
		Short:              short,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && (args[len(args)-1] == "--help" || args[len(args)-1] == "-h" || args[len(args)-1] == "help") {
				// Legacy leaves still own their flag parsing. Ask the leaf to
				// render its native flag package usage, but translate help into a
				// successful Cobra command so users get the real flags instead of
				// Cobra's flag-less parent summary.
				var helpOut, helpErr bytes.Buffer
				run(args, &helpOut, &helpErr)
				nativeHelp := append(helpOut.Bytes(), helpErr.Bytes()...)
				if !bytes.Contains(nativeHelp, []byte("Usage of ")) {
					return cmd.Help()
				}
				fmt.Fprintln(stdout, cmd.Short)
				if helpOut.Len() > 0 {
					_, _ = stdout.Write(helpOut.Bytes())
				}
				if helpErr.Len() > 0 {
					_, _ = stdout.Write(helpErr.Bytes())
				}
				return nil
			}
			if code := run(args, stdout, stderr); code != 0 {
				return commandExitError{code: code}
			}
			return nil
		},
	}
}
