// This file contains the terminal interface to the program.
// It responds to the UCI commands, and will call into other files
// to setup the chessboard and run move searches.

package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
)

func handleUci() {
  fmt.Println("uciok")
}

func handleIsReady() {
  fmt.Println("readyok")
}

func handleNewGame() {
  // Clear chessboard and initialize new game
}

func handlePosition(position string) {
  board, _ := NewChessboard(position)
  fmt.Println(board.boardSquares)
}

func handleInput(input string) {
  switch strings.Split(input, " ")[0] {
  case "uci":
    handleUci()
  case "isready":
    handleIsReady()
  case "position":
    handlePosition("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
  }
}

func main() {
  fmt.Println("BrainyEngine by Vignesh Varadarajan v0.0")

  for {

    buf := bufio.NewReader(os.Stdin)
    sentence, err := buf.ReadBytes('\n')

    if err != nil {
      fmt.Println(err)
    } else {
      handleInput(strings.TrimSpace(string(sentence)))
    }

  }
}
