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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.board.Valid(); got != tt.valid {
				t.Fatalf("got %v, want %v", got, tt.valid)
			}
		})
	}
}
