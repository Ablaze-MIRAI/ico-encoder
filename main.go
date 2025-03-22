/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	os.Exit(run())
}

func run() int {
	var input_filepath string
	flag.StringVar(&input_filepath, "i", "", "Specify input file path (.png)")

	var output_filepath string
	flag.StringVar(&output_filepath, "o", "", "Specify output file path (.ico)")

	var sizes string
	flag.StringVar(&sizes, "s", "16,24,32,48,256", "Specify sizes of icon (maximum is 256)")

	flag.Parse()

	if input_filepath == "" || output_filepath == "" || sizes == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return 1
	}

	var icons [][]byte

	for _, size := range strings.Split(sizes, ",") {
		size_uint, err := strconv.ParseUint(size, 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", "invalid size: \""+size+"\"")
			return 1
		}

		icon, err := PngResize(input_filepath, int(size_uint), int(size_uint))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			return 1
		}

		icons = append(icons, icon)
	}

	if err := IcoWrite(icons, output_filepath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return 1
	}

	return 0
}
