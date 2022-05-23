all: build

build: build_utils build_kvdb build_crypt build_vault build_cli build_server build_term build_wasm build_ui

build_utils:
	echo "Building utils package ..."
	cd utils; go mod tidy; go build; go test .; cd -

build_kvdb:
	echo "Building kvdb package ..."
	cd kvdb; go mod tidy; go build; go test .; cd -

build_crypt:
	echo "Building crypt package ..."
	cd crypt; go mod tidy; go build; go test .; cd -

build_vault:
	echo "Building vault package ..."
	cd vault; go mod tidy; go build; go test .; cd -

build_cli:
	echo "Building cli package ..."
	cd cli; go mod tidy; make; make test; cd -

build_server:
	echo "Building server package ..."
	cd server; go mod tidy; make; make test; cd -

build_term:
	echo "Building term package ..."
	cd term; go mod tidy; make; make test; cd -

build_wasm:
	echo "Building wasm package ..."
	cd wasm; go mod tidy; make; cd -

build_ui:
	echo "Building ui package ..."
	cd ui; go mod tidy; make; make test; cd -
