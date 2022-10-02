package helper

import (
	"flag"
	"fmt"
)

func GetLevel(level string) int8 {
	switch level {
	case "trace":
		return -1
	case "debug":
		return 0
	case "info":
		return 1
	case "warn":
		return 2
	case "error":
		return 3
	case "fatal":
		return 4
	case "panic":
		return 5
	case "disabled":
		return 7
	default:
		return 6
	}
}

func EnumFlag(targetVar *string, flagName string, safeList []string, usage string) {
	flag.Func(flagName, usage, func(flagValue string) error {
		for _, safeValue := range safeList {
			if flagValue == safeValue {
				*targetVar = flagValue
				return nil
			}
		}

		return fmt.Errorf("must be one of %v", safeList)
	})
}
