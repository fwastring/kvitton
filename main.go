package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	// "code.sajari.com/docconv"
	// "gopkg.in/gographics/imagick.v3/imagick"
	// "github.com/otiai10/gosseract/v2"
	// "fmt"

	"github.com/fwastring/kvitton/database"
	"github.com/fwastring/kvitton/shell"
	"github.com/gin-gonic/gin"
)

func addReceipt(c *gin.Context) {

	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	fileName := "./file.pdf"
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, fileName)
	shell.Shellout(fmt.Sprintf("pdftoppm -png %s receipt", fileName))
	shell.Shellout("tesseract receipt-1.png - -l eng >> file.txt")
	data, _ := ioutil.ReadFile("file.txt")
	content := string(data)

    // Split the content into paragraphs using double newlines
    paragraphs := strings.Split(content, "\n\n")

    // Check if there are at least two paragraphs
    if len(paragraphs) < 2 {
        fmt.Println("The file does not contain enough paragraphs.")
        return
    }

    // Print the second paragraph
	block := paragraphs[1]

	lines := strings.Split(block, "\n")
	var number int

	for i, line := range lines {
		blocks := strings.Split(line, " ")
		price := blocks[len(blocks)-1]
		item := strings.Join(blocks[0:len(blocks)-1], " ")
		database.Set(item, price)
		number = i
	}

	os.Remove("receipt-1.png")
	os.Remove(fileName)
	os.Remove("file.txt")

	c.String(http.StatusOK, fmt.Sprintf("%s items added!", fmt.Sprint(number+1)))
}

func getAllItems(c *gin.Context) {
	items, err := database.GetAll()
	if err != nil {

	}
	c.IndentedJSON(http.StatusOK, items)
}


func main() {
	router := gin.Default()
	router.POST("/receipt", addReceipt)
	router.GET("/receipt", getAllItems)

    router.Run("localhost:8080")
}
