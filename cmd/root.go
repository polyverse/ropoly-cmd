package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/Gadget"
    "github.com/spf13/cobra"
    "os"
)

const MIN_GADGET_LENGTH_DEFAULT = 0
const MAX_GADGET_LENGTH_DEFAULT = 6

var gadgetSpec Gadget.UserSpec

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

func init() {
    rootCmd.PersistentFlags().UintVarP(&gadgetSpec.MinSize, "gadget-min-length", "s", MIN_GADGET_LENGTH_DEFAULT, "minimum gadget length to include in fingerprint")
    rootCmd.PersistentFlags().UintVarP(&gadgetSpec.MaxSize, "gadget-max-length", "l", MAX_GADGET_LENGTH_DEFAULT, "maximum gadget length to include in fingerprint")
}