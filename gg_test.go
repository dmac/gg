package gg

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func testVersion(t *testing.T, version string) {
	t.Parallel()
	cmd := exec.Command("go", "test")
	cmd.Dir = filepath.Join("testdata", version)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error testing %s: %s\nOutput:\n%s", version, err, string(output))
	}
}

func Test21(t *testing.T) {
	testVersion(t, "v2.1")
}

func Test32Core(t *testing.T) {
	testVersion(t, "v3.2-core")
}
