package cmd

import (
    "errors"
    "fmt"
    "github.com/polyverse/ropoly-cmd/lib/eqi"
    "github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
    "github.com/spf13/cobra"
    "strconv"
)

var eqiFuncString string
var monteCarloNumGadgets uint
var monteCarloTrials uint

func init() {
    eqiCmd.PersistentFlags().StringVarP(&eqiFuncString, "eqi-func", "f", "shared-offsets", "EQI calculation to use")
    eqiCmd.PersistentFlags().UintVarP(&monteCarloNumGadgets, "num-gadgets", "g", 3, "number of gadgets to use each trial for \"--eqi-func=monte-carlo\"")
    eqiCmd.PersistentFlags().UintVarP(&monteCarloTrials, "trials", "t", 100000, "number of trials to perform for \"--eqi-func=monte-carlo\"")

    rootCmd.AddCommand(eqiCmd)
}

var eqiCmd = &cobra.Command{
    Use:        "eqi --binary=<path>|--pid=<PID>|--dump=<path>|--fingerprint=<path>",
    Short:      "Compares two binaries and generates an EQI representing how well-scrambled the second is relative to the first.",
    Args:       cobra.ExactArgs(2),
    RunE:       func(cmd *cobra.Command, args []string) error {
        f0, input0 := positionalArgAsFormAndValue(args[0])
        fingerprint0, err := getInputAsFingerprint(input0, f0, gadgetSpec)
        if err != nil {
       		return err
       	}

       	f1, input1 := positionalArgAsFormAndValue(args[1])
       	fingerprint1, err := getInputAsFingerprint(input1, f1, gadgetSpec)
       	if err != nil {
       	    return err
       	}

        eqiFunc, err := getEqiFunc()
        if err != nil {
            return err
        }

        eqi := eqiFunc(fingerprint0, fingerprint1)

        eqiString := strconv.FormatFloat(eqi, 'f', -1, 64)
        fmt.Print(eqiString)
        return nil
    },
}

type eqiFunc func(f1, f2 Fingerprint.Fingerprint) float64

func getEqiFunc() (eqiFunc, error) {
	switch eqiFuncString {
	case "shared-offsets":
		return eqi.SharedOffsetsPerGadgetEqi, nil
	case "kill-rate":
		return eqi.KillRateEqi, nil
	case "highest-offset-count":
		return eqi.HighestOffsetCountEqi, nil
	case "kill-rate-without-movement":
		return eqi.KillRateWithoutMovementEqi, nil
	case "monte-carlo":
		return eqi.SharedOffsetExistsMonteCarloEqi(monteCarloNumGadgets, monteCarloTrials), nil
	default:
		return nil, errors.New(eqiFuncString + " is not a valid EQI function.")
	}
}