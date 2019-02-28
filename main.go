package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sort"

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

	resPath = map[string]string{
		"a": "./result/a_example.txt",
		"b": "./result/b_lovely_landscapes.txt",
		"c": "./result/c_memorable_moments.txt",
		"d": "./result/d_pet_pictures.txt",
		"e": "./result/e_shiny_selfies.txt",
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
	fileKeys := []string{"a", "b", "c", "d", "e"}

	for _, fileKey := range fileKeys {
		images, err := parseFile(imagePath[fileKey])
		if err != nil {
			fmt.Printf("fatal: %v\n", err)
		}

		//fmt.Printf("Result : %#v\n", images)
		tags := GenerateDictionnaireTag(images)

		curSlideShow := processSlideShow(tags)

		//fmt.Printf("Result : %#v\n", curSlideShow)

		writeFile(curSlideShow, fileKey)
	}
}

func processSlideShow(tags map[string][]image) slideShow {
	curSlideShow := slideShow{}
	//fmt.Printf("\nthe lenght of tags is %d", len(tags))

	sortedKeys := make([]string, 0, len(tags))
	for k, _ := range tags {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		//fmt.Printf("\ntag is %d", v)
		//fmt.Printf("\nthe lenght of v is %d", len(v))
		for _, image := range tags[key] {
			if !intInSlice(image.id, curSlideShow.ids) {
				lastSlide := slide{}
				if len(curSlideShow.slides) > 0 {
					lastSlide = curSlideShow.slides[len(curSlideShow.slides)-1]
				}

				if len(lastSlide) == 1 && lastSlide[0].vertical && image.vertical {
					lastSlide = append(lastSlide, image)
					curSlideShow.slides[len(curSlideShow.slides)-1] = lastSlide
					curSlideShow.ids = append(curSlideShow.ids, image.id)
				} else {
					temporySlide := make(slide, 1)
					temporySlide[0] = image
					curSlideShow.slides = append(curSlideShow.slides, temporySlide)
					curSlideShow.ids = append(curSlideShow.ids, image.id)
				}
			}
		}
	}
	return curSlideShow
}

func writeFile(curSlideShow slideShow, fileKey string) {
	resultFileContent := fmt.Sprintf("%d\n", len(curSlideShow.slides))
	for _, slide := range curSlideShow.slides {
		if len(slide) == 1 {
			resultFileContent += fmt.Sprintf("%d\n", slide[0].id)
		} else {
			resultFileContent += fmt.Sprintf("%d %d\n", slide[0].id, slide[1].id)
		}
	}
	if err := afero.WriteFile(fs, resPath[fileKey], []byte(resultFileContent), 0644); err != nil {
		panic(err)
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
