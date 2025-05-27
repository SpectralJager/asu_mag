// The control does not make full use of the ev3dev API where it could.
package main

import (
	"asu"
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	P        = flag.Float64("p", .0, "set pid's p coefficient")
	I        = flag.Float64("i", .0, "set pid's i coefficient")
	D        = flag.Float64("d", .0, "set pid's d coefficient")
	Duration = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"

	DistanceSensor   = "/sys/class/lego-sensor/sensor1"
	ColorSensorLeft  = "/sys/class/lego-sensor/sensor2"
	ColorSensorRight = "/sys/class/lego-sensor/sensor0"
	GyroSensor       = "/sys/class/lego-sensor/sensor3"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	defer asu.Command(LeftMotor, "stop")
	defer asu.Command(RightMotor, "stop")

	after := time.After(time.Second * time.Duration(*Duration))

	correctionPID := asu.NewPIDController(*P, *I, *D, 0)

	asu.SetSpeed(LeftMotor, int(0), 0, 0)
	asu.SetSpeed(RightMotor, int(0), 0, 0)

	state := 0
	initAngle := 0
	counter := 0
	rot := 50
	for {
		select {
		case <-after:
			return nil
		default:
		}

		fmt.Printf("state: %d, ", state)
		switch state {
		case 0:
			// follow line
			distance, _ := asu.GetDistance(DistanceSensor)
			fmt.Printf("distance: %d, ", distance)
			if distance < 100 {
				state = 1
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				leftColor, err := asu.GetColor(ColorSensorLeft)
				if err != nil {
					return fmt.Errorf("can't get left color: %w", err)
				}
				rightColor, err := asu.GetColor(ColorSensorRight)
				if err != nil {
					return fmt.Errorf("can't get right color: %w", err)
				}
				correction := correctionPID.Update(float64(leftColor - rightColor))
				fmt.Printf("correction: %f\n", correction)

				newLeftSpeed := int(80 + (-1 * correction))
				newRightSpeed := int(80 + correction)

				asu.SetSpeed(LeftMotor, newLeftSpeed, 0, 160)
				asu.SetSpeed(RightMotor, newRightSpeed, 0, 160)
			}
		case 1:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 85 || dif < -85 {
				state = 2
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, -50, -160, 160)
				asu.SetSpeed(RightMotor, 50, -160, 160)
			}
			// case 8:
			// 	leftColor, err := asu.GetColor(ColorSensorLeft)
			// 	if err != nil {
			// 		return fmt.Errorf("can't get left color: %w", err)
			// 	}
			// 	rightColor, err := asu.GetColor(ColorSensorRight)
			// 	if err != nil {
			// 		return fmt.Errorf("can't get right color: %w", err)
			// 	}
			// 	correction := correctionPID.Update(float64(leftColor - rightColor))
			// 	fmt.Printf("correction: %f\n", correction)
			// 	if correction >= 15 || correction <= -20 {
			// 		state = 0
			// 		asu.SetSpeed(LeftMotor, 0, 0, 160)
			// 		asu.SetSpeed(RightMotor, 0, 0, 160)
			// 	} else {
			// 		asu.SetSpeed(LeftMotor, 140, 0, 160)
			// 		asu.SetSpeed(RightMotor, 120, 0, 160)
			// 	}

		case 2:
			counter += 1
			if counter == rot {
				counter = 0
				state = 3
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 3:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 85 || dif < -85 {
				state = 4
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 50, -160, 160)
				asu.SetSpeed(RightMotor, -50, -160, 160)
			}
		case 4:
			counter += 1
			if counter == 80 {
				counter = 0
				state = 5
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 5:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 85 || dif < -85 {
				state = 6
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 50, -160, 160)
				asu.SetSpeed(RightMotor, -50, -160, 160)
			}
		case 6:
			counter += 1
			if counter == rot {
				counter = 0
				state = 7
				initAngle, _ = asu.GetAngle(GyroSensor)
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, 100, 0, 160)
				asu.SetSpeed(RightMotor, 100, 0, 160)
			}
		case 7:
			angle, _ := asu.GetAngle(GyroSensor)
			dif := initAngle - angle
			log.Printf("angle: %d\n", dif)
			if dif > 85 || dif < -85 {
				state = 0
				asu.SetSpeed(LeftMotor, 0, 0, 160)
				asu.SetSpeed(RightMotor, 0, 0, 160)
			} else {
				asu.SetSpeed(LeftMotor, -50, -160, 160)
				asu.SetSpeed(RightMotor, 50, -160, 160)
			}
		}

		asu.Command(LeftMotor, "run-forever")
		asu.Command(RightMotor, "run-forever")
		time.Sleep(time.Millisecond * 100)
	}
}
