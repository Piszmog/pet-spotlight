package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"pet-spotlight/http"
	"pet-spotlight/io"
	"strings"
)

const (
	baseURL                = "https://www.petstablished.com"
	twoBlondesURL          = "/organization/80925"
	adoptionText           = "Adoption fee includes the following"
	showLessText           = "show less"
	fosterText             = "Foster"
	petLinkClass           = ".pet-link"
	header3                = "h3"
	clientsId              = "#oc-clients"
	petDescriptionClass    = ".pet-description-full"
	petContainerClass      = ".pet-container"
	urlLink                = "href"
	petGalleryClass        = ".pet-gallery-thumb"
	petGalleryURLAttribute = "data-pet-gallery-url"
)

// RunDogDownloads starts scrapping the description and the pictures of the specified dogs to the specified directory.
// First, it must match the specified dog names against all available dogs on the web page. When it finds a match
// it will grab the description of the dog and visit the dog's personal information page.
// On the personal page, it will download all images there are of the dog.
func RunDogDownloads(dogs []string, baseDirectory string) error {
	// Create the scrappers
	availableDogs := colly.NewCollector()
	dogPictures := availableDogs.Clone()

	// Save the current dog to use when downloading pictures
	var currentDog string
	// Handle when the page of all the available dogs is loaded
	availableDogs.OnHTML(petLinkClass, func(e *colly.HTMLElement) {
		dogName := e.ChildText(header3)
		currentDog = strings.ToLower(dogName)
		dogMatch := isMatch(dogs, currentDog)
		// If a match then create dir and description.txt file
		if dogMatch {
			if err := io.MakeDir(baseDirectory + "/" + currentDog); err != nil {
				log.Println(err)
				return
			}
			log.Println("Found", dogName)
			fullDescription := e.ChildText(petDescriptionClass)
			// Remove the adoption fee part
			index := strings.Index(fullDescription, adoptionText)
			var desc string
			if index < 0 {
				desc = fullDescription[:strings.Index(fullDescription, showLessText)]
			} else {
				desc = fullDescription[:index]
			}
			// Add the link for adopting
			desc += "\n"
			desc += `👇👇SUBMIT AN APPLICATION HERE: 👇👇
https://2babrescue.com/adoption-fees-info`
			descFile := baseDirectory + "/" + currentDog + "/description.txt"
			if err := io.WriteFile(desc, descFile); err != nil {
				log.Println(err)
				return
			}
			// Get the link to the dog's page to download pictures
			link := e.Attr(urlLink)
			if err := dogPictures.Visit(baseURL + link); err != nil {
				log.Println(err)
				return
			}
		}
	})

	// When the dog page is loaded, download pictures
	dogPictures.OnHTML(clientsId, func(e *colly.HTMLElement) {
		imageURLs := e.ChildAttrs(petGalleryClass, petGalleryURLAttribute)
		// Save all the images
		for index, imageURL := range imageURLs {
			imagePath := fmt.Sprintf("%s/%s", baseDirectory, currentDog)
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
	return availableDogs.Visit(baseURL + twoBlondesURL)
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

// RunFosters looks up all the dogs that are foster-able and prints a comma separated list of the dogs.
func RunFosters() error {
	// Create the scrappers
	availableDogs := colly.NewCollector()

	// List of dogs to be fostered
	var fosters []string

	// Handle when the page of all the available dogs is loaded
	availableDogs.OnHTML(petContainerClass, func(e *colly.HTMLElement) {
		var dogName string
		dom := e.DOM
		dom.Find(petLinkClass).Each(func(i int, selection *goquery.Selection) {
			dogName = strings.TrimSpace(selection.Find(header3).Text())
		})
		var buttonName string
		dom.Find(".actions").Each(func(i int, selection *goquery.Selection) {
			buttonName = selection.Find(".button").Text()
		})
		if strings.Contains(buttonName, fosterText) {
			fosters = append(fosters, dogName)
		}
	})

	// Log when the request is being made
	availableDogs.OnRequest(func(request *colly.Request) {
		log.Println("Starting foster lookup...")
	})

	// Start scrapping
	if err := availableDogs.Visit(baseURL + twoBlondesURL); err != nil {
		return err
	}
	printDogs(fosters)
	return nil
}

func printDogs(fosters []string) {
	log.Printf("A total of %d dogs needs to be fostered", len(fosters))
	log.Printf("Dogs to foster: %s", strings.Join(fosters, ","))
}
