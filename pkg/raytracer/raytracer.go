package raytracer

import (
	"fmt"
	"image"
	"sync"
	"gotracing/pkg/geometry"
	"gotracing/pkg/material"
	"gotracing/pkg/scene"
	"gotracing/pkg/utility"
)

type Raytracer struct {
	Scene    *scene.Scene
	Width    int
	Height   int
	MaxDepth int
	Samples  int
}

func NewRaytracer(scene *scene.Scene, width, height, maxDepth, samples int) *Raytracer {
	return &Raytracer{
		Scene:    scene,
		Width:    width,
		Height:   height,
		MaxDepth: maxDepth,
		Samples:  samples,
	}
}

func (r *Raytracer) Render() (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))

	var wg sync.WaitGroup
	pixels := r.Width * r.Height
	wg.Add(pixels)

	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			go func(x, y int) {
				defer wg.Done()

				color := &material.Color{}

				for s := 0; s < r.Samples; s++ {
					u := (float64(x) + utility.Random()) / float64(r.Width-1)
					v := (float64(y) + utility.Random()) / float64(r.Height-1)

					ray, err := r.Scene.Camera.GetRay(u, v)
					if err != nil {
						// handle error
						return
					}

					sampleColor := r.traceRay(ray, r.MaxDepth)

					color = color.Add(sampleColor)
				}

				color = color.Scale(1 / float64(r.Samples))

				rgba := color.ToRGBA()
				img.Set(x, y, rgba)
			}(x, y)
		}
		fmt.Printf("\rRendering... %d%% complete", 100*(r.Width*y+x)/pixels)
	}

	wg.Wait()

	fmt.Println("\rRendering... done")

	return img, nil
}


func (r *Raytracer) traceRay(ray *geometry.Ray, depth int, inside bool) *material.Color {
	if depth <= 0 {
		return &material.Color{0, 0, 0} 
	}

	hit, object := r.Scene.FindClosestIntersection(ray)
	if hit == nil {
		return &material.Color{0, 0, 0} 
	}

	hitColor := object.Material.ComputeColor(hit, r.Scene, ray)

	if object.Material.Reflectivity > 0 {
		reflectionRay := ray.Reflect(hit.Normal)
		reflectedColor := r.traceRay(reflectionRay, depth-1, inside)
		hitColor = hitColor.Add(reflectedColor.Scale(object.Material.Reflectivity))
	}

	if object.Material.Transparency > 0 {
		refractionRay, totalInternalReflection := ray.Refract(hit.Normal, object.Material.RefractiveIndex, inside)
		if totalInternalReflection {
			reflectedColor := r.traceRay(refractionRay, depth-1, inside)
			hitColor = hitColor.Add(reflectedColor.Scale(object.Material.Transparency))
		} else {
			refractedColor := r.traceRay(refractionRay, depth-1, !inside)
			hitColor = hitColor.Add(refractedColor.Scale(object.Material.Transparency))
		}
	}

	return hitColor
}


