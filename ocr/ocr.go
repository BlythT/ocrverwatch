package ocr

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func ReadTextFromImg(filepath string, whitelist string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImage(filepath); err != nil {
		return "", fmt.Errorf("setting image %q, %w", filepath, err)
	}
	client.Languages = []string{"eng"}
	
	if err := client.SetWhitelist(whitelist); err != nil {
		return "", fmt.Errorf("setting whitelist %q, %w", whitelist, err)
	}
	if err := client.SetPageSegMode(gosseract.PSM_SINGLE_WORD); err != nil {
		return "", fmt.Errorf("setting pageseg mode %q, %w", gosseract.PSM_SINGLE_WORD, err)
	}
	if err := client.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK); err != nil {
		return "", fmt.Errorf("setting pageseg mode %q, %w", gosseract.PSM_SINGLE_BLOCK, err)
	}
	if err := client.SetPageSegMode(gosseract.PSM_RAW_LINE); err != nil {
		return "", fmt.Errorf("setting pageseg mode %q, %w", gosseract.PSM_SINGLE_BLOCK, err)
	}

	text, err := client.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}