package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func dirTree(out *os.File, path string, printFiles bool) error {
	//_, file := filepath.Split(path)
	//fmt.Println("└───" + file)
	return scanDir(out, path, printFiles, 0, "")
}

func scanDir(out *os.File, path string, printFiles bool, shift int, prefix string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	filesLeftInDirectory := len(files) - 0
	//fmt.Println("Files in " + path + " list: ", filesLeftInDirectory, " shift: ", shift)

	//err1 := printDir(file, shift)

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
	//fmt.Println(prefix, " Files to display in " + path + ": ", filesLeftInDirectory, ", shift: ", shift)

	for _, file := range files {

		//fmt.Printf("%d ->%s<-\t", n, file.Name())

		if file.IsDir() == false { // if file is not a Directory
			if file.Name() == ".DS_Store" || file.Name() == ".dockerignore" || file.Name() == "dockerfile" || file.Name() == "hw1.md" {
				//fmt.Println(" skipping file...")
				//filesLeftInDirectory--
				//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
				continue
			}
			if printFiles { // if we need to Print files
				err := printFile(file, shift, filesLeftInDirectory, prefix)
				if err != nil {
					return err
				}
			}
			filesLeftInDirectory--
			//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
		} else { // if file is Directory
			//fmt.Println()
			if file.Name() == ".git" || file.Name() == ".idea" {
				//fmt.Println(" skipping directory...")
				//filesLeftInDirectory--
				//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
				continue
			}
			err1 := printDir(file, shift, filesLeftInDirectory, prefix)
			if err1 != nil {
				return err1
			}
			filesLeftInDirectory--
			//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
			//fmt.Println( "Entering new directory; >" + path + "/" + file.Name() + "<")
			if filesLeftInDirectory > -1 {
				err2 := scanDir(out, path+"/"+file.Name(), printFiles, shift+1, prefix+"│   ")
				if err2 != nil {
					return err2
				}
			} else {
				err2 := scanDir(out, path+"/"+file.Name(), printFiles, shift+1, prefix+"    ")
				if err2 != nil {
					return err2
				}
			}
		}
	}
	return err
}

func printFile(file fs.FileInfo, shift int, filesLeftInDirectory int, prefix string) error {
	//fmt.Println("printFile with shift ", shift, " filesLeftInDirectory ", filesLeftInDirectory)
	//for i := 0; i < shift*1+0; i++ {
	//	fmt.Printf("│   ")
	//}
	fmt.Printf(prefix)
	if filesLeftInDirectory > 0 {
		_, err := fmt.Println("├───"+file.Name(), "\t files to display ", filesLeftInDirectory, "\t shift ", shift)
		return err
	} else {
		_, err := fmt.Println("└───"+file.Name(), "\t files to display ", filesLeftInDirectory, "\t shift ", shift)
		return err
	}
}

func printDir(file fs.FileInfo, shift int, filesLeftInDirectory int, prefix string) error {
	//fmt.Println("printDir with shift ", shift, " filesLeftInDirectory ", filesLeftInDirectory)
	//for i := 0; i < shift*1+0; i++ {
	//	fmt.Printf("│   ")
	//}
	fmt.Printf(prefix)
	if filesLeftInDirectory > 1 {
		_, err := fmt.Println("├───"+file.Name(), "\t files to display ", filesLeftInDirectory, "\t shift ", shift)
		return err
	} else {
		_, err := fmt.Println("└───"+file.Name(), "\t files to display ", filesLeftInDirectory, "\t shift ", shift)
		return err
	}
}

func main() {
	// go run main.go . -f
	//fmt.Println("os.Args len =",len(os.Args))
	//fmt.Println("os.Args cap =",cap(os.Args))

	//for i := range os.Args {
	//	fmt.Println(i, " ", os.Args[i])
	//}

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	//scanDir(os.Args[1])
	//	fmt.Println("Good number of arguments: ", len(os.Args))
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	out := os.Stdout
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
