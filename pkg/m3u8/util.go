package m3u8

func AppendLine(sl []byte, line line) []byte {
	return append(sl, append(line, byte('\n'))...)
}
