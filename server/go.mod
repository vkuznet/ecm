module github.com/vkuznet/ecm-server

go 1.18

require (
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/disintegration/imaging v1.6.2
	github.com/gorilla/context v1.1.1
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/ulule/limiter/v3 v3.10.0
	github.com/vkuznet/ecm/crypt v0.0.0-20220524114141-5e230e2ed56b
	github.com/vkuznet/ecm/kvdb v0.0.0-20220524114141-5e230e2ed56b
	github.com/vkuznet/ecm/utils v0.0.0-20220524114141-5e230e2ed56b
	github.com/vkuznet/ecm/vault v0.0.0-20220524114141-5e230e2ed56b
	github.com/vkuznet/http-logging v0.0.0-20210729230351-fc50acd79868
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898
	rsc.io/qr v0.2.0
)

require (
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.2 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v2.0.6+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/klauspost/compress v1.15.4 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/exp/errors v0.0.0-20220518171630-0b5c67f07fdf // indirect
	golang.org/x/image v0.0.0-20220413100746-70e8d0d3baa9 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

// replace github.com/vkuznet/ecm/crypt => /Users/vk/Work/Languages/Go/ecm/crypt

// replace github.com/vkuznet/ecm/utils => /Users/vk/Work/Languages/Go/ecm/utils

// replace github.com/vkuznet/ecm/vault => /Users/vk/Work/Languages/Go/ecm/vault
