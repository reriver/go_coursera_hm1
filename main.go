package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	return scanDir(out, path, printFiles, 0, []rune{})
}

func scanDir(out io.Writer, path string, printFiles bool, shift int, prefix []rune) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	filesLeftInDirectory := len(files)
	for _, file := range files {
		if file.IsDir() == false { // if file is not a Directory
			if file.Name() == ".DS_Store" || file.Name() == ".dockerignore" || file.Name() == "dockerfile" || file.Name() == "hw1.md" {
				filesLeftInDirectory--
				continue
			}
			if !printFiles { // if we need to Print files
				filesLeftInDirectory--
			}
		} else { // if file is Directory
			if file.Name() == ".git" || file.Name() == ".idea" {
				filesLeftInDirectory--
				continue
			}
		}
	}
	for _, file := range files {
		if file.IsDir() == false { // if file is not a Directory
			if file.Name() == ".DS_Store" || file.Name() == ".dockerignore" || file.Name() == "dockerfile" || file.Name() == "hw1.md" {
				continue
			}
			if printFiles { // if we need to Print files
				err := printFile(out, file, &filesLeftInDirectory, prefix)
				if err != nil {
					return err
				}
			}
		} else { // if file is Directory
			if file.Name() == ".git" || file.Name() == ".idea" {
				continue
			}
			err1 := printDir(out, file, &filesLeftInDirectory, prefix)
			if err1 != nil {
				return err1
			}
			if filesLeftInDirectory > 0 {
				err2 := scanDir(out, path+"/"+file.Name(), printFiles, shift+1, append(prefix, []rune("│\t")...))
				if err2 != nil {
					return err2
				}
			} else {
				err2 := scanDir(out, path+"/"+file.Name(), printFiles, shift+1, append(prefix, []rune("\t")...))
				if err2 != nil {
					return err2
				}
			}
		}
	}
	return err
}

func printFile(out io.Writer, file fs.FileInfo, filesLeftInDirectory *int, prefix []rune) error {
	_, err := fmt.Fprintf(out, string(prefix))
	if err != nil {
		return err
	}
	var size string
	if file.Size() > 0 {
		size = strconv.Itoa(int(file.Size())) + "b"
	} else {
		size = "empty"
	}
	if *filesLeftInDirectory > 1 {
		*filesLeftInDirectory--
		_, err := fmt.Fprintf(out, "├───%s (%s)\n", file.Name(), size)
		return err
	} else {
		*filesLeftInDirectory--
		_, err := fmt.Fprintf(out, "└───%s (%s)\n", file.Name(), size)
		return err
	}
}

func printDir(out io.Writer, file fs.FileInfo, filesLeftInDirectory *int, prefix []rune) error {
	_, err := fmt.Fprintf(out, string(prefix))
	if err != nil {
		return err
	}
	if *filesLeftInDirectory > 1 {
		*filesLeftInDirectory--
		_, err := fmt.Fprintf(out, "├───%s\n", file.Name())
		return err
	} else {
		*filesLeftInDirectory--
		_, err := fmt.Fprintf(out, "└───%s\n", file.Name())
		return err
	}
}

func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(os.Stdout, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
