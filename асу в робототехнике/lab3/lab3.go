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
	defer asu.Command(LeftMotor, "stop")
	defer asu.Command(RightMotor, "stop")

	pid := asu.NewPIDController(*P, *I, *D, 0)

	after := time.After(time.Second * time.Duration(*Duration))

	asu.SetSpeed(LeftMotor, int(0), 0, 0)
	asu.SetSpeed(RightMotor, int(0), 0, 0)

	for {
		select {
		case <-after:
			return nil
		default:
		}
		leftColor, err := asu.GetColor(Sensor1)
		if err != nil {
			return fmt.Errorf("can't get left color: %w", err)
		}

		rightColor, err := asu.GetColor(Sensor2)
		if err != nil {
			return fmt.Errorf("can't get right color: %w", err)
		}

		correction := pid.Update(float64(leftColor - rightColor))

		newLeftSpeed := int(80 + (-1 * correction))
		newRightSpeed := int(80 + correction)

		asu.SetSpeed(LeftMotor, newLeftSpeed, 0, 160)
		asu.SetSpeed(RightMotor, newRightSpeed, 0, 160)
		fmt.Printf(
			"Correction: %.2f, left_color: %d, right_color: %d, left speed: %d, right speed: %d\n",
			correction,
			leftColor,
			rightColor,
			asu.Must(asu.GetSpeed(LeftMotor)),
			asu.Must(asu.GetSpeed(RightMotor)),
		)

		asu.Command(LeftMotor, "run-forever")
		asu.Command(RightMotor, "run-forever")

		time.Sleep(time.Millisecond * 100)
	}
}
