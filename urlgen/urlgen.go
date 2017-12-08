package urlgen

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	adjectivesFilename = "adjectives.txt"
	animalsFilename    = "animals.txt"
	adverbsFilename    = "adverbs.txt"
	verbsFilename      = "verbs.txt"
)

var (
	// use our own rand.Rand so we don't disturb a calling program potentially using the rand package
	rng *rand.Rand

	adjectives []string
	animals    []string
	adverbs    []string
	verbs      []string
)

type Error struct {
	File string
	error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.File, e.error)
}

func Gen() string {
	var (
		adj = rng.Intn(len(adjectives))
		an  = rng.Intn(len(animals))
		adv = rng.Intn(len(adverbs))
		v   = rng.Intn(len(verbs))
	)

	return adjectives[adj] + animals[an] + adverbs[adv] + verbs[v]
}

func LoadDir(dir string) error {
	if err := fill(&adjectives, filepath.Join(dir, adjectivesFilename)); err != nil {
		return Error{adjectivesFilename, err}
	}

	if err := fill(&animals, filepath.Join(dir, animalsFilename)); err != nil {
		return Error{animalsFilename, err}
	}

	if err := fill(&adverbs, filepath.Join(dir, adverbsFilename)); err != nil {
		return Error{adverbsFilename, err}
	}

	if err := fill(&verbs, filepath.Join(dir, verbsFilename)); err != nil {
		return Error{verbsFilename, err}
	}

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	return nil
}

func fill(arr *[]string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()

		if line != "" {
			*arr = append(*arr, line)
		}
	}

	return nil
}
