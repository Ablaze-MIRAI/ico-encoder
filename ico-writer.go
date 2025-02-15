/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"encoding/binary"
	"errors"
	"os"
)

type IcoHeader struct {
	Reserved   uint16 // Reserved
	Type       uint16 // ICO Type, .ico: 1, .cur: 2
	ImageCount uint16 // Number of images
}

type IcoDirectory struct {
	Width             uint8  // width
	Height            uint8  // height
	ColorPaletteCount uint8  // Color palette count
	Reserved          uint8  // Reserved
	ColorPlanesCount  uint16 // Color planes count
	BitsPerPixel      uint16 // Bits per pixel
	DataLength        uint32 // Data length
	DataOffset        uint32 // Data offset
}

func uint16ToBytesLittle(v uint16) []byte {
	return binary.LittleEndian.AppendUint16([]byte{}, v)
}

func uint32ToBytesLittle(v uint32) []byte {
	return binary.LittleEndian.AppendUint32([]byte{}, v)
}

func IcoWrite(icons [][]byte, path string) error {
	var ico []byte

	if len(icons) > 65535 {
		return errors.New("The maximum number of images is 65535")
	}

	header := IcoHeader{
		Reserved:   0,
		Type:       1,
		ImageCount: uint16(len(icons)),
	}

	var directories []IcoDirectory
	for i, icon := range icons {
		if err := PngValidate(icon); err != nil {
			return errors.New("png reader: " + err.Error())
		}

		width, height := PngReadSize(icon)

		bit_depth := PngReadBitDepth(icon)
		color_type := PngReadColorType(icon)
		bits_per_pixel := bit_depth * PngGetColorChannelCount(color_type)

		if width > 256 || height > 256 {
			return errors.New("The maximum width/height is 256")
		}

		if width == 256 {
			width = 0 // if width is 256, write 0x00
		}
		if height == 256 {
			height = 0 // if height is 256, write 0x00
		}

		var offset uint32 = uint32( // offset without data length
			6 /* ico file header size */ + 16*len(icons), /* directories size */
		)
		for j := 0; j < i; j++ {
			offset += uint32(len(icons[j]))
		}

		directories = append(directories,
			IcoDirectory{
				Width:             uint8(width),
				Height:            uint8(height),
				ColorPaletteCount: 0,
				Reserved:          0,
				ColorPlanesCount:  1,
				BitsPerPixel:      uint16(bits_per_pixel),
				DataLength:        uint32(len(icon)),
				DataOffset:        offset,
			},
		)
	}

	// Write .ico

	// header
	ico = append(ico,
		uint16ToBytesLittle(header.Reserved)..., // Reserved
	)
	ico = append(ico,
		uint16ToBytesLittle(header.Type)..., // ICO Type, .ico: 1, .cur: 2
	)
	ico = append(ico,
		uint16ToBytesLittle(header.ImageCount)..., // Number of images
	)
	// directories
	for _, directory := range directories {
		ico = append(ico,
			directory.Width,             // width
			directory.Height,            // height
			directory.ColorPaletteCount, // Color palette count
			directory.Reserved,          // Reserved
		)
		ico = append(ico,
			uint16ToBytesLittle(directory.ColorPlanesCount)..., // Color planes
		)
		ico = append(ico,
			uint16ToBytesLittle(directory.BitsPerPixel)..., // Bits per pixel
		)
		ico = append(ico,
			uint32ToBytesLittle(directory.DataLength)..., // Data length
		)
		ico = append(ico,
			uint32ToBytesLittle(directory.DataOffset)..., // Data offset
		)
	}
	// icons
	for _, icon := range icons {
		ico = append(ico,
			icon...,
		)
	}

	if err := os.WriteFile(path, ico, 0644); err != nil {
		return err
	}

	return nil
}
