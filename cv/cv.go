package cv

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"gocv.io/x/gocv"
)

func FindColouredRects(img gocv.Mat, targetColor color.RGBA, tolerance int) ([]image.Rectangle, error) {
	// Convert the target color to the gocv Scalar format
	// NOTE: BGR not RGB!
	lowerBound := gocv.NewScalar(

		float64(max(0, int(targetColor.B)-tolerance)),
		float64(max(0, int(targetColor.G)-tolerance)),
		float64(max(0, int(targetColor.R)-tolerance)),
		0,
	)
	upperBound := gocv.NewScalar(
		float64(min(255, int(targetColor.B)+tolerance)),
		float64(min(255, int(targetColor.G)+tolerance)),
		float64(min(255, int(targetColor.R)+tolerance)),
		0,
	)

	// Create a mask for the desired color range
	mask := gocv.NewMat()
	defer mask.Close()
	gocv.InRangeWithScalar(img, lowerBound, upperBound, &mask)

	// Dilate the mask to remove holes
	mask = dilateMask(mask)

	// Save the mask for debugging
	if ok := gocv.IMWrite("mask.jpg", mask); !ok {
		return nil, fmt.Errorf("error: could not write mask image")
	}

	// Find contours in the mask
	contours := gocv.FindContours(mask, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// Copy the original image into a separate debug image
	debugImg := img.Clone()
	defer debugImg.Close()

	// Draw rectangles around detected contours on the debug image
	var boundingBoxes []image.Rectangle
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		rect := gocv.BoundingRect(contour)
		boundingBoxes = append(boundingBoxes, rect)

		gocv.Rectangle(&debugImg, rect, color.RGBA{0, 255, 0, 0}, 2) // Green rectangle
	}

	// Save the debug image with the rectangles
	if ok := gocv.IMWrite("output.jpg", debugImg); !ok {
		return nil, fmt.Errorf("error: could not write output image")
	}

	return boundingBoxes, nil
}

func CropBoundingBox(img gocv.Mat, boundingBox image.Rectangle, outputName string) error {
	croppedImg := img.Region(boundingBox)
	if ok := gocv.IMWrite(outputName, croppedImg); !ok {
		return fmt.Errorf("error: could not write cropped image to %q", filepath.Join(outputName))
	}

	return nil
}


// OptimiseForTextClarity pre-processes the image to give the best environment for OCR.
func OptimiseForTextClarity(img gocv.Mat) {
	// Convert to grayscale
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	
	// Apply Gaussian Blur to smooth out noise
	gocv.GaussianBlur(img, &img, image.Pt(15, 15), 0, 0, gocv.BorderDefault)

	// Apply binary thresholding (black and white only)
	gocv.Threshold(img, &img, 200, 255, gocv.ThresholdBinary)

	// Invert for white text on black background
	gocv.BitwiseNot(img, &img)
}

func dilateMask(mask gocv.Mat) gocv.Mat {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(10, 10))
	defer kernel.Close()

	gocv.Dilate(mask, &mask, kernel)
	return mask
}

// Utility functions for clamping values
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

