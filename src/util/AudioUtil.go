package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/wav"
	"math"
	"math/cmplx"
)

type VoiceEffect struct {
	Name     string
	Factor   float64
	Shift    float64
	LowCutHz float64
}

var effects = []VoiceEffect{
	{Name: "Robot", Factor: 1.0, LowCutHz: 300},
	{Name: "Loli", Factor: 1.25, LowCutHz: 700},
	{Name: "Uncle", Factor: 0.8, LowCutHz: 100},
}

func Change(data []byte, effectIndex int) ([]byte, float64, error) {
	//effect := effects[effectIndex]
	w, e := wav.New(bytes.NewReader(data))
	if e != nil {
		return nil, 0, e
	}
	res := make([]byte, 0, w.Samples*2)
	for {
		fmt.Println("处理")
		buffer, e := w.ReadFloats(int(w.SampleRate * 2))
		if e != nil {
			break
		}
		// 将样本从float32转换为float64
		samples := make([]float64, len(buffer))
		for i := range buffer {
			samples[i] = float64(buffer[i])
		}
		// 对音频样本应用频率变换
		//changePitch(samples, effect.Factor)
		//cutLowFreq(samples, float64(w.SampleRate), effect.LowCutHz)
		//freqShift(samples, float64(w.SampleRate), effect.Shift)
		//将处理后的样本转换为float32
		for i := range samples {
			buffer[i] = float32(samples[i])
		}
		// 将float32类型的音频样本转换为int16类型
		us := make([]int16, len(buffer))
		for i := range samples {
			if buffer[i] > 1 {
				buffer[i] = 1
			} else if buffer[i] < -1 {
				buffer[i] = -1
			}
			us[i] = int16(buffer[i] * float32(math.MaxInt16))
		}
		// 将int16类型的样本转换为字节数据
		bytesData := bytes.NewBuffer(make([]byte, 0, w.Samples*2))
		for _, sample := range us {
			err := binary.Write(bytesData, binary.LittleEndian, sample)
			if err != nil {
				panic(err)
			}
		}
		res = append(res, bytesData.Bytes()...)
	}
	return res, w.Duration.Seconds(), nil
}

// changePitch 改变音频样本的基频
func changePitch(samples []float64, factor float64) {
	for i := range samples {
		samples[i] = samples[i] * factor
	}
}

// cutLowFreq 剪切低频部分
func cutLowFreq(samples []float64, sampleRate float64, lowCutHz float64) {
	cutOff := int(sampleRate / lowCutHz)
	for i := range samples {
		if i < cutOff {
			samples[i] = 0
		}
	}
}

// freqShift 将输入的PCM样本数据提高或降低指定的半音数
func freqShift(samples []float64, sampleRate float64, semitones float64) {
	N := len(samples)
	halfN := N / 2
	original := samples[0:halfN]
	reversed := samples[halfN:]
	copy(reversed, original)
	for i := range reversed {
		reversed[i] *= -1
	}
	samples = append(original, reversed...)

	// 计算傅立叶变换
	transformed := fft.FFTReal(samples)

	// 提高或降低频率
	for i := range transformed {
		transformed[i] *= cmplx.Rect(math.Pow(2, semitones/12.0), 0)
	}

	// 计算逆傅立叶变换
	result := fft.IFFT(transformed)

	// 复原复数表示并将其转换回PCM格式
	for i := range result {
		samples[i] = real(result[i])
	}
}
