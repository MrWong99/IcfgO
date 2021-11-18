package main

import (
	"flag"
	"log"
	"sync"

	"github.com/MrWong99/IcfgO/console"
	"github.com/MrWong99/filemanager/fman"
)

type Reader interface {
	Read() ([]string, error)
}

type Writer interface {
	Write([]string) error
}

// Just an example to show interface composition.
type ReaderWriter interface {
	Reader
	Writer
}

// This WaitGroup is used to wait for all routines to finish before exiting.
// See https://gobyexample.com/waitgroups
var wg sync.WaitGroup

func main() {
	inputFile := flag.String("input-file", "", "A file that should be read to be scrambled with console.")
	outputFile := flag.String("output-file", "./test.out", "A file that should contain the scrambled result.")
	flag.Parse()

	if len(*inputFile) == 0 {
		log.Fatalln("Please provide a valid --input-file.")
	}

	// Initialize used in- and outputs
	var consoleIO ReaderWriter = &console.Console{}
	inFile := fman.File{
		Path: *inputFile,
	}
	outFile := fman.File{
		Path: *outputFile,
	}

	// Initialize input slices and use direct type declarations
	var readers []Reader = []Reader{consoleIO, &inFile}
	var writers []Writer = []Writer{consoleIO, &outFile}

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
		go func(writy Writer, inputChan <-chan []string) {
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

	scrambledOutput := scramble(scramblePreparation)

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

/* Scrambles the given inputs together. It will append to the result alternating evenly from the given inputs, e.g.:
toScramble := [][]string{
	[]string{"Hello"},
	[]string{"World", "is big"},
    []string{"!", "even a", "english sentence?"}
}
will produce
[]string{
	"Hello",
	"World",
	"!",
	"is big",
	"even a",
	"english sentence?",
}
*/
func scramble(toScramble [][]string) []string {
	var result []string
	lastIndex := 0
	for {
		// Retrieve the smallest length of all of the given slices. In the example above this would return 1,
		// but if the first slice []string{"Hello"} wouldn't exist it would return 2
		upToIndex := smallestLength(toScramble)

		// Iterate over the arrays by the index and append each string to the result slice, starting by the index from the last loop iteration or 0 if there was none yet,
		// up to the shortest length of all of the input arrays
		for i := lastIndex; i < upToIndex; i++ {
			for _, theArray := range toScramble {
				result = append(result, theArray[i])
			}
		}
		lastIndex = upToIndex

		// This removes all of the slices from toScramble that are as long as the last calculated smallest length.
		// These slices are too short to provide any new values for the output.
		var biggerSlices [][]string = [][]string{}
		for _, arr := range toScramble {
			if len(arr) > upToIndex {
				biggerSlices = append(biggerSlices, arr)
			}
		}
		toScramble = biggerSlices

		// Jump out of the infinite loop if there are no more slices left to process
		lengthAfterFilter := len(biggerSlices)
		if lengthAfterFilter <= 0 {
			break
		}
	}
	return result
}

// Retrieve the smallest length of all of the given slices or 0 if an empty slice is passed in
func smallestLength(toCompare [][]string) int {
	if len(toCompare) <= 0 {
		return 0
	}
	smallestLength := len(toCompare[0])
	for _, comp := range toCompare {
		length := len(comp)
		if length < smallestLength {
			smallestLength = length
		}
	}
	return smallestLength
}
