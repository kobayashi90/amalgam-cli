package detektivconanch

import (
	"amalgamDCLoader/pkg/models"
	"amalgamDCLoader/pkg/web"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/mholt/archiver"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	LANGUAGE_DE          = "DE"
	LANGUAGE_JP          = "JP"
	GermanStartEpisode   = 0
	JapaneseStartEpisode = 400
	BaseUrl              = "https://www.detektiv-conan.ch/episodeOverviewPage/showPageTableContent"
)

func fetchEpisodePage(startNumber string, language string) ([]byte, error) {
	form := url.Values{}
	form.Add("startNumber", startNumber)
	form.Add("type", language)
	req, err := http.NewRequest("POST", BaseUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("returned status code was not 200 but %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, err
}

func parseEpisodeTableRow(tr *html.Node) (*models.Episode, error) {
	episode := models.NewEmptyEpisode()
	tds := htmlquery.Find(tr, "//td")
	if len(tds) < 4 {
		return nil, fmt.Errorf("could not parse episode number, title and manga")
	}
	// parse episode number
	episodeNr := strings.TrimSpace(htmlquery.InnerText(tds[0]))
	episode.EpisodeNr = episodeNr

	// parse episode title
	titleLink := htmlquery.FindOne(tds[2], "//a")
	episode.Title = strings.TrimSpace(htmlquery.InnerText(titleLink))

	// parse manga
	mangaLinks := htmlquery.Find(tds[3], "//a")
	if len(mangaLinks) == 0 {
		episode.Note = strings.TrimSpace(htmlquery.InnerText(tds[3]))
	} else if len(mangaLinks) == 1 {
		mangaStr := strings.TrimSpace(htmlquery.InnerText(mangaLinks[0]))
		mangaNrStr := strings.Split(mangaStr, " ")[1]
		mangaNr, err := strconv.Atoi(mangaNrStr)
		if err != nil {
			return nil, err
		}
		episode.Mangas = []int{mangaNr}

	} else if len(mangaLinks) > 1 {
		mangaStart, err := strconv.Atoi(htmlquery.InnerText(mangaLinks[0]))
		if err != nil {
			return nil, err
		}
		mangaEnd, err := strconv.Atoi(htmlquery.InnerText(mangaLinks[1]))
		if err != nil {
			return nil, err
		}
		episode.Mangas = []int{mangaStart, mangaEnd}
	}

	episode.Language = "de"
	return episode, nil
}

func getEpisodesInfos(language string) ([]*models.Episode, error) {
	var lastEpisode int
	episodesFound := true
	if language == LANGUAGE_DE {
		lastEpisode = GermanStartEpisode
	} else if language == LANGUAGE_JP {
		lastEpisode = JapaneseStartEpisode
	} else {
		return nil, fmt.Errorf("wrong language, supported are %q and %q", LANGUAGE_DE, LANGUAGE_JP)
	}
	episodes := make([]*models.Episode, 0)

	for episodesFound {

		body, err := fetchEpisodePage(fmt.Sprintf("%d", lastEpisode+1), language)
		if err != nil {
			return nil, err
		}
		bodyStr := string(body)
		bodyStr = fmt.Sprintf("<html><head></head><body><table>\n%s\n</table></body></html>", bodyStr)

		doc, err := htmlquery.Parse(strings.NewReader(bodyStr))
		if err != nil {
			return nil, err
		}
		// get all tr with data-episode-id

		episodeTrs := htmlquery.Find(doc, "//tr[@data-episode-id]")
		if len(episodeTrs) == 0 {
			episodesFound = false
		} else {
			// parse and add episodes to list
			for _, tr := range episodeTrs {
				ep, err := parseEpisodeTableRow(tr)
				if err != nil {
					return nil, err
				}
				episodes = append(episodes, ep)
			}
			lastEpisode, _ = strconv.Atoi(episodes[len(episodes)-1].EpisodeNr)
		}

	}

	return episodes, nil
}

func GetGermanEpisodeInfos() ([]*models.Episode, error) {
	return getEpisodesInfos(LANGUAGE_DE)
}

func GetJapaneseEpisodeInfos() ([]*models.Episode, error) {
	return getEpisodesInfos(LANGUAGE_JP)
}

func GetAllEpisodeInfos() ([]*models.Episode, error) {
	episodesDE, err := GetGermanEpisodeInfos()
	if err != nil {
		return nil, err
	}
	episodesJP, err := GetJapaneseEpisodeInfos()
	if err != nil {
		return nil, err
	}
	episodesDE = append(episodesDE, episodesJP...)
	return episodesDE, nil
}

func FetchMusic() ([]*models.Music, error) {
	musicUrl := "https://www.detektiv-conan.ch/index.php?page=aplayer/musik.php"
	doc, err := htmlquery.LoadURL(musicUrl)
	if err != nil {
		return nil, err
	}

	musicDivs := htmlquery.Find(doc, "//div[@class='album_content']")

	var musics []*models.Music
	for _, musicDiv := range musicDivs {
		htmlLink := htmlquery.FindOne(musicDiv, "//a")
		relativeDownloadLink := htmlquery.SelectAttr(htmlLink, "href")
		if !strings.Contains(relativeDownloadLink, ".zip") {
			continue
		}
		downloadLink := fmt.Sprintf("https://www.detektiv-conan.ch%v", relativeDownloadLink)

		filename := filepath.Base(relativeDownloadLink)
		title := strings.TrimSuffix(filename, ".zip")

		musics = append(musics, &models.Music{
			Title:        title,
			Filename:     filename,
			DownloadLink: downloadLink,
		})
	}

	return musics, nil
}

func DownloadMusic(music *models.Music, unzip, keepArchive bool, outputDir string) error {

	fp := fmt.Sprintf("%v/%v", outputDir, music.Filename)

	err := web.DownloadFile(fp, music.DownloadLink)

	if unzip {
		// create directory for extraction
		archivePath := filepath.Dir(fp)
		extractionDirPath := fmt.Sprintf("%s/%s/", archivePath, strings.TrimSuffix(music.Filename, ".zip"))
		err = os.Mkdir(extractionDirPath, 0775)
		if err != nil {
			return err
		}

		err = archiver.Unarchive(fp, extractionDirPath)
		if err != nil {
			return err
		}

		if !keepArchive {
			err = os.Remove(fp)
			if err != nil {
				return err
			}
		}
	}

	return err
}
