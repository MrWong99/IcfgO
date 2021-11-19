package scramble

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
func Scramble(toScramble [][]string) []string {
	var result []string

	// Retrieve the longest length of all of the given slices. In the example above this would return 3
	upToIndex := longestSize(toScramble)

	// Iterate over all the indices and append one value of each slice to the result, as long as the given slice
	// has values for the given index.
	for i := 0; i < upToIndex; i++ {
		for _, theArray := range toScramble {
			if len(theArray) > i {
				result = append(result, theArray[i])
			}
		}
	}

	return result
}

// Retrieve the smallest length of all of the given slices or 0 if an empty slice is passed in
func longestSize(toCompare [][]string) int {
	longest := len(toCompare[0])
	for _, comp := range toCompare {
		length := len(comp)
		if length > longest {
			longest = length
		}
	}
	return longest
}
