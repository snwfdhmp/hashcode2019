package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/afero"
)

var (
	fs = afero.NewOsFs()

	imageChoice string

	imagePath = map[string]string{
		"a": "./data/a_example.txt",
		"b": "./data/b_lovely_landscapes.txt",
		"c": "./data/c_memorable_moments.txt",
		"d": "./data/d_pet_pictures.txt",
		"e": "./data/e_shiny_selfies.txt",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type image struct {
	id       int
	vertical bool
	tags     []string
}

func main() {
	flag.StringVar(&imageChoice, "image", "a", "image to choose")
	flag.Parse()

	fmt.Printf("Processing image %s\n", strings.ToUpper(imageChoice))

	images, err := parseFile(imagePath[imageChoice])
	if err != nil {
		fmt.Printf("fatal: %v\n", err)
		return
	}

	slides := make([]slide, 0)
	for i := 0; i < len(images); i++ {
		imgIndex := randomIndex(images)
		img := images[imgIndex]
		images = append(images[:imgIndex], images[imgIndex+1:]...)
		if !img.vertical {
			slides = append(slides, slide{images: []image{img}})
			continue
		}

		vIndex := findFirstVertical(images)
		if vIndex == -1 {
			continue
		}
		slides = append(slides, slide{images: []image{img, images[vIndex]}})
		images = append(images[:vIndex], images[vIndex+1:]...)
		i++
	}

	if err := afero.WriteFile(fs, fmt.Sprintf("ouput_%s.txt", imageChoice), Marshal(slides), 0760); err != nil {
		fmt.Printf("PANIIIC: cannot write file: %s\n", err)
		return
	}

	fmt.Println("Done :D")
}

func randomIndex(items []image) int {
	return rand.Intn(len(items))
}

func findFirstVertical(items []image) int {
	for i := range items {
		if items[i].vertical {
			return i
		}
	}
	return -1
}

type slide struct {
	images []image
}

func Marshal(slides []slide) []byte {
	output := fmt.Sprintf("%d\n", len(slides))
	for i := range slides {
		if len(slides[i].images) == 2 {
			output += fmt.Sprintf("%d %d\n", slides[i].images[0].id, slides[i].images[1].id)
			continue
		}

		output += fmt.Sprintf("%d\n", slides[i].images[0].id)
	}

	return []byte(output)
}

func (s *slide) isHealthy() bool {
	switch len(s.images) {
	case 2:
		if !s.images[0].vertical || !s.images[1].vertical {
			return false
		}
		return true
	case 1:
		if s.images[0].vertical {
			return false
		}
		return true
	default:
		return false
	}
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
