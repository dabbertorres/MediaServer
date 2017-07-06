package urlgen

import (
	"bufio"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultDir         = "dat"
	adjectivesFilename = "adjectives.txt"
	animalsFilename    = "animals.txt"
	adverbsFilename    = "adverbs.txt"
	verbsFilename      = "verbs.txt"
)

var (
	// use our rand.Rand so we don't disturb a calling program potentially using the rand package
	rng *rand.Rand

	adjectives []string
	animals    []string
	adverbs    []string
	verbs      []string
)

func Gen() string {
	var (
		adj = rng.Intn(len(adjectives))
		an  = rng.Intn(len(animals))
		adv = rng.Intn(len(adverbs))
		v   = rng.Intn(len(verbs))
	)

	return adjectives[adj] + animals[an] + adverbs[adv] + verbs[v]
}

func Load() error {
	return LoadDir(defaultDir)
}

func LoadDir(dir string) error {
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	var (
		adjectivesFound = false
		adverbsFound    = false
		animalsFound    = false
		verbsFound      = false
	)

	for _, cd := range contents {
		name := cd.Name()
		switch {
		case name == adjectivesFilename:
			adjectivesFound = true

		case name == adverbsFilename:
			adverbsFound = true

		case name == animalsFilename:
			animalsFound = true

		case name == verbsFilename:
			verbsFound = true
		}
	}

	if !adjectivesFound {
		return errors.New("adjectives " + os.ErrNotExist.Error())
	}

	if !adverbsFound {
		return errors.New("adverbs " + os.ErrNotExist.Error())
	}

	if !animalsFound {
		return errors.New("animals " + os.ErrNotExist.Error())
	}

	if !verbsFound {
		return errors.New("verbs " + os.ErrNotExist.Error())
	}

	if err := fill(&adjectives, filepath.Join(dir, adjectivesFilename)); err != nil {
		return err
	}

	if err := fill(&animals, filepath.Join(dir, animalsFilename)); err != nil {
		return err
	}

	if err := fill(&adverbs, filepath.Join(dir, adverbsFilename)); err != nil {
		return err
	}

	if err := fill(&verbs, filepath.Join(dir, verbsFilename)); err != nil {
		return err
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
