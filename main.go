package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	return scanDir(out, path, printFiles, 0, []rune{})
}

func scanDir(buf *bytes.Buffer, path string, printFiles bool, shift int, prefix []rune) error {
	out := bufio.NewWriter(buf)
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
	//fmt.Println(prefix, " Files to display in " + path + ": ", filesLeftInDirectory, ", shift: ", shift, "prefix ->", prefix, "<-")

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
				err := printFile(out, file, &filesLeftInDirectory, prefix)
				if err != nil {
					return err
				}
			}
			//filesLeftInDirectory--
			//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
		} else { // if file is Directory
			//fmt.Println()
			if file.Name() == ".git" || file.Name() == ".idea" {
				//fmt.Println(" skipping directory...")
				//filesLeftInDirectory--
				//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
				continue
			}
			err1 := printDir(out, file, &filesLeftInDirectory, prefix)
			if err1 != nil {
				return err1
			}
			//filesLeftInDirectory--
			//fmt.Println("filesLeftInDirectory: ", filesLeftInDirectory)
			//fmt.Println( "Entering new directory; >" + path + "/" + file.Name() + "<")
			if filesLeftInDirectory > 0 {
				err2 := scanDir(buf, path+"/"+file.Name(), printFiles, shift+1, append(prefix, []rune("│   ")...))
				if err2 != nil {
					return err2
				}
			} else {
				err2 := scanDir(buf, path+"/"+file.Name(), printFiles, shift+1, append(prefix, []rune("    ")...))
				if err2 != nil {
					return err2
				}
			}
		}
	}
	return err
}

func printFile(out io.Writer, file fs.FileInfo, filesLeftInDirectory *int, prefix []rune) error {
	//fmt.Println("printFile with shift ", shift, " filesLeftInDirectory ", filesLeftInDirectory)
	//for i := 0; i < shift*1+0; i++ {
	//	fmt.Printf("│   ")
	//}
	fmt.Printf(string(prefix))
	if *filesLeftInDirectory > 1 {
		//_, err := fmt.Println("├───"+file.Name(), "         left ", filesLeftInDirectory, "\t shift ", shift)
		//_, err := fmt.Printf("├───%.10s         left %d shift %d\n", file.Name(), *filesLeftInDirectory, shift)
		//_, err := fmt.Printf("├───%s\n", file.Name())
		_, err := fmt.Fprintf(out, "├───%s\n", file.Name())
		*filesLeftInDirectory--
		return err
	} else {
		//_, err := fmt.Println("└───"+file.Name(), "         left ", filesLeftInDirectory, "\t shift ", shift)
		//_, err := fmt.Printf("└───%.10s         left %d shift %d\n", file.Name(), *filesLeftInDirectory, shift)
		_, err := fmt.Fprintf(out, "└───%s\n", file.Name())
		*filesLeftInDirectory--
		return err
	}
}

func printDir(out io.Writer, file fs.FileInfo, filesLeftInDirectory *int, prefix []rune) error {
	//fmt.Println("printDir with shift ", shift, " filesLeftInDirectory ", filesLeftInDirectory)
	//for i := 0; i < shift*1+0; i++ {
	//	fmt.Printf("│   ")
	//}
	fmt.Printf(string(prefix))
	if *filesLeftInDirectory > 1 {
		//_, err := fmt.Printf("├───%.10s         left %d shift %d\n", file.Name(), *filesLeftInDirectory, shift)
		_, err := fmt.Fprintf(out, "├───%s\n", file.Name())
		*filesLeftInDirectory--
		return err
	} else {
		_, err := fmt.Fprintf(out, "└───%s\n", file.Name())
		*filesLeftInDirectory--
		return err
	}
}

func main() {
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
	//out := *os.File(os.Stdout)
	out := new(bytes.Buffer)
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
