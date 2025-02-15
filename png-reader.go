/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var true_png_magic []byte = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
var true_png_ihdr_header []byte = []byte{0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52}

func PngValidate(png []byte) error {
	magic := png[0:8]
	if !bytes.Equal(magic, true_png_magic) {
		return errors.New("Bad magic number")
	}

	ihdr_header := png[8:16]
	if !bytes.Equal(ihdr_header, true_png_ihdr_header) {
		return errors.New("Bad IHDR header")
	}

	return nil
}

func PngReadSize(png []byte) (uint32, uint32) {
	width_raw := png[16:20]
	height_raw := png[20:24]

	width := binary.BigEndian.Uint32(width_raw)
	height := binary.BigEndian.Uint32(height_raw)

	return width, height
}

func PngReadBitDepth(png []byte) uint8 {
	return png[24]
}

func PngReadColorType(png []byte) uint8 {
	return png[25]
}

func PngGetColorChannelCount(color_type uint8) uint8 {
	switch color_type {
	case 0: // gray scale
		return 1
	case 2: // RGB
		return 3
	case 3: // indexed color
		return 1
	case 4: // gray scale + alpha
		return 2
	case 6: // RGBA
		return 4
	default:
		return 1 // unknown color type (return 1 as a safe default)
	}
}
