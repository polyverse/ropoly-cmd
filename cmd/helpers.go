package cmd

import (
	"errors"
	"github.com/polyverse/ropoly-cmd/lib/types/BinDump"
	"github.com/polyverse/ropoly-cmd/lib/types/Fingerprint"
	"github.com/polyverse/ropoly-cmd/lib/types/Gadget"
	"strconv"
	"strings"
)

type form int

const FORM_FILEPATH = form(1)
const FORM_PID = form(2)
const FORM_BINDUMP = form(3)
const FORM_FINGERPRINT = form(4)
const FORM_EQI = form(5)

const FORM_FILEPATH_STRING = "binary"
const FORM_PID_STRING = "pid"
const FORM_BINDUMP_STRING = "dump"
const FORM_FINGERPRINT_STRING = "fingerprint"
const FORM_EQI_STRING = "eqi"

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

func positionalArgAsFormAndValue(arg string) (form, string) {
    split := strings.SplitN(arg, "=", 2)
    return STRING_FORM_MAP[split[0]], split[1]
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
		err := errors.New("Bad input form " + FORM_STRING_MAP[form])
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