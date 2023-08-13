package main

import (
	"fmt"
	"github.com/daog1/envcli"
	"github.com/urfave/cli/v2"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"log"
	"os"
	"strings"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}
	app := &cli.App{
		Name:  "envcli",
		Usage: "edit .env",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Value:   ".env",
				Aliases: []string{"f"},
				Usage:   "-f filename",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "read",
				Aliases: []string{"r"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "keys",
						Aliases: []string{"k"},
						Usage:   "-k a,b,c",
					},
				},
				Usage: "read key",
				Action: func(cCtx *cli.Context) error {
					filename := cCtx.String("file")
					keystrs := cCtx.String("keys")
					keys := strings.Split(keystrs, ",")
					file, err := os.Open(filename)
					if err != nil {
						return nil
					}
					defer file.Close()
					envkeys, err := envcli.Parse(file)
					if err != nil {
						return nil
					}
					for _, key := range keys {
						value, b := envkeys.Get(key)
						if b {
							println(value)
						} else {
						}
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add key=value",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "-k a,b,c",
					},
				},
				Action: func(cCtx *cli.Context) error {
					filename := cCtx.String("file")
					var envkeys *orderedmap.OrderedMap[string, string]
					{
						file, err := os.Open(filename)
						if err != nil {
							return nil
						}
						defer file.Close()
						envkeys, err = envcli.Parse(file)
						if err != nil {
							return nil
						}
					}
					keystr := cCtx.String("key")
					keys := strings.Split(keystr, "=")
					envkeys.Set(keys[0], keys[1])
					envcli.Write(envkeys, filename)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list key=value",
				Action: func(cCtx *cli.Context) error {
					filename := cCtx.String("file")
					var envkeys *orderedmap.OrderedMap[string, string]
					{
						file, err := os.Open(filename)
						if err != nil {
							return nil
						}
						defer file.Close()
						envkeys, err = envcli.Parse(file)
						if err != nil {
							return nil
						}
					}
					for pair := envkeys.Oldest(); pair != nil; pair = pair.Next() {
						fmt.Printf("%s=%s\n", pair.Key, pair.Value)
					}
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
