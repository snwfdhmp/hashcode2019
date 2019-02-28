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

	imageChoice       string
	maxTries          int
	minimumScoreToLog int

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
	flag.StringVar(&imageChoice, "image", "d", "image to choose")
	flag.IntVar(&maxTries, "max_tries", 100, "maximum tries to make")
	flag.IntVar(&minimumScoreToLog, "minimum", 100000, "minimum score to save")
	flag.Parse()

	fmt.Printf("Processing image %s\n", strings.ToUpper(imageChoice))

	images, err := parseFile(imagePath[imageChoice])
	if err != nil {
		fmt.Printf("fatal: %v\n", err)
		return
	}

	maxScore := 0
	for try := 0; try < maxTries; try++ {
		saveImg := make([]image, len(images))
		copy(saveImg, images)

		fmt.Printf("Test n°%d ...", try)
		slides := makeSlideShow(saveImg)
		if slides == nil {
			fmt.Printf(" -- watchdog\n")
			continue
		}
		score := getScore(slides)

		fmt.Printf("\rTest n°%d: score %d", try, score)
		if score > maxScore {
			fmt.Print(" HIGH SCORE !!!!")
			maxScore = score
			if score >= minimumScoreToLog {
				afero.WriteFile(fs, fmt.Sprintf("./output/image_%s_score_%d.txt", imageChoice, score), Marshal(slides), 0760)
			}
		}
		fmt.Print("\n")
	}

	fmt.Printf("Done :D Best score: %d\n", maxScore)
}

func makeSlideShow(images []image) []slide {
	slides := make([]slide, 0)
	lengthImg := len(images)

	attempt := 0
	for i := 0; i < lengthImg; i++ {
		if attempt > len(images) {
			return nil
		}
		imgIndex := randomIndex(images)
		img := images[imgIndex]
		if !img.vertical {
			slide := slide{images: []image{img}}
			if i == 0 {
				slides = append(slides, slide)
				continue
			}
			scorePrevision := getTagScore(getTags(slides[i-1]), getTags(slide))
			if scorePrevision > 5 || (scorePrevision > 4 && attempt > 10) || (scorePrevision > 3 && attempt > 20) || (scorePrevision > 2 && attempt > 30) || (scorePrevision >= 1 && attempt > 40) {
				slides = append(slides, slide)
				attempt = 0
				images[imgIndex] = images[len(images)-1]
				images = images[:len(images)-1]
			} else {
				attempt++
				i--
			}
			continue
		}

		vIndex := findFirstVertical(images)
		if vIndex == -1 {
			fmt.Printf("\nWARN: found -1, len(images)=%d\n", len(images))
			continue
		}
		slide := slide{images: []image{img, images[vIndex]}}
		if i == 0 {
			slides = append(slides, slide)
			continue
		}
		scorePrevision := getTagScore(getTags(slides[i-1]), getTags(slide))
		if scorePrevision > 5 || (scorePrevision > 4 && attempt > 10) || (scorePrevision > 3 && attempt > 20) || (scorePrevision > 2 && attempt > 30) || (scorePrevision >= 1 && attempt > 40) {
			slides = append(slides, slide)
			attempt = 0
			images[vIndex] = images[len(images)-1]
			images = images[:len(images)-1]
		} else {
			attempt++
			i--
		}
		slides = append(slides, slide)
		i++
	}

	return slides
}

func randomIndex(items []image) int {
	if len(items) == 1 {
		return 0
	}
	return rand.Intn(len(items) - 1)
}

func findFirstVertical(items []image) int {
	for i := range items {
		if items[i].vertical {
			return i
		}
	}
	return -1
}

func getScore(slides []slide) int {
	score := 0
	for i := 0; i < len(slides)-1; i++ {
		tagsA := getTags(slides[i])
		tagsB := getTags(slides[i+1])
		score += getTagScore(tagsA, tagsB)
	}
	return score
}

func getTags(slide slide) map[string]bool {
	tags := make(map[string]bool)
	for i := range slide.images {
		for j := range slide.images[i].tags {
			tags[slide.images[i].tags[j]] = true
		}
	}

	return tags
}

// returns score from tag arrays comparison
func getTagScore(a, b map[string]bool) int {
	common := 0
	for n := range a {
		if _, ok := b[n]; ok {
			common++
		}
	}

	return min(min(len(a)-common, common), len(b)-common)
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
