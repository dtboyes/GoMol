package main

import (
	"math"
)

func RenderScene(camera *Camera, light *Light, atoms []*Atom) {
	// Initialize pixel data
	pixels = make([]uint8, 4*imageWidth*imageHeight)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			// pixel_center = pixel00Location + pixel_delta_u * i + pixel_delta_v * j
			pixel_center := camera.pixel00.Add(camera.pixelDeltaU.Scale(float64(i))).Add(camera.pixelDeltaV.Scale(float64(j)))
			// ray_direction = pixel_center - camera.position
			ray_direction := pixel_center.Subtract(camera.position)
			// create a ray object
			ray := &Ray{camera.position, ray_direction, Color{0, 0, 0, 1}}
			// calculate the color of the ray
			pixel_color := RayColor(ray, light, atoms)

			color := colorToRGBA(pixel_color)
			pixels[4*(j*imageWidth+i)] = color[0]
			pixels[4*(j*imageWidth+i)+1] = color[1]
			pixels[4*(j*imageWidth+i)+2] = color[2]
			pixels[4*(j*imageWidth+i)+3] = color[3]

			// pixels[4*(j*imageWidth+i)] = uint8(pixel_color.x)
			// pixels[4*(j*imageWidth+i)+1] = uint8(pixel_color.y)
			// pixels[4*(j*imageWidth+i)+2] = uint8(pixel_color.z)
			// pixels[4*(j*imageWidth+i)+3] = 255
		}
	}
}
func RayColor(r *Ray, light *Light, atoms []*Atom) vec3 {
	for i := 0; i < len(atoms); i++ {
		t := RaySphereCollision(r, atoms[i])
		if t > 0.0 {
			surfaceNormal := r.getAt(t).Subtract(vec3{atoms[i].x, atoms[i].y, atoms[i].z}).UnitVector()
			color := LambertianShading(surfaceNormal, light, vec3{0.2, 0.8, 0.2})
			return color
			return vec3{66, 135, 245}
		}
	}
	return vec3{0, 0, 0}
}

func RaySphereCollision(r *Ray, atom *Atom) float64 {
	oc := r.getOrigin().Subtract(vec3{atom.x, atom.y, atom.z})
	a := r.getDirection().Dot(r.getDirection())
	b := 2.0 * oc.Dot(r.getDirection())
	c := oc.Dot(oc) - atom.radius*atom.radius
	discriminant := b*b - 4*a*c
	return discriminant
}

func LambertianShading(surfaceNormal vec3, light *Light, color vec3) vec3 {
	lightDirection := light.position.Subtract(surfaceNormal)
	cosTheta := math.Max(0, surfaceNormal.Dot(lightDirection))
	diffuse := color.Scale(cosTheta)
	ambient := color.Scale(0.2)
	color = ambient.Add(diffuse)
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
