package main

import (
	"errors"
	"github.com/polyverse/ropoly-cmd/cmd"
	"github.com/polyverse/ropoly-cmd/lib/eqi"
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"strconv"
)

type form int

const FORM_FILEPATH = form(1)
const FORM_PID = form(2)
const FORM_BINDUMP = form(3)
const FORM_FINGERPRINT = form(4)
const FORM_EQI = form(5)

const FORM_FILEPATH_STRING = "file"
const FORM_PID_STRING = "pid"
const FORM_BINDUMP_STRING = "bindump"
const FORM_FINGERPRINT_STRING = "fingerprint"
const FORM_EQI_STRING = "eqi"

const MIN_GADGET_LENGTH_DEFAULT = 0
const MAX_GADGET_LENGTH_DEFAULT = 6

func main() {
    cmd.Execute()
}

func parseGadgetSpec(args []string, form form) (Gadget.UserSpec, error) {
	var minGadgetLength uint
	specified, minGadgetLengthString := getArgValue(args, "min-gadget-length")
	if specified {
		if form == FORM_FINGERPRINT {
			return Gadget.UserSpec{}, errors.New("Cannot specify min gadget length with input format fingerprint.")
		}
		minGadgetLength64, err := strconv.ParseUint(minGadgetLengthString, 10, 64)
		if err != nil {
			return Gadget.UserSpec{}, err
		}
		minGadgetLength = uint(minGadgetLength64)
	} else {
		minGadgetLength = MIN_GADGET_LENGTH_DEFAULT
	}

	var maxGadgetLength uint
	specified, maxGadgetLengthString := getArgValue(args, "max-gadget-length")
	if specified {
		if form == FORM_FINGERPRINT {
			return Gadget.UserSpec{}, errors.New("Cannot specify max gadget length with input format fingerprint.")
		}
		maxGadgetLength64, err := strconv.ParseUint(maxGadgetLengthString, 10, 64)
		if err != nil {
			return Gadget.UserSpec{}, err
		}
		maxGadgetLength = uint(maxGadgetLength64)
	} else {
		maxGadgetLength = MAX_GADGET_LENGTH_DEFAULT
	}

	return Gadget.UserSpec{
		minGadgetLength,
		maxGadgetLength,
	}, nil
}

func getArgValue(args []string, name string) (bool, string) {
	for index, token := range args {
		if token == name {
			return true, args[index+1]
		}
	}
	return false, ""
}

func getUintArg(args []string, name string, defaultArgValue uint) (uint, error) {
	argStringExists, argString := getArgValue(args, name)
	if argStringExists {
		argUint, err := strconv.ParseUint(argString, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(argUint), nil
	} else {
		return defaultArgValue, nil
	}
}

type eqiFunc func(f1, f2 Fingerprint.Fingerprint) float64

func getEqiFunc(args []string) (eqiFunc, error) {
	argExists, arg := getArgValue(args, "eqi-func")
	if !argExists {
		arg = "shared-offsets"
	}

	switch arg {
	case "shared-offsets":
		return eqi.SharedOffsetsPerGadgetEqi, nil
	case "kill-rate":
		return eqi.KillRateEqi, nil
	case "highest-offset-count":
		return eqi.HighestOffsetCountEqi, nil
	case "kill-rate-without-movement":
		return eqi.KillRateWithoutMovementEqi, nil
	case "monte-carlo":
		trials, err := getUintArg(args, "trials", 100000)
		if err != nil {
			return nil, err
		}
		numGadgets, err := getUintArg(args, "num-gadgets", 3)
		if err != nil {
			return nil, err
		}
		return eqi.SharedOffsetExistsMonteCarloEqi(numGadgets, trials), nil
	default:
		return nil, errors.New(arg + " is not a valid EQI function.")
	}
}