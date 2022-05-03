package http

import (
	"github.com/gorilla/websocket"
	"snake/game"
	"snake/generators"
)

var Initial = []byte{
	157, 79, 170, 8, 108, 234, 163, 16, 251, 181, 23, 148, 55, 162, 211, 186, 194, 222, 152, 207, 57, 97, 87, 45, 245, 141, 142, 40, 13, 92, 89, 64, 191, 102, 247, 178, 28, 138, 118, 68, 226, 24, 151, 103, 15, 139, 154, 244, 180, 83, 82, 196, 171, 167, 31, 155, 63, 246, 38, 200, 228, 120, 218, 204, 10, 238, 47, 56, 146, 185, 172, 158, 133, 53, 117, 42, 193, 241, 206, 86, 161, 0, 77, 243, 149, 239, 121, 129, 2, 85, 159, 59, 96, 164, 81, 220, 114, 18, 214, 65, 60, 125, 188, 201, 104, 174, 153, 75, 240, 223, 126, 35, 189, 113, 27, 236, 122, 143, 124, 73, 227, 43, 49, 67, 187, 48, 99, 250, 39, 20, 165, 115, 1, 177, 93, 232, 202, 249, 116, 54, 6, 242, 252, 69, 255, 22, 176, 197, 110, 5, 61, 169, 254, 183, 19, 229, 109, 150, 111, 131, 156, 253, 208, 145, 58, 179, 76, 7, 91, 78, 37, 233, 212, 9, 215, 192, 62, 209, 33, 32, 198, 168, 17, 195, 136, 166, 98, 130, 71, 248, 90, 217, 25, 30, 112, 34, 231, 3, 237, 21, 80, 224, 100, 66, 52, 84, 106, 4, 101, 205, 26, 105, 128, 225, 210, 135, 137, 175, 95, 70, 132, 203, 182, 29, 219, 190, 199, 44, 235, 140, 147, 74, 144, 46, 123, 216, 221, 14, 94, 127, 119, 36, 184, 88, 107, 12, 41, 134, 213, 72, 173, 160, 50, 51, 11, 230,
}
var Secret = []byte("ababubazzz")
var Counter = uint64(123)

type GameConn struct {
	conn   *websocket.Conn
	level  *game.Level
	gameId int
}

type MoveMsg struct {
	Direction string `json:"direction"`
	CloseGame bool   `json:"closeGame"`
	NewGame   bool   `json:"newGame"`
}

type EndGameAnsw struct {
	Permutation [256]int `json:"permutation"`
	Counter     int      `json:"counter"`
	GameResult  string   `json:"gameResult"`
	Prize       string   `json:"prize"`
}

type MoveAnsw struct {
}

type ErrAnsw struct {
	msg string
}

func NewGameConn(conn *websocket.Conn, gameId int) GameConn {
	gameConn := GameConn{conn: conn, gameId: gameId}
	return gameConn
}

func (gameConn *GameConn) startGame() {
	//todo: получение игры из базы
	seed, err := generators.GenerateSeed(Initial, Secret, Counter)
	if err != nil {
		gameConn.conn.WriteJSON(&ErrAnsw{msg: err.Error()})
		return
	}

	level, err := generators.GenerateLevel(seed)
	if err != nil {
		gameConn.conn.WriteJSON(&ErrAnsw{msg: err.Error()})
		return
	}
	gameConn.level = level
	_ = gameConn.level.Step(game.DIRECTION_RIGHT)
	for {
		var gameMsg MoveMsg
		err = gameConn.conn.ReadJSON(&gameMsg) //тут бы сделать ожидание сообщения
		if err != nil {
			gameConn.conn.WriteJSON(&ErrAnsw{msg: err.Error()})
			return
		}
		moveAnsw := gameConn.handleGame(gameMsg)
		if gameConn.level.Status() != game.STATUS_UNFINISHED {
			_ = gameConn.conn.WriteJSON(EndGameAnsw{})
		}
		_ = gameConn.conn.WriteJSON(moveAnsw)
	}
}

func (gameConn *GameConn) handleGame(msg MoveMsg) MoveAnsw {
	direction := gameConn.level.Snake.Direction()
	switch msg.Direction {
	case UP:
		direction = game.DIRECTION_UP
	case DOWN:
		direction = game.DIRECTION_DOWN
	case LEFT:
		direction = game.DIRECTION_LEFT
	case RIGHT:
		direction = game.DIRECTION_RIGHT
	}
	_ = gameConn.level.Step(direction)
	if gameConn.level.Status() == game.STATUS_WIN {
		return MoveAnsw{}
	}
	if gameConn.level.Status() == game.STATUS_LOSE {
		return MoveAnsw{}
	}
	return MoveAnsw{}
}
