package cmd

import (
	"errors"
	"github.com/polyverse/ropoly-cmd/lib/types/BinDump"
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

var STRING_FORM_MAP = map[string]form{
	FORM_FILEPATH_STRING:    FORM_FILEPATH,
	FORM_PID_STRING:         FORM_PID,
}

var FORM_STRING_MAP = map[form]string{
	FORM_FILEPATH:    FORM_FILEPATH_STRING,
	FORM_PID:         FORM_PID_STRING,
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
	default:
		err := errors.New("Bad input form " + FORM_STRING_MAP[form])
		return BinDump.BinDump{}, err
	}
}