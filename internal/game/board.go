package game

type Board [3][3]string

func NewBoard() Board {
	return Board{}
}

func (b Board) IsFull() bool {
	for _, row := range b {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func (b Board) Winner() string {
	// Check rows, columns, and diagonals
	for i := 0; i < 3; i++ {
		if b[i][0] != "" && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			return b[i][0]
		}
		if b[0][i] != "" && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			return b[0][i]
		}
	}
	if b[0][0] != "" && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return b[0][0]
	}
	if b[0][2] != "" && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return b[0][2]
	}
	return ""
}
