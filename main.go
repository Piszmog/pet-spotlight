package main

import (
	"flag"
	"fmt"
	"log"
	"pet-spotlight/io"
	"time"
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
	defer runtime(time.Now())

	if len(f.dogs) != 0 {
		// Create directory where the dog info will go
		if err := io.MakeDir(f.baseDirectory); err != nil {
			log.Fatalln(err)
		}
		if err := RunDogDownloads(f.dogs, f.baseDirectory); err != nil {
			log.Fatalln(err)
		}
	}
	if f.determineFosters {
		if len(f.baseDirectory) > 0 {
			// Create directory where the dog info will go
			if err := io.MakeDir(f.baseDirectory); err != nil {
				log.Fatalln(err)
			}
		}
		if err := RunGetFosters(f.baseDirectory); err != nil {
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

func runtime(t time.Time) {
	fmt.Printf("Application ran in %fsec\n", time.Since(t).Seconds())
}
