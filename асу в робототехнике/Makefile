.PHONY: lab1
lab1:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab1/.

.PHONY: lab2
lab2:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab2/.

.PHONY: lab3
lab3:
	GOOS=linux GOARCH=arm GOARM=5 go build -o build/app ./lab3/.

send:
	scp build/app robot@192.168.175.78:labs/