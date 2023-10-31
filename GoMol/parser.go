package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

func ParsePDB(pdbFile string) []Atom {
	atoms := make([]Atom, 0)
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
		newAtom := Atom{number, element, residue, chain, sequence, x, y, z, 5.0}
		atoms = append(atoms, newAtom)
	}
	return atoms
}
