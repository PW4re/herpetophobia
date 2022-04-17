package server

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"snake/game"
	"snake/generators"
	"strconv"
	"strings"
)

var Initial = []byte{
	157, 79, 170, 8, 108, 234, 163, 16, 251, 181, 23, 148, 55, 162, 211, 186, 194, 222, 152, 207, 57, 97, 87, 45, 245, 141, 142, 40, 13, 92, 89, 64, 191, 102, 247, 178, 28, 138, 118, 68, 226, 24, 151, 103, 15, 139, 154, 244, 180, 83, 82, 196, 171, 167, 31, 155, 63, 246, 38, 200, 228, 120, 218, 204, 10, 238, 47, 56, 146, 185, 172, 158, 133, 53, 117, 42, 193, 241, 206, 86, 161, 0, 77, 243, 149, 239, 121, 129, 2, 85, 159, 59, 96, 164, 81, 220, 114, 18, 214, 65, 60, 125, 188, 201, 104, 174, 153, 75, 240, 223, 126, 35, 189, 113, 27, 236, 122, 143, 124, 73, 227, 43, 49, 67, 187, 48, 99, 250, 39, 20, 165, 115, 1, 177, 93, 232, 202, 249, 116, 54, 6, 242, 252, 69, 255, 22, 176, 197, 110, 5, 61, 169, 254, 183, 19, 229, 109, 150, 111, 131, 156, 253, 208, 145, 58, 179, 76, 7, 91, 78, 37, 233, 212, 9, 215, 192, 62, 209, 33, 32, 198, 168, 17, 195, 136, 166, 98, 130, 71, 248, 90, 217, 25, 30, 112, 34, 231, 3, 237, 21, 80, 224, 100, 66, 52, 84, 106, 4, 101, 205, 26, 105, 128, 225, 210, 135, 137, 175, 95, 70, 132, 203, 182, 29, 219, 190, 199, 44, 235, 140, 147, 74, 144, 46, 123, 216, 221, 14, 94, 127, 119, 36, 184, 88, 107, 12, 41, 134, 213, 72, 173, 160, 50, 51, 11, 230,
}
var Secret = []byte("ababubazzz")
var Counter = uint64(123)

type GameConn struct {
	conn  net.Conn
	stage ConnStage
	level *game.Level
}

func NewGameConn(conn net.Conn) GameConn {
	gameConn := GameConn{conn: conn, stage: INIT}
	return gameConn
}

func (gameConn *GameConn) startGame() error {
	seed, err := generators.GenerateSeed(Initial, Secret, Counter)
	if err != nil {
		return err
	}

	level, err := generators.GenerateLevel(seed)
	if err != nil {
		return err
	}
	gameConn.level = level
	_ = gameConn.level.Step(game.DIRECTION_RIGHT)
	return nil
}

func (gameConn *GameConn) handlePlay(msgBytes []byte) []byte {
	msg := strings.TrimSpace(string(bytes.Trim(msgBytes, "\x00")))
	direction := gameConn.level.Snake.Direction()
	switch msg {
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
	if gameConn.level.Status() != game.STATUS_UNFINISHED {
		gameConn.stage = INIT
	}
	if gameConn.level.Status() == game.STATUS_WIN {
		return []byte("flag")
	}
	if gameConn.level.Status() == game.STATUS_LOSE {
		return []byte("you lose")
	}
	return []byte(gameConn.level.Str())
}

func (gameConn *GameConn) handleInitMsg(msgBytes []byte) []byte {
	msg := strings.Split(string(bytes.Trim(msgBytes, "\x00")), " ")
	switch strings.TrimSpace(msg[0]) {
	case START:
		return gameConn.handleStart()
	case STOP:
		_ = gameConn.conn.Close()
		break
	case CREATE:
		return gameConn.handleCreate(msg)
	}
	return nil
}

func (gameConn *GameConn) handleCreate(msg []string) []byte {
	secret := []byte(msg[1])
	init, err := StringToBytes(msg[2])
	if err != nil {
		return []byte("Incorrect init")
	}
	counter := 1
	seed, err := generators.GenerateSeed(init, secret, uint64(counter))
	if err != nil {
		return []byte("Can't create seed")
	}
	gameConn.level, err = generators.GenerateLevel(seed)
	if err != nil {
		return []byte("Can't create level")
	}
	return []byte("Level created")
}

func (gameConn *GameConn) handleStart() []byte {
	gameConn.stage = PLAYING_GAME
	err := gameConn.startGame()
	if err != nil {
		_ = gameConn.conn.Close()
	}
	_, _ = gameConn.conn.Write([]byte(gameConn.level.Str()))
	return []byte(gameConn.level.Str())
}

func (gameConn *GameConn) handleConnection() {
	fmt.Println("handling")
	for {
		msgBytes := make([]byte, 2000)
		n, err := gameConn.conn.Read(msgBytes)
		if err != nil {
			fmt.Println("err")
			_ = gameConn.conn.Close()
			return
		}
		if n == 0 {
			continue
		}
		var answ []byte
		switch gameConn.stage {
		case INIT:
			answ = gameConn.handleInitMsg(msgBytes)
		case PLAYING_GAME:
			answ = gameConn.handlePlay(msgBytes)
		}
		_, _ = gameConn.conn.Write(answ)
	}
}

func StringToBytes(str string) ([]byte, error) {
	strs := strings.Split(str, ",")
	var bytes_ []byte
	for _, strEl := range strs {
		el, err := strconv.ParseUint(strEl, 10, 16)
		if err != nil {
			return nil, err
		}
		if el > 255 {
			return nil, errors.New("Incorrect num")
		}
		bytes_ = append(bytes_, byte(el))
	}
	return bytes_, nil
}
