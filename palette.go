package main

import (
	"image/color"
	"strconv"
)

var PALETTE = []string{"#060608", "#141013", "#3b1725", "#73172d", "#b4202a", "#df3e23", "#fa6a0a", "#f9a31b", "#ffd541", "#fffc40", "#d6f264", "#9cdb43", "#59c135", "#14a02e", "#1a7a3e", "#24523b", "#122020", "#143464", "#285cc4", "#249fde", "#20d6c7", "#a6fcdb", "#ffffff", "#fef3c0", "#fad6b8", "#f5a097", "#e86a73", "#bc4a9b", "#793a80", "#403353", "#242234", "#221c1a", "#322b28", "#71413b", "#bb7547", "#dba463", "#f4d29c", "#dae0ea", "#b3b9d1", "#8b93af", "#6d758d", "#4a5462", "#333941", "#422433", "#5b3138", "#8e5252", "#ba756a", "#e9b5a3", "#e3e6ff", "#b9bffb", "#849be4", "#588dbe", "#477d85", "#23674e", "#328464", "#5daf8d", "#92dcba", "#cdf7e2", "#e4d2aa", "#a08662", "#796755", "#5a4e44", "#423934"}
var TABLE = func() map[rune](color.RGBA) {
	lookup := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '+', '-'}

	table := make(map[rune](color.RGBA))

	table['0'] = color.RGBA{0, 0, 0, 0}
	for i, v := range PALETTE {
		table[lookup[i+1]] = hexToRgba(v)
	}
	return table
}()

func hexToRgba(hex string) color.RGBA {
	values, err := strconv.ParseUint(hex[1:], 16, 32)
	if err != nil {
		panic(err)
	}

	return color.RGBA{
		R: uint8(values >> 16),
		G: uint8(values >> 8),
		B: uint8(values),
		A: 255,
	}
}

func Lookup(c rune) color.RGBA {
	return TABLE[c]
}
