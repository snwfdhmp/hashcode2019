package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

var (
	fs = afero.NewOsFs()
)

type image struct {
	vertical bool
	tags     []string
}

func NewImage() *image {
	return &image{
		vertical: false,
		tags:     make([]string, 0),
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
	images := make([]*image, 0)

	for i := 0; i < nbrImages; i++ {
		items := strings.Split(lines[i+1], " ")
		img := NewImage()

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
