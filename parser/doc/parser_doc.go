package doc

import (
	"fmt"
	"os/exec"
	"os"

)

func ParseAndWriteDoc(inputPath string, outputPath string) error{
	if _, err := os.Stat(inputPath); err != nil {
		return err
	}
	antiword := fmt.Sprintf("antiword -f %v > %v", inputPath, outputPath)
	cmd := exec.Command("bash", "-c", antiword)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
	
}