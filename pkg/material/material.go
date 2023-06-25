package material

import (
	"gotracing/pkg/geometry"
	"gotracing/pkg/scene"
)

type Color struct {
	R, G, B float64
}

type Material struct {
	Color              Color
	Reflectivity       float64
	Transparency       float64
	RefractiveIndex    float64
	AmbientCoefficient float64
	DiffuseCoefficient float64
	SpecularCoefficient float64
	Shininess          float64
}

func (m *Material) ComputeColor(hit *geometry.Hit, s *scene.Scene, ray *geometry.Ray) Color {
	color := m.Color.Scale(m.AmbientCoefficient)

	for _, light := range s.Lights {
		lightDirection := light.Position.Subtract(hit.Position)
		distance := lightDirection.Length()
		lightDirection = lightDirection.Normalize()

		attenuation := 1 / (1 + 0.1*distance)

		diffuse := m.Color.Scale(m.DiffuseCoefficient * max(0, hit.Normal.Dot(lightDirection)))

		reflectDirection := lightDirection.Negate().Reflect(hit.Normal)
		viewDirection := ray.Direction.Negate()
		specular := light.Color.Scale(m.SpecularCoefficient * math.Pow(max(0, viewDirection.Dot(reflectDirection)), m.Shininess))

		lightContribution := diffuse.Add(specular).Scale(attenuation)
		color = color.Add(lightContribution)
	}

	color.R = min(1, color.R)
	color.G = min(1, color.G)
	color.B = min(1, color.B)

	return color
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
