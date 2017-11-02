package main

import (
  "errors"
  "regexp"
  "strconv"
  "strings"
  "fmt"
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
  turn bool // false for white's move, true for black's move
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

  // Handle the turn encoded in the fen
  if len(fenParts) > 0 {
    if fenParts[1] == "w" {
      board.turn = false
    }

    if fenParts[1] == "b" {
      board.turn = true
    }
  }

  return
}

// Special Moves Convenience Checks

// Checks for promotion validity, does not take into account turn
func (c Chessboard) attemptedPromotion(from int, to int) bool {
  color := c.pieceColorOnPosition(from)

  if color == 0 && to > -1 && to < 8 {
    return true
  }

  if color == 1 && to > 55 && to < 64 {
    return true
  }

  return false
}

// Makes a move on the board, returning the success/failure as a bool

func (c Chessboard) Move(from int, to int, promopiece string) bool {
  color := c.pieceColorOnPosition(from)
  turn := 0

  if c.turn {
    turn = 1
  }

  // Incorrect turn
  if color != turn {
    return false
  }

  // TODO: Validate promopiece

  if !c.prelimValidMove(from, to) {
    return false
  }

  // TODO: Check/Checkmate validation, update board, handle en passant

  return true
}

// Validates a move on a board with a from index and a to index
// Returns a boolean indicating success of the move.
// This is a preliminary validator, and does not take into account checkmate
// or checks, which will be dealt with after this first pass confirmation

// TODO: En passant capture, castling

func (c Chessboard) prelimValidMove(from int, to int) bool {
  // Cannot move a piece to the same square.
  if from == to {
    return false
  }

  // Cannot capture your own pieces.
  if c.pieceColorOnPosition(from) == c.pieceColorOnPosition(to) {
    return false
  }

  switch c.boardSquares[from] % 10 {
  case 6: // King
    // TODO: Castling

    if c.validMoveKing(from, to) {
      return true
    }
  case 5: // Queen
    if c.validMoveQueen(from, to) && c.validateEmptySquaresBetween(from, to) {
      return true
    }
  case 4: // Rook
    if c.validMoveRook(from, to) && c.validateEmptySquaresBetween(from, to) {
      return true
    }
  case 3: // Bishop
    if c.validMoveBishop(from, to) && c.validateEmptySquaresBetween(from, to) {
      return true
    }
  case 2: // Knight
    if c.validMoveKnight(from, to) {
      return true
    }
  case 1: // Pawn
    if c.validMovePawn(from, to) {
      return true
    }
  default:
    return false
  }

  return false
}

// Piece Existence Validators:
func (c Chessboard) validPieceKing(square int) bool {
  return int(c.boardSquares[square] % 10) == 6
}

func (c Chessboard) validPieceQueen(square int) bool {
  return int(c.boardSquares[square] % 10) == 5
}

func (c Chessboard) validPieceRook(square int) bool {
  return int(c.boardSquares[square] % 10) == 4
}

func (c Chessboard) validPieceBishop(square int) bool {
  return int(c.boardSquares[square] % 10) == 3
}

func (c Chessboard) validPieceKnight(square int) bool {
  return int(c.boardSquares[square] % 10) == 2
}

func (c Chessboard) validPiecePawn(square int) bool {
  return int(c.boardSquares[square] % 10) == 1
}

func (c Chessboard) validPiece(square int) bool {
  return int(c.boardSquares[square]) != -1
}

// Candidate piece Move Validators:

// Validates a queen move, does not take into account turn.
func (c Chessboard) validMoveKnight(from int, to int) bool {
  if !c.validPieceKnight(from) {
      return false
  }

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
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
  if !c.validPieceQueen(from) {
      return false
  }

  return c.validMoveRook(from, to) || c.validMoveBishop(from, to)
}

// Validates a rook move, does not take into account turn.
func (c Chessboard) validMoveRook(from int, to int) bool {
  if !c.validPieceRook(from) {
      return false
  }

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
  toCol := colFromPosition(to)

  return fromRow == toRow || fromCol == toCol

}

// Validates a bishop move, does not take into account turn.
func (c Chessboard) validMoveBishop(from int, to int) bool {
  if !c.validPieceBishop(from) {
      return false
  }

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
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

// Validates a king move, does not take into account turn.
func (c Chessboard) validMoveKing(from int, to int) bool {
  if !c.validPieceKing(from) {
      return false
  }

  if from == to {
    return false
  }

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
  toCol := colFromPosition(to)

  if abs(fromRow - toRow) <= 1 && abs(fromCol - toCol) <= 1 {
    return true
  }

  return false

}

// Validates a pawn move, does not take into account turn.
func (c Chessboard) validMovePawn(from int, to int) bool {
  if !c.validPiecePawn(from) {
      return false
  }

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
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

// Lists the pieces threatening a position
func (c Chessboard) piecesThreateningPos(sq int, color int) []int {
  threatening := make([]int, 0, 16)

  for i := 0; i < 64; i++ {
    if c.validPiece(i) && c.prelimValidMove(i, sq) &&
      c.pieceColorOnPosition(i) == 1 - color {
      threatening = append(threatening, i)
    }
  }

  return threatening
}

// Confirms if a square is under threat
func (c Chessboard) squareThreatened(sq int, color int) bool {
  return len(c.piecesThreateningPos(sq, color)) > 0
}

// Finds the appropriately colored king
func (c Chessboard) positionForKing(color int) int {
  for i, v := range(c.boardSquares) {
    if int(v) == (6 + 10 * color) {
      return i
    }
  }

  return -1
}

// Checks if a king is in check
func (c Chessboard) kingInCheck(color int) {
  pos := c.positionForKing(color)
  c.squareThreatened(pos, color)
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

func posFromRowColumn(r int, c int) int {
  return r * 8 + c
}

// Return the row/col direction (1,-1,0) to travel "from" to "to."
func travelDirection(from int, to int) (rowDir int, colDir int) {
  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
  toCol := colFromPosition(to)

  if fromRow < toRow {
    rowDir = 1
  } else if fromRow > toRow {
    rowDir = -1
  } else {
    rowDir = 0
  }

  if fromCol < toCol {
    colDir = 1
  } else if fromCol > toCol {
    colDir = -1
  } else {
    colDir = 0
  }

  return
}

func squaresBetween(from int, to int) []int {
  sqBetween := make([]int, 0, 8)

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)

  rowDir, colDir := travelDirection(from, to)

  currPos := from
  currRow := fromRow + rowDir
  currCol := fromCol + colDir
  currPos = posFromRowColumn(currRow, currCol)

  for currPos > -1 && currPos < 64 && currPos != to {
    currRow = rowFromPosition(currPos)
    currCol = rowFromPosition(currPos)

    sqBetween = append(sqBetween, currPos)

    currRow = currRow + rowDir
    currCol := currCol + colDir

    currPos = posFromRowColumn(currRow, currCol)
  }

  return sqBetween
}

func (c Chessboard) validateEmptySquaresBetween(from int, to int) bool {
  for _, v := range(squaresBetween(from, to)) {
    if c.validPiece(v) {
      return false
    }
  }

  return true
}

func abs(n int) int {
  if n > 0 {
    return n
  }

  return -n
}

// Debug Methods
func (c Chessboard) prettyPrintedPieceOnSquare(sq int) string {
  prettyText := ""

  switch c.boardSquares[sq] % 10 {
  case 6:
    prettyText = "K"
  case 5:
    prettyText = "Q"
  case 4:
    prettyText = "R"
  case 3:
    prettyText = "B"
  case 2:
    prettyText = "N"
  case 1:
    prettyText = "P"
  default:
    prettyText = "-"
  }

  if c.boardSquares[sq] / 10 == 1 {
    prettyText = strings.ToLower(prettyText)
  }

  return prettyText
}

func (c Chessboard) PrintBoard() {
  if c.turn {
    fmt.Println("Black to play.")
  } else {
    fmt.Println("White to play.")
  }

  for i, _ := range(c.boardSquares) {
    fmt.Printf("%s ", c.prettyPrintedPieceOnSquare(i))
    if i % 8 == 7 {
      fmt.Println()
    }
  }
}
