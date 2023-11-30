package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "calculator",
		Usage: "A simple calculator app",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add two numbers",
				Action: func(ctx *cli.Context) error {
					a := ctx.Int("a")
					b := ctx.Int("b")
					result := a + b
					fmt.Printf("Result of %d + %d is %d\n", a, b, result)
					return nil
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "a",
						Usage: "First number",
					},
					&cli.IntFlag{
						Name:  "b",
						Usage: "Second number",
					},
				},
			},
			{
				Name:    "subtract",
				Aliases: []string{"s"},
				Usage:   "Subtract two numbers",
				Action: func(ctx *cli.Context) error {
					a := ctx.Int("a")
					b := ctx.Int("b")
					result := a - b
					fmt.Printf("Result of %d - %d is %d\n", a, b, result)
					return nil
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "a",
						Usage: "First number",
					},
					&cli.IntFlag{
						Name:  "b",
						Usage: "Second number",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
