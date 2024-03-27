package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/glamour"
)

func main() {
	err := Main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
	}
}

func Main() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: md PATH")
	}

	inFile, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}

	in, err := io.ReadAll(inFile)
	if err != nil {
		return err
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithWordWrap(80),
	)

	out, err := r.Render(string(in))
	if err != nil {
		return err
	}

	fmt.Print(out)

	return nil
}
