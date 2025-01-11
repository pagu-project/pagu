package main

import (
	"github.com/pagu-project/pagu/cmd"
	pagu "github.com/pagu-project/pagu/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-http",
		Version: pagu.StringVersion(),
	}

	runCommand(rootCmd)

	err := rootCmd.Execute()
	cmd.ExitOnError(rootCmd, err)
}
