package crlf

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/kortschak/utter"
)

func TestReadFile(t *testing.T) {
	b, err := ioutil.ReadFile("crlf_test.go")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	utter.Config.BytesWidth = 8
	dump := utter.Sdump(b)
	t.Log(dump)
	if bytes.Contains(b, []byte("\r\n")) {
		t.Errorf("unexpected CRLF in unix file: %s", dump)
	}
}
