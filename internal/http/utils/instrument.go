package utils

import (
	"fmt"
	"strings"
)

func SplitInstrument(instr string) (string, string, error) {
	parts := strings.Split(instr, "/")
	if len(parts) != 2 {
		parts = strings.Split(instr, "-")
	}
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid instrument: %s", instr)
	}
	return parts[0], parts[1], nil
}
