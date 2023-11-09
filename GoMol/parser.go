package main

import (
	"bufio"
	"fmt"
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
	for scanner.Scan() {
		re := regexp.MustCompile(`\s+`)
		line := scanner.Text()
		parts := re.Split(line, -1)
		if parts[0] != "ATOM" {
			continue
		}
		number, _ := strconv.Atoi(parts[1])
		element := parts[2]
		residue := parts[3]
		chain := parts[4]
		sequence := parts[5]
		x, _ := strconv.ParseFloat(parts[6], 64)
		y, _ := strconv.ParseFloat(parts[7], 64)
		z, _ := strconv.ParseFloat(parts[8], 64)
		newAtom := &Atom{number, element, residue, chain, sequence, x, y, z, 5.0}
		atoms = append(atoms, newAtom)
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
			} else if parts[i] == "focal_length" {
				cam.focalLength, _ = strconv.ParseFloat(parts[i+1], 64)
				break
			} else if parts[i] == "viewport_height" {
				cam.viewportHeight, _ = strconv.ParseFloat(parts[i+1], 64)
				break
			}
		}
	}
	fmt.Println(cam)
	return cam
}
