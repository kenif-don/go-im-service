package util

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/format/wav"
	"os"
)

func test() {
	// 打开输入音频文件
	inputFile, err := os.Open("input.wav")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	// 解码输入音频文件
	inputStream, format, err := wav.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	// 创建变声效果
	pitchShift := &effects.PitchShift{
		Shift: 3, // 萝莉音效的音调变化量，可以调整以达到所需效果
	}

	// 将输入音频流应用于变声效果
	pitchedStream := pitchShift.Stream(inputStream)

	// 创建输出音频文件
	outputFile, err := os.Create("output.wav")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// 编码并写入输出音频文件
	err = wav.Encode(outputFile, beep.Resample(4, format.SampleRate, format.SampleRate, pitchedStream))
	if err != nil {
		panic(err)
	}

	println("变声完成，输出文件：output.wav")
}
