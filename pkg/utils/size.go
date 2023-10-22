package utils

import (
	"fmt"
)

func HumanReadableSize(size int) string {
	const scale = 1024
	if size < scale {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(scale), 0
	for n := size / scale; n >= scale; n /= scale {
		div *= scale
		exp++
	}
	return fmt.Sprintf("%.2f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
