package chessboard

import (
  "fmt"
)

type Bitboards struct {
  whiteKing uint64
  blackKing uint64

  whiteQueen uint64
  blackQueen uint64

  whiteKnight uint64
  blackKnight uint64

  whiteBishop uint64
  blackBishop uint64

  whiteRook uint64
  blackRook uint64

  whitePawn uint64
  blackPawn uint64

  tables BitboardTables
}

func NewBitboards() (b Bitboards) {
  b = Bitboards{}
  b.tables = NewBitboardTables()

  return
}

func (b Bitboards) addToBitboards(piece string, square int) {

}

func (b Bitboards) PieceExistsOnSquare(square int) bool {
  return pieceOnBitboard(b.AllPieces(), square)
}

func (b Bitboards) AllPieces() uint64 {
    return b.WhitePieces() | b.BlackPieces()
}

func (b Bitboards) WhitePieces() uint64 {
  return b.whitePawn | b.whiteRook | b.whiteBishop | b.whiteKnight | b.whiteQueen | b.whiteKing
}

func (b Bitboards) BlackPieces() uint64 {
  return b.blackPawn | b.blackRook | b.blackBishop | b.blackKnight | b.blackQueen | b.blackKing
}

func pieceOnBitboard(bitboard uint64, square int) bool {
  return int64ToBool((bitboard >> uint(square)) & 0x1)
}

func reverse(s string) (ret string) {
    for _, v := range s {
        defer func(r rune) { ret += string(r) }(v)
    }
    return
}

func int64ToBool(i uint64) bool {
  if i == 1 {
    return true
  }

  return false
}

func PrettyPrintBitboard(bitboard uint64) {
  c := bitboard
  rank := ""

  for i := 0; i < 64; i++ {
    msb := c & (0x1 << 63)
    val := 0

    if msb != 0 {
      val = 1
    }

    rank += fmt.Sprintf("%d ", val)

    if ((i + 1) % 8 == 0) {
      fmt.Println(reverse(rank))
      rank = ""
    }

    c = c << 1
  }
}
