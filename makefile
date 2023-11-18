build_windows:
	go build -o app main.go

build_linux:
	go build -o app main.go

build_docker:
	python3 dockerBuild.py
docker_run:
	docker run --name email email