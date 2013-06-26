package main

import (
  "os"
  "net/http"
  "github.com/mreiferson/go-httpclient"
  "time"
  "log"
  "io"
  "fmt"
  "io/ioutil"
  "strings"
)

func main(){

  var TimeoutTransport = &httpclient.Transport{
    ConnectTimeout:        1*time.Second,
    RequestTimeout:        1*time.Second,
    ResponseHeaderTimeout: 1*time.Second,
  }

  url := "http://devimages.apple.com/iphone/samples/bipbop/gear4/prog_index.m3u8"
  destination := "./file.ts"

  out, err := os.Create(destination)
  defer out.Close()

  client := &http.Client{Transport: TimeoutTransport}
  response, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't download m3u8 url: %s\n", url)
		os.Exit(0)
	}
	defer response.Body.Close()

  contents, err := ioutil.ReadAll(response.Body)
  if err != nil { panic(err) }

  m3u8content := string(contents)
	m3u8Lines := strings.Split(strings.TrimSpace(m3u8content), "\n")

	if m3u8Lines[0] != "#EXTM3U" {
		log.Fatal("not a valid m3u8 file")
		os.Exit(0)
	}

  var segmentUrls []string
	for i, value := range m3u8Lines {
		// trim each line
		m3u8Lines[i] = strings.TrimSpace(value)
		if m3u8Lines[i] != "" && !strings.HasPrefix(m3u8Lines[i], "#") {
			segmentUrls = append(segmentUrls, "http://devimages.apple.com/iphone/samples/bipbop/gear4/" + m3u8Lines[i])
		}
	}

  for _, url := range segmentUrls {
    log.Println("downloading ", url)
    resp, err := client.Get(url)
    if err != nil {
      log.Fatal(err)
    } else {
      _, err := io.Copy(out, resp.Body)
      if err != nil { log.Fatal(err) }
    }
    resp.Body.Close()
  }

}
