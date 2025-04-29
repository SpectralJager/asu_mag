.PHONY: lab1
lab1:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab1/.

.PHONY: lab2
lab2:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab2/.

.PHONY: lab3
lab3:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab3/.

.PHONY: lab4
lab4:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab4/.

send:
	scp build/app robot@ev3dev.local:labs/