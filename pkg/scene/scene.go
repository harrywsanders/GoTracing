package scene

import (
	"math"

	"gotracing/pkg/geometry"
	"gotracing/pkg/material"
)

type Scene struct {
	Objects []geometry.Object
	Lights  []material.Light
	Camera  *geometry.Camera
}

func NewScene(camera *geometry.Camera) *Scene {
	return &Scene{
		Objects: make([]geometry.Object, 0),
		Lights:  make([]material.Light, 0),
		Camera:  camera,
	}
}

func (s *Scene) AddObject(object geometry.Object) {
	s.Objects = append(s.Objects, object)
}

func (s *Scene) AddLight(light material.Light) {
	s.Lights = append(s.Lights, light)
}

func (s *Scene) FindClosestIntersection(ray *geometry.Ray) (*geometry.Hit, *geometry.Object) {
	var closestHit *geometry.Hit
	var closestObject *geometry.Object
	closestDistance := math.Inf(1)

	for _, object := range s.Objects {
		hit := object.Intersect(ray)
		if hit != nil && hit.T < closestDistance {
			closestHit = hit
			closestObject = &object
			closestDistance = hit.T
		}
	}

	return closestHit, closestObject
}
