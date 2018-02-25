package board

type Board [][]Tile

func (b Board) Valid() bool {
	return true
}

type Tile uint64
