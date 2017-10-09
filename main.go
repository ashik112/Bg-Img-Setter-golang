package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// Desktop contains the current desktop environment on Linux.
// Empty string on all other operating systems.
var Desktop = os.Getenv("XDG_CURRENT_DESKTOP")

// ErrUnsupportedDE is thrown when Desktop is not a supported desktop environment.
var ErrUnsupportedDE = errors.New("your desktop environment is not supported")

type Image struct {
	name string
	url  string
}

func NewImage(link string) *Image {
	image := &Image{
		name: getImgName(link),
		url:  link,
	}
	return image
}

func getImgName(link string) string {
	//Split the URL
	splitURL := strings.Split(link, "/")
	//Get Divided URL length
	length := len(splitURL)
	//Get image name and type
	img := splitURL[length-1]
	//Trim Whitespaces
	img = strings.TrimSpace(img)
	return img
}

func getImgLink() string {
	//Take url from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a image url: ")
	link, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err.Error())
	}
	return link
}

// Get image for url using HTTP
func getImg(img *Image) *http.Response {
	response, err := http.Get(img.url)
	if err != nil {
		log.Fatal(err.Error())
	}
	return response
}

func createFile(img *Image) *os.File {
	//open a file for writing
	file, err := os.Create(img.name)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file
}

func saveImg(file *os.File, response *http.Response) {
	//User io.Copy to dump the response body to the file. Supports big files too
	_, err := io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	//Get Image Link from User Input
	link := getImgLink()
	//Create Img object
	img := NewImage(link)
	//Get Image from HTTP response
	response := getImg(img)
	//Check if Image available
	if response.StatusCode == 404 {
		fmt.Println("Status Code:", 404, ", Image not found")
	}
	fmt.Println("Downloading")
	// Create file to store image
	file := createFile(img)
	//Save image to file
	saveImg(file, response)
	// close response Body after finishing all task
	defer response.Body.Close()
	//close file afte finishing all task
	defer file.Close()
	fmt.Println("Successfully Downloaded")
	choice := ""
	for strings.Compare(choice, "y") != 0 || strings.Compare(choice, "n") != 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to set it as dekstop background? (y/n) ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		choice = strings.ToLower(strings.Trim(choice, " \r\n"))
		switch choice {
		case "y":
			//fmt.Println("yes")
			error := SetFromFile(basepath + "/" + img.name)
			if error != nil {
				fmt.Println(error.Error())

			} else {
				fmt.Println("Dekstop Background Set")
				return
			}

		case "n":
			return
		default:
			fmt.Println("invalid input")
		}
	}

}
