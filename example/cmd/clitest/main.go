package main

import (
	"fmt"
	"os"

	"github.com/plinx2/grepo/cli"
	"github.com/plinx2/grepo/example/internal"
)

func main() {
	api := internal.InitializeAPI()
	if err := cli.New(api, "clitest").Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
