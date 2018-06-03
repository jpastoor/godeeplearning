package main

import (
	"math/rand"
	"fmt"
)

func main() {

	fmt.Println("var (")
	fmt.Println("hashes = map[bool]map[int]map[int]uint64{")

	hashes := make(map[bool]map[int]map[int]uint64)
	for _, player := range []bool{true, false} {
		hashesPlayer := make(map[int]map[int]uint64)
		fmt.Printf("%s: {\n", player)

		for row := 1; row <= 19; row++ {
			fmt.Printf("%d: {\n", row)
			hashesRow := make(map[int]uint64)
			for col := 1; col <= 19; col++ {

				fmt.Printf("%d: %b,\n", col, rand.Uint64())

			}
			hashesPlayer[row] = hashesRow
			fmt.Println("	},")
		}
		hashes[player] = hashesPlayer
		fmt.Println("	},")
	}

	fmt.Println("	}")
	fmt.Println(")")
}

