package main

import (
	"log"
	"sync"

	"github.com/MrWong99/IcfgO/config"
	"github.com/MrWong99/IcfgO/console"
	"github.com/MrWong99/IcfgO/scramble"
	"github.com/MrWong99/filemanager/fman"
)

// This WaitGroup is used to wait for all routines to finish before exiting.
// See https://gobyexample.com/waitgroups
var wg sync.WaitGroup

func main() {

	appConfig := config.ParseAppConfig()

	// Initialize used in- and outputs
	var consoleIO config.ReaderWriter = &console.Console{}
	inFile := fman.File{
		Path: appConfig.InputFile,
	}
	outFile := fman.File{
		Path: appConfig.OutputFile,
	}

	// Initialize input slices and use direct type declarations
	var readers []config.Reader = []config.Reader{consoleIO, &inFile}
	var writers []config.Writer = []config.Writer{consoleIO, &outFile}

	// Initialize a slice of channels that can be used to send slices of strings
	var scrambleChannels []chan []string

	for idx, w := range writers {
		// For each writer add a new channel to the list of channels which awaits exactly one string slice
		scrambleChannels = append(scrambleChannels, make(chan []string, 1))

		wg.Add(1) // This tells the wait group that there is a new routine started

		// Call a new async go routine that outputs to the writer.
		// The writer and the channel have to be specifically given as function paramerter
		// since their references "w" and "scramblechannels[idx]" are changed for each iteration in the loop
		// which would result in unexpected/unwanted behaviour.
		// The type syntax "<-chan []string" indicates that the given channel can only be used to read
		// messages but not send
		go func(writy config.Writer, inputChan <-chan []string) {
			defer wg.Done() // Upon exiting the routine tell the wait group that one routine is finished

			// Wait for a message to be received on the given channel. This will block this go routine until this happens
			scrambledOutput := <-inputChan

			log.Printf("Started writing output to %T\n", writy)
			err := writy.Write(scrambledOutput)
			if err != nil {
				log.Printf(`Error while processing write to "%T": %v`+"\n", writy, err)
			}
			log.Printf("Finished writing output to %T\n", writy)
		}(w, scrambleChannels[idx])
	}

	// Printf with %#v placeholders is good for slices. See https://yourbasic.org/golang/fmt-printf-reference-cheat-sheet/
	log.Printf("Scramble adventure! Using %#v as input to scramble and saving output to %#v\n", readers, writers)

	// Get all the input from each reader store it to a two-dimensional slice in order by input
	var scramblePreparation [][]string
	for _, r := range readers {
		input, err := r.Read()
		if err != nil {
			log.Fatalf(`Error while processing read from "%T": %v`, r, err)
		}
		scramblePreparation = append(scramblePreparation, input)
	}

	scrambledOutput := scramble.Scramble(scramblePreparation)

	// Send the scrambled output via the previously created channels to all writers.
	// Start a go routine so it is non-blocking
	for _, outputChan := range scrambleChannels {
		// Initialize a new go routine and give the channel as function parameter since we are in a loop (same as above)
		// The type syntax "chan<- []string" indicates that the given channel can only be used to send
		// messages but not read
		go func(outChan chan<- []string) {
			outChan <- scrambledOutput
		}(outputChan)
	}

	wg.Wait() // Wait for all routines to finish
}
