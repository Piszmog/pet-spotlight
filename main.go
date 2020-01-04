package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	baseURL       = "https://www.petstablished.com"
	twoBlondesURL = "/organization/80925"
	adoptionText  = "Adoption fee includes the following"
)

func main() {
	dogsFlag := flag.String("d", "", "A comma separate list of availableDogs to extract")
	baseDirectory := flag.String("l", "", "Location to place data on the extracted dogs")
	flag.Parse()

	if len(*dogsFlag) == 0 {
		fmt.Println("Missing flag '-d'")
		flag.Usage()
		return
	} else if len(*baseDirectory) == 0 {
		fmt.Println("Missing flag '-l'")
		flag.Usage()
		return
	}

	if err := makeDir(*baseDirectory); err != nil {
		log.Fatalln(err)
	}

	dogs := strings.Split(*dogsFlag, ",")
	selectedDogs := make(map[string]bool)
	for _, dog := range dogs {
		selectedDogs[strings.ToLower(dog)] = true
	}

	availableDogs := colly.NewCollector()
	dogPictures := availableDogs.Clone()

	var currentDog string
	availableDogs.OnHTML(".pet-link", func(e *colly.HTMLElement) {
		dogName := e.ChildText("h3")
		currentDog = strings.ToLower(dogName)
		if selectedDogs[currentDog] {
			if err := makeDir(*baseDirectory + "/" + currentDog); err != nil {
				log.Println(err)
				return
			}
			log.Println("Found", dogName)
			fullDescription := e.ChildText(".pet-description-full")
			index := strings.Index(fullDescription, adoptionText)
			desc := fullDescription[:index]
			descFile := *baseDirectory + "/" + currentDog + "/description.txt"
			if err := ioutil.WriteFile(descFile, []byte(desc), os.ModePerm); err != nil {
				log.Println(err)
				return
			}
			link := e.Attr("href")
			if err := dogPictures.Visit(baseURL + link); err != nil {
				log.Println(err)
				return
			}
		}
	})

	dogPictures.OnHTML("#oc-clients", func(e *colly.HTMLElement) {
		imageURLs := e.ChildAttrs(".pet-gallery-thumb", "data-pet-gallery-url")
		for index, imageURL := range imageURLs {
			if err := saveImage(*baseDirectory, currentDog, imageURL, index); err != nil {
				log.Println(err)
			}
		}
	})

	availableDogs.OnRequest(func(request *colly.Request) {
		log.Println("Starting extraction...")
	})

	if err := availableDogs.Visit(baseURL + twoBlondesURL); err != nil {
		log.Fatalln(err)
	}
}

func makeDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func saveImage(baseDir string, currentDog string, imageURL string, imageNumber int) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Println(err)
	}
	defer closeResource(resp.Body)
	imagePath := fmt.Sprintf("%s/%s/image-%d.png", baseDir, currentDog, imageNumber)
	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer closeResource(f)
	if _, err = io.Copy(f, resp.Body); err != nil {
		return err
	}
	return nil
}

func closeResource(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}
