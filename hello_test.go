package gonum

import "testing"

func TestTheMatrix(t *testing.T) {
	if theMatrix != -1 {
		t.Errorf("Bad constant value")
	}
}
