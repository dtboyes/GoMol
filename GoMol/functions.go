package main

func InitializeCamera() *Camera {
	camera := ParseCamera("input/camera.txt")
	camera.viewportWidth = camera.viewportHeight * float64(imageWidth) / float64(imageHeight)
	return camera
}

func RenderScene(camera *Camera, atoms []*Atom, pixel00Location, pixel_delta_u, pixel_delta_v vec3, viewport_upper_left vec3) {
	// Initialize pixel data
	pixels = make([]uint8, 4*imageWidth*imageHeight)
	for j := 0; j < imageHeight; j++ {
		for i := 0; i < imageWidth; i++ {
			// pixel_center = pixel00Location + pixel_delta_u * i + pixel_delta_v * j
			pixel_center := pixel00Location.vectorAddition(pixel_delta_u.multiplyScalar(float64(i))).vectorAddition(pixel_delta_v.multiplyScalar(float64(j)))
			// ray_direction = pixel_center - camera.position
			ray_direction := pixel_center.vectorSubtraction(camera.position)
			// create a ray object
			ray := &Ray{camera.position, ray_direction, vec3{0, 0, 0}, Color{0, 0, 0, 1}}
			// calculate the color of the ray
			pixel_color := RayColor(ray, atoms)
			atoms[0].x = 0
			atoms[0].y = 0
			atoms[0].z = -1
			pixels[4*(j*imageWidth+i)] = uint8(pixel_color.x)
			pixels[4*(j*imageWidth+i)+1] = uint8(pixel_color.y)
			pixels[4*(j*imageWidth+i)+2] = uint8(pixel_color.z)
			pixels[4*(j*imageWidth+i)+3] = 1
		}
	}
}
func RayColor(r *Ray, atoms []*Atom) vec3 {
	for i := 0; i < len(atoms); i++ {
		if RaySphereCollision(r, atoms[i]) {
			return vec3{59, 118, 212}
		}
	}
	return vec3{0, 0, 0}
}

func RaySphereCollision(r *Ray, atom *Atom) bool {
	oc := r.getOrigin().vectorSubtraction(vec3{atom.x, atom.y, atom.z})
	a := r.getDirection().dotProduct(r.getDirection())
	b := 2.0 * oc.dotProduct(r.getDirection())
	c := oc.dotProduct(oc) - atom.radius*atom.radius
	discriminant := b*b - 4*a*c
	return discriminant > 0
}
