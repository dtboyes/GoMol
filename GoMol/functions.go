package main

import (
	"fmt"
	"math"
)

// parallelizing the rendering process by splitting up viewport into multiple sections
func RenderMultiProc(pixels []uint8, numProcs int, window1 bool) {
	// pixels = make([]uint8, 4*imageWidth*imageHeight)
	finished := make(chan bool, numProcs)
	for i := 0; i < numProcs; i++ {
		start_height := i * imageHeight / numProcs
		end_height := (i + 1) * imageHeight / numProcs
		if renderProtein1 || renderKabsch {
			go RenderScene(camera, light, atoms1, atoms1_sequence, alignedSeq2, start_height, end_height, pixels, finished)
		} else if renderProtein2 {
			go RenderScene(camera, light, atoms2, atoms2_sequence, alignedSeq1, start_height, end_height, pixels, finished)
		}
	}
	for i := 0; i < numProcs; i++ {
		<-finished
	}
}

// RenderScene renders the scene by iterating through each pixel in the viewport
// for each pixel, a ray is generated with the origin at the camera position and the direction
// then computes a color for that ray
func RenderScene(camera *Camera, light *Light, atoms []*Atom, atoms_sequence, aligned_sequence string, start, end int, pixels []uint8, finished chan bool) {
	for j := start; j < end; j++ {
		for i := 0; i < imageWidth; i++ {
			// pixel_center = pixel00Location + pixel_delta_u * i + pixel_delta_v * j
			pixel_center := camera.pixel00.Add(camera.pixelDeltaU.Scale(float64(i))).Add(camera.pixelDeltaV.Scale(float64(j)))
			// ray_direction = pixel_center - camera.position
			ray_direction := pixel_center.Subtract(camera.position)
			// create a ray object
			ray := &Ray{camera.position, ray_direction, Color{0, 0, 0, 1}}
			// calculate the color of the ray
			pixel_color := RayColor(ray, light, camera, atoms, atoms_sequence, aligned_sequence)

			color := colorToRGBA(pixel_color)
			pixels[4*(j*imageWidth+i)] = color[0]
			pixels[4*(j*imageWidth+i)+1] = color[1]
			pixels[4*(j*imageWidth+i)+2] = color[2]
			pixels[4*(j*imageWidth+i)+3] = color[3]
		}
	}
	finished <- true
}

// RayColor calculates the color of a ray based on the atoms in the scene
// for a particular ray, the function checks to see if it collides with any atoms in the scene
// it then computes the color of the ray based on the Phong shading model
func RayColor(r *Ray, light *Light, camera *Camera, atoms []*Atom, atoms_sequence, aligned_sequence string) vec3 {
	for i := 0; i < len(atoms); i++ {
		collision := RaySphereCollision(r, atoms[i])
		if !collision.getNormal().EqualsZero() {
			if colorByChain {
				if atoms[i].chain == "A" {
					collision.color = PhongShading(collision, light, camera, vec3{0.2, 0.7, 0.1})
				} else if atoms[i].chain == "B" {
					collision.color = PhongShading(collision, light, camera, vec3{0.1, 0.2, 1.0})
				} else if atoms[i].chain == "C" {
					collision.color = PhongShading(collision, light, camera, vec3{1.0, 0.1, 0.2})
				} else if atoms[i].chain == "D" {
					collision.color = PhongShading(collision, light, camera, vec3{1.0, 0.55, 0.0})
				} else {
					collision.color = PhongShading(collision, light, camera, vec3{1.0, 1.0, 1.0})
				}
			} else if colorByAtom {
				if atoms[i].element == "CA" {
					collision.color = PhongShading(collision, light, camera, vec3{0.565, 0.565, 0.565})
				} else if atoms[i].element == "N" {
					collision.color = PhongShading(collision, light, camera, vec3{0.188, 0.313, 0.9725})
				} else if atoms[i].element == "O" {
					collision.color = PhongShading(collision, light, camera, vec3{1.0, 0.051, 0.051})
				} else if atoms[i].element == "S" {
					collision.color = PhongShading(collision, light, camera, vec3{1.0, 0.784, 0.196})
				}
			} else if colorByDifferingRegions {
				if alignedSeq1[atoms[i].seqIndex] != alignedSeq2[atoms[i].seqIndex] {
					collision.color = PhongShading(collision, light, camera, vec3{0.69, 0.22, 0.188})
				} else {
					collision.color = PhongShading(collision, light, camera, vec3{0.373, 0.651, 0.286})
				}
			} else if renderKabsch {
				if i < len(atoms1_sequence) {
					collision.color = PhongShading(collision, light, camera, vec3{0.373, 0.651, 0.286})
				} else {
					collision.color = PhongShading(collision, light, camera, vec3{1.0, 0.22, 1.0})
				}
			} else if renderProtein1 {
				collision.color = PhongShading(collision, light, camera, vec3{0.373, 0.651, 0.286})
			} else if renderProtein2 {
				collision.color = PhongShading(collision, light, camera, vec3{1.0, 0.22, 1.0})
			} else {
				collision.color = PhongShading(collision, light, camera, vec3{0.373, 0.651, 0.286})
			}
			return collision.color
		}

	}
	return vec3{0, 0, 0}
}

// RaySphereCollision calculates the point of collision between a ray and a sphere
// it returns a collision object which contains the point of collision and the normal vector
// collision computations done using well defined equations for ray-sphere intersection
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

// PhongShading calculates the color of a pixel based on the Phong shading model
// it takes specific attenuation constants, light parameters, material properites, etc.
// uses all of these properties to calculate the color of a ray based on ambient, diffuse, and specular properties
// specified by the Phong Shading model
func PhongShading(collision Collision, light *Light, camera *Camera, color vec3) vec3 {
	constantAttenuation := 0.1
	linearAttenuation := 0.03
	quadraticAttenuation := 0.0001
	lightIntensity := 1.0
	specularColor := vec3{1.0, 1.0, 1.0}
	dist := light.getPosition().Subtract(collision.getPoint()).Length()
	totalAttenuation := 1.0 / (constantAttenuation + linearAttenuation*dist + quadraticAttenuation*dist*dist)
	// // lightDirection = unit vector of (light.position - collision.point)
	lightDirection := light.getPosition().Subtract(collision.point).Normalize()

	cameraDirection := camera.getPosition().Subtract(collision.point).Normalize()
	reflectDirection := lightDirection.Subtract(collision.getNormal().Scale(2.0 * lightDirection.Dot(collision.getNormal()))).Normalize()
	// diffuse = color * max(0, collision.normal dot lightDirection) * totalAttenuation * lightIntensity
	diffuse := color.Scale(math.Max(0.0, collision.getNormal().Dot(lightDirection))).Scale(totalAttenuation * lightIntensity)
	ambient := color.Scale(0.6)
	specular := specularColor.Scale(math.Pow(math.Max(0.0, reflectDirection.Dot(cameraDirection)), 5.0)).Scale(totalAttenuation * lightIntensity)
	color = diffuse.Add(ambient).Add(specular)
	return color
}

// InitializeCamera initializes the camera position and viewport
// the camera is set to point at the center of mass of all atoms through the viewport
// each pixel location in the viewport is set up to map to a physical location in space
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

// InitializeLight initializes the light position to be some set distance away from the center of mass
func InitializeLight(atoms []*Atom) *Light {
	light.position = light.position.Add(CenterOfMass(atoms))
	return light
}

// RotateAtoms rotates a slice of atoms around the x and y axes based on cursor position
// after clicking and dragging using the mouse
func RotateAtoms(atoms []*Atom, rotationX, rotationY float64) []*Atom {
	for i := 0; i < len(atoms); i++ {
		// rotate around x axis
		// y' = y*cos q - z*sin q
		// z' = y*sin q + z*cos q
		y := atoms[i].y
		z := atoms[i].z
		atoms[i].y = y*math.Cos(rotationX) - z*math.Sin(rotationX)
		atoms[i].z = y*math.Sin(rotationX) + z*math.Cos(rotationX)

		// rotate around y axis
		// x' = x*cos q - z*sin q
		// z' = x*sin q + z*cos q
		x := atoms[i].x
		z = atoms[i].z
		atoms[i].x = x*math.Cos(rotationY) - z*math.Sin(rotationY)
		atoms[i].z = x*math.Sin(rotationY) + z*math.Cos(rotationY)
	}
	return atoms
}

// converts a vec3 of float values between 0 and 1 to a vec3 of uint8 values between 0 and 255
// values were originally between 0 and 1 to perform shading computations since PhongShading expects
// values between 0 and 1
func colorToRGBA(c vec3) [4]uint8 {
	return [4]uint8{
		uint8(c.x * 255),
		uint8(c.y * 255),
		uint8(c.z * 255),
		255,
	}
}

// GetQuerySequence returns the amino acid sequence of a slice of atoms
func GetQuerySequence(atoms []*Atom) string {
	sequence := ""
	current_ind := -100
	for i := 0; i < len(atoms); i++ {
		if atoms[i].seqIndex != current_ind {
			sequence += ConvertAminoAcidToSingleChar(atoms[i].amino)
			current_ind = atoms[i].seqIndex
		}
	}
	return sequence
}

// ConvertAminoAcidToSingleChar converts a 3 letter amino acid to a single character code
// this is to make it easier to perform sequence alignment, this way each index is associated
// with a single character in the amino acid sequence
func ConvertAminoAcidToSingleChar(aa string) string {
	switch aa {
	case "MET":
		return "M"
	case "ALA":
		return "A"
	case "ARG":
		return "R"
	case "ASN":
		return "N"
	case "ASP":
		return "D"
	case "CYS":
		return "C"
	case "GLN":
		return "Q"
	case "GLU":
		return "E"
	case "GLY":
		return "G"
	case "HIS":
		return "H"
	case "ILE":
		return "I"
	case "LEU":
		return "L"
	case "LYS":
		return "K"
	case "PHE":
		return "F"
	case "PRO":
		return "P"
	case "SER":
		return "S"
	case "THR":
		return "T"
	case "TRP":
		return "W"
	case "TYR":
		return "Y"
	case "VAL":
		return "V"
	default:
		fmt.Println(aa)
		panic("Invalid amino acid")
	}

}

func MaxSeqIndex(atoms []*Atom) int {
	max := 0
	for i := 0; i < len(atoms); i++ {
		if atoms[i].seqIndex > max {
			max = atoms[i].seqIndex
		}
	}
	return max
}
