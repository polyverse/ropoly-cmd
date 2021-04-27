package cmd

import (
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/eqi"
    "github.com/polyverse/ropoly-cmd/lib/types/Gadget"
    "github.com/spf13/cobra"
    "strconv"
)

func init() {
  rootCmd.AddCommand(eqiCmd)
}

var eqiCmd = &cobra.Command{
    Use:        "eqi --binary=<path>|--pid=<PID>|--dump=<path>|--fingerprint=<path>",
    Short:      "Compares two binaries and generates an EQI representing how well-scrambled the second is relative to the first.",
    Args:       cobra.ExactArgs(2),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f0, input0 := positionalArgAsFormAndValue(args[0])
        fingerprint0, err := getInputAsFingerprint(input0, f0, Gadget.UserSpec{
            MinSize: MIN_GADGET_LENGTH_DEFAULT,
            MaxSize: MAX_GADGET_LENGTH_DEFAULT,
        })
        if err != nil {
       		return err
       	}

       	f1, input1 := positionalArgAsFormAndValue(args[0])
       	fingerprint1, err := getInputAsFingerprint(input1, f1, Gadget.UserSpec{
            MinSize: MIN_GADGET_LENGTH_DEFAULT,
            MaxSize: MAX_GADGET_LENGTH_DEFAULT,
        })

        eqiFunc := eqi.SharedOffsetsPerGadgetEqi

        eqi := eqiFunc(fingerprint0, fingerprint1)

        eqiString := strconv.FormatFloat(eqi, 'f', -1, 64)
        fmt.Print(eqiString)
        return nil
    },
}