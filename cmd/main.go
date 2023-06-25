package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"gotracing/pkg/raytracer"
	"gotracing/pkg/scene"
)

func main() {
	sceneFile := flag.String("scene", "", "Path to the scene file")
	outputFile := flag.String("output", "output.png", "Path to the output image file")

	flag.Parse()

	if *sceneFile == "" {
		fmt.Fprintln(os.Stderr, "You must specify a scene file with the -scene flag.")
		os.Exit(1)
	}

	s, err := scene.LoadScene(*sceneFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load scene file: %v\n", err)
		os.Exit(1)
	}

	width := 800
	height := 600
	maxDepth := 5
	samples := 100

	r := raytracer.NewRaytracer(s, width, height, maxDepth, samples)

	img, err := r.Render()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Rendering failed: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write image: %v\n", err)
		os.Exit(1)
	}
}
