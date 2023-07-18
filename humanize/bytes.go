package humanize

import (
	"fmt"
	"math"
)

func bytes(size uint64, base int64, units []string) string {
	if size < 10 {
		return fmt.Sprintf("%d %s", size, units[0])
	}

	e := math.Floor(math.Log(float64(size)) / math.Log(float64(base)))
	val := math.Floor(float64(size)/math.Pow(float64(base), e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}
	return fmt.Sprintf(f, val, units[int(e)])
}

// Bytes convert to human-readable representation of an SI size.
//
// e.g. Bytes(1024) => "1.0 kB"
func Bytes(size uint64) string {
	return bytes(size, 1000, []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"})
}

// IBytes convert to human-readable representation of an IEC size.
//
// e.g. IBytes(1024) => "1.0 KiB"
func IBytes(size uint64) string {
	return bytes(size, 1024, []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"})
}
