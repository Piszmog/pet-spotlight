package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"pet-spotlight/http"
	"pet-spotlight/io"
	"pet-spotlight/sync"
	"pet-spotlight/wait"
	"sort"
	"strings"
)

const (
	actionsClass           = ".actions"
	adoptionText           = "Adoption fee includes the following"
	baseURL                = "https://www.petstablished.com"
	buttonClass            = ".button"
	clientsId              = "#oc-clients"
	dogNameContext         = "dogName"
	errorClass             = ".error"
	fosterText             = "Foster"
	header3                = "h3"
	linkAttribute          = "href"
	maxPages               = 100
	petDescriptionClass    = ".pet-description-full"
	petContainerClass      = ".pet-container"
	petGalleryClass        = ".thumb-img"
	petGalleryURLAttribute = "data-pet-gallery-url"
	petLinkClass           = ".pet-link"
	showLessText           = "show less"
	twoBlondesPath         = "/organization/80925"
	urlLink                = "href"
	widgetPage             = "/widget/dogs?page=%d"
)

const defaultDescription = `👇👇SUBMIT AN APPLICATION HERE: 👇👇
https://2babrescue.com/adoption-fees-info`

// RunDogDownloads starts scrapping the description and the pictures of the specified dogs to the specified directory.
// First, it must match the specified dog names against all available dogs on the web page. When it finds a match
// it will grab the description of the dog and visit the dog's personal information page.
// On the personal page, it will download all images there are of the dog.
func RunDogDownloads(dogs string, baseDirectory string) error {
	// Convert the comma sep list of dogs to a map
	dogMap := createDogMap(dogs)
	// Create the scrappers
	availableDogs := colly.NewCollector(colly.Async(true))
	if err := availableDogs.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10}); err != nil {
		return fmt.Errorf("failed to set parallel limit: %w", err)
	}
	dogPictures := availableDogs.Clone()

	// Save the current dog to use when downloading pictures
	isDone := sync.AtomicBoolean{}

	// Handle when last page is reached
	availableDogs.OnHTML(errorClass, func(e *colly.HTMLElement) {
		isDone.Set(true)
	})

	// Handle when the page of all the available dogs is loaded
	availableDogs.OnHTML(petLinkClass, func(e *colly.HTMLElement) {
		if dogMap.IsCompete() {
			if !isDone.Get() {
				isDone.Set(true)
			}
			return
		}
		name := e.ChildText(header3)
		dogName := strings.ToLower(name)
		dogMatch := dogMap.IsMatch(dogName)
		// If a match then create dir and description.txt file
		if dogMatch {
			if err := io.MakeDir(baseDirectory + "/" + dogName); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Found", name)
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
			desc += defaultDescription
			descFile := baseDirectory + "/" + dogName + "/description.txt"
			if err := io.WriteFile(desc, descFile); err != nil {
				fmt.Println(err)
				return
			}
			// Get the link to the dog's page to download pictures
			link := e.Attr(urlLink)
			// Add the dog name to the context of the request
			dogPictures.OnRequest(func(request *colly.Request) {
				request.Ctx.Put(dogNameContext, dogName)
			})
			if err := dogPictures.Visit(link); err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	// When the dog page is loaded, download pictures
	dogPictures.OnHTML(clientsId, func(e *colly.HTMLElement) {
		dogName := e.Request.Ctx.Get(dogNameContext)
		imageURLs := e.ChildAttrs(petGalleryClass, petGalleryURLAttribute)
		videoURLs := e.ChildAttrs(petGalleryClass, linkAttribute)
		// Save all the images
		wg := wait.NewBoundedWaitGroup(5)
		for index, imageURL := range imageURLs {
			imageFile := fmt.Sprintf("image-%d.png", index)
			wg.Add(1)
			go download(baseDirectory, dogName, imageFile, imageURL, &wg)
		}
		for index, videoURL := range videoURLs {
			videoFile := fmt.Sprintf("video-%d.mp4", index)
			wg.Add(1)
			go download(baseDirectory, dogName, videoFile, videoURL, &wg)
		}
		wg.Wait()
	})

	// Handle errors
	availableDogs.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s, Status Code %d, Error %+v\n", r.Request.URL, r.StatusCode, err)
	})

	// Handle errors
	dogPictures.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s, Status Code %d, Error %+v\n", r.Request.URL, r.StatusCode, err)
	})

	// Start scrapping
	fmt.Println("Starting extraction...")
	for i := 1; i < maxPages && !isDone.Get(); i++ {
		page := fmt.Sprintf(widgetPage, i)
		if err := availableDogs.Visit(baseURL + twoBlondesPath + page); err != nil {
			return err
		}
	}
	availableDogs.Wait()
	dogPictures.Wait()
	dogMap.PrintMissing()
	return nil
}

func createDogMap(dogsList string) *sync.DogMap {
	selectedDogs := strings.Split(dogsList, ",")
	return sync.InitializeMap(selectedDogs)
}

func download(baseDirectory string, dogName string, fileName string, url string, b *wait.BoundedWaitGroup) {
	defer b.Done()
	directoryPath := fmt.Sprintf("%s/%s", baseDirectory, dogName)
	if strings.HasSuffix(fileName, "png") {
		if err := http.Download(url, directoryPath, fileName); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := http.DownloadVideo(url, directoryPath, fileName); err != nil {
			fmt.Println(err)
		}
	}
}

// RunGetFosters looks up all the dogs that are foster-able and prints a comma separated list of the dogs.
func RunGetFosters() error {
	// Create the scrappers
	availableDogs := colly.NewCollector(colly.Async(true))
	if err := availableDogs.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10}); err != nil {
		return fmt.Errorf("failed to set parallel limit: %w", err)
	}

	// List of dogs to be fostered
	fosters := sync.DogList{}
	isDone := sync.AtomicBoolean{}

	// Handle when last page is reached
	availableDogs.OnHTML(errorClass, func(e *colly.HTMLElement) {
		isDone.Set(true)
	})

	// Handle when the page of all the available dogs is loaded
	availableDogs.OnHTML(petContainerClass, func(e *colly.HTMLElement) {
		var dogName string
		dom := e.DOM
		dom.Find(petLinkClass).Each(func(i int, selection *goquery.Selection) {
			dogName = strings.TrimSpace(selection.Find(header3).Text())
		})
		var buttonName string
		dom.Find(actionsClass).Each(func(i int, selection *goquery.Selection) {
			buttonName = selection.Find(buttonClass).Text()
		})
		if strings.Contains(buttonName, fosterText) {
			fosters.Add(dogName)
		}
	})

	// Handle errors
	availableDogs.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s, Status Code %d, Error %+v\n", r.Request.URL, r.StatusCode, err)
	})

	// Start scrapping
	fmt.Println("Starting foster lookup...")
	for i := 1; i < maxPages && !isDone.Get(); i++ {
		page := fmt.Sprintf(widgetPage, i)
		if err := availableDogs.Visit(baseURL + twoBlondesPath + page); err != nil {
			return err
		}
	}
	availableDogs.Wait()
	printDogs(fosters.Get())
	return nil
}

func printDogs(fosters []string) {
	sort.Sort(sort.StringSlice(fosters))
	fmt.Printf("A total of %d dogs needs to be fostered\n", len(fosters))
	fmt.Println("Dogs to foster:", strings.Join(fosters, ","))
}
