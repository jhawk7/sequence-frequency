# Sequence Frequecy ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white) 
This golang application parses a given text file (or a list of text files) and returns the top 100 most frequent sequences - comprised of 3 words.

## Setting up environment
Download the go runtime from the link to the page below. This page will also contain further instructions to setup your local go env.
`https://go.dev/doc/install`
This application can also be ran via docker (see instructions for running with docker)

## Running Application with Go
The application can accept arguments via the command line or stdin. You can choose to run the application with `go run main.go` providing it with agrguments, or you can use `go build` to compile your own executable and run it directly from the binary. The pre-compiled binary `sequence-frequency` (compiled using macOS x86_64) is provided and can be found in the `bin` directory of this repo.
- `go run main.go <your_file_path1> <your_file_path2>` - runs application with list of files (the results will be returned in sequential order)
- `cat <your_file_path> | go run main.go` - pipes file contents into the application via stdin
- `./bin/sequence-frequency <your_file_path>` - example using bin file with cmdline arguments
- `cat <your_file_path> | ./bin/sequence-frequency` - example using bin file with stdin

## Running Application with Docker
The application is dockerized and can be ran using `docker-compose run` or `docker run`. Before running using docker, insert your text files into the `texts/` directory of the application. Then, pass the path to the text files in the docker command (pass via stdin or commandline argument).
- `docker-compose run sequence-frequency texts/<your_filename1> texts/<your_filename2>` - to run via docker and cmdline arguments
- `cat texts/<your_filename> | docker-compose run sequence-frequency` - to run via docker and stdin

## Running Tests
Tests in this application are split up by packages (pkg). Each package has its own test that can be found in the pkg subfolder with the exception of the test file for the main.go file which can be found in the root directory of the application. 
 - Use `go test ./...` to run all tests in the application at once
 - To run tests individually, navigate to the directory of the test and run `go test`

## Sample Output
    Top 100 Seqences for File: texts/moby-dick.txt
    the sperm whale - 86
    of the whale - 77
    one of the - 64
    the white whale - 59
    ...

## Epilogue
- I was able to get this to process text files pretty quickly by taking advantage of Go's channels and goroutines, as well as the underlying heap-map data structure. This application could be more resilient if implemented as a worker scaled horizontally and subscribed to a queue
- The output frequencies were not the exact same as what was given in the sample output (86 vs 85 for "the sperm whale")
