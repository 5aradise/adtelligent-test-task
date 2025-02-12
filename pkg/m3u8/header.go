package m3u8

import (
	"strconv"
)

func AdHeader(durationInMs int) line {
	durInSec := float64(durationInMs) / 1000
	durStr := strconv.FormatFloat(durInSec, 'f', 3, 64)

	header := make(line, len(AdStartPrefix), len(AdStartPrefix)+len(durStr)+1)
	copy(header, AdStartPrefix)
	header = append(header, []byte(durStr)...)
	header = append(header, '\n')

	return header
}
