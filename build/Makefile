.PHONY: proto_go
#for mac
GENERATE_DIR := generated
BUILD_TARGET_MAC := ../build

#for windows
BUILD_TARGET_WIN := ./build
ANDROID_SDK_DIR=android

#--------------------------------------------------------------------IOS打包----------------------------------------------------------------------------------
proto_go_mac:
	rm -rf $(GENERATE_DIR)
	mkdir -p $(GENERATE_DIR)
	protoc proto/*.proto -Iproto --go_out=$(GENERATE_DIR)/
   	--go_opt=paths=import --experimental_allow_proto3_optional

proto_oc_mac:
	rm -rf $(BUILD_TARGET_MAC)/build_ios/proto
	mkdir -p $(BUILD_TARGET_MAC)/build_ios/proto
	protoc proto/*.proto -Iproto \
		--objc_out=$(BUILD_TARGET_MAC)/build_ios/proto

ios_mac:proto_go_mac proto_oc_mac
	rm -rf $(BUILD_TARGET_MAC)/build_ios/WorldIMSDK.xcframework
	mkdir -p $(BUILD_TARGET_MAC)/build_ios
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
 	-o ${BUILD_TARGET_MAC}/build_ios/WorldIMSDK.xcframework -target=ios ../src/api

#--------------------------------------------------------------------Windows Android打包----------------------------------------------------------------------------------
proto_go_windows:
	rmdir /s /q ${GENERATE_DIR}
	mkdir $(GENERATE_DIR)
	protoc proto/*.proto -Iproto --go_out=$(GENERATE_DIR) --go_opt=paths=import --experimental_allow_proto3_optional
proto_kt_windows:
	cd ${ANDROID_SDK_DIR}/ && gradlew :bean:clean && gradlew :bean:assembleRelease
	copy ${ANDROID_SDK_DIR}\\bean\\build\\outputs\\aar\\bean-release.aar  build_android
android_windows: proto_go_windows
	rmdir /s /q build_android
	mkdir build_android
	go get golang.org/x/mobile/cmd/gomobile
	gomobile init
	go mod download golang.org/x/exp
	gomobile bind -v -trimpath -ldflags "-s -w" \
	-o build_android\\im-service.aar -target=android ..\\src\\api
android_all: android_windows proto_kt_windows
#--------------------------------------------------------------------Mac Android打包----------------------------------------------------------------------------------
proto_kt_mac:
	cd ${BUILD_TARGET_MAC}/android/ && chmod +x gradlew && ./gradlew -Dhttp.proxyHost :bean:clean && ./gradlew -Dhttp.proxyHost :bean:assembleRelease
	cp ${BUILD_TARGET_MAC}/android/bean/build/outputs/aar/bean-release.aar  ${BUILD_TARGET_MAC}/android/
android_mac:proto_go_mac
	rm -rf ${BUILD_TARGET_MAC}/build_android
	mkdir -p ${BUILD_TARGET_MAC}/build_android
	go get golang.org/x/mobile
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" \
	-o ${BUILD_TARGET_MAC}/build_android/im-service.aar -target=android ../src/api
	unzip -d $(BUILD_TARGET_MAC)/build_android/sources $(BUILD_TARGET_MAC)/build_android/im-service-sources.jar