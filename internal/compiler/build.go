package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// Build translates the Go package in srcDir to C and compiles it into outFile.
// Uses CC (default "cc"), CFLAGS, and LDFLAGS environment variables.
func Build(srcDir, outFile string) error {
	tmpDir, err := os.MkdirTemp("", "solod_build")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := Translate(srcDir, tmpDir); err != nil {
		return err
	}

	cFiles, err := findCFiles(tmpDir)
	if err != nil {
		return err
	}

	return compileC(tmpDir, cFiles, outFile)
}

// Run translates and compiles the Go package in srcDir, then executes it.
// Returns an *exec.ExitError if the program exits with a non-zero status.
func Run(srcDir string) error {
	tmpFile, err := os.CreateTemp("", "solod_run")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err := Build(srcDir, tmpFile.Name()); err != nil {
		return err
	}

	cmd := exec.Command(tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Version returns the compiler version string to embed into compiled
// programs via -Dso_version. It uses the module version from
// runtime/debug.BuildInfo when available (e.g. go install ...@vx.y.z),
// falling back to "(devel)" (e.g. go run during development).
func Version() string {
	const devel = "(devel)"
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return devel
	}
	if v := info.Main.Version; v != "" {
		return v
	}
	return devel
}

// findCFiles returns all .c files under dir, recursively.
func findCFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".c") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("find C files: %w", err)
	}
	return files, nil
}

// compileC invokes the C compiler to produce an executable.
func compileC(includeDir string, cFiles []string, outFile string) error {
	cc := os.Getenv("CC")
	if cc == "" {
		cc = "cc"
	}

	args := []string{"-I" + includeDir}
	args = append(args, fmt.Sprintf(`-Dso_version="%s"`, Version()))
	args = append(args, splitFlags(os.Getenv("CFLAGS"))...)
	args = append(args, cFiles...)
	args = append(args, "-o", outFile)
	args = append(args, splitFlags(os.Getenv("LDFLAGS"))...)

	cmd := exec.Command(cc, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("C compiler: %w", err)
	}
	return nil
}

// splitFlags splits a space-separated flags string into individual args.
func splitFlags(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Fields(s)
}
