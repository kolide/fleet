package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	debDir      = "deb"
	rpmDir      = "rpm"
	binDir      = "bin"
	repoBaseURL = "https://dl.kolide.co/"
)

type pkg struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Kind        string `json:"kind"`
}
type metadata struct {
	Current  []pkg `json:"current"`
	Previous []pkg `json:"previous"`
}

func main() {
	var (
		flRepoPath   = flag.String("repo", "", "path to binary packages repo")
		flCurrentTag = flag.String("git-tag", "", "the tag for the latest kolide release")
	)
	flag.Parse()
	m, err := getMetadata(*flRepoPath, *flCurrentTag)
	if err != nil {
		log.Fatal(err)
	}
	metadataFilePath := filepath.Join(*flRepoPath, "metadata.json")
	os.Remove(metadataFilePath)
	f, err := os.Create(metadataFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		log.Fatal(err)
	}
}

func getMetadata(repoPath, current string) (*metadata, error) {
	var m metadata
	walkFn := func(dir string) filepath.WalkFunc {
		return func(path string, info os.FileInfo, err error) error {
			switch ext := filepath.Ext(path); ext {
			case ".rpm", ".deb", ".zip":
				p := pkg{
					Name:        info.Name(),
					DownloadURL: repoBaseURL + dir + "/" + info.Name(),
					Kind:        dir,
				}
				if isCurrent(info.Name(), current, dir) {
					m.Current = append(m.Current, p)
					return nil
				}
				m.Previous = append(m.Previous, p)
			}
			return nil
		}
	}
	dirs := []string{debDir, rpmDir, binDir}
	for _, dir := range dirs {
		err := filepath.Walk(filepath.Join(repoPath, dir), walkFn(dir))
		if err != nil {
			return nil, errors.Wrapf(err, "walking %s", repoPath)
		}
	}
	return &m, nil
}

// determines wether the file is the current version
// parses the filename based on the conventions for rpms and debs
// set by `fpm`. Unfortunately it doesn't seem possible to keep
// the filename format the same for the different filetypes.
func isCurrent(have, current, kind string) bool {
	switch kind {
	case "bin":
		binSplit := strings.SplitN(have, "_", 2)[1]
		binSplit = strings.TrimSuffix(binSplit, ".zip")
		return binSplit == current
	case "deb":
		debSplit := strings.SplitN(have, "_", 3)[1]
		return debSplit == current
	case "rpm":
		rpmSplit := strings.SplitN(have, "-", 3)[1]
		rpmSplit = strings.Replace(rpmSplit, "_", "-", -1)
		return rpmSplit == current
	default:
		return false
	}
}
