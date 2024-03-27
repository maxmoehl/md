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

func Main() (err error) {
	inReader := os.Stdin

	if len(os.Args) == 2 {
		inReader, err = os.Open(os.Args[1])
		if err != nil {
			return err
		}
	}

	in, err := io.ReadAll(inReader)
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
