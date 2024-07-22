package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/wmattei/minceraft/pkg/engine"
)

type TextureSide string

const (
	TopText    TextureSide = "top"
	BottomText TextureSide = "bottom"
	SideText   TextureSide = "side"
)

type Texture struct {
	ref uint32

	ColorStr string `json:"color"`
	Color    *Color `json:"-"`
	Path     string `json:"path"`
	Index    int    `json:"-"`
}

type TextureFile map[BlockType]map[TextureSide]Texture

func LoadTextures() map[string]Texture {
	texturesFile, err := os.Open("assets/textures/texture_atlas.json")
	if err != nil {
		panic(err)
	}

	var file TextureFile

	json.NewDecoder(texturesFile).Decode(&file)

	var result = make(map[string]Texture)

	index := 0

	for blockType, textures := range file {
		for side, texture := range textures {
			ref, err := engine.LoadTexture(texture.Path)
			if err != nil {
				panic(err)
			}

			texName := string(blockType) + string(side)
			if texture.ColorStr != "" {
				colorParts := strings.Split(texture.ColorStr, ",")
				r, _ := strconv.Atoi(colorParts[0])
				g, _ := strconv.Atoi(colorParts[1])
				b, _ := strconv.Atoi(colorParts[2])

				texture.Color = &Color{R: r, G: g, B: b}

			}
			texture.ref = ref
			texture.Index = index
			result[texName] = texture

			index++
		}
	}

	return result
}
