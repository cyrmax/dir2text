package main

import (
	"dir2text/internal/app"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

var version string

func main() {
	workingDirArg := flag.String("dir", "", "Specify custom working directory. Current working directory is used by default")
	outputFilenameArg := flag.String("output", "", "Specify custom output file name. If empty, output file will have the same name as the working directory")
	overwriteArg := flag.Bool("overwrite", false, "Set to true to automatically overwrite the output file if it already exists. Defaults to false for safety")
	versionArg := flag.Bool("version", false, "If specified, outputs the program version and exits")
	flag.Parse()

	if *versionArg {
		fmt.Printf("dir2text version %s", version)
		return
	}

	if *workingDirArg != "" {
		if err := os.Chdir(*workingDirArg); err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
	}
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	fmt.Printf("Working directory: %s\n", workingDir)
	var outputFileName string
	workingDirName := filepath.Base(workingDir)
	if *outputFilenameArg != "" {
		outputFileName = *outputFilenameArg
	} else {
		outputFileName = path.Join(workingDir, workingDirName) + ".txt"
	}
	if _, err := os.Stat(outputFileName); err == nil {
		if !*overwriteArg {
			fmt.Printf("Error! The output file %s already exists\n", outputFileName)
			return
		} else {
			fmt.Printf("Output file %s already exists, overwriting\n", outputFileName)
		}
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error opening output file: %s\n", err.Error())
		return
	}
	defer outputFile.Close()

	ignoreProcessor := app.NewIgnoreProcessor(outputFileName)
	ignoreProcessor.AddCommonIgnores()
	err = filepath.WalkDir(workingDir, func(currentPath string, dirInfo fs.DirEntry, walkError error) error {
		if walkError != nil {
			return walkError
		}
		isIgnored := ignoreProcessor.IsPathShouldBeIgnored(currentPath)
		if dirInfo.IsDir() {
			if isIgnored {
				return filepath.SkipDir
			}
			return nil
		} else {
			if isIgnored {
				return nil
			}
			isBinary, err := app.IsFileBinary(currentPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", currentPath, err)
			}
			if isBinary {
				fmt.Printf("Skipping binary file %s\n", currentPath)
				return nil
			}
			fmt.Printf("Adding file %s\n", currentPath)
			return app.WriteFileWithTitle(currentPath, outputFile)
		}
	})
	if err != nil {
		fmt.Printf("Error walking the directory %s: %s\n", workingDir, err.Error())
		return
	}
	fmt.Println("Finished!")
}
