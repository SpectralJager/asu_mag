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
	Duration = flag.Int("duration", 10, "set duration in seconds")
)

const (
	LeftMotor  = "/sys/class/tacho-motor/motor0"
	RightMotor = "/sys/class/tacho-motor/motor1"

	DistanceSensor = "/sys/class/lego-sensor/sensor0"
	ColorSensor1   = "/sys/class/lego-sensor/sensor1"
	ColorSensor2   = "/sys/class/lego-sensor/sensor2"
	GyroSensor     = "/sys/class/lego-sensor/sensor3"
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

	correctionPID := asu.NewPIDController(.0, .0, .0, 0)
	distancePID := asu.NewPIDController(.0, .0, .0, 0)
	gyroPID := asu.NewPIDController(.0, .0, .0, 0)

	asu.SetSpeed(LeftMotor, int(0), 0, 0)
	asu.SetSpeed(RightMotor, int(0), 0, 0)

	for {
		select {
		case <-after:
			return nil
		default:
		}

		leftColor, err := asu.GetColor(ColorSensor1)
		if err != nil {
			return fmt.Errorf("can't get left color: %w", err)
		}
		rightColor, err := asu.GetColor(ColorSensor2)
		if err != nil {
			return fmt.Errorf("can't get right color: %w", err)
		}
		correction := correctionPID.Update(float64(leftColor - rightColor))

		newLeftSpeed := int(80 + (-1 * correction))
		newRightSpeed := int(80 + correction)

		distance, err := asu.GetDistance(DistanceSensor)
		if err != nil {
			return err
		}

		if distance < 200 {

		} else {
		}

		asu.SetSpeed(LeftMotor, newLeftSpeed, 0, 160)
		asu.SetSpeed(RightMotor, newRightSpeed, 0, 160)
		asu.Command(LeftMotor, "run-forever")
		asu.Command(RightMotor, "run-forever")

		time.Sleep(time.Millisecond * 100)
	}
}
