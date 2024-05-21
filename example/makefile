
ver:
	date +"%H:%M:%S" > ver.txt

build:ver
	go build  -o example .

run:build
	./example

build-new:ver
	go build -o example.new .

http-hot_reload:build-new
	curl  http://127.0.0.1:8087/upload  -F "file=@./example.new" -v

new:
	./example.new
