package pdf

import (
	"fmt"
	
	"github.com/dslipak/pdf"
)

func ValidatePdfArgs(params ...bool) error {
	sum := 0
	for _, param := range params {
		if param {
			sum++
		}
	}
	if sum > 1 {
		return fmt.Errorf("more than one pdf parser's parameter set!")
	}
	if sum == 0 {
		return fmt.Errorf("one pdf parser's parameter must be set!")
	}
	return nil
}

func isSameSentence(t1, t2 pdf.Text) bool {
	if t1.Y != t2.Y {
		return false
	}
 	return true
}