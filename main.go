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
	return scanDir(out, path, printFiles, 0)
}

func scanDir(out *os.File, path string, printFiles bool, shift int) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	filesLeftInDirectory := len(files)
	//fmt.Println("\n", " Files in " + path + " list: ", filesLeftInDirectory, " shift: ", shift)

	//err1 := printDir(file, shift)

	for _, file := range files {

		//fmt.Printf("%d ->%s<-\t", n, file.Name())

		if file.IsDir() == false { // if file is not a Directory
			if file.Name() == ".DS_Store" || file.Name() == ".dockerignore" || file.Name() == "dockerfile" || file.Name() == "hw1.md" {
				//fmt.Println(" skipping file...")
				filesLeftInDirectory--
				fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
				continue
			}
			if printFiles { // if we need to Print files
				err := printFile(file, shift)
				if err != nil {
					return err
				}
				filesLeftInDirectory--
				fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
			}
		} else { // if file is Directory
			//fmt.Println()
			if file.Name() == ".git" || file.Name() == ".idea" {
				//fmt.Println(" skipping directory...")
				filesLeftInDirectory--
				fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
				continue
			}
			err1 := printDir(file, shift)
			if err1 != nil {
				return err1
			}
			filesLeftInDirectory--
			fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
			//fmt.Println( "Entering new directory; >" + path + "/" + file.Name() + "<")
			err2 := scanDir(out, path+"/"+file.Name(), printFiles, shift+1)
			if err2 != nil {
				return err2
			}
		}
	}
	return err
}

func printFile(file fs.FileInfo, shift int) error {
	//fmt.Println("printFile with shift ", shift)
	for i := 0; i < shift*1+0; i++ {
		fmt.Printf("│   ")
	}
	_, err := fmt.Println("├───" + file.Name()) //, "\t is Dir = ", file.IsDir())
	return err
}

func printDir(file fs.FileInfo, shift int) error {
	//fmt.Println("printDir with shift ", shift)
	for i := 0; i < shift*1+0; i++ {
		fmt.Printf("│   ")
	}
	_, err := fmt.Println("├───" + file.Name()) //, "\t is Dir = ", file.IsDir())
	return err
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
