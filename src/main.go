package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type UpdockConfig struct {
	dryRun    bool
	accept    bool
	inputFile string
}

type UpdockCli struct {
	config UpdockConfig
	services []ComposeService
}

type ComposeService struct {
	name string
	tag string
	latest string
	registry string
}

func (compose ComposeService) String() string {
	return fmt.Sprintf("ComposeService{ name: %s, tag: %s, latest: %s, registry: %s}",
		compose.name,
		compose.tag,
		compose.latest,
		compose.registry)
}

type ComposeConfig struct {
	Version string `yaml:"version"`
	Services map[string]ComposeConfigService `yaml:"services"`
}

type ComposeConfigService struct {
	Image string `yaml:"image"`
}

func (cli UpdockCli) String() string {
	svcs := make([]string, len(cli.services))
	for i, svc := range cli.services {
		svcs[i] = svc.String()
	}
	return fmt.Sprintf("UpdockCli{services: %s}",
		strings.Join(svcs, ", "))
}

func (cli *UpdockCli) readComposeConfig(compose ComposeConfig) {
	for _, service := range compose.Services {
		nameParts := strings.Split(service.Image, ":")
		name := nameParts[0]
		tag := nameParts[1]
		reg := "docker.io"

		slashCount := strings.Count(name, "/")
		switch slashCount {
		case 2:
			names := strings.Split(name, "/")
			reg = names[0]
			name = names[1] + "/" + names[2]
		}
		cli.services = append(cli.services, ComposeService{
			name: name,
			tag: tag,
			registry: reg,
		})
	}
}

func (cli *UpdockCli) getLatestTag() {
	for _, svc := range cli.services {
		fmt.Printf("%s\t%s\t\t:: %s\n", svc.name, svc.tag, svc.registry)
	}
}

func main() {
	dryRun := flag.Bool("d", false, "Dry run")
	accept := flag.Bool("y", false, "Accept changes")
	inputFile := flag.String("f", "podman-compose.yaml", "Input file")

	flag.Parse()

	cmd := UpdockConfig{
		dryRun:    *dryRun,
		accept:    *accept,
		inputFile: *inputFile,
	}

	file, err := os.ReadFile(cmd.inputFile); if err != nil {
		fmt.Println("Error reading file", err)
		return
	}

	var compose ComposeConfig
	err = yaml.Unmarshal(file, &compose); if err != nil {
		fmt.Println("Error unmarshalling file", err)
		return
	}

	updock := UpdockCli{
		config: cmd,
		services: make([]ComposeService, 0),
	}

	updock.readComposeConfig(compose)

	updock.getLatestTag()
}