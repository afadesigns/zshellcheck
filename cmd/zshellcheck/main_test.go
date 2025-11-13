package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestProcessFile(t *testing.T) {
	// Create a temporary file with a violation
	tmpfile, err := ioutil.TempFile("", "test.zsh")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte(`echo $my_array[1]`)
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	var outBuf, errBuf bytes.Buffer
	processFile(tmpfile.Name(), &outBuf, &errBuf)

	if errBuf.Len() > 0 {
		t.Errorf("Expected no errors, but got: %s", errBuf.String())
	}

	expected := "ZC1001: Use ${} for array element access"
	if !strings.Contains(outBuf.String(), expected) {
		t.Errorf("Expected output to contain %q, but got %q", expected, outBuf.String())
	}
}