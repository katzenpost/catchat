linux: ccdeps
	docker run --rm -v $(shell readlink -f ..):/go/katzenpost/ -it katzenpost/ccdeps bash -c 'cd /go/katzenpost/catchat/; go build'

ccdeps:
	if ! docker images|grep katzenpost/ccdeps; then \
		docker run --name katzenpost_ccdeps -it katzenpost/deps bash -c 'apt install -y --no-install-recommends build-essential libgles2 libgles2-mesa-dev libglib2.0-dev libxkbcommon-dev libxkbcommon-x11-dev libglu1-mesa-dev libxcursor-dev libwayland-dev libx11-xcb-dev' \
		&& docker commit katzenpost_ccdeps katzenpost/ccdeps \
		&& docker rm katzenpost_ccdeps; \
	fi
