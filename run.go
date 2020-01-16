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
	adoptionText           = "Adoption fee includes the following"
	baseURL                = "https://www.petstablished.com"
	clientsId              = "#oc-clients"
	errorClass             = ".error"
	fosterText             = "Foster"
	header3                = "h3"
	maxPages               = 100
	petDescriptionClass    = ".pet-description-full"
	petContainerClass      = ".pet-container"
	petGalleryClass        = ".pet-gallery-thumb"
	petGalleryURLAttribute = "data-pet-gallery-url"
	petLinkClass           = ".pet-link"
	showLessText           = "show less"
	twoBlondesPath         = "/organization/80925"
	urlLink                = "href"
	widgetPage             = "/widget/dogs?page=%d"
)

// RunDogDownloads starts scrapping the description and the pictures of the specified dogs to the specified directory.
// First, it must match the specified dog names against all available dogs on the web page. When it finds a match
// it will grab the description of the dog and visit the dog's personal information page.
// On the personal page, it will download all images there are of the dog.
func RunDogDownloads(dogs map[string]bool, baseDirectory string) error {
	// Create the scrappers
	availableDogs := colly.NewCollector()
	dogPictures := availableDogs.Clone()

	// Save the current dog to use when downloading pictures
	var currentDog string
	dogsDownloaded := 0
	lastPageReached := false

	// Handle when last page is reached
	availableDogs.OnHTML(errorClass, func(e *colly.HTMLElement) {
		lastPageReached = true
	})

	// Handle when the page of all the available dogs is loaded
	availableDogs.OnHTML(petLinkClass, func(e *colly.HTMLElement) {
		if dogsDownloaded == len(dogs) {
			return
		}
		dogName := e.ChildText(header3)
		currentDog = strings.ToLower(dogName)
		dogMatch := isMatch(dogs, currentDog)
		// If a match then create dir and description.txt file
		if dogMatch {
			dogsDownloaded++
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
			if err := dogPictures.Visit(link); err != nil {
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

	// Start scrapping
	log.Println("Starting extraction...")
	for i := 1; i < maxPages && !lastPageReached; i++ {
		page := fmt.Sprintf(widgetPage, i)
		if err := availableDogs.Visit(baseURL + twoBlondesPath + page); err != nil {
			return err
		}
	}
	for name, found := range dogs {
		if !found {
			log.Printf("Failed to find %s", name)
		}
	}
	return nil
}

func isMatch(dogs map[string]bool, currentDog string) bool {
	dogMatch := false
	for dog, alreadyDownloaded := range dogs {
		if !alreadyDownloaded {
			if strings.Contains(currentDog, dog) {
				dogMatch = true
				dogs[dog] = true
				break
			}
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
	lastPageReached := false

	// Handle when last page is reached
	availableDogs.OnHTML(errorClass, func(e *colly.HTMLElement) {
		lastPageReached = true
	})

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

	// Start scrapping
	log.Println("Starting foster lookup...")
	for i := 1; i < maxPages && !lastPageReached; i++ {
		page := fmt.Sprintf(widgetPage, i)
		if err := availableDogs.Visit(baseURL + twoBlondesPath + page); err != nil {
			return err
		}
	}
	printDogs(fosters)
	return nil
}

func printDogs(fosters []string) {
	log.Printf("A total of %d dogs needs to be fostered", len(fosters))
	log.Printf("Dogs to foster: %s", strings.Join(fosters, ","))
}
