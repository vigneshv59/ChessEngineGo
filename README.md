# BrainyEngine

## UCI Interface
BrainyEngine currently responds to the following UCI commands. More commands will be added as the engine is completed. The full list of commands and longer descriptions of them can be seen in [this](http://wbec-ridderkerk.nl/html/UCIProtocol.html) document published by Stefan-Meyer Khalen.

### Standard UCI Commands (Currently Implemented)
- `uci`
	- **Usage**: Tells the engine to use the Universal Chess Interface protocol.
	- **Expected Response**: The engine should respond with the `uciok` message back to the caller.

- `debug [on | off]`
	- **Usage**: Toggles debug mode.
	- **Expected Response**: No response.

- `isready`
	- **Usage**: Used to query for engine intialization.
	- **Expected Response**: `readyok` after the engine is ready.

- `position [fen | startpos] moves ....`
	- **Usage**: Used to initialize the engine at a position.
	- **Expected Response**: No response.

- `go`
	- **Usage**: Start the calculation of the current position with some of the following subcommands.
		- `infinite`: Seach until the `stop` command is sent.
		- `go depth [depth]`: Seach moves for a specific depth.
		- Other subcommands can be seen on the official UCI documentation, and will be added if time permits.
	- **Expected Response**:
		- Various `info` lines.
		- `bestmove [move]` when the command terminates.
		- Sends periodic update messages with explored lines at a depth: `info depth [depth] nodes [nodes] nps [nodes/sec] score cp [score] time [time(ms)] pv [moves]`
- `stop`
	- **Usage**: Stop the calculations.
	- **Expected Response**: `bestmove [move]` if a calculation is running, nothing otherwise.

### Custom Debugging Commands (Currently Implemented)
These commands are only valid when debug mode is enabled, otherwise the engine will not respond to these commands.

- `dump`
	- **Usage**: Dump the board to the command line (with relevant information).
	- **Expected Response**: The engine should respond with a visual representation of the chessboard.
- `legalmoves [square]`
 	- **Usage**: Dump the board to the command line with legal moves of the piece on the given square indicated by `x` or `c` depending on whether the move will be a capture.
	- **Expected Response**: The engine will respond with a visual representation of the chessboard with legal moves indicated.

## Project Organization
The project consists of the following files, and the files planned in the future:

- `chessboard.go`: Consists of the logic for the chess game. Legal moves, board representation, etc. is in this file.

- `uci.go`: Consists of the UCI interface implementation. Calls into functions in `chessboard.go` to handle legality checking and board representation handling.

- `brain.go`: Contains the alpha-beta pruning/minimax algorithm implementation for the engine. Calls into `chessboard.go` to handle the board representations and legality.

- `evaluate.go`: Contains the shallow evaluation functions for the board that eventually feed into the alpha-beta search algorithm. Uses primarily point values for pieces, along with bonuses for centralization (piece movement potential) and pawn structure. **TODO:** It would be nice to have some sort of smart way to handle king safety.

- `book.go`: Contains the code to read from a polyglot opening book and lookup moves from the given book.

- `bitboards.go`: Contains the logic for basic bitboard operations and handling bitboards for the chessboard itself. Does not handle any bitboard table lookup logic -- this is handled by `tables.go`

- `tables.go`: Handles the precomputation of bitboard lookup tables, such as the lookup for the attack tables for all the pieces, to allow for fast access for the engine.


## Project Milestone Goals
- November 9
	- Complete the chessboard library with all the rules of chess, including en passant, castling, and promotions.
	- Add a legal move list to the chessboard library, developing a way to conveniently determine all the valid moves from a given position on the board.
- November 16
	- Develop a basic way of evaluating a given position, most likely purely based on naive piece/point designations.
	- Implement a basic Alpha-Beta algorithm to search the move tree for the best move based on a given evaluation function.
- November 22
	- Improve positional evaluation function manually (give points to passed pawns, centralized pieces, pawns about to promote, etc.
	- Begin integration of an opening book for
- November 30
	- Complete opening book integration
	- Experiment with multithreading support
- December 7
	- Complete modified multithreaded alpha-beta if possible
	- Write unit tests for the chessboard libraries, and tests for the engine search algorithms.

## Resources
- [UCI Protocol](http://wbec-ridderkerk.nl/html/UCIProtocol.html): A complete description of the chess UCI protocol.
- [Alpha-Beta Pruning](https://chessprogramming.wikispaces.com/Alpha-Beta): A description of the alpha-beta algorithm, as it relates to chess.
- [Multithreaded Alpha-Beta on GPUs](https://hrcak.srce.hr/file/114783): A paper which demonstrates a method of developing an Alpha-Beta algorithm that works with GPU multithreading.
