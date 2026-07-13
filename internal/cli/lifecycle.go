package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/JonatasFreireDev/spec-framework/internal/clifecycle"
	"github.com/spf13/cobra"
)

func (app App) updateCommand(stdout, stderr io.Writer) *cobra.Command {
	var requested string
	var check, yes bool
	command := &cobra.Command{
		Use:   "update",
		Short: "Check for or install a CLI release.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := clifecycle.Default(app.version)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return commandExitError{code: 1}
			}
			release, err := manager.Check(context.Background(), requested)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return commandExitError{code: 1}
			}
			fmt.Fprintf(stdout, "CLI update\n- Current: %s\n- Available: %s\n- Platform archive: %s\n", release.Current, release.Latest, release.Archive)
			if !release.UpdateAvailable {
				fmt.Fprintln(stdout, "The CLI is already at the requested release.")
				return nil
			}
			if check || !yes {
				if !check {
					fmt.Fprintln(stdout, "Re-run with --yes to download, verify, and replace the CLI binary.")
				}
				return nil
			}
			result, err := manager.Update(context.Background(), requested)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return commandExitError{code: 1}
			}
			if result.Updated {
				if result.Scheduled {
					fmt.Fprintf(stdout, "Verified spec-framework %s; Windows replacement is scheduled after this process exits.\n", result.Release.Latest)
				} else {
					fmt.Fprintf(stdout, "Updated spec-framework to %s.\n", result.Release.Latest)
				}
			}
			return nil
		},
	}
	command.Flags().StringVar(&requested, "version", "", "release version; defaults to latest")
	command.Flags().BoolVar(&check, "check", false, "check without changing the installed binary")
	command.Flags().BoolVar(&yes, "yes", false, "confirm the binary replacement")
	return command
}

func (app App) uninstallCommand(stdout, stderr io.Writer) *cobra.Command {
	var purge, yes bool
	command := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the local CLI installation.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := clifecycle.Default(app.version)
			if err != nil {
				fmt.Fprintln(stderr, err)
				return commandExitError{code: 1}
			}
			plan := manager.PlanUninstall(purge)
			fmt.Fprintf(stdout, "CLI uninstall\n- Binary: %s\n- Install PATH entry: %s\n- Managed installation: %t\n", plan.Executable, plan.InstallDir, plan.Managed)
			if plan.Managed {
				fmt.Fprintf(stdout, "- Install manifest: %s\n", plan.Manifest)
			}
			fmt.Fprintln(stdout, "- Product repositories: never touched")
			if purge {
				fmt.Fprintf(stdout, "- Cache: %s\n", plan.CacheRoot)
				for _, dispatcher := range plan.Dispatchers {
					fmt.Fprintf(stdout, "- Dispatcher: %s\n", dispatcher)
				}
			}
			if !yes {
				fmt.Fprintln(stdout, "Re-run with --yes to remove these paths.")
				return nil
			}
			if _, err := manager.Uninstall(purge); err != nil {
				fmt.Fprintln(stderr, err)
				return commandExitError{code: 1}
			}
			fmt.Fprintln(stdout, "Spec Framework CLI removal scheduled/completed. Open a new terminal to refresh PATH.")
			return nil
		},
	}
	command.Flags().BoolVar(&purge, "purge", false, "also remove versioned runtime caches and Spec Framework dispatchers")
	command.Flags().BoolVar(&yes, "yes", false, "confirm removal")
	return command
}
