package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func makePlot(path string, data plotter.XYs) {

	// xticks defines how we convert and display time.Time values.
	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p := plot.New()
	p.Title.Text = "Weekly Closing Price (LINKUSD)"
	p.X.Tick.Marker = xticks
	p.Y.Label.Text = "LINK Price ($)"
	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = color.RGBA{G: 255, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, points)

	f, err := os.Create(path)
	if err != nil {
		log.Panicf("could not create %s: %v", path, err)
	}
	err = p.Save(10*vg.Centimeter, 5*vg.Centimeter, path)
	if err != nil {
		log.Panic(err)
	}
	if err = f.Close(); err != nil {
		log.Panicf("could not close %s: %v", path, err)

	}
}

func readData(path string) (plotter.XYs, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var xys plotter.XYs
	i := -1
	s := bufio.NewScanner(f)
	for s.Scan() {
		// skip header row
		i++
		if i == 0 {
			continue
		}

		var x, y float64

		// time_period_start;time_period_end;time_open;time_close;price_close
		values := strings.Split(s.Text(), ";")

		time_period_end := values[1]
		price_close := values[4]

		var xtime time.Time
		xtime, err = time.Parse(time.RFC3339Nano, time_period_end)
		if err != nil {
			log.Printf("error parsing time %q: %v", time_period_end, err)
			continue
		}
		x = float64(xtime.Unix())
		y, err = strconv.ParseFloat(price_close, 64)
		if err != nil {
			log.Printf("error parsing closing price %q: %v", price_close, err)
			continue
		}
		xys = append(xys, struct{ X, Y float64 }{x, y})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("could not scan: %v", err)
	}
	return xys, nil
}

func main() {
	pxys, err := readData("7DAY-link-data.csv")

	if err != nil {
		log.Fatalf("could not read data.txt: %v", err)
	}
	path := "timeseries.png"
	makePlot(path, pxys)
}
