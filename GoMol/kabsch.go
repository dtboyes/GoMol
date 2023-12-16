package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// RunKabsch takes as input two slices of atoms pointers and returns the rotated versions of these atoms slices
// such that the RMSD is minimized
func RunKabsch(samp1, samp2 []*Atom) ([]*Atom, []*Atom, float64) {
	var s1 *mat.Dense
	var s2 *mat.Dense

	s1 = GenerateMatrix(samp1)
	s2 = GenerateMatrix(samp2)

	result1, result2, rmsd := kabsch(s1, s2)

	r1 := GenerateAtomSlice(result1, samp1)
	r2 := GenerateAtomSlice(result2, samp2)

	return r1, r2, rmsd
}

// Kabsch takes as input two dense matrices and performs the linear algebra calcuations and returns the matrices
// of the rotated points and the RMSD
func kabsch(p, q *mat.Dense) (*mat.Dense, *mat.Dense, float64) {
	a := CopyMatrix(p)
	b := CopyMatrix(q)

	bColAvgs := AvgColumns(b)
	aColAvgs := AvgColumns(a)

	a = CenterAtOrigin(a)
	b = CenterAtOrigin(b)

	var E0 float64
	for i := 0; i < a.RawMatrix().Rows; i++ {
		for j := 0; j < a.RawMatrix().Cols; j++ {
			E0 += a.At(i, j) * a.At(i, j)
			E0 += b.At(i, j) * b.At(i, j)
		}
	}

	var h mat.Dense
	h.Mul(b.T(), a)

	fmt.Println("Covariance matrix h: ")

	matPrint(&h)

	var svd mat.SVD
	if ok := svd.Factorize(&h, mat.SVDFull); !ok {
		fmt.Println("SVD failed")
		return nil, nil, 0.0
	}

	S := svd.Values(nil)

	var U mat.Dense
	svd.UTo(&U)

	fmt.Println("SVD U matrix: ")
	matPrint(&U)

	var VT mat.Dense
	svd.VTo(&VT)

	fmt.Println("SVD VT: ")
	matPrint(VT.T())

	reflect := mat.Det(&U) * mat.Det(VT.T())
	if reflect < 0 {
		S[len(S)-1] = -S[len(S)-1]
		for i := range U.RawMatrix().Data {
			if i%U.RawMatrix().Cols == U.RawMatrix().Cols-1 {
				U.RawMatrix().Data[i] = -U.RawMatrix().Data[i]
			}
		}
	}

	fmt.Println("Reflect value: ", reflect)

	RMSD := E0 - 2*mat.Sum(mat.NewVecDense(len(S), S))
	RMSD = math.Sqrt(math.Abs(RMSD / float64(a.RawMatrix().Rows)))

	var r mat.Dense

	r.Mul(&U, VT.T())

	for i := 0; i < b.RawMatrix().Rows; i++ {
		for j := 0; j < b.RawMatrix().Cols; j++ {
			b.Set(i, j, b.At(i, j)+bColAvgs[j])
		}
	}

	var bRotated mat.Dense

	bRotated.Mul(b, &r)

	for i := 0; i < a.RawMatrix().Rows; i++ {
		for j := 0; j < a.RawMatrix().Cols; j++ {
			a.Set(i, j, a.At(i, j)+aColAvgs[j])
		}
	}

	return a, &bRotated, RMSD
}

// GenerateMatrix takes as input a slice of atom pointers and returns a dense matrix of the coordinates
func GenerateMatrix(atoms []*Atom) *mat.Dense {
	n := len(atoms)

	data := make([]float64, 3*n)

	for i, atom := range atoms {
		data[3*i] = atom.x
		data[3*i+1] = atom.y
		data[3*i+2] = atom.z
	}

	matrix := mat.NewDense(n, 3, data)

	return matrix
}

// GenerateAtomSlice takes as input a dense matrix and atom slice and returns an atom slice with
// with the correct pdb info incorporated
func GenerateAtomSlice(matrix *mat.Dense, pdbInfo []*Atom) []*Atom {
	rows, _ := matrix.Dims()

	atoms := make([]*Atom, rows)

	for i := 0; i < rows; i++ {
		atoms[i] = &Atom{
			number:   pdbInfo[i].number,
			element:  pdbInfo[i].element,
			amino:    pdbInfo[i].amino,
			chain:    pdbInfo[i].chain,
			seqIndex: pdbInfo[i].seqIndex,
			x:        matrix.At(i, 0),
			y:        matrix.At(i, 1),
			z:        matrix.At(i, 2),
			radius:   pdbInfo[i].radius,
		}
	}

	return atoms
}

// CenterAtOrigin takes as input a dense matrix and returns another dense matrix with the coordinates
// centered at the origin
func CenterAtOrigin(a *mat.Dense) *mat.Dense {
	numRows, numCols := a.Dims()
	colAvgs := AvgColumns(a)
	for j := 0; j < numCols; j++ {
		for i := 0; i < numRows; i++ {
			a.Set(i, j, a.At(i, j)-colAvgs[j])
		}
	}
	return a
}

// SumColumns takes an input a dense matrix and returns a slice of floats with the sum of the columns
func SumColumns(a *mat.Dense) []float64 {
	numRows, numCols := a.Dims()
	sums := make([]float64, numCols)

	for j := 0; j < numCols; j++ {
		for i := 0; i < numRows; i++ {
			sums[j] += a.At(i, j)
		}
	}
	return sums
}

// AvgColumns takes as input a dense matrix and returns a slice of floats that correspond to the column averages
func AvgColumns(a *mat.Dense) []float64 {
	numRows, _ := a.Dims()
	avgs := SumColumns(a)

	for i := range avgs {
		avgs[i] = avgs[i] / float64(numRows)
	}
	return avgs
}

// CopyMatrix takes as input a dense matrix and returns a exact copy as a dense matrix
func CopyMatrix(a *mat.Dense) *mat.Dense {
	r, c := a.Dims()
	data := make([]float64, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = a.At(i, j)
		}
	}
	newMat := mat.NewDense(r, c, data)
	return newMat
}

// matPrint takes as input a matrix and prints it
func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}
