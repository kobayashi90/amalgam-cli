package gdrive

import (
	"amalgamDCLoader/pkg/web"
	"fmt"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetConfirmCodeAndCookies(exportUrl string) ([]*http.Cookie, string, error) {
	resp, err := http.Get(exportUrl)
	if err != nil {
		return nil, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, "", err
	}

	ucHtmlLink := htmlquery.FindOne(doc, "//a[@id='uc-download-link']")
	ucLink := htmlquery.SelectAttr(ucHtmlLink, "href")

	ucConfirmCode := strings.Split(strings.Split(ucLink, "&")[1], "=")[1]

	return resp.Cookies(), ucConfirmCode, nil
}

func GetTotalFileSize(exportUrl string) (int, error) {
	doc, err := htmlquery.LoadURL(exportUrl)
	if err != nil {
		return -1, err
	}
	sizeSpan := htmlquery.FindOne(doc, "//span[@class = 'uc-name-size']")
	sizeSpanText := htmlquery.InnerText(sizeSpan)

	r, err := regexp.Compile("\\(([0-9]+.?[0-9]*)([kKmMgG])\\)")
	if err != nil {
		return -1, err
	}
	matches := r.FindStringSubmatch(sizeSpanText)
	// fmt.Printf("Matches: %v\n", matches)

	size, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return -1, err
	}
	sizeEntity := strings.ToLower(matches[2])
	switch sizeEntity {
	case "g":
		size *= 1024 * 1024 * 1024
	case "m":
		size *= 1024 * 1024
	case "k":
		size *= 1024
	}

	return int(size), err
}

func GdriveDownload(url, filePath string) error {
	// get file id
	splitted := strings.Split(url, "/")
	fileId := splitted[5]

	exportUrl := fmt.Sprintf("https://drive.google.com/uc?id=%v&export=download", fileId)

	cookies, confirmCode, err := GetConfirmCodeAndCookies(exportUrl)
	if err != nil {
		return err
	}

	totalFileSize, err := GetTotalFileSize(exportUrl)
	if err != nil {
		return err
	}

	confirmUrl := fmt.Sprintf("%v&confirm=%v", exportUrl, confirmCode)

	err = DownloadFile(filePath, confirmUrl, cookies, totalFileSize)

	return err
}

func DownloadFile(filepath string, url string, cookies []*http.Cookie, totalFileSize int) error {
	client := http.Client{}

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Limit totalFileSize lowest value to 0
	if totalFileSize < 0 {
		totalFileSize = 0
	}

	// Write the body to file
	counter := &web.WriteCounter{Filename: filepath, Total: uint64(totalFileSize)}

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	return err
}
