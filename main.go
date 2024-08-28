package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

var (
	flagHelp     = flag.Bool("h", false, "show help")
	flagVersion  = flag.Bool("v", false, "show version")
	flagDebug    = flag.Bool("d", false, "print debug output")
	flagNewLines = flag.Bool("n", false, "preserve newlines")
	flagPager    = flag.Bool("p", false, "send output to a pager")
	flagWidth    = flag.Int("w", 0, "specify a maximum width")
)

func debugf(format string, a ...any) {
	if *flagDebug {
		fmt.Fprintf(os.Stderr, "debug: "+format+"\n", a...)
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

	if *flagHelp {
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		return nil
	}

	if *flagVersion {
		fmt.Println(version())
		return nil
	}

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
		width, _, err = term.GetSize(outFd)
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
	}

	if *flagNewLines {
		opts = append(opts, glamour.WithPreservedNewLines())
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

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("tt must be built with go module support")
	}

	v := info.Main.Version

	var revision, modified string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}

	if revision != "" && len(revision) > 7 {
		v += "." + revision[:7]
	} else if revision != "" {
		v += "." + revision
	}

	if modified == "true" {
		v += "." + "modified"
	}

	return v + " built using " + info.GoVersion
}
