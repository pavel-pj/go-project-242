package main

import (
	si "code/sizeinformer"
	"fmt"
)

func main() {
	/*
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
	*/
	size, err := si.GetSize("/var/www/go1/proj1/go-project-242/testdata/dir200")
	//size, err := si.GetSize("/tmp/file2.pdf")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(size)

}
