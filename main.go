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
	flagPager = flag.Bool("p", false, "send output to a pager")
)

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
	width := 80
	if term.IsTerminal(outFd) {
		width, _, err = term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return err
		}
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithWordWrap(width),
	)

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
