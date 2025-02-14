package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
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
	services := make([]string, len(cli.services))
	for i, service := range cli.services {
		services[i] = service.String()
	}
	return fmt.Sprintf("UpdockCli{services: %s}",
		strings.Join(services, ", "))
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
type DockerTagResponse struct {
	Results []struct {
		Name   string `json:"name"`
		Images []struct {
			Digest    string `json:"digest"`
			OS        string `json:"os"`
			Arch      string `json:"architecture"`
			CreatedAt string `json:"last_pushed"`
		} `json:"images"`
	} `json:"results"`
}

func (cli *UpdockCli) getLatestTag() {
	for i, svc := range cli.services {
		fmt.Printf("[%s] %s:%s\n", svc.registry, svc.name, svc.tag)

		// TODO: support other registries
		if svc.registry != "docker.io" {
			fmt.Println("!!~ Only docker.io is supported for now")
			continue
		}

		url := "https://hub.docker.com/v2/repositories/%s/tags/?page_size=50"
		fullUrl := fmt.Sprintf(url, svc.name)

		fmt.Printf("url: %s\n", fullUrl)

		res, err := http.Get(fullUrl); if err != nil {
			fmt.Printf("Error getting tags for %s: %v\n", svc.name, err)
			os.Exit(1)
		}

		body, err := io.ReadAll(res.Body); if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			os.Exit(1)
		}

		var tagData DockerTagResponse
		if err := json.Unmarshal(body, &tagData); err != nil {
			fmt.Println("Error decoding JSON:", err)
			os.Exit(1)
		}

		for _, tag := range tagData.Results {
			if tag.Name != svc.tag || len(tag.Images) == 0 {
				continue
			}

			fmt.Println("Tag Digest:", tag.Images[0].Digest)
			fmt.Println("Platform:", tag.Images[0].OS, tag.Images[0].Arch)
			fmt.Println("Created At:", tag.Images[0].CreatedAt)

			latestDigest := tag.Images[0].Digest

			if tag.Name == "latest" {
				fmt.Printf("Looking up for latest tag\n")
				// Step 2: Find another tag with the same digest
				for _, secondTag := range tagData.Results {
					if secondTag.Name != "latest" &&
					   len(secondTag.Images) > 0 &&
					   secondTag.Images[0].Digest == latestDigest {
						fmt.Println("The 'latest' tag points to version:", secondTag.Name)
						cli.services[i].latest = secondTag.Name
					}
				}
			}
		}
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
		fmt.Printf("Error reading file %s: %v\n", cmd.inputFile, err)
		return
	}

	var compose ComposeConfig
	err = yaml.Unmarshal(file, &compose); if err != nil {
		fmt.Println("Error decoding YAML", err)
		return
	}

	updock := UpdockCli{
		config: cmd,
		services: make([]ComposeService, 0),
	}

	updock.readComposeConfig(compose)

	updock.getLatestTag()
}