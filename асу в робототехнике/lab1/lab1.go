// The control does not make full use of the ev3dev API where it could.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ev3go/ev3dev"
)

var (
	TargetSpeed = flag.Int("target", 500, "set target speed")
	P           = flag.Float64("p", .01, "set pid's p coefficient")
	I           = flag.Float64("i", .1, "set pid's i coefficient")
	D           = flag.Float64("d", .1, "set pid's d coefficient")
	Duration    = flag.Int("duration", 10, "set duration in seconds")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	leftMotor, err := ev3dev.TachoMotorFor("ev3-ports:outA", "lego-ev3-l-motor")
	if err != nil {
		return fmt.Errorf("failed to find motor on port A: %w", err)
	}
	defer leftMotor.Command("stop")

	rightMotor, err := ev3dev.TachoMotorFor("ev3-ports:outB", "lego-ev3-l-motor")
	if err != nil {
		return fmt.Errorf("failed to find motor on port B: %w", err)
	}
	defer rightMotor.Command("stop")

	pid1 := NewPIDController(*P, *I, *D, float64(*TargetSpeed))
	pid2 := NewPIDController(*P, *I, *D, float64(*TargetSpeed))

	after := time.After(time.Second * time.Duration(*Duration))

	log.Printf("P=%.2f, I=%.2f, D=%.2f\nTarget=%d, Duration=%d sec\n", pid1.P, pid1.I, pid1.D, *TargetSpeed, *Duration)
	log.Printf("P=%.2f, I=%.2f, D=%.2f\nTarget=%d, Duration=%d sec\n", pid2.P, pid2.I, pid2.D, *TargetSpeed, *Duration)

	leftMotor.SetSpeedSetpoint(0).Command("run-forever")
	rightMotor.SetSpeedSetpoint(0).Command("run-forever")

	for {
		select {
		case <-after:
			return nil
		default:
		}
		leftSpeed, err := leftMotor.Speed()
		if err != nil {
			return fmt.Errorf("can't get speed of left motor: %w", err)
		}

		rightSpeed, err := rightMotor.Speed()
		if err != nil {
			return fmt.Errorf("can't get speed of right motor: %w", err)
		}

		newSpeed_left := pid1.Update(float64(leftSpeed))
		newSpeed_right := pid2.Update(float64(rightSpeed))
		if newSpeed_left >= 1050 {
			newSpeed_left = 1050
		}
		if newSpeed_right >= 1050 {
			newSpeed_right = 1050
		}
		log.Printf("New speed: %0.2f, %0.2f\n", newSpeed_left, newSpeed_right)

		leftMotor.SetSpeedSetpoint(int(newSpeed_left)).Command("run-forever")
		rightMotor.SetSpeedSetpoint(int(newSpeed_right)).Command("run-forever")

		time.Sleep(time.Millisecond * 500)
	}
}

// func weels() {
// 	m1, err := ev3dev.TachoMotorFor("ev3-ports:outA", "lego-ev3-l-motor")
// 	if err != nil {
// 		log.Fatalf("failed to find motor on outA: %v", err)
// 	}
// 	err = m1.SetStopAction("brake").Err()
// 	if err != nil {
// 		log.Fatalf("failed to set brake stop for left large motor on outB: %v", err)
// 	}

// 	m2, err := ev3dev.TachoMotorFor("ev3-ports:outB", "lego-ev3-l-motor")
// 	if err != nil {
// 		log.Fatalf("failed to find motor on outB: %v", err)
// 	}
// 	err = m2.SetStopAction("brake").Err()
// 	if err != nil {
// 		log.Fatalf("failed to set brake stop for left large motor on outB: %v", err)
// 	}

// 	m1.SetSpeedSetpoint(-10 * m1.MaxSpeed() / 100).Command("run-forever")
// 	m2.SetSpeedSetpoint(10 * m2.MaxSpeed() / 100).Command("run-forever")

// 	fmt.Println("Suslic")
// 	time.Sleep(time.Second * 2)

// 	m1.Command("stop")
// 	m2.Command("stop")
// }

// func sound() {
// 	const SoundPath = "/dev/input/by-path/platform-sound-event"
// 	speaker := ev3dev.NewSpeaker(SoundPath)
// 	err := speaker.Init()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer speaker.Close()
// 	// Play tone at 440Hz for 200ms...
// 	err = speaker.Tone(440)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	time.Sleep(200 * time.Millisecond)
// 	// play tone at 220Hz for 200ms...
// 	err = speaker.Tone(220)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	time.Sleep(200 * time.Millisecond)
// 	// then stop tone playback.
// 	err = speaker.Tone(0)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
