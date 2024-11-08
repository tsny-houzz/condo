package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tsny-houzz/condo/pkg/jbd/jobs"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

func main() {
	klient := newClient()

	commands := []*cli.Command{
		{
			Name:  "whoami",
			Usage: "Show the user's email",
			Action: func(c *cli.Context) error {
				email, err := LoadEmail()
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
				}
				fmt.Printf("User email: %s\n", email)
				return nil
			},
		},

		// TODO: Fix the client [needs oauth :( ])
		{
			Name:   "jbd",
			Usage:  "jukwaa-tools (JBD) util commands",
			Hidden: true,
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List JBD jobs",
					Action: func(ctx *cli.Context) error {
						j, err := jobs.NewJobClient().FetchJobs()
						if err != nil {
							return err
						}
						for _, job := range j.Jobs {
							println(job.ID, job.Status)
						}
						return nil
					},
				},
			},
		},

		{
			Name:  "cs",
			Usage: "Codespace-related commands",
			Before: func(ctx *cli.Context) error {
				return klient.validateCluster()
			},
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "List all Codespace",
					Action: func(c *cli.Context) error {
						email := c.String("owner")
						if email == "" {
							e, err := LoadEmail()
							if err != nil {
								return cli.Exit(err.Error(), 1)
							}
							email = e
						}
						if err := klient.ListNamespacesWithEmail(email); err != nil {
							return cli.Exit(err.Error(), 1)
						}
						return nil
					},
				},
				{
					Name:  "create",
					Usage: "Create a codespace",
					Action: func(c *cli.Context) error {
						email, err := LoadEmail()
						if err != nil {
							return cli.Exit(err.Error(), 1)
						}
						name := c.Args().Get(0)
						if err := klient.CreateNamespace(name, email); err != nil {
							return cli.Exit(err.Error(), 1)
						}
						fmt.Println("Created codespace:", name)
						return nil
					},
				},
				{
					Name:  "desc",
					Usage: "Describe a codespace",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							return cli.Exit("arg 1: 'name' is required", 1)
						}
						_, err := klient.selectResource(name)
						return err
					},
				},
			},
		},
	}

	app := &cli.App{
		Name:     "condo",
		Usage:    "CLI tool to manage codespaces",
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// LoadEmail tries to retrieve an email from ~/.quorumrc or ~/.gitconfig
func LoadEmail() (string, error) {
	paths := []string{"~/.condorc", "~/.gitconfig"}

	for _, path := range paths {
		configPath := os.ExpandEnv(strings.Replace(path, "~", "$HOME", 1))
		cfg, err := ini.Load(configPath)
		if err == nil {
			email := cfg.Section("user").Key("email").String()
			if email != "" {
				return email, nil
			}
		}
	}

	return "", fmt.Errorf("email not found in any configuration files")
}
