package main

import (
  "errors"
  "regexp"
  "strconv"
  "strings"
  "fmt"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Map fen piece strings to piece values. This makes for easy color/piece
// checking by dividing or mod 10 operations.
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

// Struct to represent chessboard state.
type Chessboard struct {
  boardSquares []int8
  enpassantPos int // The position for an enpassant capture, -1 if it doesnt exist.
  ksCanCastle []bool // Can players castle king-side? (0 white, 1 black)
  qsCanCastle []bool // Can players castle queen-side? (0 white, 1 black)
  turn bool // false for white's move, true for black's move
}

// Creates a new chessboard from a given fen position. Either returns the
// chessboard in board, or the creation error in err.
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
  if len(fenParts) > 1 {
    if fenParts[1] == "w" {
      board.turn = false
    }

    if fenParts[1] == "b" {
      board.turn = true
    }
  }

  board.ksCanCastle = make([]bool, 2, 2)
  board.qsCanCastle = make([]bool, 2, 2)

  board.ksCanCastle[0] = board.validCastlePositionKingside(0)
  board.ksCanCastle[1] = board.validCastlePositionKingside(1)
  board.qsCanCastle[0] = board.validCastlePositionQueenside(0)
  board.qsCanCastle[1] = board.validCastlePositionQueenside(1)

  if len(fenParts) > 2 {
    ksw, _ := regexp.MatchString("^.*K.*$", fenParts[2])
    ksb, _ := regexp.MatchString("^.*k.*$", fenParts[2])
    qsw, _ := regexp.MatchString("^.*Q.*$", fenParts[2])
    qsb, _ := regexp.MatchString("^.*q.*$", fenParts[2])

    if !ksw {
      board.ksCanCastle[0] = false
    }

    if !ksb {
      board.ksCanCastle[1] = false
    }

    if !qsw {
      board.qsCanCastle[0] = false
    }

    if !qsb {
      board.qsCanCastle[1] = false
    }
  }

  board.enpassantPos = -1
  if len(fenParts) > 3 {
    board.enpassantPos = alToPos(fenParts[3])
  }

  return
}

// Checks for promotion validity, does not take into account turn
func (c Chessboard) attemptedPromotion(from int, to int) bool {
  color := c.pieceColorOnPosition(from)

  if !c.validPiecePawn(from) {
    return false
  }

  if color == 0 && to > -1 && to < 8 {
    return true
  }

  if color == 1 && to > 55 && to < 64 {
    return true
  }

  return false
}

// Checks if a move is legal without altering the board.
func (c Chessboard) moveIsLegal(from int, to int, promopiece string) bool {
  return c.Move(from, to, promopiece, true)
}

// Makes a move using algebraic descriptive notation.
// Example: e2e4
func (c *Chessboard) MoveAlDescriptive(notation string) bool {
  fromSquare := alToPos(notation[0:2])
  toSquare := alToPos(notation[2:])

  return c.MakeMove(fromSquare, toSquare, "")
}

// Moves from->to if the move is legal.
func (c *Chessboard) MakeMove(from int, to int, promopiece string) bool {
  return c.Move(from, to, promopiece, false)
}

// Makes a move on the board, returning the legaility of the move
// as a boolean.
func (c *Chessboard) Move(from int, to int, promopiece string, dryrun bool) bool {
  color := c.pieceColorOnPosition(from)
  turn := 0
  restoreMap := make(map[int]int8)

  if c.turn {
    turn = 1
  }

  // Incorrect turn
  if color != turn {
    return false
  }

  // TODO: Validate promopiece correctness.

  if !c.prelimValidMove(from, to) {
    return false
  }

  preKsCanCastle := make([]bool, 2)
  preQsCanCastle := make([]bool, 2)

  copy(preKsCanCastle, c.ksCanCastle)
  copy(preQsCanCastle, c.qsCanCastle)

  // The king cannot castle after it has moved.
  if c.validPieceKing(from) {
    c.ksCanCastle[color] = false
    c.qsCanCastle[color] = false
  }

  // The king cannot castle on the side which the rook has moved.
  if c.validPieceRook(from) {
    if from == c.rookCastleKingsidePosition(color) {
      c.ksCanCastle[color] = false
    }

    if from == c.rookCastleQueensidePosition(color) {
      c.qsCanCastle[color] = false
    }
  }

  // Castling
  if c.castlingAttempt(color, from, to) {
    afterPos := c.rookPositionAfterCastle(color, from, to)
    beforePos := c.rookPositionBeforeCastle(color, from, to)

    restoreMap[afterPos] = c.boardSquares[afterPos]
    restoreMap[beforePos] = c.boardSquares[beforePos]

    c.boardSquares[afterPos] = c.boardSquares[beforePos]
    c.boardSquares[beforePos] = -1
  }

  // En passant capture
  if c.validPiecePawn(from) && c.validColorPiece(from, color) &&
    c.enpassantPos == to {
      if color == 1 {
        restoreMap[to-8] = c.boardSquares[to-8]
        c.boardSquares[to-8] = -1
      } else {
        restoreMap[to+8] = c.boardSquares[to+8]
        c.boardSquares[to+8] = -1
      }
  }

  // Set the enpassant position
  preEpPos := c.enpassantPos
  c.enpassantPos = c.generateEpPos(color, from, to)

  // Update the underlying move.
  restoreMap[to] = c.boardSquares[to]
  restoreMap[from] = c.boardSquares[from]

  // If this is an attempted promotion, promote the pawn
  // If it is a dryrun, we don't really care what piece the pawn
  // is promoted to, therefore, we will ignore this validation.
  if c.attemptedPromotion(from, to) && !dryrun {
    if promopiece == "" {
      c.boardSquares[to] = int8(color * 10 + 5)
    } else {
      c.boardSquares[to] = int8(color) * 10 + (pieceVals[promopiece] % 10)
    }

    c.boardSquares[from] = -1
  } else {
    c.boardSquares[to] = c.boardSquares[from]
    c.boardSquares[from] = -1
  }

  kingInCheck := c.kingInCheck(color)

  // If the king is in check, or this is a dry run, reset the board.
  if kingInCheck || dryrun {
    for k, v := range restoreMap {
      c.boardSquares[k] = v
    }

    c.ksCanCastle[0] = preKsCanCastle[0]
    c.qsCanCastle[0] = preQsCanCastle[0]

    c.ksCanCastle[1] = preKsCanCastle[1]
    c.qsCanCastle[1] = preQsCanCastle[1]

    c.enpassantPos = preEpPos

    return !kingInCheck
  }

  c.turn = !c.turn

  return true
}

// Returns the legal moves of a piece on a particular
// square, represented as an array of legal destination
// squares.
func (c Chessboard) LegalMovesFromSquare(from int) []int {
  candMoves := c.candSquares(from)
  legalMoves := make([]int, 0, 32)

  for _, v := range(candMoves) {
    if c.moveIsLegal(from, v, "") {
      legalMoves = append(legalMoves, v)
    }
  }

  return legalMoves
}

// Validates a move on a board with a from index and a to index
// Returns a boolean indicating success of the move.
// This is a preliminary validator, and does not take into account checkmate
// or checks, which will be dealt with after this first pass confirmation
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
    if c.validMoveKing(from, to) {
      return true
    }

    if c.validCastle(from, to) {
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
    fmt.Println("Default")
    return false
  }

  return false
}

// Generates candidate squares for a piece on a given
// from square.
func (c Chessboard) candSquares(from int) []int {
  switch c.boardSquares[from] % 10 {
  case 6: // King
    return c.candSquaresKing(from)
  case 5: // Queen
    return c.candSquaresQueen(from)
  case 4: // Rook
    return c.candSquaresRook(from)
  case 3: // Bishop
    return c.candSquaresBishop(from)
  case 2: // Knight
    return c.candSquaresKnight(from)
  case 1: // Pawn
    return c.candSquaresPawn(from)
  default:
    return make([]int, 0, 1)
  }

  return make([]int, 0, 1)
}

// Piece Existence Validators:

// Checks if a king exists on a square.
func (c Chessboard) validPieceKing(square int) bool {
  return int(c.boardSquares[square] % 10) == 6
}

// Checks if a queen exists on a square.
func (c Chessboard) validPieceQueen(square int) bool {
  return int(c.boardSquares[square] % 10) == 5
}

// Checks if a rook exists on a square.
func (c Chessboard) validPieceRook(square int) bool {
  return int(c.boardSquares[square] % 10) == 4
}

// Checks if a bishop exists on a square.
func (c Chessboard) validPieceBishop(square int) bool {
  return int(c.boardSquares[square] % 10) == 3
}

// Checks if a knight exists on a square.
func (c Chessboard) validPieceKnight(square int) bool {
  return int(c.boardSquares[square] % 10) == 2
}

// Checks if a pawn exists on a square.
func (c Chessboard) validPiecePawn(square int) bool {
  return int(c.boardSquares[square] % 10) == 1
}

// Checks if a piece exists on a square.
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

// Generates all possible candidate moves for a knight at position from.
func (c Chessboard) candSquaresKnight(from int) []int {
  if !c.validPieceKnight(from) {
    return make([]int, 0, 1)
  }

  candSquares := make([]int, 0, 16)
  color := c.pieceColorOnPosition(from)

  directions := [][]int{{2,1},{-2,1},{2,-1},{-2,-1}, {1,2},{-1,2},{1,-2},{-1,-2}}

  for _,v := range(directions) {
    currRow := rowFromPosition(from) + v[0]
    currCol := colFromPosition(from) + v[1]

    currPos := posFromRowColumn(currRow, currCol)

    if currRow > -1 && currRow < 8 && currCol > -1 && currCol < 8 {
      if c.pieceColorOnPosition(currPos) != color {
        candSquares = append(candSquares, currPos)
      }
    }
  }

  return candSquares
}

// Validates a queen move, does not take into account turn.
func (c Chessboard) validMoveQueen(from int, to int) bool {
  if !c.validPieceQueen(from) {
      return false
  }

  return c.validMoveRook(from, to) || c.validMoveBishop(from, to)
}

// Generates all possible candidate moves for a queen at position from.
func (c Chessboard) candSquaresQueen(from int) []int {
  return append(c.candSquaresRook(from), c.candSquaresBishop(from)...)
}

// Validates a rook move, does not take into account turn.
func (c Chessboard) validMoveRook(from int, to int) bool {
  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
  toCol := colFromPosition(to)

  return fromRow == toRow || fromCol == toCol

}

// Generates all possible candidate moves for a rook at position from.
// Should only be used for queens (through candSquaresQueen), or rooks.
func (c Chessboard) candSquaresRook(from int) []int {
  candSquares := make([]int, 0, 16)

  color := c.pieceColorOnPosition(from)

  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)

  // Squares below the rook
  for i := 1; i < 8 - fromRow; i++ {
    to := posFromRowColumn(fromRow + i, fromCol)

    if c.pieceColorOnPosition(to) == color {
      break
    }

    candSquares = append(candSquares, to)

    if c.pieceColorOnPosition(to) != -1 {
      break
    }
  }

  // Squares above the rook
  for i := 1; i < fromRow + 1; i++ {
    to := posFromRowColumn(fromRow - i, fromCol)

    if c.pieceColorOnPosition(to) == color {
      break
    }

    candSquares = append(candSquares, to)

    if c.pieceColorOnPosition(to) != -1 {
      break
    }
  }

  // Squares to the right of the rook
  for i := 1; i < 8 - fromCol; i++ {
    to := posFromRowColumn(fromRow, fromCol + i)

    if c.pieceColorOnPosition(to) == color {
      break
    }

    candSquares = append(candSquares, to)

    if c.pieceColorOnPosition(to) != -1 {
      break
    }
  }

  // Squares to the left of the rook
  for i := 1; i < fromCol + 1; i++ {
    to := posFromRowColumn(fromRow, fromCol - i)

    if c.pieceColorOnPosition(to) == color {
      break
    }

    candSquares = append(candSquares, to)

    if c.pieceColorOnPosition(to) != -1 {
      break
    }
  }

  return candSquares
}

// Validates a bishop move, does not take into account turn.
func (c Chessboard) validMoveBishop(from int, to int) bool {
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

// Generates the candidate squares for a bishop on the given
// square. Returns the squares as an array of possible destinations.
func (c Chessboard) candSquaresBishop(from int) []int {
  candSquares := make([]int, 0, 16)

  color := c.pieceColorOnPosition(from)

  for _,v := range([][]int{{1,1},{1,-1},{-1,1},{-1,-1}}) {
    currRow := rowFromPosition(from) + v[0]
    currCol := colFromPosition(from) + v[1]

    currPos := posFromRowColumn(currRow, currCol)

    for currRow > -1 && currRow < 8 && currCol > -1 && currCol < 8 {
      if c.pieceColorOnPosition(currPos) == color {
        break
      }

      candSquares = append(candSquares, currPos)

      if c.pieceColorOnPosition(currPos) != -1 {
        break
      }

      currRow = rowFromPosition(currPos) + v[0]
      currCol = colFromPosition(currPos) + v[1]

      currPos = posFromRowColumn(currRow, currCol)
    }
  }

  return candSquares
}

// Checks if the move from -> to is a valid kingside castling attempt
// for the given color.
func (c Chessboard) kingsideCastlingAttempt(color int, from int, to int) bool {
  if color == 0 && from == 60 && to == 62 {
    return true
  }

  if color == 1 && from == 4 && to == 6 {
    return true
  }

  return false
}

// Checks if the move from -> to is a valid queenside castling attempt
// for the given color.
func (c Chessboard) queensideCastlingAttempt(color int, from int, to int) bool {
  if color == 0 && from == 60 && to == 58 {
    return true
  }

  if color == 1 && from == 4 && to == 2 {
    return true
  }

  return false
}

// Checks if the move from -> to is a castling attempt for the given color.
func (c Chessboard) castlingAttempt(color int, from int, to int) bool {
  return (c.kingsideCastlingAttempt(color, from, to) ||
    c.queensideCastlingAttempt(color, from, to))
}

// Checks if the move from -> to is a legal castling move.
func (c Chessboard) validCastle(from int, to int) bool {
  if !c.validPieceKing(from) {
    return false
  }

  color := c.pieceColorOnPosition(from)

  if !c.castlingAttempt(color, from, to) {
    return false
  }

  if c.kingInCheck(color) {
    return false
  }

  rookPos := c.rookCastleKingsidePosition(color)

  if c.kingsideCastlingAttempt(color, from, to) {
    if !c.ksCanCastle[color] {
      return false
    }
  }

  if c.queensideCastlingAttempt(color, from, to) {
    if !c.qsCanCastle[color] {
      return false
    }

    rookPos = c.rookCastleQueensidePosition(color)
  }

  if !c.validateEmptySquaresBetween(from, rookPos) {
    return false
  }

  for _, v := range(squaresBetween(from, rookPos)) {
    if c.squareThreatened(v, color) {
      return false
    }
  }

  return true
}

// Validates a king move, does not take into account turn or castling.
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

// Generates the candidate squares for a king on the given
// square. Returns the squares as an array of possible destinations.
func (c Chessboard) candSquaresKing(from int) []int {
  if !c.validPieceKing(from) {
    return make([]int, 0, 1)
  }

  candSquares := make([]int, 0, 16)
  color := c.pieceColorOnPosition(from)
  directions := [][]int{{1,1},{1,-1},{-1,1},{-1,-1},{0,1},{0,-1},{-1,0},{1,0}}
  for _,v := range(directions) {
    currRow := rowFromPosition(from) + v[0]
    currCol := colFromPosition(from) + v[1]

    if currRow < 0 || currRow > 7 || currCol < 0 || currCol > 7 {
      continue
    }

    currPos := posFromRowColumn(currRow, currCol)

    if c.pieceColorOnPosition(currPos) != color {
      candSquares = append(candSquares, currPos)
    }
  }

  if c.validCastle(from, from + 2) {
    candSquares = append(candSquares, from + 2)
  }

  if c.validCastle(from, from - 2) {
    candSquares = append(candSquares, from - 2)
  }

  return candSquares
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
      if (toRow == fromRow - 1) &&
          ((toCol == fromCol + 1) || (toCol == fromCol - 1)) {

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
      if (toRow == fromRow + 1) &&
          ((toCol == fromCol + 1) || (toCol == fromCol - 1)) {

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

// Generates the candidate squares for a pawn on the given
// square. Returns the squares as an array of possible destinations.
func (c Chessboard) candSquaresPawn(from int) []int {
  if !c.validPiecePawn(from) {
    return make([]int, 0, 1)
  }

  candSquares := make([]int, 0, 16)
  color := c.pieceColorOnPosition(from)

  forwardDir := color*2 - 1

  directions := [][]int{{forwardDir,0}}

  row := rowFromPosition(from)

  if color == 0 && row == 6 {
    directions = append(directions, []int{-2, 0})
  } else if color == 1 && row == 1 {
    directions = append(directions, []int{2, 0})
  }

  for _,v := range(directions) {
    currRow := rowFromPosition(from) + v[0]
    currCol := colFromPosition(from) + v[1]

    currPos := posFromRowColumn(currRow, currCol)

    if c.pieceColorOnPosition(currPos) != color {
      candSquares = append(candSquares, currPos)
    } else {
      break
    }
  }

  directions = [][]int{{forwardDir,-1},{forwardDir,1}}
  for _,v := range(directions) {
    currRow := rowFromPosition(from) + v[0]
    currCol := colFromPosition(from) + v[1]

    currPos := posFromRowColumn(currRow, currCol)

    if c.validMovePawn(from, currPos) {
      candSquares = append(candSquares, currPos)
    }
  }

  return candSquares
}

// Lists the pieces threatening a position
func (c Chessboard) piecesThreateningPos(sq int, color int) []int {
  threatening := make([]int, 0, 16)

  for i := 0; i < 64; i++ {
    if c.pieceColorOnPosition(i) == 1 - color && c.prelimValidMove(i, sq) &&
       c.validPiece(i) {
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
func (c Chessboard) kingInCheck(color int) bool {
  pos := c.positionForKing(color)
  return c.squareThreatened(pos, color)
}

// Returns 0 for white pieces and 1 for
// black pieces, -1 for empty.
func (c Chessboard) pieceColorOnPosition(pos int) int {
  if c.boardSquares[pos] < 0 {
      return -1
  }

  return int(c.boardSquares[pos] / 10)
}

// Validates that the color of the piece at position pos is color.
func (c Chessboard) validColorPiece(pos int, color int) bool {
  return c.pieceColorOnPosition(pos) == color
}

// Helper methods

// Returns the row of a given pos.
func rowFromPosition(pos int) int {
  return pos / 8
}

// Returns the column of a given pos.
func colFromPosition(pos int) int {
  return pos % 8
}

// Returns the position from a given row and column.
func posFromRowColumn(r int, c int) int {
  return r * 8 + c
}

// Takes an algebraic square and converts it into a 0-63 square.
func alToPos(al string) int {
  r, _ := strconv.Atoi(al[1:])
  r = 8 - r

  c := ([]byte(al[0:1]))[0] - []byte("a"[0:1])[0]

  return posFromRowColumn(r, int(c))
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

// Generates the en-passant square given the from/to square of
// a move, and the color of the move.
func (c Chessboard) generateEpPos(color int, from int, to int) int {
  fromRow := rowFromPosition(from)
  fromCol := colFromPosition(from)
  toRow := rowFromPosition(to)
  toCol := colFromPosition(to)

  if !(c.validPiecePawn(from) && c.validColorPiece(from, color)) {
    return -1
  }

  if (c.validColorPiece(from, 0) && toRow == 4 && fromRow == 6 &&
    toCol == fromCol && c.validateEmptySquaresBetween(from, to)) {
    return from - 8
  }

  if (c.validColorPiece(from, 1) && toRow == 3 && fromRow == 1 &&
    toCol == fromCol && c.validateEmptySquaresBetween(from, to)) {
    return from + 8
  }

  return -1
}

// Returns an array of the squares directly in between two squares.
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
    currCol = colFromPosition(currPos)

    sqBetween = append(sqBetween, currPos)

    currRow = currRow + rowDir
    currCol = currCol + colDir

    currPos = posFromRowColumn(currRow, currCol)
  }

  return sqBetween
}

// Validates that there aren't pieces between from and to.
func (c Chessboard) validateEmptySquaresBetween(from int, to int) bool {
  for _, v := range(squaresBetween(from, to)) {
    if c.validPiece(v) {
      return false
    }
  }

  return true
}

// Checks if the kingside position is castleable for the given color.
func (c Chessboard) validCastlePositionKingside(color int) bool {
  kingPos := c.kingCastlePosition(color)
  rookPos := c.rookCastleKingsidePosition(color)

  if c.pieceColorOnPosition(kingPos) != color {
    return false
  }

  if c.pieceColorOnPosition(rookPos) != color {
    return false
  }

  return (c.validPieceKing(kingPos) && c.validPieceRook(rookPos))
}

// Checks if the queenside position is castleable for the given color.
func (c Chessboard) validCastlePositionQueenside(color int) bool {
  kingPos := c.kingCastlePosition(color)
  rookPos := c.rookCastleQueensidePosition(color)

  if c.pieceColorOnPosition(kingPos) != color {
    return false
  }

  if c.pieceColorOnPosition(rookPos) != color {
    return false
  }

  return (c.validPieceKing(kingPos) && c.validPieceRook(rookPos))
}

// Returns the king's position before castling for a given color.
func (c Chessboard) kingCastlePosition(color int) int {
  if color == 0 {
    return 60
  }

  return 4
}

// Returns the rook's position before castling for a given color and move
// (identifies whether the move is a ks or qs castle).
func (c Chessboard) rookPositionBeforeCastle(color int, from int, to int) int {
  if color == 0 && from == 60 && to == 62 {
    return 63
  }

  if color == 0 && from == 60 && to == 58 {
    return 56
  }

  if color == 1 && from == 4 && to == 6 {
    return 7
  }

  if color == 1 && from == 4 && to == 2 {
    return 0
  }

  return -1
}

// Returns the rook's position after castling for a given color and move
// (identifies whether the move is a ks or qs castle).
func (c Chessboard) rookPositionAfterCastle(color int, from int, to int) int {
  if color == 0 && from == 60 && to == 62 {
    return 61
  }

  if color == 0 && from == 60 && to == 58 {
    return 59
  }

  if color == 1 && from == 4 && to == 6 {
    return 5
  }

  if color == 1 && from == 4 && to == 2 {
    return 3
  }

  return -1
}

// Returns the rook's initial position in a kingside castle.
func (c Chessboard) rookCastleKingsidePosition(color int) int {
  if color == 0 {
    return 63
  }

  return 7
}

// Returns the rook's initial position in a queenside castle.
func (c Chessboard) rookCastleQueensidePosition(color int) int {
  if color == 0 {
    return 56
  }

  return 0
}


// Returns the absolute value of an int.
func abs(n int) int {
  if n > 0 {
    return n
  }

  return -n
}

// Check if a given element exists in a slice.
func intInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }

    return false
}

// Debug Methods

// Returns the pretty printed symbol of a piece on a square, used in
// debugging board prints.
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


// Prints the board along with any necessary information to fully
// describe a position.
func (c Chessboard) PrintBoard() {
  if c.turn {
    fmt.Println("Black to play.")
  } else {
    fmt.Println("White to play.")
  }

  fmt.Print("Enpassant position: ")
  if c.enpassantPos == -1 {
    fmt.Println("-")
  } else {
    fmt.Println(c.enpassantPos)
  }

  for i, _ := range(c.boardSquares) {
    fmt.Printf("%s ", c.prettyPrintedPieceOnSquare(i))
    if i % 8 == 7 {
      fmt.Println()
    }
  }
}


// Prints a board with the legal moves of a piece on source indicated by
// algebraic notation (alsource).
func (c Chessboard) PrintLegalMoves(alsource string) {
  source := alToPos(alsource)

  if c.turn {
    fmt.Println("Black to play.")
  } else {
    fmt.Println("White to play.")
  }

  fromSquares := c.LegalMovesFromSquare(source)

  for i, _ := range(c.boardSquares) {
    if intInSlice(i, fromSquares) {
      if (c.boardSquares[i] % 10) < 0 {
        fmt.Printf("x ")
      } else {
        fmt.Printf("c ")
      }
    } else {
      fmt.Printf("%s ", c.prettyPrintedPieceOnSquare(i))
    }

    if i % 8 == 7 {
      fmt.Println()
    }
  }
}
