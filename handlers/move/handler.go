package move

import (
	"errors"
	"fmt"

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
		return errors.New("cannot convert speed to float64")
	}

	direction, ok := m["direction"].(string)
	if !ok {
		return errors.New("cannot convert direction to string")
	}

	goMove := motors.GoMove{
		Direction: direction,
		Speed:     byte(speed), // watch out: speed can be only in range 0-255 (inclusive)
	}

	err := motors.ExecuteGoMove(goMove)
	if err != nil {
		return fmt.Errorf("execute go move: %v", err)
	}

	return nil
}

// handleGoMove handles "turn moves".
func handleTurnMove(m map[string]interface{}) error {
	degrees, ok := m["degrees"].(float64)
	if !ok {
		return errors.New("cannot convert degrees to float64")
	}

	turnMove := servos.TurnMove{
		Degrees: byte(degrees),
	}

	err := servos.ExecuteTurnMove(turnMove)
	if err != nil {
		return fmt.Errorf("execute turn move: %v", err)
	}

	return nil
}
