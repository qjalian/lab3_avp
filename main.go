package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

var (
	board         [3][3]string
	currentPlayer string
	gameActive    bool
	scoreX, scoreO int
)

type GameResponse struct {
	Board         [3][3]string `json:"board"`
	CurrentPlayer string       `json:"currentPlayer"`
	Message       string       `json:"message"`
	GameActive    bool         `json:"gameActive"`
	ScoreX        int          `json:"scoreX"`
	ScoreO        int          `json:"scoreO"`
}

func main() {
	initializeBoard()
	currentPlayer = "X"
	gameActive = true

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/makeMove", makeMoveHandler)
	http.HandleFunc("/reset", resetHandler)

	port := 5000
	fmt.Printf("Server is running at http://localhost:5000 ")
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	gameResponse := GameResponse{
		Board:         board,
		CurrentPlayer: currentPlayer,
		Message:       "",
		GameActive:    gameActive,
		ScoreX:        scoreX,
		ScoreO:        scoreO,
	}

	err = tmpl.Execute(w, gameResponse)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if !gameActive {
		http.Error(w, "Game over", http.StatusBadRequest)
		return
	}

	row, col := r.URL.Query().Get("row"), r.URL.Query().Get("col")
	if row == "" || col == "" {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	rowInt, colInt := parseInt(row), parseInt(col)
	if !isValidMove(rowInt, colInt) {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	board[rowInt][colInt] = currentPlayer
	if checkWinner() {
		gameActive = false
		updateScores(currentPlayer)
		sendJSONResponse(w, "Player "+currentPlayer+" wins!")
		return
	}

	if isBoardFull() {
		gameActive = false
		sendJSONResponse(w, "It's a draw!")
		return
	}

	switchPlayer()
	sendJSONResponse(w, "")
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	initializeBoard()
	currentPlayer = "X"
	gameActive = true

	sendJSONResponse(w, "")
}

func initializeBoard() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = " "
		}
	}
}

func isValidMove(row, col int) bool {
	return row >= 0 && row < 3 && col >= 0 && col < 3 && board[row][col] == " "
}

func checkWinner() bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == currentPlayer && board[i][1] == currentPlayer && board[i][2] == currentPlayer {
			return true
		}
		if board[0][i] == currentPlayer && board[1][i] == currentPlayer && board[2][i] == currentPlayer {
			return true
		}
	}

	if board[0][0] == currentPlayer && board[1][1] == currentPlayer && board[2][2] == currentPlayer {
		return true
	}
	if board[0][2] == currentPlayer && board[1][1] == currentPlayer && board[2][0] == currentPlayer {
		return true
	}

	return false
}

func isBoardFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == " " {
				return false
			}
		}
	}
	return true
}

func switchPlayer() {
	if currentPlayer == "X" {
		currentPlayer = "O"
	} else {
		currentPlayer = "X"
	}
}

func updateScores(player string) {
	if player == "X" {
		scoreX++
	} else {
		scoreO++
	}
}

func sendJSONResponse(w http.ResponseWriter, message string) {
	gameResponse := GameResponse{
		Board:         board,
		CurrentPlayer: currentPlayer,
		Message:       message,
		GameActive:    gameActive,
		ScoreX:        scoreX,
		ScoreO:        scoreO,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameResponse)
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

const indexHTML = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tic-Tac-Toe</title>
    <style>
        body {
            text-align: center;
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
        }

        h1 {
            color: #4646be;
            margin: 60px;
        }

        table {
            border-collapse: collapse;
            margin: 40px 40px;
        }

        .cell {
            width: 70px;
            height: 70px;
            border: 2px solid #ccc;
            display: inline-block;
            text-align: center;
            font-size: 24px;
            cursor: pointer;
            vertical-align: middle;
            line-height: 70px;
            color: #333;
        }

        .reset-btn {
            margin-top: 50px;
            padding: 10px 20px;
            font-size: 18px;
            cursor: pointer;
            background-color: #4646be;
            color: white;
            border: none;
            border-radius: 5px;
            position: fixed;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            transition: opacity 0.5s ease;
        }

        .reset-btn:hover {
            opacity: 0.8;
        }

        p {
            font-size: 20px;
            color: #333;
            margin: 20px;
        }

        .winning-cell-x {
            background-color: #a1e9a1;
        }

        .winning-cell-o {
            background-color: #f2a0a0;
        }

        
        @keyframes resetBoard {
            0% {
                transform: scale(1);
            }
            50% {
                transform: scale(0.1);
            }
            100% {
                transform: scale(1);
            }
        }
		.game-container {
	        text-align: center;

	    }

	    .board-and-scores {
	        display: flex;
	        align-items: center;
	        margin-top: 20px;
	        max-width: 500px; 
	        margin: 0 auto; 
	        
	    }
		
		.score-container {
		    text-align: left;
		    font-size: 18px;
		    font-weight: bold;
		}

		
		.player-x {
		    color: #2ecc71; 
		}

		.player-o {
		    color: #e74c3c; 
		}


		.winning-cell-x {
		    background-color: #a1e9a1; 
		}

		.winning-cell-o {
		    background-color: #f2a0a0; 
		}
    </style>
</head>

<body>
    <div class="game-container">
        <h1>Tic-Tac-Toe</h1>
        <p>Player's Turn: {{.CurrentPlayer}}</p>

        <div class="board-and-scores">
            <table>
                {{range $i, $row := .Board}}
                <tr>
                    {{range $j, $cell := $row}}
                    <td class="cell" onclick="makeMove({{$i}}, {{$j}})" style="border: 2px solid black;">{{$cell}}</td>
                    {{end}}
                </tr>
                {{end}}
            </table>

			<div class="score-container">
			    <p id="info">Score:</p>
			    <p class="score-x" id="scoreX">Player <span class="player-x">X</span>: {{.ScoreX}}</p>
			    <p class="score-o" id="scoreO">Player <span class="player-o">O</span>: {{.ScoreO}}</p>
			</div>
        </div>

        <p id="resultMessage"></p>
        <button class="reset-btn" onclick="resetGame()">Restart</button>
    </div>
    <script>
        function makeMove(row, col) {
            fetch("/makeMove?row=" + row + "&col=" + col)
                .then(response => response.json())
                .then(data => {
                    updateBoard(data);
                });
        }

        function resetGame() {
            document.querySelectorAll('.cell').forEach(cell => {
                cell.style.animation = 'resetBoard 1s ease';
                setTimeout(() => {
                    cell.style.animation = '';
                }, 1000);
            });

            
            fetch("/reset")
                .then(response => response.json())
                .then(data => {
                    updateBoard(data);
                });
        }

        function updateBoard(data) {
		    const cells = document.querySelectorAll('.cell');
		    let index = 0;

		    data.board.forEach((row, i) => {
		        row.forEach((cell, j) => {
		            cells[index].textContent = cell;
		            cells[index].style.color = cell === 'X' ? 'green' : 'red';

		            if (isWinningCell(i, j, data)) {
		                cells[index].classList.add(data.currentPlayer === 'X' ? 'winning-cell-x' : 'winning-cell-o');
		            } else {
		                cells[index].classList.remove('winning-cell-x', 'winning-cell-o');
		            }
		            index++;
		        });
		    });

		    const resultMessageElement = document.getElementById('resultMessage');
		    resultMessageElement.textContent = data.message;

		    if (data.message !== "") {
		        resultMessageElement.style.color = data.message.includes("X") ? 'green' : (data.message.includes("O") ? 'red' : 'black');
		        resultMessageElement.style.fontWeight = 'bold';
		        resultMessageElement.style.fontSize = '24px';
		    }

		    const currentPlayerElement = document.querySelector('p');
		    currentPlayerElement.textContent = "Player's Turn: ";
		    const playerIndicatorElement = document.createElement('span');
		    playerIndicatorElement.textContent = data.currentPlayer;
		    playerIndicatorElement.style.color = data.currentPlayer === 'X' ? 'green' : 'red';
		    currentPlayerElement.appendChild(playerIndicatorElement);

		   
		    const scoreXElement = document.getElementById('scoreX');
		    const scoreOElement = document.getElementById('scoreO');
		    scoreXElement.innerHTML = "Player <span class='player-x'>X</span>: " + data.scoreX;
		    scoreOElement.innerHTML = "Player <span class='player-o'>O</span>: " + data.scoreO;
		}

        function isWinningCell(row, col, data) {
            // Определение выигрышной стратегии
            return data.message.includes(data.currentPlayer) && (
                (checkRow(row, data) && data.board[row][col] === data.currentPlayer) ||
                (checkColumn(col, data) && data.board[row][col] === data.currentPlayer) ||
                (checkDiagonals(data) && (row === col || row + col === 2 || (row === 1 && col === 1)) && data.board[row][col] === data.currentPlayer)
            );
        }

        function checkRow(row, data) {
            return data.board[row][0] === data.currentPlayer &&
                data.board[row][1] === data.currentPlayer &&
                data.board[row][2] === data.currentPlayer;
        }

        function checkColumn(col, data) {
            return data.board[0][col] === data.currentPlayer &&
                data.board[1][col] === data.currentPlayer &&
                data.board[2][col] === data.currentPlayer;
        }

        function checkDiagonals(data) {
            return (data.board[0][0] === data.currentPlayer &&
                data.board[1][1] === data.currentPlayer &&
                data.board[2][2] === data.currentPlayer) ||
                (data.board[0][2] === data.currentPlayer &&
                    data.board[1][1] === data.currentPlayer &&
                    data.board[2][0] === data.currentPlayer);
        }
    </script>
</body>

</html>
`
