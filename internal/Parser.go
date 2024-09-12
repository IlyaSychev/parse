package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"

	"github.com/amsokol/go-grib2"
)

const filename = "gdas.t06z.pgrb2.1p00 (1).anl"

func Parse() ([][]float64, [][2]float64) {
	log.Println("Downloading data from HTTP server...")
	grib2file := filename
	infile, err := os.Open(grib2file)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	data, err := io.ReadAll(infile)
	if err != nil {
		log.Fatal(err)
	}

	gribs, err := grib2.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Source package contains %d GRIB2 file(s)\n", len(gribs))

	newfile := "NewTxtFile.txt"

	outfile, err := os.Create(newfile)

	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	ArrPoints := make([][2]float64, 0)
	arrVelocities := make([][]float64, 0)
	tmp := []float64{}
	for _, g := range gribs {
		log.Printf("Published='%s', Forecast='%s', Parameter='%s', Unit='%s', Description='%s'\n",
			g.RefTime.Format("2006-01-02 15:04:05"), g.VerfTime.Format("2006-01-02 15:04:05"), g.Name, g.Unit, g.Description)

		//refTime := g.RefTime
		//verfTime := g.VerfTime

		//name := g.Name
		level := g.Level
		for _, v := range g.Values {
			lon := v.Longitude
			// if lon > 180.0 {
			// 	lon -= 360.0
			// }
			// if !Compare(v.Longitude, v.Latitude) {
			// 	continue
			// }

			_, err := fmt.Fprintf(outfile, "Date: %s, Parameter: %s, %s, Level: %s, Value: %g, Lat: %g, Lon: %g\n",
				//g.RefTime, //.Format("2006-01-02 15:04:05"),
				g.VerfTime.Format("2006-01-02 15:04:05"),
				//name,
				g.Description,
				g.Unit,
				level,
				v.Value,
				v.Latitude,
				lon,
			)

			if g.Name == "UGRD" {
				ArrPoints = append(ArrPoints, [2]float64{v.Longitude, v.Latitude})
			}
			if g.Name == "UGRD" {
				tmp = append(tmp, float64(v.Value))
			}
			if g.Name == "VGRD" {
				tmp = append(tmp, float64(v.Value))
			}
			if len(tmp) == 2 {
				arrVelocities = append(arrVelocities, tmp)
				tmp = []float64{}
			}

			if err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
	}
	fmt.Println(len(ArrPoints), len(arrVelocities))
	//CreateImage(ArrPoints, arrVelocities)

	//fmt.Println(arrVelocities)
	return arrVelocities, ArrPoints
}

func CreateImage(arr [][2]float64, arrOfVel [][]float32) error {
	// Открыть исходный файл
	f, err := os.Open("grid_map2.png")
	if err != nil {
		return err
	}
	defer f.Close()

	// Декодировать исходное изображение
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	// Создать копию изображения
	newImg := image.NewRGBA(img.Bounds())
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// Нарисовать точки на копии изображения
	for i := range arr {
		x1 := int(arr[i][0] * 3)
		y1 := int(math.Abs(arr[i][1]*3 - 270))
		//x2 := x1 + int(float64(arrOfVel[i][0]))
		//y2 := y1 - int((float64(arrOfVel[i][1])))
		newImg.Set(x1, y1, color.RGBA{255, 0, 0, 255}) // Рисуем красную точку
		//drawLine(newImg, x1, y1, x2, y2, color.RGBA{255, 0, 0, 255})
	}

	// Сохранить копию изображения в новый файл
	out, err := os.Create("results.png")
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, newImg)
}

func Compare(long, lat float64) bool {
	return long == 37 && lat == 68
}
