package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

var (
	fs = afero.NewOsFs()

	imagePath = map[string]string{
		"a": "./data/a_example.txt",
		"b": "./data/b_lovely_landscapes.txt",
		"c": "./data/c_memorable_moments.txt",
		"d": "./data/d_pet_pictures.txt",
		"e": "./data/e_shiny_selfies.txt",
	}
)

type image struct {
	id       int
	vertical bool
	tags     []string
}

type tag map[string][]image

func main() {
	images, err := parseFile(imagePath["a"])
	if err != nil {
		fmt.Printf("fatal: %v\n", err)
	}

	//fmt.Printf("Result : %#v\n", images)

	tags := GenerateDictionnaireTag(images)

	for k, v := range tags {
		for _, image := range v {
			fmt.Println(k, ": ", image)
		}
	}

	//fmt.Printf("Result : %#v\n", tags)
}

func parseFile(filepath string) ([]image, error) {
	b, err := afero.ReadFile(fs, filepath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	if len(lines) < 1 {
		return nil, errors.New("no lines in file")
	}

	nbrImages, _ := strconv.Atoi(lines[0])
	images := make([]image, 0)

	for i := 0; i < nbrImages; i++ {
		items := strings.Split(lines[i+1], " ")
		img := NewImage()
		img.id = i

		if items[0] == "V" {
			img.vertical = true
		}

		for j := 0; j < len(items)-2; j++ {
			img.tags = append(img.tags, items[j+2])
		}

		images = append(images, img)
	}

	return images, nil
}

func NewImage() image {
	return image{
		vertical: false,
		tags:     make([]string, 0),
	}
}

func GenerateDictionnaireTag(images []image) map[string][]image {
	tag := make(map[string][]image)
	for i := range images {
		for j := range images[i].tags {
			//fmt.Printf("\ntag is : %s and image is : %s", images[i].tags[j], images[i])
			tag[images[i].tags[j]] = append(tag[images[i].tags[j]], images[i])
		}

	}
	return tag
}
