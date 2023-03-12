package parser

import (
	"fmt"
	"os"
)

func WriteContent[T string | []string](content T, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	switch v := any(content).(type) {
	case string:
		_, err = f.WriteString(v)
		if err != nil {
			return err
		}
	case []string:
		for _, line := range v {
			_, err := f.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid type of input")
	}
	return nil
}
