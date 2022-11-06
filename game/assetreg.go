package game

import (
	"embed"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// handles asset loading
// (C) Luminoso 2022 / All Rights Reserved

type AssetRegistry struct {
	Img map[string]*ebiten.Image
}

func LoadAssets(DEV_MODE bool, embedFS embed.FS) AssetRegistry {
	var result AssetRegistry

	//== Image Loading ==
	var imgfiles []string
	result.Img = make(map[string]*ebiten.Image)
	if DEV_MODE {
		filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if strings.Contains(path, "data\\sprite\\") {
					imgfiles = append(imgfiles, strings.Replace(path, "\\", "/", -1))
				}
				return nil
			})
	} else {
		fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, err error) error {
			if strings.Contains(path, "data/sprite/") {
				imgfiles = append(imgfiles, path)
			}
			return nil
		})
	}
	log.Print("[Info] [AssetReg] Loading Images...")
	for _, file := range imgfiles {
		//sanity check
		if !strings.Contains(file, ".png") {
			continue
		}
		//create our reg id
		id := strings.Replace(file, ".png", "", -1)
		id = strings.Replace(id, "data/sprite/", "", -1)
		log.Printf("[Info] [AssetReg] Loading Image w/ ID %s\n", id)
		var img image.Image
		if DEV_MODE {
			data, _ := os.Open(file)
			img, _, _ = image.Decode(data)
		} else {
			data, _ := embedFS.Open(file)
			img, _, _ = image.Decode(data)
		}
		result.Img[id] = ebiten.NewImageFromImage(img)
	}
	log.Printf("[Info] [AssetReg] Loaded %d images.\n", len(result.Img))
	return result
}
