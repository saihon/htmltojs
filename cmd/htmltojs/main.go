package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	flag "github.com/saihon/flags"
	htmltojs "github.com/saihon/htmltojs"
)

const (
	NAME    = "htmltojs"
	VERSION = "v0.0.1"
)

var (
	input         string
	output        string
	defaultParent string
)

func init() {
	flag.CommandLine.Init(NAME, flag.ContinueOnError, true)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "\nUsage: %s [options] [arguments]\n\nOptions:\n\n", NAME)
		flag.PrintCustom()
	}

	flag.Bool("version", 'v', false, "Output version information and exit\n",
		func(_ flag.Getter) error {
			fmt.Fprintf(flag.CommandLine.Output(), "%s: %s\n", NAME, VERSION)
			return flag.ErrHelp
		})

	flag.StringVar(&input, "input", 'i', "", "Specify the input file name.\n", nil)
	flag.StringVar(&output, "output", 'o', "", "Specify the output file name.\n", nil)
	flag.StringVar(&defaultParent, "default-parent", 'p', "document.body",
		"Specify a default appending target element.\n", nil)
}

func run() error {
	var r io.Reader = os.Stdin
	var w io.Writer = os.Stdout

	if terminal.IsTerminal(0) {
		f, err := os.Open(input)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	if output != "" {
		f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	h := htmltojs.New()

	if defaultParent != "" {
		h.DefaultParent = defaultParent
	}

	if err := h.Parse(r); err != nil {
		return err
	}

	_, err := h.WriteTo(w)
	return err
}

func _main() int {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if flag.IsIgnorableError(err) {
			return 2
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
	}()

	os.Exit(_main())
}
