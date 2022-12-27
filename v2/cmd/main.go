package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	stdownloader "github.com/city404/simple-http-downloader/v2"
)

func main() {
	stdownloader.NewClient()
	client := &stdownloader.Client{
		UserAgent: "grab",
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
	req, err := stdownloader.NewRequest("/Users/herui/vscode/v6/simple-http-downloader/v2/test/xdd.dump", "https://www.golang-book.com/public/pdf/gobook.pdf")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	headReceived := false
	go func() {
		<-time.After(time.Microsecond)
		if headReceived {
			return
		}
		cancel()
	}()
	resp := client.Do(req.WithContext(ctx))
	headReceived = true
	if resp.Err() != nil {
		panic(resp.Err())
	}
	// fmt.Printf("  %v\n", resp.HTTPResponse.Status)
	t := time.NewTicker(2000 * time.Millisecond)
	defer t.Stop()
	lastDownloadBytes := int64(0)
Loop:
	for {
		select {
		case <-t.C:
			downloaded := resp.BytesComplete()
			if lastDownloadBytes == downloaded {
				cancel()
			}
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			t.Stop()
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)
}
