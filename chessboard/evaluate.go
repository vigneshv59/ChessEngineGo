package chessboard

// Evaluate the position with some shallow parameters.
func (c Chessboard) Evaluate() int {
  evaluation := 0
  turn := 0

  if !c.turn {
    turn = 1
  } else {
    turn = -1
  }

  for i := 0; i < 64; i++ {
    evaluation += c.pointsAtSquare(i)
  }

  evaluation += c.pawnStructureBonus()

  return turn * evaluation
}

// A simple evaluation function, assigning points to each square.
func (c Chessboard) pointsAtSquare(pos int) int {
  score := 0
  switch c.boardSquares[pos] % 10 {
  case 5:
    score = c.evaluateQueen(pos)
  case 4:
    score = c.evaluateRook(pos)
  case 3:
    score = c.evaluateBishop(pos)
  case 2:
    score = c.evaluateKnight(pos)
  case 1:
    score = c.pawnPointsAtSquare(pos)
  }

  return (1 - 2*(int(c.boardSquares[pos]) / 10)) * score
}

func (c Chessboard) evaluateQueen(pos int) int {
  return 900 + 2 * (len(c.LegalMovesFromSquare(pos)) - 11)
}

func (c Chessboard) evaluateRook(pos int) int {
  return 500 + 4 * (len(c.LegalMovesFromSquare(pos)) - 7)
}

func (c Chessboard) evaluateBishop(pos int) int {
  return 330 + 7 * (len(c.LegalMovesFromSquare(pos)) - 7)
}

func (c Chessboard) evaluateKnight(pos int) int {
  return 300 + 20 * (len(c.LegalMovesFromSquare(pos)) - 6)
}

func (c Chessboard) pawnStructureBonus() int {
  whiteCols := make([][]int, 8, 8)
  blackCols := make([][]int, 8, 8)

  for i, _ := range(c.boardSquares) {
    row := rowFromPosition(i)
    col := colFromPosition(i)

    if c.validPiecePawn(i) {
      if c.pieceColorOnPosition(i) == 0 {
        whiteCols[col] = append(whiteCols[col], row)
      } else {
        blackCols[col] = append(blackCols[col], row)
      }
    }
  }

  score := scorePawnColsArray(whiteCols, 0) - scorePawnColsArray(blackCols, 1)

  return score
}

// Takes an array of pawn positions by column and scores the connectedness.
func scorePawnColsArray(colorCols [][]int, color int) int {
  bonus := 0

  for i, col := range(colorCols) {
    if len(col) > 1 {
      bonus -= 5 * (len(col) - 1)
    }

    for _, r := range(col) {
      if i > 0 {
        for _, rp := range(colorCols[i-1]) {
          if r - rp == (color * 2) - 1 {
            bonus += 2
          }
        }
      }

      if i < 7 {
        for _, rp := range(colorCols[i+1]) {
          if r - rp == 1 {
            bonus += 2
          }
        }
      }
    }
  }

  return bonus
}

// Basic optimizations of pawns, centered pawns > points.
func (c Chessboard) pawnPointsAtSquare(s int) int {
  row := rowFromPosition(s)
  col := colFromPosition(s)

  color := int(c.boardSquares[s] / 10)
  backRank := (1 - color) * 7

  advanced := 6 - (backRank - row) * (color * 2 - 1)
  score := 100 + 8*(advanced)

  if (row == 3 || row == 4) && (col == 3 || col == 4) {
    score += 10
  }

  return score
}
