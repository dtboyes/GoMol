package main

import (
	"math"
)

func qRes(atoms1, atoms2 []*Atom) []float64 {
	if len(atoms1) != len(atoms2) {
		panic("residue length mismatch")
	}

	n := len(atoms1)

	qRes := make([]float64, n)

	cMap1 := GenerateContactMap(atoms1)
	cMap2 := GenerateContactMap(atoms2)

	for i := range qRes {
		var k float64
		if i == 0 || i == n-1 {
			k = 2.0
		} else {
			k = 3.0
		}
		sum := 0.0
		for j := range cMap2 {
			if j == i-1 || j == i || j == i+1 {
				continue
			}
			varIJ := math.Pow(math.Abs(float64(i-j)), 0.15)
			deltaDists := cMap1[i][j] - cMap2[i][j]
			expression := (deltaDists * deltaDists) / (2 * varIJ)
			sum += math.Exp(-expression)
		}

		qi := (1 / (float64(n) - k)) * sum

		qRes[i] = qi
	}
	return qRes
}

func GenerateContactMap(atoms []*Atom) [][]float64 {
	n := len(atoms)

	M := make([][]float64, n)

	for i := range M {
		M[i] = make([]float64, n)
	}

	for i, atom1 := range atoms {
		for j, atom2 := range atoms {
			dist := Distance(atom1, atom2)

			M[i][j] = dist
		}
	}
	return M

}

func FilterAlignedAtoms(seq1, seq2, align1, align2 string, atoms1, atoms2 []*Atom) ([]*Atom, []*Atom) {
	/*
		kAtoms1 := []*Atom{}
		kAtoms2 := []*Atom{}
		for i := range atoms1 {
			if atoms1[i].element == "CA" {
				kAtoms1 = append(kAtoms1, atoms1[i])
			}
		}

		for i := range atoms2 {
			if atoms2[i].element == "CA" {
				kAtoms2 = append(kAtoms2, atoms2[i])
			}
		}
	*/

	alignedAtoms1 := []*Atom{}
	alignedAtoms2 := []*Atom{}
	seqIndex1, seqIndex2 := 0, 0

	for i := 0; i < len(align1); i++ {
		// Check if the current position is not a gap in either sequence
		if align1[i] != '-' && align2[i] != '-' {
			// Add the atoms corresponding to the current aligned position
			alignedAtoms1 = append(alignedAtoms1, atoms1[seqIndex1])
			alignedAtoms2 = append(alignedAtoms2, atoms2[seqIndex2])
		}

		// Increment sequence indices if not a gap
		if align1[i] != '-' {
			seqIndex1++
		}
		if align2[i] != '-' {
			seqIndex2++
		}
	}

	return alignedAtoms1, alignedAtoms2
}

func Distance(atom1, atom2 *Atom) float64 {
	deltaX := atom2.x - atom1.x
	deltaY := atom2.y - atom1.y
	deltaZ := atom2.z - atom1.z

	dist := math.Sqrt(deltaX*deltaX + deltaY*deltaY + deltaZ*deltaZ)

	return dist
}
