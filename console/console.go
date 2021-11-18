package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Console struct{}

func (*Console) Read() ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	var result []string
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		// remove line breaks
		text = strings.Replace(text, "\r\n", "", -1)
		text = strings.Replace(text, "\n", "", -1)
		if text == "quit" {
			break
		}
		result = append(result, text)
	}
	return result, nil
}

func (*Console) Write(content []string) error {
	for _, line := range content {
		fmt.Println(line)
	}
	return nil
}
