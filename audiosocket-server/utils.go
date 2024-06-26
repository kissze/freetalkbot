package audiosocketserver

import (
	"encoding/binary"
	"log"
	"math"
	"net"
	"time"

	"github.com/CyCoreSystems/audiosocket"
)

// filterSilence filters out initial silence from the audio data
func filterSilence(data []byte, format string) []byte {
	threshold := 400                                          // Define el umbral de silencio según el formato
	silenceDuration := time.Duration(3000 * time.Millisecond) // Ajusta la duración según sea necesario

	var volumeFunc func([]byte) float64
	if format == "g711" {
		volumeFunc = calculateVolumeG711
	} else {
		volumeFunc = calculateVolumePCM16
	}

	var filteredData []byte
	silenceStart := time.Now()

	for len(data) > 0 {
		var chunkSize int
		if format == "g711" {
			chunkSize = 160 // Aproximadamente 20ms de audio para g711
		} else {
			chunkSize = 340 // Aproximadamente 20ms de audio para PCM
		}

		if len(data) < chunkSize {
			chunkSize = len(data)
		}

		chunk := data[:chunkSize+20]
		data = data[chunkSize:]

		volume := volumeFunc(chunk)
		log.Printf("Volume: %f\n", volume)
		if volume > float64(threshold) {
			filteredData = append(filteredData, chunk...)
		} else if time.Since(silenceStart) >= silenceDuration {
			filteredData = append(filteredData, chunk...)
		}
	}

	return filteredData
}

// Calculate the volume of the audio data. This is done by calculating the amplitude of the audio data wave.
// We are receiving 16-bit signed linear audio data.
func calculateVolumePCM16(buffer []byte) float64 {
	// Check if the buffer length is a multiple of 2
	if len(buffer)%2 != 0 {
		log.Println("Buffer length is not a multiple of 2")
		return 0
	}

	var sum float64

	// Iterate on the buffer by 2 bytes at a time
	for i := 0; i < len(buffer); i += 2 {
		// Takes two bytes of the buffer and converts them to a 16-bit signed integer in little-endian format
		// convert from unsigned int to signed int. This is the sample to be used for calculating the amplitude
		sample := int16(binary.LittleEndian.Uint16(buffer[i:]))
		// The amplitude of the audio data is calculated by squaring the sample and adding it to the sum
		sum += float64(sample) * float64(sample)
	}

	// And finally, the square root of the average, which is the sum of the samples divided by the number of samples.
	// This is the amplitude of the audio wave.
	return math.Sqrt(sum / float64(len(buffer)/2))
}

// ulawToLinear decodes a byte coded in g711 u-law format to a 16-bit signed linear PCM value.
func ulawToLinear(ulaw byte) int16 {
	ulaw ^= 0xFF
	sign := int16(ulaw & 0x80)
	exponent := int16((ulaw >> 4) & 0x07)
	mantissa := int16(ulaw & 0x0F)
	value := (mantissa << 4) + 0x08
	if exponent != 0 {
		value += 0x100
		value <<= (exponent - 1)
	}
	if sign != 0 {
		value = -value
	}
	return value
}

// Calculate volume data for G711 audio data
func calculateVolumeG711(buffer []byte) float64 {
	var sum float64
	sampleCount := len(buffer)
	for _, ulaw := range buffer {
		sample := ulawToLinear(ulaw)
		sum += float64(sample) * float64(sample)
	}
	return math.Sqrt(sum / float64(sampleCount))
}

// sendHangupSignal sends a hangup signal to the client
func sendHangupSignal(c net.Conn) {
	hangupMessage := audiosocket.HangupMessage()
	if _, err := c.Write(hangupMessage); err != nil {
		log.Println("Failed to send hangup signal:", err)
	} else {
		log.Println("Hangup signal sent successfully")
	}
}
