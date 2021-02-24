package mo2utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func CopyAllFiles(srcPath string, tgtPath string) {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		//os.Mkdir(rootPath, 0755)
		fmt.Println("Directory created")
	} else {
		os.Mkdir(tgtPath, 0755)
		copyAllFiles(srcPath, tgtPath)
		fmt.Println("Directory already exists")
	}
	return

}
func copyAllFiles(curPath string, tgtPath string) {
	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
		srcNewPath := path.Join(curPath, f.Name())
		tgtNewPath := path.Join(tgtPath, f.Name())
		if f.IsDir() {
			//mkdir
			os.Mkdir(tgtNewPath, 0755)
			copyAllFiles(srcNewPath, tgtNewPath)
		} else {
			// Part 1: open input file.
			inputFile, _ := os.Open(srcNewPath)

			// Part 2: call ReadAll to get contents of input file.
			data, _ := ioutil.ReadAll(inputFile)

			// Part 3: write data to copy file.
			ioutil.WriteFile(tgtNewPath, data, 0)
			fmt.Println("DONE")

		}
	}
}
