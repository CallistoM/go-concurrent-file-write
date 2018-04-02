package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type imageFile struct {
	gorm.Model
	FileImage string `gorm:"size:255"`
	Extension string `gorm:"size:255"`
}

func main() {

	// open db
	db, err := gorm.Open("mysql", "root:admin@tcp(localhost:3306)/images?charset=utf8&parseTime=True&loc=Local")

	// check if db throws error
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// migrate using struct
	db.AutoMigrate(&imageFile{})

	// open sync wait group
	wg := sync.WaitGroup{}

	// check location
	originalLocation := ""

	// read directory
	files, err := ioutil.ReadDir(originalLocation)

	// check if read dir has error
	if err != nil {
		log.Fatal(err)
	}

	// loop over all files
	for index, file := range files {

		// set sync group to one
		wg.Add(1)

		// set file variable
		f := file

		// set files concurrent
		go func(index int) {

			// check if sync group is done
			defer wg.Done()

			// set image location
			setImageLocation := originalLocation + "/" + f.Name()

			// set extension
			ext := filepath.Ext(setImageLocation)

			// set struct
			image := imageFile{FileImage: setImageLocation, Extension: ext}

			// create image record
			db.Create(&image)

		}(index)

		// get width and height
		width, height := getImageDimension(image.FileImage)

	}

	// close db
	defer db.Close()

	// wait for sync group to finish
	defer wg.Wait()

}

// getImageDimension return image dimensions of a image
func getImageDimension(imagePath string) (int, int) {
	// open file
	file, err := os.Open(imagePath)

	// check if os gives error when opening image path
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	// decode image config
	imageData, err := image.DecodeConfig(file)

	// check if decoding is failing
	if err != nil {
		os.Exit(2)
	}

	return imageData.Width, imageData.Height
}
