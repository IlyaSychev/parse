package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/amsokol/go-grib2"
)

const URL = "https://nomads.ncep.noaa.gov/cgi-bin/filter_fnl.pl?dir=%2Fgdas.20240907%2F18%2Fatmos&file=gdas.t18z.pgrb2.1p00.anl&var_ABSV=on&var_DZDT=on&var_HGT=on&var_O3MR=on&var_RH=on&var_SPFH=on&var_TMP=on&var_UGRD=on&var_VGRD=on&var_VVEL=on&lev_70_mb=on&lev_50_mb=on&lev_40_mb=on&lev_30_mb=on&lev_20_mb=on&lev_15_mb=on&lev_10_mb=on&lev_7_mb=on&lev_5_mb=on&lev_3_mb=on&lev_2_mb=on&lev_1_mb=on&subregion=&toplat=90&leftlon=20&rightlon=200.00&bottomlat=50"

func main() {
	log.Println("Downloading data from HTTP server...")
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data downloaded")

	newfile := "./gfs.t00z.pgrb2.0p25.f001.txt"

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	gribs, err := grib2.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Source package contains %d GRIB2 file(s)\n", len(gribs))

	outfile, err := os.Create(newfile)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	ArrPoints := make([][2]float64, 0)

	for _, g := range gribs {
		log.Printf("Published='%s', Forecast='%s', Parameter='%s', Unit='%s', Description='%s'\n",
			g.RefTime.Format("2006-01-02 15:04:05"), g.VerfTime.Format("2006-01-02 15:04:05"), g.Name, g.Unit, g.Description)

		//refTime := g.RefTime
		//verfTime := g.VerfTime
		name := g.Name
		level := g.Level
		for _, v := range g.Values {
			lon := v.Longitude
			if lon > 180.0 {
				lon -= 360.0
			}
			// if !Compare(v.Longitude, v.Latitude) {
			// 	continue
			// }
			_, err := fmt.Fprintf(outfile, "\"%s\",\"%s\",\"%s\", Parametr %s, Means %s, Unit %s, Долгота %g, Широта %g, %g\n",
				g.RefTime, //.Format("2006-01-02 15:04:05"),
				//g.VerfTime, //.Format("2006-01-02 15:04:05"),
				name,
				level,
				g.Name,
				g.Description,
				g.Unit,
				lon,
				v.Latitude,
				v.Value)
			ArrPoints = append(ArrPoints, [2]float64{v.Longitude, v.Latitude})
			if err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}
	}
	CreateImage(ArrPoints)
}

func CreateImage(arr [][2]float64) error {
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
		newImg.Set(int(arr[i][0]*3), int(math.Abs(arr[i][1]*3-270)), color.RGBA{255, 0, 0, 255}) // Рисуем красную точку
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
