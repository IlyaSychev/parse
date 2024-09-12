package main

import (
	"fmt"
	"math"
	"parse/internal"
)

func main() {
	arrOfVel, arrOfPoints := internal.Parse()

	// мапа координата - 2 компоненты скорости
	InptData := make(map[[2]float64][]float64, 0)
	for i := range arrOfVel {
		InptData[[2]float64{arrOfPoints[i][0], arrOfPoints[i][1]}] = []float64{arrOfVel[i][0], arrOfVel[i][1]}
	}
	//fmt.Println(InptData)
	for i := range arrOfPoints {
		InptData[[2]float64{arrOfPoints[i][0], arrOfPoints[i][1]}] = []float64{arrOfVel[i][0], arrOfVel[i][1]}
	}

	//fmt.Println(len(InptData))
	//internall2.ParseTxt()

	var starty, startx float64 = 68, 37
	aerostat := NewPoint(float64(startx), starty)

	for t := 0; t <= 216000; t += 600 {
		aerostat.Polinom(InptData, 600)
	}
	fmt.Println(aerostat.path)
}

type Point struct {
	x    float64
	y    float64
	path [][]float64
}

func NewPoint(startX, startY float64) *Point {
	tmp := make([][]float64, 0)
	tmp = append(tmp, []float64{startX, startY})
	return &Point{startX, startY, tmp}
}

func (aerostat *Point) Polinom(vel map[[2]float64][]float64, dt float64) {
	fmt.Println("coordinates x and y = ", aerostat.x, aerostat.y)
	x1, x2 := round(aerostat.x)
	y1, y2 := round(aerostat.y)
	fmt.Println("local x1 x2 = ", x1, x2)
	fmt.Println("local y1 y2 = ", y1, y2)

	xCenter := (x2 + x1) / 2
	yCenter := (y2 + y1) / 2

	nu := xCenter - aerostat.x
	nu *= 2
	ksi := yCenter - aerostat.y
	ksi *= 2

	fnForm1 := float64(-1 * (ksi + 1) * (nu - 1) / 4)
	fnForm2 := float64((ksi + 1) * (nu + 1) / 4)
	fnForm3 := float64((ksi - 1) * (nu - 1) / 4)
	fnForm4 := float64(-1 * (ksi - 1) * (nu + 1) / 4)

	u := fnForm1*vel[[2]float64{x1, y1}][0] + fnForm2*vel[[2]float64{x2, y1}][0] + fnForm3*vel[[2]float64{x1, y2}][0] + fnForm4*vel[[2]float64{x2, y2}][0]
	v := fnForm1*vel[[2]float64{x1, y1}][1] + fnForm2*vel[[2]float64{x2, y1}][1] + fnForm3*vel[[2]float64{x1, y2}][1] + fnForm4*vel[[2]float64{x2, y2}][1]

	fmt.Println("u, v = ", u, v)
	aerostat.x = aerostat.x + metresToDegrees(aerostat.y, u*dt)
	aerostat.y = aerostat.y + metresToDegrees(0, v*dt)
	aerostat.path = append(aerostat.path, []float64{aerostat.x, aerostat.y})
	//fmt.Println(aerostat.x, aerostat.y)
}

// функция округления
func round(coordinate float64) (float64, float64) {
	coordinateMin := float64(math.Floor(float64(coordinate)))
	coordinateMax := coordinateMin + 1
	return coordinateMin, coordinateMax
}

func metresToDegrees(y, S float64) float64 {
	return S / (math.Pi * 6378000 / 180 * math.Cos(y/180*math.Pi))
}

// func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
// 	dx := float64(x2 - x1)
// 	dy := float64(y2 - y1)
// 	length := math.Sqrt(dx*dx + dy*dy)
// 	dx /= length
// 	dy /= length

// 	x := float64(x1)
// 	y := float64(y1)

// 	for i := 0; i <= int(length); i++ {
// 		img.Set(int(x), int(y), c)
// 		x += dx
// 		y += dy
// 	}
// }
