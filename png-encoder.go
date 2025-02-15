/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"bytes"
	"image"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func PngResize(png_path string, width int, height int) ([]byte, error) {
	f, err := os.Open(png_path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	new_img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.BiLinear.Scale(new_img, new_img.Bounds(), img, img.Bounds(), draw.Over, nil)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, new_img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
