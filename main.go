package main

import (
	"flag"
	"fmt"
)

/*
TODO:
 1. Define the flags:
    - dry run
    - input file
    - auto-accept changes (--yes|-y)
 2. Get a list of services from the compose file
 3. Check what version of an the images is
 4. Check what version of the image is available
    - is it possible to get the repo used for the image?
 5. Ask user if they want to update the image
 6. Write changes to the compose file
*/

type UpdockConfig struct {
	dryRun    bool
	accept    bool
	inputFile string
}

type UpdockCli struct {
	config UpdockConfig
}

func (cli UpdockCli) String() string {
	return fmt.Sprintf("UpdockCli{dryRun: %t, accept: %t, inputFile: %s}",
		cli.config.dryRun,
		cli.config.accept,
		cli.config.inputFile)
}

func main() {
	dryRun := flag.Bool("d", false, "Dry run")
	accept := flag.Bool("y", false, "Accept changes")
	inputFile := flag.String("f", "podman-compose.yaml", "Input file")

	flag.Parse()

	fmt.Println("Hello, World from Updock!")

	cmd := UpdockConfig{
		dryRun:    *dryRun,
		accept:    *accept,
		inputFile: *inputFile,
	}

	updock := UpdockCli{
		config: cmd,
	}

	fmt.Println(updock)
}
