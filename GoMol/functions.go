package main

import (
	"fmt"
	"math"
)

func RenderMultiProc(pixels []uint8, numProcs int, window1 bool) {
	// pixels = make([]uint8, 4*imageWidth*imageHeight)
	finished := make(chan bool, numProcs)
	for i := 0; i < numProcs; i++ {
		start_height := i * imageHeight / numProcs
		end_height := (i + 1) * imageHeight / numProcs
		if renderProtein1 || renderBoth {
			go RenderScene(camera, light, atoms1, atoms1_sequence, alignedSeq2, start_height, end_height, pixels, finished)
		} else if renderProtein2 {
			go RenderScene(camera, light, atoms2, atoms2_sequence, alignedSeq1, start_height, end_height, pixels, finished)
		}
	}
	for i := 0; i < numProcs; i++ {
		<-finished
	}
}

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
func RayColor(r *Ray, light *Light, camera *Camera, atoms []*Atom, atoms_sequence, aligned_sequence string) vec3 {
	for i := 0; i < len(atoms); i++ {
		collision := RaySphereCollision(r, atoms[i])
		if !collision.getNormal().EqualsZero() {
			if colorByChain {
				if atoms[i].chain == "A" {
					collision.color = LambertianShading(collision, light, camera, vec3{0.2, 0.7, 0.1})
				} else if atoms[i].chain == "B" {
					collision.color = LambertianShading(collision, light, camera, vec3{0.1, 0.2, 1.0})
				} else if atoms[i].chain == "C" {
					collision.color = LambertianShading(collision, light, camera, vec3{1.0, 0.1, 0.2})
				} else if atoms[i].chain == "D" {
					collision.color = LambertianShading(collision, light, camera, vec3{1.0, 0.55, 0.0})
				} else {
					collision.color = LambertianShading(collision, light, camera, vec3{1.0, 1.0, 1.0})
				}
			} else if colorByAtom {
				if atoms[i].element == "CA" {
					collision.color = LambertianShading(collision, light, camera, vec3{0.565, 0.565, 0.565})
				} else if atoms[i].element == "N" {
					collision.color = LambertianShading(collision, light, camera, vec3{0.188, 0.313, 0.9725})
				} else if atoms[i].element == "O" {
					collision.color = LambertianShading(collision, light, camera, vec3{1.0, 0.051, 0.051})
				} else if atoms[i].element == "S" {
					collision.color = LambertianShading(collision, light, camera, vec3{1.0, 0.784, 0.196})
				}
			} else if colorByDifferingRegions {
				// fmt.Println(MaxSeqIndex(atoms))
				if alignedSeq1[atoms[i].seqIndex] != alignedSeq2[atoms[i].seqIndex] {
					collision.color = LambertianShading(collision, light, camera, vec3{0.69, 0.22, 0.188})
				} else {
					collision.color = LambertianShading(collision, light, camera, vec3{0.373, 0.651, 0.286})
				}
			} else {
				collision.color = LambertianShading(collision, light, camera, vec3{0.373, 0.651, 0.286})
			}
			return collision.color
		}

	}
	return vec3{0, 0, 0}
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

func InitializeLight(atoms []*Atom) *Light {
	light.position = light.position.Add(CenterOfMass(atoms))
	return light
}

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

func colorToRGBA(c vec3) [4]uint8 {
	return [4]uint8{
		uint8(c.x * 255),
		uint8(c.y * 255),
		uint8(c.z * 255),
		255,
	}
}

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
