.PHONY: proto_go

GENERATE_DIR := generated
BUILD_TARGET := ./build
ANDROID_SDK_DIR=./android
proto_go_windows:
	rmdir /s /q ${GENERATE_DIR}
	mkdir $(GENERATE_DIR)
	protoc proto/*.proto -Iproto --go_out=$(GENERATE_DIR) --go_opt=paths=import --experimental_allow_proto3_optional

proto_go:
	rm -rf $(GENERATE_DIR)
	mkdir -p $(GENERATE_DIR)
	protoc proto/*.proto -Iproto --go_out=$(GENERATE_DIR)/
   	--go_opt=paths=import --experimental_allow_proto3_optional

proto_oc:
	rm -rf $(BUILD_TARGET)/build_ios/proto
	mkdir -p $(BUILD_TARGET)/build_ios/proto
	protoc proto/*.proto -Iproto \
		--objc_out=$(BUILD_TARGET)/build_ios/proto

proto_kt:
	cd ${ANDROID_SDK_DIR}/ && chmod +x gradlew && ./gradlew -Dhttp.proxyHost :bean:clean && ./gradlew -Dhttp.proxyHost :bean:assembleRelease
	cp ${ANDROID_SDK_DIR}/bean/build/outputs/aar/bean-release.aar  ${BUILD_TARGET}/android/

ios:proto_go proto_oc
	rm -rf $(BUILD_TARGET)/build_ios/Wallet.xcframework
	mkdir -p $(BUILD_TARGET)/build_ios
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
 	-o ${BUILD_TARGET}/build_ios/Wallet.xcframework -target=ios ./api

android:proto_go
	rm -rf ${BUILD_TARGET}/build_android
	mkdir -p ${BUILD_TARGET}/build_android
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
	-o ${BUILD_TARGET}/build_android/im-service.aar -target=android ./api
	unzip -d $(BUILD_TARGET)/build_android/sources $(BUILD_TARGET)/build_android/im-service-sources.jar
repo_android:android
	cd  ${ANDROID_SDK_DIR}/ && chmod +x gradlew
	cd  ${ANDROID_SDK_DIR}/ && gradlew -Dhttp.proxyHost generateLib

repo_android_push: proto_kt repo_android