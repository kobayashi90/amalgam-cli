package main

import (
	"amalgamDCLoader/pkg/amalgam"
	dch "amalgamDCLoader/pkg/detektivconanch"
	"amalgamDCLoader/pkg/models"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
)

func CmdApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Amalgam CLI"
	app.Usage = "Small CLI Download Tool written in Go to Download Episodes from amalgam-fansubs.moe or Music from detektiv-conan.ch"
	app.Version = "0.3.4"

	app.Commands = []cli.Command{
		{
			Name:  "episodes",
			Usage: "list and download episodes from amalgam",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list available episodes",
					Action:  ListEpisodes,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:     "dlink,d",
							Usage:    "List episodes with download links",
							Required: false,
							Hidden:   false,
						},
						cli.BoolFlag{
							Name:     "gdrive,g",
							Usage:    "Show if episodes can be downloaded via google drive",
							Required: false,
							Hidden:   false,
						},
						cli.StringFlag{
							Name:     "format",
							Usage:    "available values: csv, html, md",
							Required: false,
							Hidden:   false,
							Value:    "",
						},
					},
				},
				{
					Name:      "download",
					Aliases:   []string{"d"},
					Usage:     "download episodes",
					ArgsUsage: "episode list: 1 2 3  episode range: 4-10, combined: 1 2-6 8",
					Action:    DownloadEpisodes,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:     "gdrive,g",
							Usage:    "Download episode from google drive",
							Required: false,
							Hidden:   false,
						},
						cli.BoolFlag{
							Name:     "dlink,d",
							Usage:    "Download episode via direct amalgam link. This is deprecated an now replaced via default m3u8 download.",
							Required: false,
							Hidden:   false,
						},
						cli.StringFlag{
							Name:     "output-dir,o",
							Usage:    "Specify the directory where to the downloaded episodes shell be stored",
							Required: false,
							Hidden:   false,
							Value:    "",
						},
					},
				},
			},
		},
		{
			Name:  "music",
			Usage: "list and download music from detektiv-conan.ch",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Action:  ListMusic,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "format",
							Usage:    "available values: csv, html, md",
							Required: false,
							Hidden:   false,
							Value:    "",
						},
					},
				},
				{
					Name:      "download",
					Aliases:   []string{"d"},
					Action:    DownloadMusic,
					ArgsUsage: "music id list: 1 2 3  music id range: 4-10, combined: 1 2-6 8",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:     "unzip,u",
							Usage:    "directly unzip the downloaded in file into a new directory",
							Required: false,
							Hidden:   false,
						},
						cli.BoolFlag{
							Name:     "keepArchive,k",
							Usage:    "keep archive after extraction",
							Required: false,
							Hidden:   false,
						},
						cli.StringFlag{
							Name:     "output-dir,o",
							Usage:    "Specify the directory where to the downloaded episodes shell be stored",
							Required: false,
							Hidden:   false,
							Value:    "",
						},
					},
				},
			},
		},
	}

	return app
}

func handleOutputDir(outputDirFlagVal string) (string, error) {
	var outpuDir string
	// if outputDirectory is not specified, set it to the current working directory
	if outputDirFlagVal == "" {
		wcdir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		outpuDir = wcdir
	} else {
		// check if outputDirectory exists, if not create it
		if _, err := os.Stat(outputDirFlagVal); os.IsNotExist(err) {
			err = os.Mkdir(outputDirFlagVal, 0755)
			if err != nil {
				return "", err
			}
		}
		outpuDir = outputDirFlagVal
	}
	return outpuDir, nil
}

func parseArgsList(args cli.Args) ([]string, error) {
	var argList []string
	for _, s := range args {
		if strings.Contains(s, "-") {
			// handle ranges
			splitted := strings.Split(s, "-")
			start, err := strconv.Atoi(splitted[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(splitted[1])
			if err != nil {
				return nil, err
			}
			for i := start; i <= end; i++ {
				argList = append(argList, strconv.Itoa(i))
				// handle ,5 episodes (there are episodes with numbers like 704,5)
				//if _, ok := episodes[fmt.Sprintf("%v,5", i)]; ok {
				//	episodeArgList = append(episodeArgList, fmt.Sprintf("%v,5", i))
				//}
			}
		} else {
			// handle simple comma separation
			argList = append(argList, s)
		}
	}
	return argList, nil
}

func ListEpisodes(c *cli.Context) error {
	episodes, err := amalgam.FetchEpisodes()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)

	if c.Bool("dlink") {
		t.AppendHeader(table.Row{"Nr.", "Title", "Download Link", "Google Drive Link"})
	} else if c.Bool("gdrive") {
		t.AppendHeader(table.Row{"Nr.", "Title", "Google Drive"})
	} else {
		t.AppendHeader(table.Row{"Nr.", "Title"})
	}

	for _, e := range episodes {
		gdAvailable := "✓"
		if e.GDriveLink == "" {
			gdAvailable = "✘"
		}

		if strings.Contains(e.EpisodeNr, ",") {
			e.Title = fmt.Sprintf("%v (Combined Episode)", e.Title)
		}

		// Skip if available flag and episode is not downloadable
		if c.Bool("gdrive") && e.GDriveLink == "" {
			continue
		}

		if c.Bool("dlink") {
			t.AppendRow(table.Row{e.EpisodeNr, e.Title, e.DownloadLink, e.GDriveLink})
		} else if c.Bool("gdrive") {
			t.AppendRow(table.Row{e.EpisodeNr, e.Title, gdAvailable})
		} else {
			t.AppendRow(table.Row{e.EpisodeNr, e.Title})
		}

	}

	t.AppendFooter(table.Row{fmt.Sprintf("Total: %v", len(episodes))})

	if c.String("format") == "csv" {
		t.RenderCSV()
	} else if c.String("format") == "html" {
		t.RenderHTML()
	} else if c.String("format") == "md" {
		t.RenderMarkdown()
	} else {
		t.Render()
	}

	return nil
}

func DownloadEpisodes(c *cli.Context) error {
	var episodeArgList []string
	episodesArgs := c.Args()

	// Fetch episodes and create map for easier download
	episodesList, err := amalgam.FetchEpisodes()
	if err != nil {
		return err
	}
	episodes := make(map[string]*models.Episode)
	for _, e := range episodesList {
		episodes[e.EpisodeNr] = e
	}

	episodeArgList, err = parseArgsList(episodesArgs)
	if err != nil {
		return err
	}

	if c.Bool("gdrive") {
		fmt.Println("Downloading Episodes:", strings.Join(episodeArgList, " "), "via google drive")
	} else {
		fmt.Println("Downloading Episodes:", strings.Join(episodeArgList, " "))
	}

	outputDirectory, err := handleOutputDir(c.String("output-dir"))
	if err != nil {
		return err
	}

	fmt.Println()

	// download episodes
	for _, episodeNr := range episodeArgList {
		// check if episode is available
		episode, ok := episodes[episodeNr]
		if !ok {
			fmt.Printf("Episode %v is not available\nPlease checkout https://amalgam-fansubs.moe/ for more informations.\n", episodeNr)
			continue
		}

		if c.Bool("gdrive") && episode.GDriveLink == "" {
			fmt.Printf("Episode %v is not available via Google Drive\n", episodeNr)
			continue
		}

		if c.Bool("gdrive") {
			err = amalgam.DownloadEpisodeFromGDrive(episode, outputDirectory)
		} else if c.Bool("dlink") {
			err = amalgam.DownloadEpisode(episode, outputDirectory)
		} else {
			err = amalgam.DownloadEpisodeFromM3U8(episode, outputDirectory)

			if err != nil {
				fmt.Println("Could not download episode via m3u8.\nRetrying with normal download link...")
				err = amalgam.DownloadEpisode(episode, outputDirectory)
			}
		}

		if err != nil {
			fmt.Printf("Episode %v could not be downloaded, %v\n", episodeNr, err)
		}

		fmt.Println()
	}

	return nil
}

func ListMusic(c *cli.Context) error {
	musics, err := dch.FetchMusic()
	if err != nil {
		return err
	}
	fmt.Println(musics[0].DownloadLink)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.AppendHeader(table.Row{"ID", "Title"})

	for i, m := range musics {
		t.AppendRow(table.Row{i + 1, m.Title})
	}

	t.AppendFooter(table.Row{fmt.Sprintf("Total: %v", len(musics))})

	if c.String("format") == "csv" {
		t.RenderCSV()
	} else if c.String("format") == "html" {
		t.RenderHTML()
	} else if c.String("format") == "md" {
		t.RenderMarkdown()
	} else {
		t.Render()
	}

	return nil
}

func DownloadMusic(c *cli.Context) error {
	var musicArgList []string
	episodesArgs := c.Args()

	outputDirectory, err := handleOutputDir(c.String("output-dir"))
	if err != nil {
		return err
	}

	// Fetch episodes and create map for easier download
	musicList, err := dch.FetchMusic()
	if err != nil {
		return err
	}

	musicArgList, err = parseArgsList(episodesArgs)
	if err != nil {
		return err
	}

	for _, musicIndex := range musicArgList {
		index, err := strconv.Atoi(musicIndex)
		if err != nil {
			fmt.Printf("Could not download music on index %v\n", musicIndex)
		}

		index -= 1 // shift index back (input indices are from 1-end)
		err = dch.DownloadMusic(musicList[index], c.Bool("unzip"), c.Bool("keepArchive"), outputDirectory)

		if err != nil {
			fmt.Printf("Could not download music %v\n%v\n", musicList[index].Title, err)
		}
	}

	return nil
}
