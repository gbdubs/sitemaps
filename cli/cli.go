package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gbdubs/sitemaps"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Sitemap Downloader",
		Usage:   "A CLI for downloading a sitemap to a local cache.",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Usage: "the location of the sitemap to download",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Whether to print a description of the sitemap, if the command succeeds.",
			},
		},
		Action: func(c *cli.Context) error {
			if c.String("url") == "" {
				return errors.New("url must be provided")
			}
			url := c.String("url")
			v := c.Bool("verbose")
			s, err := sitemaps.GetSitemapFromURL(url)
			if err != nil {
				return err
			}
			if v {
				fmt.Printf("Sucessfully downloaded sitemap with %d entries at %s\n", len(s.LastUpdated), s.MemoPath())
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
