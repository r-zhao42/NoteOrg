package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defaultDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	} // cli args
	opts := flag.NewFlagSet("options", flag.ExitOnError)

	defaultDir += "/Documents/Obsidian Vault"

	wd := opts.String("dir", defaultDir, "Directory to organize, defaults to obsidian vault directory")
	organizeAll := opts.Bool("all", false, "Organize all files in directory, defaults to false")

	opts.Parse(os.Args[1:])
	fmt.Println("Working Directory:", *wd)

	filepath.WalkDir(*wd, func(path string, d fs.DirEntry, err error) error {

		// For daily use, we just want to organize the top level loose notes
		if d.IsDir() && !*organizeAll && path != *wd {
			return filepath.SkipDir
		}

		if notHidden(d.Name()) && !d.IsDir() {
			moveNote(path)
		}
		return nil
	})
}

func moveNote(fp string) {
	filePathSplit := strings.Split(fp, "/")

	filename := filePathSplit[len(filePathSplit)-1]
	baseDir := strings.Join(filePathSplit[:len(filePathSplit)-1], "/")

	filenameSplit := strings.Split(filename, " ")

	err := os.MkdirAll(filepath.Join(baseDir, strings.Join(filenameSplit[:len(filenameSplit)-1], "/")), 0750)

	if err != nil {
		log.Fatal("error making directories:", err)
	}

	err = os.Rename(fp, filepath.Join(baseDir, strings.Join(filenameSplit, "/")))

	if err != nil {
		log.Fatal("error renaming file:", err)
	}
}

func notHidden(filename string) bool {
	if len(filename) < 1 {
		log.Fatal("Should not call function with empty path")
	}

	return filename[0] != byte('.')
}
