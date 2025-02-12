package m3u8

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	AdStartPrefix = []byte("#EXT-X-CUE-OUT:")
	AdEnd         = line("#EXT-X-CUE-IN")
)

var (
	ErrAdLineIsTooShort = errors.New("ad line is too short")
)

func SearchAdStart(playlist Lines, dst []byte) (updated []byte, adStartLine line) {
	for line := range playlist.Lines() {
		if bytes.HasPrefix(line, AdStartPrefix) {
			return dst, line
		}

		dst = AppendLine(dst, line)
	}
	return
}

func SearchAdEnd(playlist Lines, dst []byte) (updated []byte, adEndLine line) {
	for line := range playlist.Lines() {
		if bytes.HasPrefix(line, AdEnd) {
			return dst, line
		}

		dst = AppendLine(dst, line)
	}
	return
}

func ExtractAdDuration(adLine line) (time.Duration, error) {
	const op = "m3u8.ExtractAdDuration"

	adLine = bytes.TrimSpace(adLine)
	if len(adLine) < len(AdStartPrefix) {
		return 0, fmt.Errorf("%s: %w", op, ErrAdLineIsTooShort)
	}

	durStr := string(adLine[len(AdStartPrefix):])
	durInSec, err := strconv.ParseFloat(durStr, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return time.Duration(durInSec*1000) * time.Millisecond, nil
}
