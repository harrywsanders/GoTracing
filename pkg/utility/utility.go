package utility

import (
	"math"
	"math/rand"

	"gotracing/pkg/geometry"
	"gotracing/pkg/material"
)

// Random returns a random float64 between 0 and 1.
func Random() float64 {
	return rand.Float64()
}

// RandomInRange returns a random float64 between min and max.
func RandomInRange(min, max float64) float64 {
	return min + (max-min)*Random()
}

// Clamp clamps a value between a minimum and maximum.
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Lerp performs linear interpolation between a and b, based on t.
func Lerp(a, b, t float64) float64 {
	return (1-t)*a + t*b
}

// Reflect returns the reflection of a vector v against a normal n.
func Reflect(v, n *geometry.Vector) *geometry.Vector {
	return v.Subtract(n.Scale(2 * v.Dot(n)))
}

// Refract returns the refraction of a vector v through a normal n, using the given refractive index.
// Also returns whether total internal reflection occurred.
func Refract(v, n *geometry.Vector, refractiveIndex float64) (*geometry.Vector, bool) {
	dt := v.Dot(n)
	discriminant := 1.0 - refractiveIndex*refractiveIndex*(1-dt*dt)
	if discriminant > 0 {
		refracted := v.Scale(refractiveIndex).Subtract(n.Scale(refractiveIndex*dt + math.Sqrt(discriminant)))
		return refracted, false
	}
	return nil, true
}

// ConvertColorToUint8 converts a material.Color to RGBA values in the uint8 format.
func ConvertColorToUint8(color *material.Color) (r, g, b, a uint8) {
	r = uint8(255.999 * Clamp(color.R, 0.0, 1.0))
	g = uint8(255.999 * Clamp(color.G, 0.0, 1.0))
	b = uint8(255.999 * Clamp(color.B, 0.0, 1.0))
	a = 255 // fully opaque
	return
}
