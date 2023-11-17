package main

import (
	"math"
)

func RenderMultiProc(pixels []uint8, numProcs int) {
	// pixels = make([]uint8, 4*imageWidth*imageHeight)
	finished := make(chan bool, numProcs)
	for i := 0; i < numProcs; i++ {
		start_height := i * imageHeight / numProcs
		end_height := (i + 1) * imageHeight / numProcs
		go RenderScene(camera, light, atoms, start_height, end_height, pixels, finished)
	}
	for i := 0; i < numProcs; i++ {
		<-finished
	}
}

func RenderScene(camera *Camera, light *Light, atoms []*Atom, start, end int, pixels []uint8, finished chan bool) {
	// Initialize pixel data
	light.position = light.position.Add(CenterOfMass(atoms))
	for j := start; j < end; j++ {
		for i := 0; i < imageWidth; i++ {
			// pixel_center = pixel00Location + pixel_delta_u * i + pixel_delta_v * j
			pixel_center := camera.pixel00.Add(camera.pixelDeltaU.Scale(float64(i))).Add(camera.pixelDeltaV.Scale(float64(j)))
			// ray_direction = pixel_center - camera.position
			ray_direction := pixel_center.Subtract(camera.position)
			// create a ray object
			ray := &Ray{camera.position, ray_direction, Color{0, 0, 0, 1}}
			// calculate the color of the ray
			pixel_color := RayColor(ray, light, camera, atoms)

			color := colorToRGBA(pixel_color)
			pixels[4*(j*imageWidth+i)] = color[0]
			pixels[4*(j*imageWidth+i)+1] = color[1]
			pixels[4*(j*imageWidth+i)+2] = color[2]
			pixels[4*(j*imageWidth+i)+3] = color[3]
		}
	}
	finished <- true
}
func RayColor(r *Ray, light *Light, camera *Camera, atoms []*Atom) vec3 {
	for i := 0; i < len(atoms); i++ {
		collision := RaySphereCollision(r, atoms[i])
		if !collision.normal.EqualsZero() {
			if atoms[i].chain == "A" {
				collision.color = LambertianShading(collision, light, camera, vec3{0.2, 1.0, 0.1})
			} else if atoms[i].chain == "B" {
				collision.color = LambertianShading(collision, light, camera, vec3{0.1, 0.2, 1.0})
			} else if atoms[i].chain == "C" {
				collision.color = LambertianShading(collision, light, camera, vec3{1.0, 0.1, 0.2})
			} else if atoms[i].chain == "D" {
				collision.color = LambertianShading(collision, light, camera, vec3{1.0, 0.55, 0.0})
			} else {
				collision.color = LambertianShading(collision, light, camera, vec3{1.0, 1.0, 1.0})
			}
			return collision.color
		}
	}
	return vec3{0.0, 0.0, 0.0}
}

func RaySphereCollision(r *Ray, atom *Atom) Collision {
	var collision Collision
	oc := r.getOrigin().Subtract(vec3{atom.x, atom.y, atom.z})
	a := r.getDirection().Dot(r.getDirection())
	b := 2.0 * oc.Dot(r.getDirection())
	c := oc.Dot(oc) - atom.radius*atom.radius
	discriminant := b*b - 4*a*c
	var min_t float64
	if discriminant < 0.0 {
		zero := vec3{0, 0, 0}
		collision.point = zero
		collision.normal = zero
		return collision
	} else {
		var tval1 float64 = (-b - math.Sqrt(discriminant)) / (2.0 * a)
		var tval2 float64 = (-b + math.Sqrt(discriminant)) / (2.0 * a)
		if tval2 < tval1 {
			min_t = tval2
		} else {
			min_t = tval1
		}
	}
	origin := r.getOrigin()
	direction := r.getDirection()
	collision.point = direction.Scale(min_t).Add(origin)
	collision.normal = collision.point.Subtract(vec3{atom.x, atom.y, atom.z}).Normalize()
	return collision
}

func LambertianShading(collision Collision, light *Light, camera *Camera, color vec3) vec3 {
	constantAttenuation := 0.02
	linearAttenuation := 0.005
	quadraticAttenuation := 0.0005
	lightIntensity := 2.0
	dist := light.getPosition().Subtract(collision.point).Length() * 0.1
	totalAttenuation := 1.0 / (constantAttenuation + linearAttenuation*dist + quadraticAttenuation*dist*dist)

	// // lightDirection = unit vector of (light.position - collision.point)
	lightDirection := light.getPosition().Subtract(collision.point).Normalize()

	// diffuse = color * max(0, collision.normal dot lightDirection) * totalAttenuation * lightIntensity
	diffuse := color.Scale(math.Max(0.0, collision.normal.Dot(lightDirection))).Scale(totalAttenuation * lightIntensity)
	ambient := color.Scale(0.5)
	color = diffuse.Add(ambient)
	return color
}

func InitializeCamera(atoms []*Atom) *Camera {
	camera := ParseCamera("input/camera.txt")
	// makes it so that the camera always points at the center of mass of all atoms
	camera.position = camera.position.Add(CenterOfMass(atoms))

	// viewportWidth is viewportHeight * aspectRatio
	camera.viewportWidth = camera.viewportHeight * float64(imageWidth) / float64(imageHeight)
	viewport_u := vec3{camera.viewportWidth, 0, 0}
	viewport_v := vec3{0, -camera.viewportHeight, 0}

	// Initializing viewport, pixel delta, and top left pixel location

	// pixel_delta_u = viewport_u / imageWidth
	pixel_delta_u := viewport_u.Scale(1.0 / float64(imageWidth))
	// pixel_delta_u = viewport_v / imageHeight
	pixel_delta_v := viewport_v.Scale(1.0 / float64(imageHeight))

	// uppper left of viewport is the camera position minus half of the viewport width and height minus the focal Length
	viewport_upper_left := camera.position.Subtract(viewport_u.Scale(0.5)).Subtract(viewport_v.Scale(0.5)).Subtract(vec3{0, 0, camera.focalLength})

	// top left pixel location is the upper left viewport location plus half of the pixel width and height
	pixel00Location := viewport_upper_left.Add(pixel_delta_u.Scale(0.5).Add(pixel_delta_v.Scale(0.5)))

	camera.pixelDeltaU = pixel_delta_u
	camera.pixelDeltaV = pixel_delta_v
	camera.pixel00 = pixel00Location

	return camera
}

func colorToRGBA(c vec3) [4]uint8 {
	return [4]uint8{
		uint8(c.x * 255),
		uint8(c.y * 255),
		uint8(c.z * 255),
		255,
	}
}
