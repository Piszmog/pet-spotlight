# Pet Spotlight
This is a CLI tool used to quickly download description and all images related to dogs provided as a flag to 
the CLI tool. Description and images are saved to the specified location.

Currently the tool only works for the Foster Organization [2 Blondes All Breed Rescue](https://2babrescue.com/).

There is a UI version available (only for Windows). It can be found at [Pet Spotlight UI](https://github.com/Piszmog/pet-spotlight-ui)

## Download
Visit the [Releases](https://github.com/Piszmog/pet-spotlight/releases) page to download the binary for your system.

## Running
To run the CLI the following flags are required,

* `-d` - a comma separated list of dogs to download information for -- e.g. `Boomer,Zodiac`
* `-l` - the base directory to save the information on the dogs -- e.g. `dogs`
* `-f` - determines if the list of dogs to foster should be printed

### Usage
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

### Running
1. Download the latest version from the [Releases Page](https://github.com/Piszmog/pet-spotlight/releases)
    * Download the correct binary for your operating machine
        * `pet-spotlight.exe` - Windows
        * `pet-spotlight-linux` - Linux
        * `pet-spotlight-mac` - Mac/OSX
    * Place the downloaded binary in an easily accessible location (e.g. `C:\Users\bob`)
2. Open up the Terminal
    * On Windows, search for `cmd` and open it
    * On Mac/Osx/Linux, open the `Terminal` application
3. In the Terminal, change directories to where you placed the binary - e.g. `cd C:\Users\bob`
4. In the Terminal, run the program
    * Listing Boarding dogs
        * e.g. on Windows `pet-spotlight.exe -f`
        * e.g. on Linux `./pet-spotlight-linux -f`
        * e.g. on Mac/OSX `./pet-spotlight-mac -f`
    * Download dogs
        * e.g. on Windows `pet-spotlight.exe -l dogs -d "Aladdin,Sugar Baby, Frankie"`
        * e.g. on Linux `./pet-spotlight-linux -l dogs -d "Aladdin,Sugar Baby, Frankie"`
        * e.g. on Mac/OSX `./pet-spotlight-mac -l dogs -d "Aladdin,Sugar Baby, Frankie"`

#### Examples

##### List Foster Dogs
```text
C:\Users\rande>cd goprojects

C:\Users\rande\goprojects>cd pet-spotlight

C:\Users\rande\goprojects\pet-spotlight>pet-spotlight.exe -f
Starting foster lookup...
A total of 64 dogs needs to be fostered
Dogs to foster: Aldis,Aspen,Augustus,Beamer,Beau,Biscuit,Bonnie,Boots,Brawny aka Ace,Brody,Buster aka Aspen,Chance,Cher,Chico fka Duncan,Cleo,Cody,Datsun,Delilah,Dobie,Drake,Dude,Emma,Felix,Forrest,Franny,Freckles,Gerard,Gypsy,Hershel,Jack,Jenny,Kail,Lenny,Lily,Lucille aka Lucy,Lucy,Lyla,Macaron,Marc Antony,Marcos,Matilda,Miko,Minnie,Miyako,Moonie,Murdoch,Murphy,Nickel,Niko,Nougat,Opal,Pancake,Penelope (aka Penny),Red,Rico,Roseanne aka Rosie,Saturn,Sugar,Sugar Baby,Talia,Theo,Tiffany,Ziggy,Zodiac
Application ran in 1.297434sec
```

##### Download Dogs
```text
C:\Users\rande>cd goprojects

C:\Users\rande\goprojects>cd pet-spotlight

C:\Users\rande\goprojects\pet-spotlight>pet-spotlight.exe -l dogs -d "Aldis,Sugar Baby"
Starting extraction...
Found Aldis
Found Sugar Baby
Application ran in 4.643689sec
```

## Building
To build the CLI tool, there is a `makefile` provided. However, to run the `makefile` required Windows and `nmake`.

e.g. `nmake all`