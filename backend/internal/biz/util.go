package biz

import (
	"math"
	"strings"
)

func MaskPhone(in string) string {
	const (
		TailLen = 4
		MaskLen = 4
	)

	char := []rune(in)
	if len(char) < TailLen+MaskLen+1 {
		return strings.Repeat("*", len(char))
	}

	head := char[:len(char)-(TailLen+MaskLen)]
	tail := char[len(char)-TailLen:]
	starplaceholder := strings.Repeat("*", MaskLen)
	return string(head) + starplaceholder + string(tail)
}

// 返回单位为：米
func GetDistanceReturnMeter(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}
