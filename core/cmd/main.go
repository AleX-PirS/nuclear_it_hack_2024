package main

import (
	"fmt"
	"os"

	"github.com/AleX-PirS/nuclear_it_hack_2024/services/core"
	"github.com/urfave/cli/v2"
)

func main(){
	app := &cli.App{
		Name: "Graph atributor",
		Usage: "Road graph processing",

		
		Flags: []cli.Flag{
            &cli.StringFlag{
                Name:     "af",
                Usage:    "Specifies the name of atritube file",
                Required: true,
            },
            &cli.StringFlag{
                Name:     "gf",
                Usage:    "Specifies the name of geography file",
                Required: true,
            },
            &cli.StringFlag{
                Name:     "rf",
                Usage:    "Specifies the name of output file",
                Required: true,
            },
            &cli.IntFlag{
                Name:    "acc",
                Usage:   "Specifies the accuracy",
                Value:   15,
				Required: true,
            },
        },
		Action:  func(c *cli.Context) error {
			fmt.Println("Start utility")
			core.Serve(c.Int("acc"), c.String("af"), c.String("gf"), c.String("rf"))
			fmt.Println("Utility stopped")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
