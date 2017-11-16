// This file contains the terminal interface to the program.
// It responds to the UCI commands, and will call into other files
// to setup the chessboard and run move searches.

package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "strconv"
)

type uciConfig struct {
  debug bool
}

func handleUci() {
  fmt.Println("uciok")
}

func (ec *uciConfig) setDebug(desiredState bool) {
  ec.debug = desiredState
}

func handleIsReady() {
  fmt.Println("readyok")
}

func handleNewGame() {
  // TODO: Clear chessboard and initialize new game
}

func handlePosition(position string) Chessboard {
  board, _ := NewChessboard(position)

  return board
}

func handleInput(input string, engineConfig *uciConfig, b Chessboard) Chessboard {
  cmdArr := strings.Split(input, " ")

  switch cmdArr[0] {
  case "uci":
    handleUci()
  case "isready":
    handleIsReady()
  case "debug":
    debugState := false

    if cmdArr[1] == "on" {
      debugState = true
    }

    engineConfig.setDebug(debugState)
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

    for i, v := range cmdArr {
      if v == "moves" {
        moves := cmdArr[(i+1):]

        for _, m := range moves {
          success := b.MoveAlDescriptive(m)

          if !success {
            break
          }
        }

        break
      }
    }
  case "go":
    depth := 4
    if len(cmdArr) == 3 {
      depth, _ = strconv.Atoi(cmdArr[2])
    }

    _, move := b.alphaBeta(depth)
    fmt.Println("bestmove " + posToAl(move[0]) + posToAl(move[1]))
  case "sevaluate":
    if !engineConfig.debug {
      fmt.Println("Unknown command.")

      break
    }

    fmt.Println(b.Evaluate())
  case "legalmoves":
    if !engineConfig.debug {
      fmt.Println("Unknown command.")

      break
    }

    if len(cmdArr) <= 1 {
      fmt.Println("Command takes exactly one argument.")

      break
    }

    b.PrintLegalMoves(cmdArr[1])

  case "dump":
    if !engineConfig.debug {
      fmt.Println("Unknown command.")

      break
    }

    b.PrintBoard()

    break
  default:
    fmt.Println("Unknown command.")
  }

  return b
}

func main() {
  fmt.Println("BrainyEngine by Vignesh Varadarajan v0.0")
  engineConfig := uciConfig{false}
  var board Chessboard

  for {

    buf := bufio.NewReader(os.Stdin)
    sentence, err := buf.ReadBytes('\n')

    if err != nil {
      fmt.Println(err)
    } else {
      board = handleInput(strings.TrimSpace(string(sentence)),
                &engineConfig,
                board)
    }

  }
}
