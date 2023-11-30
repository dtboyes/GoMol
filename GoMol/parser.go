package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

// ParsePDB takes as input a pdb file and returns a list of Atom objects
func ParsePDB(pdbFile string) []*Atom {
	atoms := make([]*Atom, 0)
	f, _ := os.Open(pdbFile)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	currIndex := 0
	for scanner.Scan() {
		re := regexp.MustCompile(`\s+`)
		line := scanner.Text()
		parts := re.Split(line, -1)
		if parts[0] != "ATOM" {
			continue
		}
		number, _ := strconv.Atoi(parts[1])
		element := parts[2]
		amino := parts[3]
		chain := parts[4]
		seqIndex, _ := strconv.Atoi(parts[5])
		x, _ := strconv.ParseFloat(parts[6], 64)
		y, _ := strconv.ParseFloat(parts[7], 64)
		y *= -1.0
		z, _ := strconv.ParseFloat(parts[8], 64)
		// excludes atoms that are not part of the protein backbone
		if element == "H" || element == "C" || element == "N" || element == "O" || element == "S" {
			radius := 0.0
			// radii based on Pauling radii
			if element == "H" {
				radius = 1.2
			} else if element == "C" {
				radius = 1.7
			} else if element == "N" {
				radius = 1.55
			} else if element == "O" {
				radius = 1.52
			} else if element == "S" {
				radius = 1.8
			}
			newAtom := &Atom{number, element, amino, chain, currIndex, x, y, z, radius}
			if seqIndex != currIndex {
				currIndex++
			}
			atoms = append(atoms, newAtom)
		}
	}
	return atoms
}

func ParseCamera(cameraFile string) *Camera {
	f, _ := os.Open(cameraFile)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	cam := &Camera{}
	for scanner.Scan() {
		re := regexp.MustCompile(`\s+`)
		line := scanner.Text()
		parts := re.Split(line, -1)
		if parts[0] == "CAMERA" {
			continue
		}
		for i := 0; i < len(parts); i++ {
			if parts[i] == "pos" {
				cam.position.x, _ = strconv.ParseFloat(parts[i+1], 64)
				cam.position.y, _ = strconv.ParseFloat(parts[i+2], 64)
				cam.position.z, _ = strconv.ParseFloat(parts[i+3], 64)
				break
			} else if parts[i] == "focal_Length" {
				cam.focalLength, _ = strconv.ParseFloat(parts[i+1], 64)
				break
			} else if parts[i] == "viewport_height" {
				cam.viewportHeight, _ = strconv.ParseFloat(parts[i+1], 64)
				break
			} else if parts[i] == "camera_radius" {
				cam.radius, _ = strconv.ParseFloat(parts[i+1], 64)
			} else if parts[i] == "camera_yaw" {
				cam.yaw, _ = strconv.ParseFloat(parts[i+1], 64)
			} else if parts[i] == "camera_pitch" {
				cam.pitch, _ = strconv.ParseFloat(parts[i+1], 64)
			} else if parts[i] == "camera_speed" {
				cam.speed, _ = strconv.ParseFloat(parts[i+1], 64)
			}
		}
	}
	return cam
}

func ParseLight(lightFile string) *Light {
	f, _ := os.Open(lightFile)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	light := &Light{}
	for scanner.Scan() {
		re := regexp.MustCompile(`\s+`)
		line := scanner.Text()
		parts := re.Split(line, -1)
		if parts[0] == "LIGHT" {
			continue
		}
		for i := 0; i < len(parts); i++ {
			if parts[i] == "pos" {
				light.position.x, _ = strconv.ParseFloat(parts[i+1], 64)
				light.position.y, _ = strconv.ParseFloat(parts[i+2], 64)
				light.position.z, _ = strconv.ParseFloat(parts[i+3], 64)
				break
			}
		}
	}
	return light
}
