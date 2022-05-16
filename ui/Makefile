# GITTAG=`git describe --tags`
# VERSION=`git rev-parse --short HEAD`
# flags=-ldflags="-s -w -X main.gitVersion=${VERSION} -X main.gitTag=${GITTAG} -extldflags -static"
flags=-ldflags="-s -w"

all: build

vet:
	go vet .

build:
	go clean; rm -rf pkg; go build -o ecm ${flags}

build_debug:
	go clean; rm -rf pkg; go build -o ecm ${flags} -gcflags="-m -m"

build_amd64: build_linux

build_darwin:
	go clean; rm -rf pkg ecm; GOOS=darwin go build -o ecm ${flags}

build_linux:
	go clean; rm -rf pkg ecm; GOOS=linux go build -o ecm ${flags}

build_power8:
	go clean; rm -rf pkg ecm; GOARCH=ppc64le GOOS=linux go build -o ecm ${flags}

build_arm64:
	go clean; rm -rf pkg ecm; GOARCH=arm64 GOOS=linux go build -o ecm ${flags}

build_windows:
	go clean; rm -rf pkg ecm; GOARCH=amd64 GOOS=windows go build -o ecm ${flags}

build_iphone:
	# for macOS provisioning file see instructions here:
	# https://developer.apple.com/documentation/xcode/distributing-your-app-to-registered-devices
	fyne release -os ios -appID com.example.myapp -appVersion 1.0 -appBuild 1 -profile development

build_android:
	# see additional instructions for Android: https://developer.android.com/studio/publish/app-signing
	export ANDROID_HOME=/Users/vk/Library/Android/sdk
	export ANDROID_NDK_HOME=/Users/vk/Library/Android/sdk/ndk/23.1.7779620
	fyne release -os android -appID com.example.myapp -appVersion 1.0 -appBuild 1 -keyStore /Users/vk/AndroidKeystore

install:
	go install

clean:
	go clean; rm -rf pkg

test : test1

test1:
	if ( test -d /tmp/kvdb ); then rm -rf /tmp/kvdb; fi
	go test -v -bench=.
