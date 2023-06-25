package geometry

import (
	"math"
	"gotracing/pkg/material"
)

// Vector represents a 3D vector or point.
type Vector struct {
	X, Y, Z float64
}

// Methods for Vector

func (v *Vector) Add(w *Vector) *Vector {
	return &Vector{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v *Vector) Subtract(w *Vector) *Vector {
	return &Vector{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v *Vector) Scale(s float64) *Vector {
	return &Vector{v.X * s, v.Y * s, v.Z * s}
}

func (v *Vector) Dot(w *Vector) float64 {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v *Vector) Cross(w *Vector) *Vector {
	return &Vector{
		v.Y*w.Z - v.Z*w.Y,
		v.Z*w.X - v.X*w.Z,
		v.X*w.Y - v.Y*w.X,
	}
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector) Normalize() *Vector {
	len := v.Length()
	return &Vector{v.X / len, v.Y / len, v.Z / len}
}

// Ray represents a ray.
type Ray struct {
	Origin, Direction *Vector
}

// Methods for Ray

func (r *Ray) At(t float64) *Vector {
	return r.Origin.Add(r.Direction.Scale(t))
}

// Object represents a generic object in the scene.
type Object interface {
	Intersect(ray *Ray) *Hit
	Material() *material.Material
}

// Hit represents a ray-object intersection.
type Hit struct {
	Position, Normal *Vector
	T float64
	Object Object
}

// Sphere represents a sphere.
type Sphere struct {
	Center *Vector
	Radius float64
	Mat    *material.Material
}

// Intersection method for Sphere

func (s *Sphere) Intersect(ray *Ray) *Hit {
	oc := ray.Origin.Subtract(s.Center)
	a := ray.Direction.Dot(ray.Direction)
	b := 2 * oc.Dot(ray.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return nil // no intersection
	}

	t := (-b - math.Sqrt(discriminant)) / (2 * a)
	position := ray.At(t)
	normal := position.Subtract(s.Center).Normalize()

	return &Hit{position, normal, t, s}
}

func (s *Sphere) Material() *material.Material {
	return s.Mat
}

// Plane represents a plane.
type Plane struct {
	Point  *Vector
	Normal *Vector
	Mat    *material.Material
}

// Intersection method for Plane
func (p *Plane) Intersect(ray *Ray) *Hit {
	denom := p.Normal.Dot(ray.Direction)
	if math.Abs(denom) > 0.0001 { // prevent division by zero
		t := p.Point.Subtract(ray.Origin).Dot(p.Normal) / denom
		if t >= 0 {
			position := ray.At(t)
			return &Hit{position, p.Normal, t, p}
		}
	}
	return nil // no intersection
}

func (p *Plane) Material() *material.Material {
	return p.Mat
}

// Triangle represents a triangle.
type Triangle struct {
	V0, V1, V2 *Vector // vertices
	Mat        *material.Material
}

// Intersection method for Triangle using the Moller-Trumbore algorithm
func (t *Triangle) Intersect(ray *Ray) *Hit {
	edge1 := t.V1.Subtract(t.V0)
	edge2 := t.V2.Subtract(t.V0)
	h := ray.Direction.Cross(edge2)
	a := edge1.Dot(h)

	if math.Abs(a) < 0.0001 { // this ray is parallel to this triangle
		return nil
	}

	f := 1.0 / a
	s := ray.Origin.Subtract(t.V0)
	u := f * s.Dot(h)

	if u < 0.0 || u > 1.0 { // the intersection lies outside of the triangle
		return nil
	}

	q := s.Cross(edge1)
	v := f * ray.Direction.Dot(q)

	if v < 0.0 || u+v > 1.0 { // the intersection lies outside of the triangle
		return nil
	}

	// at this stage we can compute t to find out where the intersection point is on the line
	tt := f * edge2.Dot(q)

	if tt > 0.0001 { // ray intersection
		position := ray.At(tt)
		normal := edge1.Cross(edge2).Normalize()
		return &Hit{position, normal, tt, t}
	}

	return nil // this means that there is a line intersection but not a ray intersection
}

func (t *Triangle) Material() *material.Material {
	return t.Mat
}

// Cylinder represents an infinite cylinder along the y-axis.
type Cylinder struct {
	Center *Vector
	Radius float64
	Mat    *material.Material
}

// Intersection method for Cylinder
func (c *Cylinder) Intersect(ray *Ray) *Hit {
	oc := ray.Origin.Subtract(c.Center)
	oc.Y = 0 // ignore y component
	a := ray.Direction.X*ray.Direction.X + ray.Direction.Z*ray.Direction.Z
	b := 2 * oc.Dot(ray.Direction)
	cval := oc.Dot(oc) - c.Radius*c.Radius
	discriminant := b*b - 4*a*cval

	if discriminant < 0 {
		return nil // no intersection
	}

	t := (-b - math.Sqrt(discriminant)) / (2 * a)
	position := ray.At(t)
	normal := position.Subtract(c.Center)
	normal.Y = 0 // ignore y component
	normal = normal.Normalize()

	return &Hit{position, normal, t, c}
}

func (c *Cylinder) Material() *material.Material {
	return c.Mat
}

// Cube represents a cube.
type Cube struct {
	Center *Vector
	Length float64
	Mat    *material.Material
}

// Intersection method for Cube
func (cube *Cube) Intersect(ray *Ray) *Hit {
	min := cube.Center.Subtract(&Vector{cube.Length / 2, cube.Length / 2, cube.Length / 2})
	max := cube.Center.Add(&Vector{cube.Length / 2, cube.Length / 2, cube.Length / 2})

	tmin := (min.X - ray.Origin.X) / ray.Direction.X
	tmax := (max.X - ray.Origin.X) / ray.Direction.X

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	tymin := (min.Y - ray.Origin.Y) / ray.Direction.Y
	tymax := (max.Y - ray.Origin.Y) / ray.Direction.Y

	if tymin > tymax {
		tymin, tymax = tymax, tymin
	}

	if (tmin > tymax) || (tymin > tmax) {
		return nil
	}

	if tymin > tmin {
		tmin = tymin
	}

	if tymax < tmax {
		tmax = tymax
	}

	tzmin := (min.Z - ray.Origin.Z) / ray.Direction.Z
	tzmax := (max.Z - ray.Origin.Z) / ray.Direction.Z

	if tzmin > tzmax {
		tzmin, tzmax = tzmax, tzmin
	}

	if (tmin > tzmax) || (tzmin > tmax) {
		return nil
	}

	if tzmin > tmin {
		tmin = tzmin
	}

	if tzmax < tmax {
		tmax = tzmax
	}

	t := tmin

	if t < 0 {
		t = tmax
		if t < 0 {
			return nil
		}
	}
		position := ray.At(t)
		normal := cube.NormalAt(position)
	
		return &Hit{position, normal, t, cube}
	}
	
func (cube *Cube) NormalAt(point *Vector) *Vector {
		// Determine the closest face of the cube that the point is on, and return that face's normal.
		min := cube.Center.Subtract(&Vector{cube.Length / 2, cube.Length / 2, cube.Length / 2})
		max := cube.Center.Add(&Vector{cube.Length / 2, cube.Length / 2, cube.Length / 2})
	
		normals := []*Vector{
			&Vector{-1, 0, 0},
			&Vector{1, 0, 0},
			&Vector{0, -1, 0},
			&Vector{0, 1, 0},
			&Vector{0, 0, -1},
			&Vector{0, 0, 1},
		}
	
		sides := []float64{
			math.Abs(min.X - point.X),
			math.Abs(max.X - point.X),
			math.Abs(min.Y - point.Y),
			math.Abs(max.Y - point.Y),
			math.Abs(min.Z - point.Z),
			math.Abs(max.Z - point.Z),
		}
	
		minIndex := 0
		minValue := sides[0]
		for i := 1; i < len(sides); i++ {
			if sides[i] < minValue {
				minValue = sides[i]
				minIndex = i
			}
		}
	
		return normals[minIndex]
	}
	
	func (cube *Cube) Material() *material.Material {
		return cube.Mat
	}
	
