package main

func (b Board) IsPointAnEye(point Point, player Player) bool {
	_, exists := b.get(point)
	if exists {
		return false
	}

	// If one of the neighbors is the other player, the Point is not an eye
	for _, neighbor := range point.Neighbors() {
		neighborPlayer, neighborExists := b.get(neighbor)
		if neighborExists && !neighborPlayer.equals(player) {
			return false
		}
	}

	friendlyCorners := 0
	offboardCorners := 0

	corners := []Point{
		{point.Row - 1, point.Col - 1},
		{point.Row - 1, point.Col + 1},
		{point.Row + 1, point.Col - 1},
		{point.Row + 1, point.Col + 1},
	}

	for _, corner := range corners {
		if b.isOnGrid(corner) {
			cornerPlayer, cornerExists := b.get(corner)
			if cornerExists && cornerPlayer.equals(player) {
				friendlyCorners++
			}

		} else {
			offboardCorners++
		}
	}

	if offboardCorners > 0 {
		return offboardCorners+friendlyCorners == 4
	}

	return friendlyCorners >= 3
}
