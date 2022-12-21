package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func getFileNames(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		log.Fatal("ERROR getFileNames os.Open: ", err)
		return []string{}, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		log.Fatal("ERROR loadImage f.Readdir: ", err)
		return []string{}, err
	}
	filenames := make([]string, len(files))
	for i, v := range files {
		if !v.IsDir() {
			filenames[i] = v.Name()
		}
	}
	return filenames, nil
}

func loadImage(filePath string) (*ebiten.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("ERROR loadImage os.Open: ", err)
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal("ERROR loadImage image.Decode: ", err, filePath)
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
