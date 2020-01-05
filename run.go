package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"pet-spotlight/http"
	"pet-spotlight/io"
	"strings"
)

// Run starts scrapping the description and the pictures of the specified dogs to the specified directory.
// First, it must match the specified dog names against all available dogs on the web page. When it finds a match
// it will grab the description of the dog and visit the dog's personal information page.
// On the personal page, it will download all images there are of the dog.
func Run(dogs []string, baseDirectory *string) {
	// Create the scrappers
	availableDogs := colly.NewCollector()
	dogPictures := availableDogs.Clone()

	// Save the current dog to use when downloading pictures
	var currentDog string
	// Handle when the page of all the available dogs is loaded
	availableDogs.OnHTML(".pet-link", func(e *colly.HTMLElement) {
		dogName := e.ChildText("h3")
		currentDog = strings.ToLower(dogName)
		dogMatch := isMatch(dogs, currentDog)
		// If a match then create dir and description.txt file
		if dogMatch {
			if err := io.MakeDir(*baseDirectory + "/" + currentDog); err != nil {
				log.Println(err)
				return
			}
			log.Println("Found", dogName)
			fullDescription := e.ChildText(".pet-description-full")
			// Remove the adoption fee part
			index := strings.Index(fullDescription, adoptionText)
			desc := fullDescription[:index]
			// Add the link for adopting
			desc += "\n"
			desc += `👇👇SUBMIT AN APPLICATION HERE: 👇👇
https://2babrescue.com/adoption-fees-info`
			descFile := *baseDirectory + "/" + currentDog + "/description.txt"
			if err := io.WriteFile(desc, descFile); err != nil {
				log.Println(err)
				return
			}
			// Get the link to the dog's page to download pictures
			link := e.Attr("href")
			if err := dogPictures.Visit(baseURL + link); err != nil {
				log.Println(err)
				return
			}
		}
	})

	// When the dog page is loaded, download pictures
	dogPictures.OnHTML("#oc-clients", func(e *colly.HTMLElement) {
		imageURLs := e.ChildAttrs(".pet-gallery-thumb", "data-pet-gallery-url")
		// Save all the images
		for index, imageURL := range imageURLs {
			imagePath := fmt.Sprintf("%s/%s/image-%d.png", *baseDirectory, currentDog, index)
			imageFile := fmt.Sprintf("image-%d.png", index)
			if err := http.DownloadImage(imageURL, imagePath, imageFile); err != nil {
				log.Println(err)
			}
		}
	})

	// Log when the request is being made
	availableDogs.OnRequest(func(request *colly.Request) {
		log.Println("Starting extraction...")
	})

	// Start scrapping
	if err := availableDogs.Visit(baseURL + twoBlondesURL); err != nil {
		log.Fatalln(err)
	}
}

func isMatch(dogs []string, currentDog string) bool {
	dogMatch := false
	for _, dog := range dogs {
		if strings.Contains(currentDog, dog) {
			dogMatch = true
			break
		}
	}
	return dogMatch
}