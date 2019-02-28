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

type slide []image

type slideShow struct {
	slides []slide
	ids    []int
}

func main() {
	images, err := parseFile(imagePath["a"])
	if err != nil {
		fmt.Printf("fatal: %v\n", err)
	}

	//fmt.Printf("Result : %#v\n", images)
	tags := GenerateDictionnaireTag(images)

	slideShow := slideShow{}

	//fmt.Printf("\nthe lenght of tags is %d", len(tags))
	for _, v := range tags {
		//fmt.Printf("\ntag is %d", v)
		//fmt.Printf("\nthe lenght of v is %d", len(v))
		for _, image := range v {
			if !intInSlice(image.id, slideShow.ids) {
				lastSlide := slide{}
				if len(slideShow.slides) > 0 {
					lastSlide = slideShow.slides[len(slideShow.slides)-1]
				}

				if len(lastSlide) == 1 && lastSlide[0].vertical && image.vertical {
					lastSlide = append(lastSlide, image)
				} else {
					temporySlide := make(slide, 1)
					temporySlide[0] = image
					slideShow.slides = append(slideShow.slides, temporySlide)
					slideShow.ids = append(slideShow.ids, image.id)
				}
			}
		}
	}

	//fmt.Printf("Result : %#v\n", slideShow)
	fmt.Printf("%d\n", len(slideShow.slides))
	for _, slide := range slideShow.slides {
		fmt.Printf("%d\n", slide[0].id)
	}
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
