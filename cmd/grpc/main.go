package main

import (
	"github.com/pagu-project/pagu/cmd"
	"github.com/pagu-project/pagu/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-grpc",
		Version: version.StringVersion(),
	}

	runCommand(rootCmd)

	err := rootCmd.Execute()
	cmd.ExitOnError(rootCmd, err)
}
