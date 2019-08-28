package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// TODO Two issues:
//
// First, small, is that data should be able to come from stdin as well as file
//
// Second, bigger, is that using crazy data like an image in an argument (e.g. "--data cRaZy_ImAgE_dAtA") ends up including
// escape characters and other shit, so it needs to come from a file or stdin, not a flag

var version string

func main() {
	app := cli.App{
		Name:        "stegosaurus",
		HelpName:    "stegosaurus",
		Version:     version,
		Description: "stegonography tool",
		Author:      "Gavin Cabbage",
		Commands: []cli.Command{
			encodeCmd,
			decodeCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(fmt.Errorf("%s: %w", app.Name, err))
		os.Exit(1)
	}
}
func reader(s string) (io.Reader, error) {
	if strings.HasPrefix(s, "@") {
		d, err := ioutil.ReadFile(strings.TrimPrefix(s, "@"))
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(d), nil
	}

	return strings.NewReader(s), nil
}
