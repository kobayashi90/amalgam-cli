package amalgam

import (
	"amalgamDCLoader/pkg/gdrive"
	"amalgamDCLoader/pkg/models"
	"amalgamDCLoader/pkg/web"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"strings"
)

//const amalgamBaseUrl = "https://cdn.amalgam-fansubs.moe"

const amalgamBaseUrl = "https://ddl.amalgam-fansubs.moe/content/Conan/1080p"

func getEpisodeLink(epNr string) string {
	return fmt.Sprintf("%s/[Totto]DetektivConan-%s-RFCT-[1080p].mp4", amalgamBaseUrl, epNr)
}

func getM3u8Link(epNr string) string {
	return fmt.Sprintf("https://cdn.amalgam-fansubs.moe/detektiv-conan/%s/master.m3u8", epNr)
}

func DownloadEpisode(episode *models.Episode, outputDir string) error {
	filepath := fmt.Sprintf("%v/%s", outputDir, episode.GetFilename())
	err := web.DownloadFile(episode.DownloadLink, filepath)

	return err
}

func DownloadEpisodeFromM3U8(episode *models.Episode, outputDir string) error {
	filepath := fmt.Sprintf("%v/%s", outputDir, episode.GetFilename())
	tsTmpDir := fmt.Sprintf("%s/tsdir", outputDir)

	dlFilename, err := web.DownloadM3U8(episode.M3U8Link, tsTmpDir)
	if err != nil {
		return err
	}
	// Rename and move the resulting video.ts file
	err = os.Rename(dlFilename, filepath)
	if err != nil {
		return err
	}

	// delete the temporary download directory
	err = os.Remove(tsTmpDir)
	return err
}

func DownloadEpisodeFromGDrive(episode *models.Episode, outputDir string) error {
	// get current working directory
	wcdir, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpDir := os.TempDir()
	rarDownloadPath := fmt.Sprintf("%v/%v.rar", tmpDir, episode.EpisodeNr)
	extractionPath := fmt.Sprintf("%v/%v-extracted", tmpDir, episode.EpisodeNr)

	// Download rar archived video into /tmp directory
	err = gdrive.GdriveDownload(episode.GDriveLink, rarDownloadPath)
	if err != nil {
		return err
	}

	// extract rar archive in /tmp directory
	fmt.Printf("  --> Extracting %v to %v\n", rarDownloadPath, extractionPath)
	err = archiver.Unarchive(rarDownloadPath, extractionPath)

	// get video filename and rename it
	files, err := ioutil.ReadDir(extractionPath)
	videoName := files[0].Name()
	err = os.Rename(fmt.Sprintf("%v/%v", extractionPath, videoName), fmt.Sprintf("%v/%v", wcdir, episode.GetFilename()))
	if err != nil {
		return err
	}
	fmt.Println("  --> Copy video to your path")

	// remove rar file and extracted directory
	fmt.Println("  --> Removing temporary files")
	err = os.Remove(rarDownloadPath)
	err = os.RemoveAll(extractionPath)

	return err
}

func FetchEpisodes() ([]*models.Episode, error) {
	defaultBaseUrl := "https://amalgam-fansubs.moe"
	backupBaseUrl := "https://amalgamsubs.lima-city.de"
	var baseUrl string

	if web.CheckIfAvailable(defaultBaseUrl) {
		baseUrl = defaultBaseUrl
	} else {
		fmt.Printf("%v is not reachable right now, trying to use backup %v:\n", defaultBaseUrl, backupBaseUrl)
		baseUrl = backupBaseUrl
	}

	urls := []string{
		baseUrl + "/detektiv-conan/",
		baseUrl + "/detektiv-conan-2018/",
	}

	var episodes []*models.Episode
	for _, url := range urls {
		doc, err := htmlquery.LoadURL(url)
		if err != nil {
			return nil, err
		}

		conanDiv := htmlquery.FindOne(doc, "//div[@id='conan']")
		if conanDiv == nil {
			continue
		}
		episodeTable := htmlquery.FindOne(conanDiv, "table")
		if episodeTable == nil {
			continue
		}

		rows := htmlquery.Find(episodeTable, "//tr")

		for i := 1; i < len(rows); i++ {
			cols := htmlquery.Find(rows[i], "//td")
			if len(cols) <= 0 {
				continue
			}
			episodeNr := htmlquery.InnerText(cols[0]) // number
			episodeNr = strings.ReplaceAll(episodeNr, ".", "")
			episodeTitle := htmlquery.InnerText(cols[1]) // title

			gdriveLink := htmlquery.SelectAttr(cols[3].FirstChild, "href") // gdrive link
			if !strings.Contains(gdriveLink, "drive.google") {
				gdriveLink = ""
			}

			downloadLink := getEpisodeLink(episodeNr)
			m3u8Link := getM3u8Link(episodeNr)

			episodes = append(episodes, &models.Episode{
				Title:        episodeTitle,
				EpisodeNr:    strings.TrimSpace(episodeNr),
				M3U8Link:     m3u8Link,
				GDriveLink:   gdriveLink,
				DownloadLink: downloadLink,
			})
		}
	}

	return episodes, nil
}
