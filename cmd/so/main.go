package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"solod.dev/internal/compiler"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var err error
	switch cmd {
	case "translate":
		err = translate(args)
	case "build":
		err = build(args)
	case "run":
		err = run(args)
	case "version":
		fmt.Printf("so version %s\n", compiler.Version())
		return
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "so %s: %s\n", cmd, err)
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `So is a tool for managing Solod source code.

Usage: so <command> [arguments]

Commands:
    build        compile package to executable
    run          compile and run a package
    translate    translate package to C
    version      print compiler version

Run 'so <command> -h' for details.
`)
}

func translate(args []string) error {
	flags := flag.NewFlagSet("translate", flag.ContinueOnError)
	outDir := flags.String("o", "", "output directory (default: package directory)")
	if err := flags.Parse(args); err != nil {
		return err
	}

	pkg := "."
	if flags.NArg() > 0 {
		pkg = flags.Arg(0)
	}

	out := *outDir
	if out == "" {
		out = pkg
	}

	return compiler.Translate(pkg, out)
}

func build(args []string) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	outFile := flags.String("o", "", "output file (default: basename of package directory)")
	if err := flags.Parse(args); err != nil {
		return err
	}

	pkg := "."
	if flags.NArg() > 0 {
		pkg = flags.Arg(0)
	}

	out := *outFile
	if out == "" {
		absDir, err := filepath.Abs(pkg)
		if err != nil {
			return err
		}
		out = filepath.Base(absDir)
	}

	return compiler.Build(pkg, out)
}

func run(args []string) error {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return err
	}

	pkg := "."
	if flags.NArg() > 0 {
		pkg = flags.Arg(0)
	}

	return compiler.Run(pkg)
}
