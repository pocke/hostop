package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/ogier/pflag"
	"github.com/pkg/errors"
)

const Hosts = "/etc/hosts"
const Footer = "###### Hostop End\n"
const Perm = 0644

func main() {
	if err := Main(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func Main(args []string) error {
	c, err := ParseArgs(args)
	if err != nil {
		return err
	}

	if c.ResetID != "" {
		return Reset(c)
	} else {
		return Stop(c)
	}
}

type CmdArgs struct {
	After   string
	ResetID string
	Hosts   []string
}

func ParseArgs(args []string) (*CmdArgs, error) {
	fs := pflag.NewFlagSet(args[0], pflag.ContinueOnError)
	c := new(CmdArgs)
	fs.StringVarP(&c.After, "after", "a", "", "")
	fs.StringVarP(&c.ResetID, "reset", "", "", "")

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Hosts = fs.Args()
	return c, nil
}

func Reset(c *CmdArgs) error {
	return nil
}

func Stop(c *CmdArgs) error {
	id := uuid.New().String()
	content := Header(id)
	for _, h := range c.Hosts {
		content += fmt.Sprintf("127.0.0.1 %s\n", h)
	}
	content += Footer

	f, err := os.OpenFile(Hosts, os.O_WRONLY|os.O_CREATE|os.O_APPEND, Perm)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	_, err = fmt.Fprint(f, content)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Header(id string) string {
	return fmt.Sprintf(`
###### Hostop Start %s
`, id)
}
