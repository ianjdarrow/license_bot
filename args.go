package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
)

type args struct {
	org    string
	force  bool
	action string
}

func parseArgs() args {
	parser := argparse.NewParser("license_bot", "a dope mass license conversion tool")
	org := parser.String("o", "org", &argparse.Options{Required: true, Help: "name of GitHub organization to convert"})
	force := parser.Flag("", "force", &argparse.Options{Required: false, Help: "☠️  automatically submit PRs for each repo"})
	action := parser.Selector("a", "action", []string{"contributors", "licenses"}, &argparse.Options{Required: true, Help: "action to take (licenses, contributors"})
	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(err)
	}
	return args{org: *org, force: *force, action: *action}
}
