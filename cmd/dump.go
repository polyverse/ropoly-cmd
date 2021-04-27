package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/BinDump"
    "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
    Use:        "dump --binary=<path>|--pid=<PID>",
    Short:      "Dumps the executable segments of a binary or running process.",
    Args:       cobra.ExactArgs(1),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f, input := positionalArgAsFormAndValue(args[0])
        binDump, err := getInputAsBinDump(input, f)
        if err != nil {
       		return err
       	}

        b, err := BinDump.EncodeBinDump(binDump)
       	if err != nil {
       		return err
       	}
       	fmt.Print(b)
       	return nil
    },
}