package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// LoadEmail tries to retrieve an email from ~/.quorumrc or ~/.gitconfig
func LoadEmail() (string, error) {
	paths := []string{"~/.quorumrc", "~/.gitconfig"}

	for _, path := range paths {
		configPath := os.ExpandEnv(strings.Replace(path, "~", "$HOME", 1))
		cfg, err := ini.Load(configPath)
		if err == nil {
			email := cfg.Section("").Key("email").String()
			if email != "" {
				return email, nil
			}
		}
	}

	return "", fmt.Errorf("email not found in any configuration files")
}

// getKubeClient initializes a Kubernetes client
func getKubeClient() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// ListNamespaces retrieves and prints all namespaces
func ListNamespaces() error {
	clientset, err := getKubeClient()
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %v", err)
	}

	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name)
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "quorum",
		Usage: "CLI tool to manage namespaces and applications in Kubernetes",
		Commands: []*cli.Command{
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
			{
				Name:  "ns",
				Usage: "Namespace-related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List all namespaces",
						Action: func(c *cli.Context) error {
							if err := ListNamespaces(); err != nil {
								return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
