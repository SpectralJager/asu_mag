// The control does not make full use of the ev3dev API where it could.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	TargetSpeed = flag.Float64("target", 100, "set target speed")
	P           = flag.Float64("p", .0, "set pid's p coefficient")
	I           = flag.Float64("i", .0, "set pid's i coefficient")
	D           = flag.Float64("d", .0, "set pid's d coefficient")
	Duration    = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"
	Sensor1    = "/sys/class/lego-sensor/sensor1"
	Sensor2    = "/sys/class/lego-sensor/sensor2"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	defer Command(LeftMotor, "stop")
	defer Command(RightMotor, "stop")

	pid := NewPIDController(*P, *I, *D, 0)

	after := time.After(time.Second * time.Duration(*Duration))

	SetSpeed(LeftMotor, int(0))
	SetSpeed(RightMotor, int(0))

	for {
		select {
		case <-after:
			return nil
		default:
		}
		leftColor, err := GetColor(Sensor1)
		if err != nil {
			return fmt.Errorf("can't get left color: %w", err)
		}

		rightColor, err := GetColor(Sensor2)
		if err != nil {
			return fmt.Errorf("can't get right color: %w", err)
		}

		correction := pid.Update(float64(leftColor - rightColor))

		newLeftSpeed := int(80 + (-1 * correction))
		newRightSpeed := int(80 + correction)

		if newLeftSpeed > 160 {
			newLeftSpeed = 160
		}
		if newRightSpeed > 160 {
			newRightSpeed = 160
		}
		if newLeftSpeed < 0 {
			newLeftSpeed = 0
		}
		if newRightSpeed < 0 {
			newRightSpeed = 0
		}

		SetSpeed(LeftMotor, newLeftSpeed)
		SetSpeed(RightMotor, newRightSpeed)
		// switch {
		// case correction < -10:
		// 	SetSpeed(LeftMotor, 80)
		// 	SetSpeed(RightMotor, 0)
		// case correction > 10:
		// 	SetSpeed(LeftMotor, 0)
		// 	SetSpeed(RightMotor, 80)
		// default:
		// 	SetSpeed(LeftMotor, 80)
		// 	SetSpeed(RightMotor, 80)
		// }
		fmt.Printf(
			"Correction: %.2f, left_color: %d, right_color: %d, left speed: %d, right speed: %d\n",
			correction,
			leftColor,
			rightColor,
			Must(GetSpeed(LeftMotor)),
			Must(GetSpeed(RightMotor)),
		)

		Command(LeftMotor, "run-forever")
		Command(RightMotor, "run-forever")

		time.Sleep(time.Millisecond * 100)
	}
}

func Command(motor string, command string) error {
	commandPath := path.Join(motor, "command")

	return os.WriteFile(commandPath, []byte(command), 0644)
}

func GetColor(sensor string) (int, error) {
	distancePath := path.Join(sensor, "value0")
	data, err := os.ReadFile(distancePath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data[:len(data)-1]))
}

func SetSpeed(motor string, speed int) error {
	if speed >= 1050 {
		speed = 1050
	}
	speedPath := path.Join(motor, "speed_sp")

	return os.WriteFile(speedPath, []byte(strconv.Itoa(speed)), 0644)
}

func GetSpeed(motor string) (int, error) {
	speedPath := path.Join(motor, "speed")
	data, err := os.ReadFile(speedPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data[:len(data)-1]))
}
