package main

import (
	"net/http"
	"log"
	"fmt"
	"os"
	"io"
	"strings"
	"sync"
	"strconv"
	"math/rand"
)

func getResponse(url string) *http.Response {
	// open url
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func getImageName(response *http.Response) string {
	// get num of image.  It may be random number
	query := strings.Split(response.Request.URL.RawQuery, "=")
	var imageNum string
	if len(query) < 2 {
		imageNum = strconv.Itoa(rand.Intn(10000))
	} else {
		imageNum = query[1]
	}
	return imageNum + ".jpg"
}

func saveResponseBodyToFile(reader io.Reader, path string) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatal(err)
	}
}

func saveImgFromUrl(url, dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	// open url
	response := getResponse(url)
	defer response.Body.Close()

	// open new file for write
	imagePath := dir + getImageName(response)

	saveResponseBodyToFile(response.Body, imagePath)

	fmt.Println(imagePath)
}

func main() {
	url := "https://picsum.photos/1920/1080/?random"
	dir := "./downloads/"
	countImages := 100

	var wg sync.WaitGroup

	for i:=0; i<countImages; i++{
		wg.Add(1)
		go saveImgFromUrl(url, dir, &wg)
	}

	wg.Wait()
}
