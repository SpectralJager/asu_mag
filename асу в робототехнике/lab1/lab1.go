// The control does not make full use of the ev3dev API where it could.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ev3go/ev3dev"
)

type PIDController struct {
	kp float64
	ki float64
	kd float64
}

func NewPIDController(kp, ki, kd float64) PIDController {
	return PIDController{
		kp: kp,
		ki: ki,
		kd: kd,
	}
}

func (pid *PIDController) Update(speed float64) {

}

func main() {
	weels()
}

func weels() {
	m1, err := ev3dev.TachoMotorFor("ev3-ports:outA", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find motor on outA: %v", err)
	}
	err = m1.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for left large motor on outB: %v", err)
	}

	m2, err := ev3dev.TachoMotorFor("ev3-ports:outB", "lego-ev3-l-motor")
	if err != nil {
		log.Fatalf("failed to find motor on outB: %v", err)
	}
	err = m2.SetStopAction("brake").Err()
	if err != nil {
		log.Fatalf("failed to set brake stop for left large motor on outB: %v", err)
	}

	m1.SetSpeedSetpoint(-10 * m1.MaxSpeed() / 100).Command("run-forever")
	m2.SetSpeedSetpoint(10 * m2.MaxSpeed() / 100).Command("run-forever")

	fmt.Println("Suslic")
	time.Sleep(time.Second * 2)

	m1.Command("stop")
	m2.Command("stop")
}

func sound() {
	const SoundPath = "/dev/input/by-path/platform-sound-event"
	speaker := ev3dev.NewSpeaker(SoundPath)
	err := speaker.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer speaker.Close()
	// Play tone at 440Hz for 200ms...
	err = speaker.Tone(440)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	// play tone at 220Hz for 200ms...
	err = speaker.Tone(220)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	// then stop tone playback.
	err = speaker.Tone(0)
	if err != nil {
		log.Fatal(err)
	}
}
