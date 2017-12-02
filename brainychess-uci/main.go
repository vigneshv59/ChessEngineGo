// This file contains the terminal interface to the program.
// It responds to the UCI commands, and will call into other files
// to setup the chessboard and run move searches.

package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "regexp"
  "strconv"
  "github.com/vigneshv59/chessboard/chessboard"
)

type uciConfig struct {
  debug bool
}

func handleUci() {
  fmt.Println("id name BrainyEngine 1.0")
  fmt.Println("id author Vignesh")
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

func handlePosition(position string) chessboard.Chessboard {
  board, _ := chessboard.NewChessboard(position)

  return board
}

func handleInput(input string,
                  engineConfig *uciConfig,
                  b *chessboard.Chessboard,
                  s *bool) {
  cmdArr := strings.Split(input, " ")

  switch cmdArr[0] {
  case "uci":
    handleUci()
  case "isready":
    handleIsReady()
  case "quit":
    os.Exit(0)
  case "seval":
    fmt.Println(b.Evaluate())
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

    if cmdArr[1] == "fen" {
      fen = strings.Join(cmdArr[2:], " ")
    } else if cmdArr[1] != "startpos" {
      return
    }

    *s = true
    *b = handlePosition(fen)

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
    *s = false
  case "go":
    depth := 4
    if len(cmdArr) == 3 {
      depth, _ = strconv.Atoi(cmdArr[2])
    }

    if len(cmdArr) == 2 {
      depth = 10000
    }

    go func ()  {
      *s = false
      _, move := b.AlphaBeta(depth, s)
      fmt.Println("bestmove " + chessboard.PosToAl(move[0]) + chessboard.PosToAl(move[1]))
      *s = false
    }()
  case "stop":
    *s = true
  case "dump":
    if !engineConfig.debug {
      fmt.Println("Unknown command.")

      break
    }

    fmt.Println(b.BookHash())
    b.PrintBoard()

    break
  default:
    fmt.Println("Unknown command.")
  }
}

func main() {
  fmt.Println("BrainyEngine by Vignesh Varadarajan v0.0")

  engineConfig := uciConfig{false}
  var board chessboard.Chessboard
  b := false
  stopped := &b

  for {

    buf := bufio.NewReader(os.Stdin)
    sentence, err := buf.ReadBytes('\n')

    if err != nil {
      fmt.Println(err)
    } else {
      m, _ := regexp.MatchString("isready", string(sentence))

      if m {
          fmt.Println("readyok")
      } else {
        handleInput(strings.TrimSpace(string(sentence)),
                &engineConfig,
                &board, stopped)
      }
    }

  }
}
