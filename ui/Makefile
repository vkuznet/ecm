GITTAG=`git describe --tags`
VERSION=`git rev-parse --short HEAD`
flags=-ldflags="-s -w -X main.appKind=desktop -X main.gitVersion=${VERSION} -X main.gitTag=${GITTAG} -X main.supportedProviders=dropbox"
# flags=-ldflags="-s -w -X main.appKind=desktop"
# flags=-ldflags="-s -w"

all: credentials build

vet:
	go vet .

credentials:
	# create bundle for credentials
	fyne bundle credentials.env > credentials.go

sign:
	sudo codesign --force --deep --sign - ${PWD}/ecm

build: credentials
	go clean; rm -rf pkg; go build -o ecm ${flags}

build_debug: credentials
	go clean; rm -rf pkg; go build -o ecm ${flags} -gcflags="-m -m"

build_amd64: build_linux

build_darwin: credentials
	go clean; rm -rf pkg ecm; GOOS=darwin go build -o ecm ${flags}

build_linux: credentials
	go clean; rm -rf pkg ecm; GOOS=linux go build -o ecm ${flags}

build_power8: credentials
	go clean; rm -rf pkg ecm; CGO_ENABLED=1 GOARCH=ppc64le GOOS=linux go build -o ecm ${flags}

build_arm64: credentials
	go clean; rm -rf pkg ecm; GOARCH=arm64 GOOS=linux go build -o ecm ${flags}

build_windows: credentials
	go clean; rm -rf pkg ecm; GOARCH=amd64 GOOS=windows go build -o ecm ${flags}

build_app: credentials
	fyne build -os darwin -verbose

build_ios: credentials
	# for macOS provisioning file see instructions here:
	# https://developer.apple.com/documentation/xcode/distributing-your-app-to-registered-devices
	# we need to have signed certificate, so setup IOS_CERT environment and use it
	fyne release -os iossimulator -verbose -appID io.github.vkuznet -appVersion 1.0 -appBuild 1 -profile developer -cert ${IOS_CERT}

release_darwin: credentials
	# for macOS provisioning file see instructions here:
	# https://developer.apple.com/documentation/xcode/distributing-your-app-to-registered-devices
	# we need to have signed certificate, so setup IOS_CERT environment and use it
	fyne release -os darwin -appID io.github.vkuznet -appVersion 1.0 -appBuild 1 -profile developer -cert ${IOS_CERT} -category utilities

clean_aab:
	rm -f ecm.aab

build_android: credentials clean_aab
	# see additional instructions for Android: https://developer.android.com/studio/publish/app-signing
	# https://stackoverflow.com/questions/25975320/create-android-keystory-private-key-command-line
	ANDROID_HOME=/Users/vk/Library/Android/sdk \
	ANDROID_NDK_HOME=/Users/vk/Library/Android/sdk/ndk/23.1.7779620 \
	fyne release -os android -appID io.github.vkuznet -appVersion 1.0 -appBuild 1 -keyStore ecm.keystore -keyName ecm-key

install:
	go install

clean:
	go clean; rm -rf pkg

test : test1

test1:
	if ( test -d /tmp/kvdb ); then rm -rf /tmp/kvdb; fi
	go test -v -bench=.
