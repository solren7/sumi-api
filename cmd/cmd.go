package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootView = &cobra.Command{
	Use:   "myapp",
	Short: "A comprehensive Fiber backend",
}

func Execute() {
	if err := rootView.Execute(); err != nil {
		os.Exit(1)
	}
}
