package models

import (
	"fmt"
	"strings"
)

type Episode struct {
	Title        string
	EpisodeNr    string
	DownloadLink string
	M3U8Link     string
	GDriveLink   string
	Note         string
	Language     string
	Mangas       []int
}

func NewEmptyEpisode() *Episode {
	return &Episode{
		Title:        "",
		EpisodeNr:    "",
		DownloadLink: "",
		GDriveLink:   "",
		Note:         "",
		Language:     "",
		Mangas:       make([]int, 0),
	}
}

func (e *Episode) GetFilename() string {
	episodeTitle := strings.ReplaceAll(e.Title, "-", ".")
	episodeTitle = strings.ReplaceAll(episodeTitle, "(", "")
	episodeTitle = strings.ReplaceAll(episodeTitle, ")", "")
	episodeTitle = strings.ReplaceAll(episodeTitle, " ", ".")
	episodeTitle = strings.ReplaceAll(episodeTitle, "..", ".")
	return fmt.Sprintf("Detektiv.Conan.E%04s-%v.mp4", e.EpisodeNr, episodeTitle)
}

type Music struct {
	Title        string
	Filename     string
	DownloadLink string
}
