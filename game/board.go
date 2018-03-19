package game

type Board [][]Tile

func (b Board) Valid() bool {
	for _, set := range b {
		if !isRun(set) && !isGroup(set) {
			return false
		}
	}
	return true
}

func isRun(set []Tile) bool {
	if len(set) < 3 {
		return false
	}
	for i := 0; i < len(set)-1; i++ {
		if set[i+1].Value-set[i].Value != 1 {
			return false
		}
		if set[i].Color != set[i+1].Color {
			return false
		}
	}
	return true
}

func isGroup(set []Tile) bool {
	if len(set) < 3 {
		return false
	}
	var seen int64
	for i := 0; i < len(set)-1; i++ {
		if set[i].Value != set[i+1].Value {
			return false
		}
		seen |= 1 << set[i].Color
		if (seen>>set[i+1].Color)&1 == 1 {
			return false
		}
	}
	return true
}

func (b Board) Add(ts ...Tile) Board {
	return append(b, ts)
}

type Tile struct {
	Value uint64
	Color
}

type Color uint64

const (
	Red Color = iota
	Green
	Blue
	Yellow
)
