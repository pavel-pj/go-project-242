package main

import (
	si "code/sizeinformer"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	cmd := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

	size := si.GetSize("/tmp/test2.txt")
	fmt.Println(size)

}
