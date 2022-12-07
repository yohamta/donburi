package fonts

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type FontName string

const (
	Excel FontName = "excel"
)

func (f FontName) Get() font.Face {
	return getFont(f)
}

var (
	fonts = map[FontName]font.Face{}
)

func LoadFont(name FontName, ttf []byte) {
	fontData, _ := truetype.Parse(ttf)
	fonts[name] = truetype.NewFace(fontData, &truetype.Options{Size: 10})
}

func getFont(name FontName) font.Face {
	f, ok := fonts[name]
	if !ok {
		panic(fmt.Sprintf("Font %s not found", name))
	}
	return f
}
