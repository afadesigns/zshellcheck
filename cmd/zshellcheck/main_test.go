package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/katas"
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
	processFile(tmpfile.Name(), &outBuf, &errBuf, "text")

	if errBuf.Len() > 0 {
		t.Errorf("Expected no errors, but got: %s", errBuf.String())
	}

	expected := "ZC1001: Use ${} for array element access"
	if !strings.Contains(outBuf.String(), expected) {
		t.Errorf("Expected output to contain %q, but got %q", expected, outBuf.String())
	}
}

func TestProcessFileJSON(t *testing.T) {
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
	processFile(tmpfile.Name(), &outBuf, &errBuf, "json")

	if errBuf.Len() > 0 {
		t.Errorf("Expected no errors, but got: %s", errBuf.String())
	}

	var violations []katas.Violation
	err = json.Unmarshal(outBuf.Bytes(), &violations)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(violations) != 2 {
		t.Fatalf("Expected 2 violations, got %d", len(violations))
	}

	foundZC1001 := false
	foundZC1037 := false
	for _, v := range violations {
		if v.KataID == "ZC1001" {
			foundZC1001 = true
		}
		if v.KataID == "ZC1037" {
			foundZC1037 = true
		}
	}

	if !foundZC1001 {
		t.Errorf("Expected KataID to be ZC1001, but it was not found")
	}

	if !foundZC1037 {
		t.Errorf("Expected KataID to be ZC1037, but it was not found")
	}
}