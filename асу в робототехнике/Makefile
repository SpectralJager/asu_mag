.PHONY: lab1
lab1:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab1/.

send:
	scp build/app robot@ev3dev.local:labs/