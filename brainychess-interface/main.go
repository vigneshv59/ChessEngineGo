// This file contains the terminal interface to the program.
// It responds to the UCI commands, and will call into other files
// to setup the chessboard and run move searches.

package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "github.com/vigneshv59/chessboard/chessboard"
)

type gameState struct {
  moves []string
  startFen string
}

func handlePosition(position string) chessboard.Chessboard {
  board, _ := chessboard.NewChessboard(position)

  return board
}

func handleInterfaceInput(input string,
                          state *gameState,
                          b chessboard.Chessboard) chessboard.Chessboard {

  cmdArr := strings.Split(input, " ")

  switch cmdArr[0] {
  case "position":
    fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

    if len(cmdArr) == 1 {
      fmt.Println("Incorrect arguments.")

      break
    }

    if cmdArr[1] != "startpos" {
      fen = strings.Join(cmdArr[1:], " ")
    }

    b = handlePosition(fen)
  default:
    if len(cmdArr) == 0 || cmdArr[0] == "" {
      return b
    }

    success := b.MoveAlDescriptive(cmdArr[0])

    if !success {
      fmt.Println("Illegal Move.")
      break
    }

    depth := 4
    b.PrintBoard()

    a := false
    _, move := b.AlphaBeta(depth, &a)

    if len(move) < 2 {
        fmt.Println("Game Over.")
        return b
    }

    moveStr := chessboard.PosToAl(move[0]) + chessboard.PosToAl(move[1])
    fmt.Println("Engine Moved: " + moveStr)
    b.MoveAlDescriptive(moveStr)

    b.PrintBoard()
  }

  return b
}

func main() {
  fmt.Println("BrainyEngine Interface by Vignesh Varadarajan v0.0")
  state := gameState{make([]string, 0, 100), "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"}
  var board chessboard.Chessboard

  for {

    buf := bufio.NewReader(os.Stdin)
    sentence, err := buf.ReadBytes('\n')

    if err != nil {
      fmt.Println(err)
    } else {
      board = handleInterfaceInput(strings.TrimSpace(string(sentence)),
                &state,
                board)
    }

  }
}
