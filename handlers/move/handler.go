package move

import (
	"errors"
	"log"

	"github.com/knei-knurow/roverd/modules/motors"
	"github.com/knei-knurow/roverd/modules/servos"
)

var (
	ErrMoveTypeNotFound = errors.New("move type not found")
	ErrInvalidMoveType  = errors.New("invalid move type")
)

// HandleMove calls function appropriate to move type ("go" and "turn").
func HandleMove(requestBody map[string]interface{}) error {
	moveType, ok := requestBody["type"]
	if !ok {
		return ErrMoveTypeNotFound
	}

	var err error
	if moveType == "go" {
		err = handleGoMove(requestBody)
	} else if moveType == "turn" {
		err = handleTurnMove(requestBody)
	} else {
		err = ErrInvalidMoveType
	}

	return err
}

// handleGoMove handles "go moves".
func handleGoMove(m map[string]interface{}) error {
	speed, ok := m["speed"].(float64)
	if !ok {
		log.Fatalln("failed to convert speed to float64")
	}

	direction, ok := m["direction"].(string)
	if !ok {
		log.Fatalln("failed to convert direction to string")
	}

	goMove := motors.GoMove{
		Direction: direction,
		Speed:     byte(speed), // watch out: speed can be only in range 0-255 (inclusive)
	}

	err := motors.ExecuteGoMove(goMove)
	if err != nil {
		log.Fatalln("failed to execute go move:", err)
	}

	return nil
}

// handleGoMove handles "turn moves".
func handleTurnMove(m map[string]interface{}) error {
	degrees, ok := m["degrees"].(float64)
	if !ok {
		log.Fatalln("failed to convert degrees to float64")
	}

	turnMove := servos.TurnMove{
		Degrees: byte(degrees),
	}

	err := servos.ExecuteTurnMove(turnMove)
	if err != nil {
		log.Fatalln("failed to execute turn move:", err)
	}

	return nil
}
