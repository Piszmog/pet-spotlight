package main

import (
	"flag"
	"fmt"
	"log"
	"pet-spotlight/io"
	"strings"
)

const (
	baseURL       = "https://www.petstablished.com"
	twoBlondesURL = "/organization/80925"
	adoptionText  = "Adoption fee includes the following"
)

func main() {
	// Setup flags
	dogsFlag := flag.String("d", "", "A comma separate list of availableDogs to extract")
	baseDirectory := flag.String("l", "", "Location to place data on the extracted dogs")
	flag.Parse()

	if validateFlags(dogsFlag, baseDirectory) {
		return
	}

	// Create directory where the dog info will go
	if err := io.MakeDir(*baseDirectory); err != nil {
		log.Fatalln(err)
	}

	dogs := createDogList(dogsFlag)

	Run(dogs, baseDirectory)
}

func validateFlags(dogsFlag *string, baseDirectory *string) bool {
	if len(*dogsFlag) == 0 {
		fmt.Println("Missing flag '-d'")
		flag.Usage()
		return true
	} else if len(*baseDirectory) == 0 {
		fmt.Println("Missing flag '-l'")
		flag.Usage()
		return true
	}
	return false
}

func createDogList(dogsFlag *string) []string {
	selectedDogs := strings.Split(*dogsFlag, ",")
	dogs := make([]string, len(selectedDogs), len(selectedDogs))
	for i, dog := range selectedDogs {
		dogs[i] = strings.ToLower(dog)
	}
	return dogs
}
