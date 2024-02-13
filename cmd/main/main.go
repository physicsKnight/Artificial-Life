package main

import (
	"Langton-Loops/internal"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Define a list of rule files to be read
	files := []string{
		"Rules/LangtonLoop-Rules.txt",
		"Rules/Byl-Rules.txt",
		"Rules/SDSR-Rules.txt",
		"Rules/Evoloop-Rules.txt",
		"Rules/sexyloop.txt",
	}

	// Iterate through each file in the list
	for _, file := range files {
		// Open the file
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", file, err)
			continue
		}
		// Ensure the file is closed after it has been read
		defer f.Close()

		// Create a scanner to read the file line by line
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			// Read a line from the file
			line := scanner.Text()

			// Append the line to the corresponding rules list based on the file being read
			switch file {
			case files[0]:
				internal.LangtonRules = append(internal.LangtonRules, line)
			case files[1]:
				internal.BylRules = append(internal.BylRules, line)
			case files[2]:
				internal.SDSRRules = append(internal.SDSRRules, line)
			case files[3]:
				internal.EvoRules = append(internal.EvoRules, line)
			case files[4]:
				internal.SexyRules = append(internal.SexyRules, line)
			}
		}

		// If an error occurred while reading the file, print the error
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
		}
	}

	internal.Run()
}
