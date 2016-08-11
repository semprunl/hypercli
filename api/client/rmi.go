package client

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	Cli "github.com/hyperhq/hypercli/cli"
	flag "github.com/hyperhq/hypercli/pkg/mflag"
)

// CmdRmi removes all images with the specified name(s).
//
// Usage: docker rmi [OPTIONS] IMAGE [IMAGE...]
func (cli *DockerCli) CmdRmi(args ...string) error {
	cmd := Cli.Subcmd("rmi", []string{"IMAGE [IMAGE...]"}, Cli.DockerCommands["rmi"].Description, true)
	force := cmd.Bool([]string{"f", "-force"}, false, "Force removal of the image")
	noprune := cmd.Bool([]string{}, false, "Do not delete untagged parents")
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	v := url.Values{}
	if *force {
		v.Set("force", "1")
	}
	if *noprune {
		v.Set("noprune", "1")
	}

	ctx := context.Background()

	var errs []string
	for _, name := range cmd.Args() {
		options := types.ImageRemoveOptions{
			Force:         *force,
			PruneChildren: !*noprune,
		}

		dels, err := cli.client.ImageRemove(ctx, name, options)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			for _, del := range dels {
				if del.Deleted != "" {
					fmt.Fprintf(cli.out, "Deleted: %s\n", del.Deleted)
				} else {
					fmt.Fprintf(cli.out, "Untagged: %s\n", del.Untagged)
				}
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}
