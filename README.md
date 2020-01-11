# Pet Spotlight
This is a CLI tool used to quickly download description and all images related to dogs provided as a flag to 
the CLI tool. Description and images are saved to the specified location.

Currently the tool only works for the Foster Organization [2 Blondes All Breed Rescue](https://2babrescue.com/).

## Download
Visit the [Releases](https://github.com/Piszmog/pet-spotlight/releases) page to download the binary for your system.

## Running
To run the CLI the following flags are required,

* `-d` - a comma separated list of dogs to download information for -- e.g. `Boomer,Zodiac`
* `-l` - the base directory to save the information on the dogs -- e.g. `dogs`
* `-f` - determines if the list of dogs to foster should be printed

### Example
Example of running the help flag `-h`

```
> pet-spotlight.exe -h
Usage of pet-spotlight.exe:
  -d string
        A comma separate list of availableDogs to extract
  -f    Determines whether to list all of the dogs that need fosters
  -l string
        Location to place data on the extracted dogs
```

Example of running

`> pet-spotlight.exe -d aladdin -l dogs`

## Building
To build the CLI tool, there is a `makefile` provided. However, to run the `makefile` required Windows and `nmake`.

e.g. `nmake all`