package main

import (
	"errors"
	"fmt"
	"github.com/polyverse/ropoly-cmd/lib/eqi"
	"github.com/polyverse/ropoly-cmd/lib/types/BinDump"
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"os"
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

func help() {
	println("Usage: " + os.Args[0] + " <input format (\"file\", \"pid\", \"bindump\", or \"fingerprint\")> <output format (\"bindump\", \"fingerprint\", or \"eqi\")> [additional args required for output format...]")
}

func main() {
	if len(os.Args) < 3 {
		help()
		println("Not enough arguments")
		os.Exit(1)
	}

	inputFormString := os.Args[1]
	inputForm := STRING_FORM_MAP[inputFormString]
	if !VALID_INPUT_FORM[inputForm] {
		help()
		println(inputFormString + " is not a valid input form.")
		os.Exit(1)
	}

	outputFormString := os.Args[2]
	outputForm := STRING_FORM_MAP[outputFormString]
	if !VALID_OUTPUT_FORM[outputForm] {
		help()
		println(outputFormString + " is not a valid output form.")
		os.Exit(1)
	}

	OUTPUT_FORM_FUNCS[outputForm](inputForm, os.Args[3:])
}

var STRING_FORM_MAP = map[string]form{
	FORM_FILEPATH_STRING:    FORM_FILEPATH,
	FORM_PID_STRING:         FORM_PID,
	FORM_BINDUMP_STRING:     FORM_BINDUMP,
	FORM_FINGERPRINT_STRING: FORM_FINGERPRINT,
	FORM_EQI_STRING:         FORM_EQI,
}

var FORM_STRING_MAP = map[form]string{
	0:                "unknown",
	FORM_FILEPATH:    FORM_FILEPATH_STRING,
	FORM_PID:         FORM_PID_STRING,
	FORM_BINDUMP:     FORM_BINDUMP_STRING,
	FORM_FINGERPRINT: FORM_FINGERPRINT_STRING,
	FORM_EQI:         FORM_EQI_STRING,
}

var VALID_INPUT_FORM = map[form]bool{
	FORM_FILEPATH:    true,
	FORM_PID:         true,
	FORM_BINDUMP:     true,
	FORM_FINGERPRINT: true,
	FORM_EQI:         false,
}

var VALID_OUTPUT_FORM = map[form]bool{
	FORM_FILEPATH:    false,
	FORM_PID:         false,
	FORM_BINDUMP:     true,
	FORM_FINGERPRINT: true,
	FORM_EQI:         true,
}

var OUTPUT_FORM_FUNCS = map[form]func(inputForm form, args []string){
	FORM_BINDUMP:     binDumpMain,
	FORM_FINGERPRINT: fingerprintMain,
	FORM_EQI:         eqiMain,
}

func helpBinDump() {
	println("Usage for bindump output: " + os.Args[0] + " <input format (\"file\", \"pid\", or \"bindump\")> " + os.Args[2] + " <input filepath or PID>")
}

func binDumpMain(inputForm form, args []string) {
	if len(args) < 1 {
		helpBinDump()
		println("Not enough arguments")
		os.Exit(1)
	}

	binDump, err := getInputAsBinDump(args[0], inputForm)
	if err != nil {
		helpBinDump()
		println("Failed to open or create bindump: " + err.Error())
		os.Exit(1)
	}

	b, err := BinDump.EncodeBinDump(binDump)
	if err != nil {
		println("Failed to encode bindump due to error: " + err.Error())
		println("If you see this, ropoly-cmd is not working correctly.")
		os.Exit(1)
	}
	fmt.Print(b)
}

func helpFingerprint() {
	println("Usage for fingerprint output: " + os.Args[0] + " <input format (\"file\", \"pid\", \"bindump\", or \"fingerprint\"> " + os.Args[2] + " <input filepath or PID> [min-gadget-length <number of instructions> (default 0 if not specified)] [max-gadget-length <number of instructions> (default 2 if not specified)]")
}

func fingerprintMain(inputForm form, args []string) {
	if len(args) < 1 {
		helpFingerprint()
		println("Not enough arguments")
		os.Exit(1)
	}

	gadgetSpec, err := parseGadgetSpec(args, inputForm)
	if err != nil {
		helpEqi()
		println("Failed to parse gadget spec: " + err.Error())
		os.Exit(1)
	}

	fingerprint, err := getInputAsFingerprint(args[0], inputForm, gadgetSpec)
	if err != nil {
		helpFingerprint()
		println("Failed to open or create fingerprint: " + err.Error())
		os.Exit(1)
	}

	b, err := Fingerprint.EncodeFingerprint(fingerprint)
	if err != nil {
		println("Failed to encode fingerprint due to error: " + err.Error())
		println("If you see this, ropoly-cmd is not working correctly.")
		os.Exit(1)
	}
	fmt.Print(b)
}

func helpEqi() {
	println("Usage for EQI output: " + os.Args[0] + " <input format > (\"file\", \"pid\", \"bindump\", or \"fingerprint\"" + os.Args[2] + " <input filepath or PID for initial binary> <input filepath or PID for modified binary> [min-gadget-length <number of instructions> (default 0 if not specified)] [max-gadget-length <number of instructions> (default 2 if not specified)]")
}

func eqiMain(inputForm form, args []string) {
	if len(args) < 2 {
		helpEqi()
		println("Not enough arguments")
		os.Exit(1)
	}

	gadgetSpec, err := parseGadgetSpec(args, inputForm)
	if err != nil {
		helpEqi()
		println("Failed to parse gadget spec: " + err.Error())
		os.Exit(1)
	}

	input1 := args[0]
	fingerprint1, err := getInputAsFingerprint(input1, inputForm, gadgetSpec)
	if err != nil {
		helpFingerprint()
		println("Failed to open or create fingerprint: " + err.Error())
		os.Exit(1)
	}

	input2 := args[1]
	fingerprint2, err := getInputAsFingerprint(input2, inputForm, gadgetSpec)
	if err != nil {
		helpFingerprint()
		println("Failed to open or create fingerprint: " + err.Error())
		os.Exit(1)
	}

	if fingerprint1.GadgetSpec != fingerprint2.GadgetSpec {
		println("Fingerprints were generated from different gadget specs; comparing the two would be meaningless.")
		println(input1 + " gadget spec: " + fingerprint1.GadgetSpec.String())
		println(input2 + " gadget spec: " + fingerprint2.GadgetSpec.String())
		os.Exit(1)
	}

	eqiFunc, err := getEqiFunc(args)
	if err != nil {
		helpFingerprint()
		println("Failed to parse EQI function: " + err.Error())
		os.Exit(1)
	}

	eqi := eqiFunc(fingerprint1, fingerprint2)

	eqiString := strconv.FormatFloat(eqi, 'f', -1, 64)
	fmt.Print(eqiString)
}

func getInputAsBinDump(inputPath string, form form) (BinDump.BinDump, error) {
	switch form {
	case FORM_FILEPATH:
		return BinDump.GenerateElfBinDumpFromFile(inputPath)
	case FORM_PID:
		pid, err := strconv.ParseUint(inputPath, 10, 64)
		if err != nil {
			return BinDump.BinDump{}, err
		}
		return BinDump.GenerateBinDumpFromPid(uint(pid))
	case FORM_BINDUMP:
		return BinDump.OpenBinDump(inputPath)
	default:
		err := errors.New("Cannot convert from " + FORM_STRING_MAP[form] + " to bindump")
		return BinDump.BinDump{}, err
	}
}

func getInputAsFingerprint(inputPath string, form form, spec Gadget.UserSpec) (Fingerprint.Fingerprint, error) {
	var binDump BinDump.BinDump
	var err error
	switch form {
	case FORM_FINGERPRINT:
		return Fingerprint.OpenFingerprint(inputPath)
	default:
		binDump, err = getInputAsBinDump(inputPath, form)
	}
	if err != nil {
		return Fingerprint.Fingerprint{}, err
	}

	return Fingerprint.GenerateFingerprintFromBinDump(binDump, spec)
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