package chessboard

// #include <stdint.h>
// int ctz(uint64_t t) {
//   return __builtin_ctzll(t);
// }
import "C"

type BitboardTables struct {
  fileTables []uint64
  rankTables []uint64
  kingTables []uint64
  knightTables []uint64

  whitePawnAttackTables []uint64
  blackPawnAttackTables []uint64
}

func NewBitboardTables() (b BitboardTables) {
    b = BitboardTables{}

    b.generateFileTables()
    b.generateRankTables()
    b.generateKingAttackTables()
    b.generateKnightAttackTables()

    return
}

// Generates the file tables for a bitboard.
func (b *BitboardTables) generateFileTables() {
  files := make([]uint64, 8)

  for i := uint16(0); i < 8; i++ {
    var a uint64 = 0x0101010101010101
    files[i] = (a << i)
  }

  b.fileTables = files
}

// Generate the rank tables for a bitboard.
func (b *BitboardTables) generateRankTables() {
  ranks := make([]uint64, 8)

  for i := uint16(0); i < 8; i++ {
    var r1 uint64 = 0xFF
    ranks[i] = r1 << (8 * i)
  }

  b.rankTables = ranks
}

func (b *BitboardTables) generateKingAttackTables() {
  kingBitboard := uint64(0x1)
  kingTables := make([]uint64, 64)

  for kingBitboard != 0 {
    kingTables[C.ctz(C.uint64_t(kingBitboard))] = b.kingAttacks(kingBitboard)
    kingBitboard = kingBitboard << 1
  }

  b.kingTables = kingTables
}

func (b *BitboardTables) generateKnightAttackTables() {
  knightBitboard := uint64(0x1)
  knightTables := make([]uint64, 64)

  for knightBitboard != 0 {
    knightTables[C.ctz(C.uint64_t(knightBitboard))] = b.kingAttacks(knightBitboard)
    knightBitboard = knightBitboard << 1
  }

  b.knightTables = knightTables
}

func (b BitboardTables) kingAttacks(kingBoard uint64) uint64 {
  bitboard := kingBoard

  bitboard = bitboard | b.oneWest(bitboard) | b.oneEast(bitboard)
  bitboard = bitboard | b.oneNorth(bitboard) | b.oneSouth(bitboard)

  return bitboard ^ kingBoard
}

func (b BitboardTables) knightAttacks(knightBoard uint64) uint64 {
  bitboard := knightBoard

  east := b.oneEast(bitboard);
  west := b.oneWest(bitboard);
  attacks := (east | west) << 16;
  attacks |= (east | west) >> 16;
  east = b.oneEast(east);
  west = b.oneWest(west);
  attacks |= (east | west) <<  8;
  attacks |= (east | west) >>  8;

  return attacks;
}

// func (b BitboardTables) pawnAttacks(pawnBoard uint64) uint64 {
//
// }

/*
Directions defined by

noWe         nort         noEa
        +7    +8    +9
            \  |  /
west    -1 <-  0 -> +1    east
            /  |  \
        -9    -8    -7
soWe         sout         soEa
*/

func (t BitboardTables) oneEast(board uint64) uint64 {
  return (^(t.fileTables[7]) & board) << 1
}

func (t BitboardTables) oneWest(board uint64) uint64 {
  return (^(t.fileTables[0]) & board) >> 1
}

func (t BitboardTables) oneNorth(board uint64) uint64 {
  return board << 8
}

func (t BitboardTables) oneSouth(board uint64) uint64 {
  return board >> 8
}

func (t BitboardTables) oneSoWest(board uint64) uint64 {
  return (^(t.fileTables[0]) & board) >> 9
}

func (t BitboardTables) oneSoEast(board uint64) uint64 {
  return (^(t.fileTables[7]) & board) >> 7
}

func (t BitboardTables) oneNoWest(board uint64) uint64 {
  return (^(t.fileTables[0]) & board) << 7
}

func (t BitboardTables) oneNoEast(board uint64) uint64 {
  return (^(t.fileTables[7]) & board) << 9
}
