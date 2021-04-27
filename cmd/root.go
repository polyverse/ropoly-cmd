package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
  Use:      "./ropoly-cmd",
  Short:    "Ropoly-cmd is a tool for evaluating entropy in scrambled binaries.",
  Long:     "Ropoly-cmd is a tool for evaluating entropy in scrambled binaries. Documentation can be found at https://github.com/polyverse/ropoly-cmd.",
  Run: func(cmd *cobra.Command, args []string) {
    println("Supported commands: dump, fingerprint, eqi")
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}