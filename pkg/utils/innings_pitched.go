package utils

import (
	"fmt"
)

func ConvertOutsToInnings(outs int) string {
	innings := outs / 3
	remainder := outs % 3
	return fmt.Sprintf("%d.%d", innings, remainder)
}
