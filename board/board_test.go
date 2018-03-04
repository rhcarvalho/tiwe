package board

import "testing"

func TestBoardValid(t *testing.T) {
	tests := []struct {
		name  string
		board Board
		valid bool
	}{
		{
			name:  "empty is valid",
			board: Board{},
			valid: true,
		},
		{
			name: "sets must contain at least 3 tiles",
			board: Board{}.
				Add(Tile{2, Red}, Tile{3, Red}).
				Add(Tile{7, Red}, Tile{7, Green}),
			valid: false,
		},
		{
			name:  "runs must be contiguous",
			board: Board{}.Add(Tile{2, Red}, Tile{3, Red}, Tile{5, Red}),
			valid: false,
		},
		{
			name:  "runs must use tiles of the same color",
			board: Board{}.Add(Tile{2, Red}, Tile{3, Red}, Tile{4, Green}),
			valid: false,
		},
		{
			name:  "sets can contain groups",
			board: Board{}.Add(Tile{7, Red}, Tile{7, Green}, Tile{7, Blue}),
			valid: true,
		},
		{
			name:  "groups must not contain repeated colors",
			board: Board{}.Add(Tile{7, Red}, Tile{7, Green}, Tile{7, Red}),
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.board.Valid(); got != tt.valid {
				t.Fatalf("got %v, want %v", got, tt.valid)
			}
		})
	}
}
