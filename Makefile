docker-build-linux: docker-go-mod
	docker run --rm -v $(shell readlink -f ../katzenpost):/go/katzenpost -v $(shell readlink -f .):/go/catchat/ -it catchat/go_mod bash -c 'cd /go/catchat/; go build'

docker-debian-base:
	if ! docker images|grep catchat/debian_base; then \
		docker run --name catchat_debian_base -it golang:buster bash -c 'apt update && apt upgrade -y && apt install -y --no-install-recommends build-essential libgles2 libgles2-mesa-dev libglib2.0-dev libxkbcommon-dev libxkbcommon-x11-dev libglu1-mesa-dev libxcursor-dev libwayland-dev libx11-xcb-dev' \
		&& docker commit catchat_debian_base catchat/debian_base \
		&& docker rm catchat_debian_base; \
	fi

docker-go-mod: docker-debian-base
	if ! docker images|grep catchat/go_mod; then \
		docker run -v $(shell readlink -f ../katzenpost):/go/katzenpost -v $(shell readlink -f .):/go/catchat --name catchat_go_mod -it catchat/debian_base \
			bash -c 'cd /go/catchat; go mod tidy' \
		&& docker commit catchat_go_mod catchat/go_mod \
		&& docker rm catchat_go_mod; \
	fi

docker-go-mod-update: docker-go-mod
	docker run -v $(shell readlink -f ../katzenpost):/go/katzenpost -v $(shell readlink -f .):/go/catchat --name catchat_go_mod -it catchat/go_mod \
			bash -c 'cd /go/catchat; go mod tidy' \
		&& docker commit catchat_go_mod catchat/go_mod \
		&& docker rm catchat_go_mod

docker-shell: docker-debian-base
	docker run -v $(shell readlink -f ../katzenpost):/go/katzenpost -v $(shell readlink -f .):/go/catchat --rm -it catchat/debian_base bash

docker-clean:
	docker rm  catchat_debian_base catchat_go_mod || true
	docker rmi catchat/debian_base catchat/go_mod || true
