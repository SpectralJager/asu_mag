.PHONY: lab1
lab1:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab1/.

.PHONY: lab2
lab2:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab2/.

send:
	scp build/app robot@ev3dev.local:labs/