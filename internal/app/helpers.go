package app

import (
	"os"
)

// TODO: Probably make a struct to not hardcode the title separator
var titleSeparator string = "================================================================================\n"

func WriteFileWithTitle(filePath string, outputFile *os.File) error {
	// TODO: Optimize memory consumption. Do not read the entire file into memory
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if _, err = outputFile.WriteString(titleSeparator); err != nil {
		return err
	}
	if _, err = outputFile.WriteString(filePath + "\n"); err != nil {
		return err
	}
	if _, err = outputFile.WriteString(titleSeparator); err != nil {
		return err
	}
	if _, err = outputFile.Write(data); err != nil {
		return err
	}
	if _, err = outputFile.WriteString("\n\n"); err != nil {
		return err
	}
	return nil
}
