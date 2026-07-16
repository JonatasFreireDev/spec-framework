package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/JonatasFreireDev/spec-framework/internal/projectserver"
)

func runServer(args []string, stdout, stderr io.Writer) int {
	action := "start"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		action, args = args[0], args[1:]
	}
	flags := flag.NewFlagSet("server "+action, flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("product-root", "product", "product root")
	port := flags.Int("port", 0, "local port; 0 selects a free port")
	noOpen := flags.Bool("no-open", false, "do not open the browser")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	productRoot := *root
	if !filepath.IsAbs(productRoot) {
		productRoot = filepath.Join(cwd, productRoot)
	}
	switch action {
	case "start":
		return startServer(productRoot, *port, *noOpen, stdout, stderr)
	case "stop":
		if err := projectserver.Stop(productRoot); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "Local project server stopped.")
		return 0
	case "status":
		descriptor, err := projectserver.Healthy(productRoot)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "Local project server: running\n- URL: %s\n", descriptor.URL)
		return 0
	default:
		fmt.Fprintln(stderr, "server requires start, stop, or status")
		return 2
	}
}

func startServer(root string, port int, noOpen bool, stdout, stderr io.Writer) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	running, err := projectserver.Start(ctx, projectserver.Config{ProductRoot: root, Port: port})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Local project server running\n- URL: %s\n- Stop: Ctrl+C or spec-framework server stop --product-root %s\n", running.URL, filepath.Clean(root))
	if !noOpen {
		if err := openBrowser(running.URL); err != nil {
			fmt.Fprintln(stderr, "Could not open browser:", err)
		}
	}
	if err := <-running.Done; err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func openBrowser(url string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		command = exec.Command("open", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	return command.Start()
}
