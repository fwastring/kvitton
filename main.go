package main

import (
	"fmt"
	// "log"
	// "net/http"

	// "code.sajari.com/docconv"
	"gopkg.in/gographics/imagick.v3/imagick"
	"github.com/otiai10/gosseract/v2"
	"github.com/fwastring/kvitton/shell"
)

func main() {
	imagick.Initialize()
    defer imagick.Terminate()
    mw := imagick.NewMagickWand()
    defer mw.Destroy()
    mw.ReadImage("data/receipt.pdf")
    mw.SetIteratorIndex(0)        // This being the page offset
    mw.SetImageFormat("png")
    mw.WriteImage("test.png")

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("test.png")
	text, _ := client.Text()
	fmt.Println(text)
}
