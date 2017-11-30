package chessboard

import "fmt"

// Runs the recursive Alpha-Beta function, returns the score, and a tuple
// representing the move.

// BUG: Very slow for depth >= 6, this shouldn't be happening, some
// optimizations probably can be added.
func (c *Chessboard) alphaBetaHelper(a int, b int, depth int) (int, []int) {
  if depth == 0 {
    // TODO: Quiesce instead of straight-up evaluation.
    return c.Evaluate(), make([]int, 0)
  }

  var moves [][]int = c.AllLegalMoves()
  var move []int
  alpha := a
  beta := b

  color := 0
  if c.turn {
    color = 1
  }

  // Checkmate
  if c.kingInCheck(color) && len(moves) == 0 {
    return 9999, make([]int, 0)
  }

  // Stalemate
  if !(c.kingInCheck(color)) && len(moves) == 0 {
    return 0, make([]int, 0)
  }

  for _, m := range(moves) {
    _, restore := c.MakeMoveWithRestore(m[0], m[1], "")
    score, _ := c.alphaBetaHelper(-beta, -alpha, depth -1)
    score = -score
    c.RestoreBoard(restore)

   if score >= beta {
      return beta, make([]int, 0)
   }

    if (score > alpha) {
      alpha = score
      move = m
    }
  }

  return alpha, move
}

// Calls the Alpha-Beta helper with a seed alpha and beta value, along with
// the given depth.
func (c Chessboard) AlphaBeta(depth int) (int, []int) {
  var cm []int

  if c.book.name != "" {
    cm = c.book.pickMove(c.BookHash())
  }

  score, m := c.alphaBetaHelper(-10000, 10000, depth)

  if len(cm) == 2 {
    m = cm
  }

  if c.turn {
    return -score, m
  }


  return score, m
}

// Uses the negamax algorithm to find a move.
func (c Chessboard) negaMax(depth int) int {
  score, move := c.negaMaxHelper(depth)

  fmt.Println(move)

  if c.turn {
    return -score
  }

  return score
}

// Uses the negamax algorithm to recursively find a move.
// This is less efficient than alpha-beta. Normally do not use.
func (c *Chessboard) negaMaxHelper(depth int) (int, []int) {
  if (depth == 0) {
    return c.Evaluate(), make([]int, 0)
  }

  var max int = -1E9

  // prevHist := hist
  var moves [][]int = c.AllLegalMoves()
  move := make([]int, 0)

  for _, m := range(moves) {
    _, restore := c.MakeMoveWithRestore(m[0], m[1], "")
    score, _ := c.negaMaxHelper(depth - 1)
    score = -score
    c.RestoreBoard(restore)

    if (score > max) {
      max = score
      move = m
    }
  }

  return max, move
}
