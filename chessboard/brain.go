package chessboard

import (
  "fmt"
  "time"
)

// Runs the recursive Alpha-Beta function, returns the score, and a tuple
// representing the move.

// BUG: Very slow for depth >= 6, this shouldn't be happening, some
// optimizations probably can be added.
func (c *Chessboard) alphaBetaHelper(a int, b int, depth int, prevMoves [][]int, searchStop *bool) (int, [][]int, int) {
  if depth == 0 || *searchStop {
    if *searchStop {
      fmt.Println("Search Stop")
    }

    // TODO: Quiesce instead of straight-up evaluation.
    return c.Evaluate(), prevMoves, 1
  }

  var moves [][]int = c.AllLegalMoves()
  alpha := a
  beta := b

  color := 0
  if c.turn {
    color = 1
  }

  // Checkmate
  if c.kingInCheck(color) && len(moves) == 0 {
    return -9999, prevMoves, 1
  }

  // Stalemate
  if !(c.kingInCheck(color)) && len(moves) == 0 {
    return 0, prevMoves, 1
  }

  var combined [][]int
  nodesSearched := 0

  for _, m := range(moves) {
    _, restore := c.MakeMoveWithRestore(m[0], m[1], "")
    mHist := append(prevMoves, m)

    score, forwardMoves, nSearch := c.alphaBetaHelper(-beta, -alpha, depth - 1, mHist, searchStop)
    nodesSearched += nSearch

    score = -score
    c.RestoreBoard(restore)

    if score >= beta {
      return beta, make([][]int, 0), nodesSearched
    }

    if (score > alpha) {
      alpha = score
      combined = make([][]int, len(forwardMoves))
      copy(combined, forwardMoves)
    }
  }

  return alpha, combined, nodesSearched
}

// Calls the Alpha-Beta helper with a seed alpha and beta value, along with
// the given depth.
func (c Chessboard) AlphaBeta(depth int, searchStop *bool) (int, []int) {
  var cm []int

  if c.book.name != "" {
    cm = c.book.pickMove(c.BookHash())
  }

  var score int
  var m []int
  var moves [][]int
  var nSearch int
  totalNodes := 0

  start := time.Now()

  for i := 1; i <= depth; i++ {
    score, moves, nSearch = c.alphaBetaHelper(-10000, 10000, i, make([][]int, 0, i), searchStop)

    if *searchStop {
        break
    }

    m = moves[0]

    pv := ""

    for j := 0; j < len(moves); j++ {
      pv += PosToAl(moves[j][0]) + PosToAl(moves[j][1]) + " "
    }

    pv = pv[0:(len(pv) - 1)]

    t := time.Now()
    elapsed := t.Sub(start)
    nps := int(float64(nSearch) / elapsed.Seconds())
    elap := int(elapsed.Seconds() * 1000.0)
    totalNodes += nSearch

    fmt.Printf("info depth %d nodes %d nps %d score cp %d time %d multipv 1 pv %s\n", i, nSearch, nps, score, elap, pv)
  }

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
