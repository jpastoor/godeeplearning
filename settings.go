package main

import (
	"strconv"
	"fmt"
	"strings"
)

type Settings struct {
	playerNames []string
	yourBot     string
	timebank    int
	timePerMove int
	yourBotId   int
	fieldWidth  int
	fieldHeight int
	maxRounds   int
}

func (s *Settings) update(key, value string) error {
	var err error
	switch key {
	case "player_names":
		s.playerNames = strings.Split(value, ",")
	case "your_bot":
		s.yourBot = value
	case "timebank":
		s.timebank, err = strconv.Atoi(value)
	case "time_per_move":
		s.timePerMove, err = strconv.Atoi(value)
	case "your_botid":
		s.yourBotId, err = strconv.Atoi(value)
	case "field_width":
		s.fieldWidth, err = strconv.Atoi(value)
	case "field_height":
		s.fieldHeight, err = strconv.Atoi(value)
	case "max_rounds":
		s.maxRounds, err = strconv.Atoi(value)
	default:
		err = fmt.Errorf("Unrecognised settings key: %s", key)
	}

	return err
}
