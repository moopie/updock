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
}

func (compose ComposeService) String() string {
	return fmt.Sprintf("ComposeService{ name: %s, tag: %s, latest: %s }",
		compose.name,
		compose.tag,
		compose.latest)
}

type ComposeConfig struct {
	services map[string]struct {
		image string `yaml:"image"`
	} `yaml:"services"`
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
	for _, service := range compose.services {
		s := strings.Split(service.image, ":")
		name := s[0]
		tag := s[1]

		cli.services = append(cli.services, ComposeService{
			name: name,
			tag: tag,
			latest: "",
		})
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

	fmt.Println(compose)

	updock := UpdockCli{
		config: cmd,
		services: make([]ComposeService, len(compose.services)),
	}

	updock.readComposeConfig(compose)

	fmt.Println(updock)
}