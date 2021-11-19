package config

import (
	"flag"
	"log"
)

type AppConfig struct {
	InputFile  string
	OutputFile string
}

func ParseAppConfig() AppConfig {
	inputFile := flag.String("input-file", "", "A file that should be read to be scrambled with console.")
	outputFile := flag.String("output-file", "./test.out", "A file that should contain the scrambled result.")
	flag.Parse()

	if len(*inputFile) == 0 {
		log.Fatalln("Please provide a valid --input-file.")
	}
	return AppConfig{
		InputFile:  *inputFile,
		OutputFile: *outputFile,
	}
}
