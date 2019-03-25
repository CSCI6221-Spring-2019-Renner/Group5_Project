package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
)


func main() {
	keywords := []string{}
	frequencyForWord.Init(keywords)

	keys := []string{}
	for _, key := range keywords {
	keys = append(keys, key)
	}
	fmt.Print(keys)

	for i :=0; i<len(keys);i++{
	K := plotter.Values{}
	for value, key := range keywords {
		K = append(K, float64(value))
		keys = append(keys, key)
	}

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	w := vg.Points(8)


	barsA, err := plotter.NewBarChart(K, w)
	if err != nil {
		log.Panic(err)
	}
	barsA.Color = color.RGBA{R: 187, A: 5}
	barsA.Offset = -w / 10


	p, err = plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "The Discussion Frequency of Common Programming Languages"
	p.Y.Label.Text = "Frequency"

	w = vg.Points(20)
	p.Add(barsA)

		p.NominalX(keys[i])


	p.Add(plotter.NewGlyphBoxes())
	err = p.Save(750, 750, "Crawler analysis.png")
	if err != nil {
	log.Panic(err)
	}
}}