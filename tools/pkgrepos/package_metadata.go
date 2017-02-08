package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	Name        string `json:"name,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
	Kind        string `json:"kind,omitempty"`
	SHA256      string `json:"sha_256,omitempty"`
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
				hash, err := shaFile(path)
				if err != nil {
					return err
				}
				p := pkg{
					Name:        info.Name(),
					DownloadURL: repoBaseURL + dir + "/" + info.Name(),
					Kind:        dir,
					SHA256:      hash,
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
	// add current release docker hub link
	p := pkg{
		Kind: "docker",
		Name: "kolide/kolide:" + current,
	}
	m.Current = append(m.Current, p)
	return &m, nil
}

func shaFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "open file %s for hashing", f.Name())
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrapf(err, "hash file %s", f.Name())
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
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
