package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
    "github.com/polyverse/ropoly-cmd/lib/types/Gadget"
    "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(fingerprintCmd)
}

var fingerprintCmd = &cobra.Command{
    Use:        "fingerprint --binary=<path>|--pid=<PID>|--dump=<path>",
    Short:      "Generates a fingerprint of a binary or running process.",
    Args:       cobra.ExactArgs(1),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f, input := positionalArgAsFormAndValue(args[0])
        fingerprint, err := getInputAsFingerprint(input, f, Gadget.UserSpec{
            MinSize: MIN_GADGET_LENGTH_DEFAULT,
            MaxSize: MAX_GADGET_LENGTH_DEFAULT,
        })
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