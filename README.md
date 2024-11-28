# OCRVerwatch
## Overview
The program:
1. Locate replay codes using OpenCV by locating areas of bright orange
2. Form a bounding rectangles around those areas
![A photo of the detected replay codes](/cmd/getreplaycodes/assets/docs/output.jpg)

3. Optimise the image for OCR
    - Resize for bigger text during OCR
    - Greyscale to remove color data
    - Threshold (full black or white only) to improve contrast
    - Invert from white text, black background, to black text, white background (better OCR results)  
![A cropped and optimized image of one overwatch replay code](/cmd/getreplaycodes/assets/docs/cropped-and-optimized.jpg)
4. Order detected areas from top of the page going down
5. Crop and save those areas for OCR
6. OCR using an allowlist based off of Overwatch replay code's specific format to improve correctness (No O, U, I, or L)
7. Print out in order

## Usage
```
$ go run cmd/getreplaycodes/main.go --help
Usage of cmd/getreplaycodes/main:
  -image string
        Path to the image file
```

e.g.
```
$ go run cmd/getreplaycodes/main.go --image cmd/getreplaycodes/assets/test/screenshot1.jpg 
859MR0
QGMZ0K
6MZ63A
QNAVY1
SNRJ59
9MKHMB
```

The output contains each Overwatch Replay Code from the source image, from top to bottom.

## Dependencies
### Tesseract
```
tesseract --version
tesseract 5.3.4
 leptonica-1.82.0
```

Installed per https://github.com/tesseract-ocr/tessdoc/blob/main/Installation.md#platforms


### OpenCV
Installed via the gocv makefile
https://github.com/hybridgroup/gocv?tab=readme-ov-file#ubuntulinux

As of writing, some dependencies in the makefile are outdated for Ubuntu 24.04 (libtbb2 and libfreetype6-dev) and it will fail to install. Remove them from `DEBS` and `make install` again.

Proper install from zero guide coming soon.

# TODO
- [ ] Train tesseract on exact font for better OCR results https://github.com/afnleaf/ow2-replaycode-ocr/pull/1#issuecomment-2505553741
- [ ] Discord bot integration
- [ ] Proper depedency install instructions, either via Makefile or better README.md

# Potential future work
- [ ] Website integration?
