package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/frk"
	"github.com/winebarrel/frk/subcmd"
)

var version string

var cli struct {
	Version  kong.VersionFlag
	Token    string             `required:"" env:"FRK_GITHUB_TOKEN" help:"GitHub token"`
	Activity subcmd.ActivityCmd `cmd:"" help:"show activity"`
	Pulls    subcmd.PullsCmd    `cmd:"" help:"show pull requests"`
}

func main() {
	parser := kong.Must(&cli, kong.Vars{"version": version})
	ctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	client, err := frk.NewGithub(cli.Token)
	parser.FatalIfErrorf(err)

	err = ctx.Run(&frk.Binds{
		Github: client,
	})

	ctx.FatalIfErrorf(err)
}
