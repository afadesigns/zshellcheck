// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDownloadedInstallation(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("requires an unprivileged Unix account")
	}
	installer, err := filepath.Abs("../../install.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, support := range []bool{false, true} {
		t.Run(fmt.Sprintf("support-files=%t", support), func(t *testing.T) {
			dir := t.TempDir()
			files := releaseFixtureFiles(t, support)
			makeReleaseArchive(t, dir, files)
			toolsDir := installerTools(t, dir)
			cmd := exec.Command("bash", installer, "-y", "-v", "v0.0.0")
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "HOME="+dir, "SHELL=/bin/sh", "TMPDIR="+dir,
				"PATH="+toolsDir+string(os.PathListSeparator)+os.Getenv("PATH"), "INSTALL_FIXTURE="+dir)
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("download installation: %v\n%s", err, out)
			}
			assertInstalledFiles(t, dir, files)
		})
	}
}

func releaseFixtureFiles(t *testing.T, support bool) map[string]string {
	t.Helper()
	files := map[string]string{"zshellcheck": "#!/bin/sh\nprintf 'fixture version\\n'\n"}
	if support {
		for _, filename := range []string{"man/man1/zshellcheck.1", "completions/zsh/_zshellcheck", "completions/bash/zshellcheck-completion.bash"} {
			files[filename] = readTestFile(t, filepath.Join("../..", filename))
		}
	}
	return files
}

func makeReleaseArchive(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		header := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(content))}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	name := "zshellcheck_Linux_x86_64.tar.gz"
	writeTestFile(t, filepath.Join(dir, name), buf.String())
	writeTestFile(t, filepath.Join(dir, "checksums.txt"), fmt.Sprintf("%x  %s\n", sha256.Sum256(buf.Bytes()), name))
}

func installerTools(t *testing.T, dir string) string {
	t.Helper()
	toolsDir := filepath.Join(dir, "tools")
	if err := os.Mkdir(toolsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	scripts := map[string]string{
		"uname": "#!/bin/sh\ncase $1 in -s) echo Linux;; -m) echo x86_64;; *) exit 1;; esac\n",
		"curl": `#!/bin/sh
case $4 in
    */checksums.txt) cp "$INSTALL_FIXTURE/checksums.txt" "$3" ;;
    */zshellcheck_Linux_x86_64.tar.gz) cp "$INSTALL_FIXTURE/zshellcheck_Linux_x86_64.tar.gz" "$3" ;;
    *) exit 22 ;;
esac
`,
	}
	for name, source := range scripts {
		filename := filepath.Join(toolsDir, name)
		writeTestFile(t, filename, source)
		if err := os.Chmod(filename, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	return toolsDir
}

func assertInstalledFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	destinations := map[string]string{
		"zshellcheck":                                  "bin/zshellcheck",
		"man/man1/zshellcheck.1":                       "share/man/man1/zshellcheck.1",
		"completions/zsh/_zshellcheck":                 "share/zsh/site-functions/_zshellcheck",
		"completions/bash/zshellcheck-completion.bash": "share/bash-completion/completions/zshellcheck",
	}
	for source, want := range files {
		if got := readTestFile(t, filepath.Join(dir, ".local", destinations[source])); got != want {
			t.Errorf("installed %s differs from the release asset", source)
		}
	}
	// #nosec G204 -- the executable is the fixed test fixture in t.TempDir.
	out, err := exec.Command(filepath.Join(dir, ".local", "bin", "zshellcheck"), "--version").CombinedOutput()
	if err != nil || string(out) != "fixture version\n" {
		t.Errorf("installed executable: %v %q", err, out)
	}
}
