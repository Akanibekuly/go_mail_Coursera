package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var testOk = `1
2
3
3
4
5
`

func TestOK(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testOk))
	out := new(bytes.Buffer)
	err := Uniq(in, out)
	if err != nil {
		t.Errorf("Test for OK Failed")
	}
}
