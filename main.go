package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

var (
	flagDebug = flag.Bool("d", false, "print debug output")
	flagPager = flag.Bool("p", false, "send output to a pager")
	flagWidth = flag.Int("w", 0, "specify a maximum width")
)

func debugf(format string, a ...any) {
	if *flagDebug {
		fmt.Fprintf(os.Stdout, "debug: "+format+"\n", a...)
	}
}

func main() {
	err := Main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
	}
}

func Main() (err error) {
	flag.Parse()

	inReader := os.Stdin
	if len(flag.Args()) > 0 {
		inReader, err = os.Open(flag.Arg(0))
		if err != nil {
			return err
		}
	}

	in, err := io.ReadAll(inReader)
	if err != nil {
		return err
	}

	outFd := int(os.Stdout.Fd())
	width := *flagWidth
	if term.IsTerminal(outFd) {
		width, _, err = term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return err
		}
	}
	debugf("term width: %d", width)

	// Limit width to the specified value.
	if *flagWidth > 0 && width > *flagWidth {
		width = *flagWidth
	}
	debugf("final width: %d", width)

	opts := []glamour.TermRendererOption{
		glamour.WithEnvironmentConfig(),
		glamour.WithPreservedNewLines(),
	}

	if width > 0 {
		opts = append(opts, glamour.WithWordWrap(width))
	}

	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return err
	}

	out, err := r.Render(string(in))
	if err != nil {
		return err
	}

	if flagPager != nil && *flagPager {
		err = page(out)
	} else {
		fmt.Print(out)
	}

	return err
}

func page(out string) error {
	var stderr bytes.Buffer

	cmd := exec.Command("less", "-R", "-F")
	cmd.Stdin = bytes.NewBuffer([]byte(out))
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.Join(err, fmt.Errorf("stderr: %s", stderr.String()))
	}

	return nil
}
