package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AleX-PirS/nuclear_it_hack_2024/services/core"
	"github.com/urfave/cli"
)

func main(){
	app := &cli.App{
		Name: "Graph atributor",
		Usage: "Road graph processing",

		Commands: []cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start processing of two graphs",
				Action:  func(c *cli.Context) error {
					core.Serve(c.Int("n"), c.String("af"), c.String("gf"), c.String("rf"))
					log.Println("Start utility")
					return nil
				},
			},
		},
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
            &cli.Float64Flag{
                Name:    "n",
                Usage:   "Specifies the accuracy",
                Value:   15,
            },
        },
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
