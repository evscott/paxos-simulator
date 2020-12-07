package util

import (
	"fmt"
	"os"
)

func WriteToFile(text string) {
	file, err := os.OpenFile("./artifacts/basic-paxos-output.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("Can't write error: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", text))
}
