
build:
	go build  -o example .

run:build
	./example

build-new:
	date +"%H:%M:%S" > ver.txt
	go build -o example.new .

http-hot_reload:build-new
	curl  http://127.0.0.1:8087/upload  -F "file=@./example.new" -v

new:
	./example.new
