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
var testOkResult = `1
2
3
4
5
`

func TestOK(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testOk))
	out := new(bytes.Buffer)
	err := Uniq(in, out)
	if err != nil {
		t.Errorf("Test for OK Failed - error")
	}

	result := out.String()
	if result != testOkResult {
		t.Errorf("Test for OK Failed - results not match: %v %v", result, testOkResult)
	}
}

var testFail = `1
2
1`

func TestError(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testFail))
	out := new(bytes.Buffer)
	err := Uniq(in, out)
	if err != nil {
		t.Errorf("Test for Error Failed - error %v", err)
	}
}
