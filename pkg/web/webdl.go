package web

import (
	"fmt"
	"github.com/canhlinh/hlsdl"
	"github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"strings"
)

type WriteCounter struct {
	Current  uint64
	Total    uint64
	Filename string
}

func NewWriteCounter(filename string, total int64) *WriteCounter {
	var utotal uint64
	if total <= 0 {
		utotal = 0
	} else {
		utotal = uint64(total)
	}

	return &WriteCounter{
		Current:  0,
		Total:    utotal,
		Filename: filename,
	}
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += uint64(n)

	wc.PrintProgress()

	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	if wc.Total == 0 {
		fmt.Printf("\rDownloading %s: %s", wc.Filename, humanize.Bytes(wc.Current))
	} else {
		fmt.Printf("\rDownloading %s: %s / %s (%v %%)", wc.Filename, humanize.Bytes(wc.Current), humanize.Bytes(wc.Total), wc.Current*100/wc.Total)
	}
}

func DownloadFile(url, filepath string) error {
	client := http.Client{}

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected response status: %v", resp.Status)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	totalFileSize := resp.ContentLength

	// Write the body to file
	counter := NewWriteCounter(filepath, totalFileSize)

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	return err
}

func DownloadM3U8(m3u8Link, directory string) (filename string, err error) {
	hlsDL := hlsdl.New(m3u8Link, directory, 1, true)
	if filename, err = hlsDL.Download(); err != nil {
		//fmt.Println("An error occurred while downloading.\nRetrying with lower worker count...")
		//hlsDL = hlsdl.New(m3u8Link, directory, 1, true)
		//if filename, err = hlsDL.Download(); err != nil {
		return filename, err
		//}
	}
	return filename, err
}
