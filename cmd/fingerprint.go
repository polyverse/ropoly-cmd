package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
    "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(fingerprintCmd)
}

var fingerprintCmd = &cobra.Command{
    Use:        "fingerprint --binary=<path>|--pid=<PID>",
    Short:      "Generates a fingerprint of a binary file or running process.",
    Args:       cobra.ExactArgs(1),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f, input := positionalArgAsFormAndValue(args[0])
        bindump, err := getInputAsBinDump(input, f)
        if err != nil {
       		return err
       	}
       	fingerprint, err := Fingerprint.GenerateFingerprintFromBinDump(bindump, gadgetSpec)
       	if err != nil {
       	    return err
       	}

        b, err := Fingerprint.EncodeFingerprint(fingerprint)
       	if err != nil {
       		return err
       	}
       	fmt.Print(b)
       	return nil
    },
}