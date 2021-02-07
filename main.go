package main

import (
	"math"
	"strings"
	"time"

	"github.com/Steven-Ireland/path-of-gamepad/config"
	"github.com/Steven-Ireland/path-of-gamepad/controllers"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-vgo/robotgo"
)

var leftMousePosition = "up"
var rightMousePosition = "up"

func safeToggleMouseLeft(toggleTo string) {
	if leftMousePosition != toggleTo {
		leftMousePosition = toggleTo
		robotgo.MouseToggle(toggleTo)
	}
}

func safeToggleMouseRight(toggleTo string) {
	if rightMousePosition != toggleTo {
		rightMousePosition = toggleTo
		robotgo.MouseToggle(toggleTo, "right")
	}
}

func someButtonHeld(input controllers.Input) bool {
	var holdables = config.Holdable()

	var holding = false
	for _, k := range holdables {

		switch k {
		case "bumper_right":
			holding = holding || input.Right.Bumper
		case "bumper_left":
			holding = holding || input.Left.Bumper
		case "a":
			holding = holding || input.A
		case "b":
			holding = holding || input.B
		case "x":
			holding = holding || input.X
		case "y":
			holding = holding || input.Y
		// TODO, fill out rest but probably refactor first
		default:
			continue
		}

	}
	return holding
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	config.Load()

	gamepad := controllers.Gamepad{glfw.Joystick1, 0.17}
	lastInput := controllers.Input{}

	for {
		glfw.PollEvents()
		input := controllers.Read(gamepad, lastInput)
		var holding = someButtonHeld(input)

		if holding && !controllers.IsDeadZone(input.Right.Direction) && !controllers.IsDeadZone(input.Left.Direction) {
			var angle = math.Atan2(input.Right.Direction.Y, input.Right.Direction.X)
			var screenAdjustmentX = math.Cos(angle) * float64(config.AttackCircleRadius())
			var screenAdjustmentY = math.Sin(angle) * float64(config.AttackCircleRadius())

			robotgo.MoveMouse(
				(int)(float64(config.ScreenWidth())/2+screenAdjustmentX),
				(int)(float64(config.ScreenHeight())/2-screenAdjustmentY)-config.CharacterOffsetY(),
			)
			time.Sleep(50 * time.Millisecond)
		}

		if input.A_PRESS || input.A_UNPRESS {
			HandleMultiActions("a", input.A_UNPRESS)
		}
		if input.B_PRESS || input.B_UNPRESS {
			HandleMultiActions("b", input.B_UNPRESS)
		}
		if input.X_PRESS || input.X_UNPRESS {
			HandleMultiActions("x", input.X_UNPRESS)
		}
		if input.Y_PRESS || input.Y_UNPRESS {
			HandleMultiActions("y", input.X_UNPRESS)
		}
		if input.Start_PRESS {
			HandleMultiActions("start", false)
		}
		if input.Back_PRESS {
			HandleMultiActions("back", false)
		}
		if input.Left.Bumper_PRESS || input.Left.Bumper_UNPRESS {
			HandleMultiActions("bumper_left", input.Left.Bumper_UNPRESS)
		}
		if input.Right.Bumper_PRESS || input.Right.Bumper_UNPRESS {
			HandleMultiActions("bumper_right", input.Right.Bumper_UNPRESS)
		}
		if input.DPad.Up_PRESS {
			HandleMultiActions("dpad_up", false)
		}
		if input.DPad.Left_PRESS {
			HandleMultiActions("dpad_left", false)
		}
		if input.DPad.Down_PRESS {
			HandleMultiActions("dpad_down", false)
		}
		if input.DPad.Right_PRESS {
			HandleMultiActions("dpad_right", false)
		}

		if !(holding && !controllers.IsDeadZone(input.Left.Direction) && !controllers.IsDeadZone(input.Right.Direction)) {
			if controllers.IsDeadZone(input.Left.Direction) && !controllers.IsDeadZone(input.Right.Direction) {
				safeToggleMouseLeft("up")

				var screenAdjustmentX = input.Right.Direction.X * float64(config.FreeMouseSensitivity())
				var screenAdjustmentY = -1 * input.Right.Direction.Y * float64(config.FreeMouseSensitivity())

				robotgo.MoveRelative((int)(screenAdjustmentX), (int)(screenAdjustmentY))
			} else if controllers.IsDeadZone(input.Left.Direction) {
				safeToggleMouseLeft("up")
				// robotgo.MoveMouse(
				// 	(int)(SCREEN_RESOLUTION_X/2),
				// 	(int)(SCREEN_RESOLUTION_Y/2),
				// )
			} else {
				var angle = math.Atan2(input.Left.Direction.Y, input.Left.Direction.X)

				var screenAdjustmentX = math.Cos(angle) * float64(config.WalkCircleRadius())
				var screenAdjustmentY = math.Sin(angle) * float64(config.WalkCircleRadius())

				safeToggleMouseLeft("down")
				robotgo.DragMouse(
					(int)(float64(config.ScreenWidth())/2+screenAdjustmentX),
					(int)(float64(config.ScreenWidth())/2-screenAdjustmentY)-config.CharacterOffsetY(),
				)
			}
		} else if holding && !controllers.IsDeadZone(input.Right.Direction) {
			var angle = math.Atan2(input.Right.Direction.Y, input.Right.Direction.X)

			var screenAdjustmentX = math.Cos(angle) * float64(config.WalkCircleRadius())
			var screenAdjustmentY = math.Sin(angle) * float64(config.WalkCircleRadius())
			robotgo.MoveMouse(
				(int)(float64(config.ScreenWidth())/2+screenAdjustmentX),
				(int)(float64(config.ScreenWidth())/2-screenAdjustmentY)-config.CharacterOffsetY(),
			)
		}

		lastInput = input
		time.Sleep(5 * time.Millisecond)
	}

}

func HandleMultiActions(button string, unpressed bool) {
	if len(config.Buttons()[button]) > 0 {
		actions := strings.Split(config.Buttons()[button], ",")
		for _, a := range actions {
			HandleAction(a, unpressed)
		}
	}
}

func HandleAction(action string, unpressed bool) {
	switch action {
	case "RightClick":
		if unpressed {
			safeToggleMouseRight("up")
		} else {
			safeToggleMouseRight("down")
		}
	case "LeftClick":
		if unpressed {
			safeToggleMouseLeft("up")
		} else {
			safeToggleMouseLeft("down")
		}
	default:
		robotgo.KeyTap(action)
	}
}
