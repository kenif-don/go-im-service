.PHONY: proto_go

GENERATE_DIR := generated
BUILD_TARGET := dist
ANDROID_SDK_DIR=./android
TIME :=$(shell date +'%Y/%m/%d-%H:%M:%S')
proto_go:
	rmdir /s /q ${GENERATE_DIR}
	mkdir $(GENERATE_DIR)
	protoc proto/*.proto -Iproto --go_out=$(GENERATE_DIR)/
   	--go_opt=paths=import --experimental_allow_proto3_optional


proto_oc:
	rmdir /s /q $(BUILD_TARGET)/build_ios/proto
	mkdir $(BUILD_TARGET)/build_ios/proto
	protoc proto/*.proto -Iproto \
		--swift_opt=Visibility=Public \
		--swift_out=$(BUILD_TARGET)/build_ios/proto

proto_kt:
	cd  ${ANDROID_SDK_DIR}/ && gradlew -Dhttp.proxyHost :bean:clean && gradlew -Dhttp.proxyHost :bean:assembleRelease

ios:proto_go proto_oc
	rmdir /s /q $(BUILD_TARGET)/build_ios/BtcApi.xcframework
	mkdir $(BUILD_TARGET)/build_ios
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
 	-o ${BUILD_TARGET}/build_ios/Wallet.xcframework -target=ios ./api

android:proto_go
	rmdir /s /q ${BUILD_TARGET}\\android
	mkdir ${BUILD_TARGET}\\android
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	SET GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
	-o ${BUILD_TARGET}/android/IM-SDK.aar -target=android api/*
	unzip -d $(BUILD_TARGET)/android/sources $(BUILD_TARGET)/android/IM-SDK-sources.jar
repo_android:android
	cd  ${ANDROID_SDK_DIR}/ && chmod +x gradlew
	cd  ${ANDROID_SDK_DIR}/ && gradlew -Dhttp.proxyHost generateLib

repo_android_push: proto_kt repo_android