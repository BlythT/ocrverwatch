package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/BlythT/ocrverwatch/cv"
	"github.com/BlythT/ocrverwatch/ocr"
	"gocv.io/x/gocv"
)

// replayCodeOrange was eye-droppered from the orange boxes that hold overwatch Replay codes.
var replayCodeOrange = color.RGBA{
	R: 241,
	G: 100,
	B: 18,
	A: 255,
}

func main() {
	// Define a command-line flag for the image file path
	filepath := flag.String("image", "", "Path to the image file")

	// Parse the flags
	flag.Parse()

	// Check if the image flag is provided
	if filepath == nil || *filepath == "" {
		log.Fatalf("Error: Image file path is required")
	}

	// Check if the file exists
	_, err := os.Stat(*filepath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Error: File does not exist at path: %s\n", *filepath)
		}

		log.Fatalf("Error: Could not check file status: %v\n", err)
	}

	codes, err := getReplayCodes(*filepath)
	if err != nil {
		log.Fatalf("Error: Could not get codes for filepath %s: %v", *filepath, err)
	}

	for _, code := range codes {
		fmt.Println(code)
	}
}

// getReplayCodes by detecting orange boxes, cropping them, pre-processing, and OCRing each one to get the text within.
func getReplayCodes(imageFilePath string) (map[int]string, error) {
	// Read the input image
	img := gocv.IMRead(imageFilePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("error: could not read image from file")
	}
	defer img.Close()

	// Calculate the scaling factor (e.g., scale by 1.5 times)
	scaleFactor := 5.0
	newWidth := int(float64(img.Cols()) * scaleFactor)
	newHeight := int(float64(img.Rows()) * scaleFactor)

	// resize
	gocv.Resize(img, &img, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationLinear)

	boundingBoxes, err := cv.FindColouredRects(img, replayCodeOrange, 50)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("finding coloured rects"), err)
	}

	cv.OptimiseForTextClarity(img)

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "cropped_images_")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up the temp directory after processing

	boundingBoxes = removeSmallBoxes(boundingBoxes, 50, 50)

	// order the bounding boxes from highest in the image to lowest
	orderedBoundingBoxes := orderByHighestPoint(boundingBoxes)
	orderedCodes := make(map[int]string)
	for i, box := range orderedBoundingBoxes {
		filename := strconv.Itoa(i) + ".jpg"
		croppedName := filepath.Join(tmpDir, filename)

		// the code section has a share symbol taking ~1/5 of the left side. Crop in to avoid OCRing it
		// and a white secion ~1/20th on the left which gets mistakenly read as an I in greyscale
		box = cropReplayCode(box)

		// croppedName := filepath.Join(tmpDir, strconv.Itoa(i) + ".jpg")
		if err := cv.CropBoundingBox(img, box, croppedName); err != nil {
			return nil, errors.Join(fmt.Errorf("cropping bounding boxes"), err)
		}

		// By scraping https://owreplays.tv/ I was able to validate a character whitelist.
		// O, I, U, and L are excluded, likely due to similarity with others..
		text, err := ocr.ReadTextFromImg(croppedName, "ABCDEFGHJKMNPQRSTVWXYZ0123456789")
		if err != nil {
			return nil, errors.Join(fmt.Errorf("error: could not read text from image"), err)
		}

		switch strings.ToLower(text) {
		case "vew", "v1ew", "mprt", "1mp0rt": // Due to character whitelist OIUL are not possible
			fmt.Print("Button in image, skipping")
			continue
		}

		orderedCodes[i] = text
	}

	return orderedCodes, nil
}

// orderByHighestPoint returns the image.Rectangles in a map with indexes starting at 0 to preserve order.
func orderByHighestPoint(rects []image.Rectangle) map[int]image.Rectangle {
	sort.Slice(rects, func(i, j int) bool {
		return rects[i].Min.Y < rects[j].Min.Y
	})

	orderedRects := make(map[int]image.Rectangle)
	for i, rect := range rects {
		orderedRects[i] = rect
	}

	return orderedRects
}

// cropReplayCode performs a crop specialized to the replay code box.
func cropReplayCode(rect image.Rectangle) image.Rectangle {
	width := rect.Dx()
	cropLeftAmount := float64(width) / 4.5 // The replay code has a share symbol on the left that takes up ~1/4-1/5th.
	cropRightAmount := width / 20          // The bounding box often leaves whitespace

	croppedBox := image.Rect(rect.Min.X+int(cropLeftAmount), rect.Min.Y, rect.Max.X-cropRightAmount, rect.Max.Y)
	return croppedBox
}

func removeSmallBoxes(rects []image.Rectangle, minWidth, minHeight int) []image.Rectangle {
	var filteredRects []image.Rectangle
	for _, rect := range rects {
		if rect.Dx() >= minWidth && rect.Dy() >= minHeight {
			filteredRects = append(filteredRects, rect)
		}
	}
	return filteredRects
}
