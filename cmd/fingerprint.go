package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
    "github.com/polyverse/ropoly-cmd/lib/types/Gadget"
    "github.com/spf13/cobra"
)

const MIN_GADGET_LENGTH_DEFAULT = 0
const MAX_GADGET_LENGTH_DEFAULT = 6

var gadgetSpec Gadget.UserSpec
var inputIsPid bool

func init() {
    fingerprintCmd.PersistentFlags().BoolVarP(&inputIsPid, "pid", "p", false, "Input argument represents the PID of a running process, rather than a path to a binary file.")
    fingerprintCmd.PersistentFlags().UintVarP(&gadgetSpec.MinSize, "gadget-min-length", "s", MIN_GADGET_LENGTH_DEFAULT, "minimum gadget length to include in fingerprint")
    fingerprintCmd.PersistentFlags().UintVarP(&gadgetSpec.MaxSize, "gadget-max-length", "l", MAX_GADGET_LENGTH_DEFAULT, "maximum gadget length to include in fingerprint")

    rootCmd.AddCommand(fingerprintCmd)
}

var fingerprintCmd = &cobra.Command{
    Use:        "fingerprint <path/to/binary (or PID with --pid)>",
    Short:      "Generates a fingerprint of a binary file (default) or running process (with --pid).",
    Args:       cobra.ExactArgs(1),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f := FORM_FILEPATH
        if inputIsPid {
            f = FORM_PID
        }

        bindump, err := getInputAsBinDump(args[0], f)
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