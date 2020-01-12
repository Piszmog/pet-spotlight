package main

import (
	"flag"
	"fmt"
	"log"
	"pet-spotlight/io"
	"strings"
)

type flags struct {
	dogs             string
	baseDirectory    string
	determineFosters bool
}

func main() {
	// Setup flags
	dogsFlag := flag.String("d", "", "A comma separate list of availableDogs to extract")
	baseDirectory := flag.String("l", "", "Location to place data on the extracted dogs")
	determineFosters := flag.Bool("f", false, "Determines whether to list all of the dogs that need fosters")
	flag.Parse()

	f := flags{
		dogs:             *dogsFlag,
		baseDirectory:    *baseDirectory,
		determineFosters: *determineFosters,
	}

	if !validateFlags(f) {
		return
	}

	if len(f.dogs) != 0 {
		// Create directory where the dog info will go
		if err := io.MakeDir(f.baseDirectory); err != nil {
			log.Fatalln(err)
		}
		dogs := createDogMap(f.dogs)
		if err := RunDogDownloads(dogs, f.baseDirectory); err != nil {
			log.Fatalln(err)
		}
	}
	if f.determineFosters {
		if err := RunFosters(); err != nil {
			log.Fatalln(err)
		}
	}
}

func validateFlags(f flags) bool {
	if !f.determineFosters && len(f.dogs) == 0 {
		fmt.Println("Missing flag '-d' or '-f'")
		flag.Usage()
		return false
	} else if !f.determineFosters && len(f.baseDirectory) == 0 {
		fmt.Println("Missing flag '-l'")
		flag.Usage()
		return false
	} else if !f.determineFosters && len(f.dogs) == 0 && len(f.baseDirectory) == 0 {
		fmt.Println("No flags provided")
		flag.Usage()
		return false
	}
	return true
}

func createDogMap(dogsFlag string) map[string]bool {
	selectedDogs := strings.Split(dogsFlag, ",")
	dogs := make(map[string]bool)
	for _, dog := range selectedDogs {
		dogs[strings.ToLower(dog)] = false
	}
	return dogs
}
