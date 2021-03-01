package Instruction

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type Octets []byte

func (o Octets) String() string {
	buffer := &strings.Builder{}
	first := true
	for _, b := range o {
		if !first {
			fmt.Fprint(buffer, " ")
		}
		str := strconv.FormatUint(uint64(b), 16)
		buffer.WriteString("0x" + strings.Repeat("0", 2-len(str)) + str)
		first = false
	}
	return buffer.String()
}

func (o *Octets) UnmarshalJSON(b []byte) error {
	if o == nil {
		return errors.Errorf("Octets Unmarshall cannot operate on a nil pointer.")
	}

	sanitizedStr := string(b)
	sanitizedStr = strings.TrimPrefix(sanitizedStr, "\"")
	sanitizedStr = strings.TrimSuffix(sanitizedStr, "\"")

	octetsStr := strings.Split(sanitizedStr, " ")
	newOctect := make([]byte, 0, len(octetsStr))
	for _, octetStr := range octetsStr {
		if !strings.HasPrefix(octetStr, "0x") {
			return errors.Errorf("Octet %s is not prefixed with 0x. Only hexadecimal Octets are allowed.", octetStr)
		}
		octetVal, err := strconv.ParseUint(strings.TrimPrefix(octetStr, "0x"), 16, 8)
		if err != nil {
			return errors.Wrapf(err, "Unable to parse octect %s into an 8-bit unsigned integer", octetStr)
		}
		newOctect = append(newOctect, byte(octetVal))
	}
	*o = newOctect
	return nil
}

func (o Octets) MarshalJSON() ([]byte, error) {
	return []byte("\"" + o.String() + "\""), nil
}

type Instruction struct {
	Octets Octets   `json:"octets"`
	DisAsm string   `json:"disasm"`
}