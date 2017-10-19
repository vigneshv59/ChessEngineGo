package main

import (
  "errors"
  "regexp"
  "strconv"
  "strings"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var pieceVals = map[string] int8 {
  "K" : 6,
  "Q" : 5,
  "R" : 4,
  "B" : 3,
  "N" : 2,
  "P" : 1,
  "k" : 16,
  "q" : 15,
  "r" : 14,
  "b" : 13,
  "n" : 12,
  "p" : 11,
}

type Chessboard struct {
  boardSquares []int8
  enpassantPos int
}

func NewChessboard(fen string) (board Chessboard, err error) {
  // Verify presence of FEN.
  if len(fen) == 0 {
    err = errors.New("A fen must be provided.")
    return
  }

  board = Chessboard{}

  // Configure the board to have -1 (no piece) on
  // every square.

  a := make([]int8, 64)

  for k := range a {
    a[k] = -1
  }

  board.boardSquares = a

  fenParts := strings.Split(fen, " ")
  square := 0

  // Begin assigning the initial FEN position to the square array
  for _, rune := range fenParts[0] {
    fenPart := string(rune)

    if square > 63 {
      err = errors.New("The FEN is invalid -- 63 squares.")
      return
    }

    padding, e := strconv.Atoi(fenPart)

    if e == nil {
      square += padding
    } else if fenPart != "/" {

      match, _ := regexp.MatchString("[rnbqkpPRNBQK]", fenPart)

      if !match {
        err = errors.New("The FEN is invalid -- invalid character.")
        return
      }

      board.boardSquares[square] = pieceVals[fenPart]

      square += 1

    }

  }

  return
}

// Piece Move Validators:

// Validates a queen move, does not take into account turn.
func (c Chessboard) validMoveKnight(from int, to int) bool {
  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(to)
  toRow := rowFromPosition(from)
  toCol := colFromPosition(to)

  if ((toRow == (fromRow - 2) || toRow == (fromRow + 2)) &&
      (toCol == (fromCol - 1) || toCol == (fromCol + 1))) {

      return true

  }

  if ((toRow == (fromRow - 1) || toRow == (fromRow + 1)) &&
      (toCol == (fromCol - 2) || toCol == (fromCol + 2))) {

      return true

  }

  return false
}

// Validates a queen move, does not take into account turn.
func (c Chessboard) validMoveQueen(from int, to int) bool {
  return c.validMoveRook(from, to) || c.validMoveBishop(from, to)
}

// Validates a rook move, does not take into account turn.
func (c Chessboard) validMoveRook(from int, to int) bool {

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(to)
  toRow := rowFromPosition(from)
  toCol := colFromPosition(to)

  return fromRow == toRow || fromCol == toCol

}

// Validates a bishop move, does not take into account turn.
func (c Chessboard) validMoveBishop(from int, to int) bool {
  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(to)
  toRow := rowFromPosition(from)
  toCol := colFromPosition(to)

  for offset := 1; offset <= 7; offset++ {
    if (toRow == (fromRow + offset) ||
        toRow == (fromRow - offset) &&
        toCol == (fromCol + offset) ||
        toCol == (fromCol - offset)) {

          return true

    }
  }

  return false

}

// Validates a pawn move, does not take into account turn.
func (c Chessboard) validMovePawn(from int, to int) bool {
  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(to)
  toRow := rowFromPosition(from)
  toCol := colFromPosition(to)

  if c.pieceColorOnPosition(from) == 0 {

    // Check for captures
    if c.validColorPiece(to, 1) || (c.enpassantPos == to) {
      if (toRow == fromRow - 1) && (toCol == fromCol + 1) {
        return true
      }

      return false
    }

    // Normal moves and first pawn moves
    if toRow == (fromRow - 1) && (toCol == fromCol) {
      return true
    } else if toRow == 4 && fromRow == 6 && (toCol == fromCol) {
      return true
    }

    return false

  } else if c.pieceColorOnPosition(from) == 1 {

    // Captures
    if c.validColorPiece(to, 0) || (c.enpassantPos == to) {
      if (toRow == fromRow + 1) && (toCol == fromCol + 1) {
        return true
      }

      return false
    }

    // Normal moves and first pawn moves
    if toRow == (fromRow + 1) && (toCol == fromCol) {
      return true
    } else if toRow == 3 && fromRow == 1 && (toCol == fromCol) {
      return true
    }

    return false

  }

  return false
}

// Returns 0 for white pieces and 1 for
// black pieces.
func (c Chessboard) pieceColorOnPosition(pos int) int {
  return int(c.boardSquares[pos] / 10)
}

func (c Chessboard) validColorPiece(pos int, color int) bool {
  return c.pieceColorOnPosition(pos) == color
}

// Helper methods
func rowFromPosition(pos int) int {
  return pos / 8
}

func colFromPosition(pos int) int {
  return pos % 8
}
