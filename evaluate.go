package main

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

  return turn * evaluation
}

// A simple evaluation function, assigning points to each square.
func (c Chessboard) pointsAtSquare(s int) int {
  score := 0
  switch c.boardSquares[s] % 10 {
  case 5:
    score = 900
  case 4:
    score = 500
  case 3:
    score = 350
  case 2:
    score = 300
  case 1:
    score = c.pawnPointsAtSquare(s)
  }

  return (1 - 2*(int(c.boardSquares[s]) / 10)) * score
}

// Basic optimizations of pawns, centered pawns > points.
func (c Chessboard) pawnPointsAtSquare(s int) int {
  row := rowFromPosition(s)
  col := colFromPosition(s)

  color := int(c.boardSquares[s] / 10)
  backRank := (1 - color) * 7

  advanced := 6 - (backRank - row) * (color * 2 - 1)
  score := 100 + 4*(advanced)

  if (row == 3 || row == 4) && (col == 3 || col == 4) {
    score += 10
  }

  return score
}
