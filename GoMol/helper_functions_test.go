package main

import (
	"bufio"
	"io/fs"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

type AddTest struct {
	v1     vec3
	v2     vec3
	result vec3
}

type SubtractTest struct {
	v1     vec3
	v2     vec3
	result vec3
}

type CrossTest struct {
	v1     vec3
	v2     vec3
	result vec3
}

type DotTest struct {
	v1     vec3
	v2     vec3
	result float64
}

type ScaleTest struct {
	v1     vec3
	s      float64
	result vec3
}

type NormalizeTest struct {
	v1     vec3
	result vec3
}

type LengthTest struct {
	v1     vec3
	result float64
}

type EqualsZeroTest struct {
	v1     vec3
	result bool
}

func TestAdd(t *testing.T) {
	// read in all tests from the Tests/Add directory and run them
	tests := ReadAddTests("Tests/Add/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Add(test.v2)
		ourAnswer.x = roundFloat(ourAnswer.x, 2)
		ourAnswer.y = roundFloat(ourAnswer.y, 2)
		ourAnswer.z = roundFloat(ourAnswer.z, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Add(%v, %v) = %v, want %v", test.v1, test.v2, ourAnswer, test.result)
		}
	}
}

func ReadAddTests(directory string) []AddTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]AddTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadAddTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadVectorFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadAddTestFile(file string) AddTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var addTest AddTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			addTest.v1 = ReadVector(line)
		} else if index == 1 {
			addTest.v2 = ReadVector(line)
		}
		index++
	}
	return addTest
}

func TestSubtract(t *testing.T) {
	// read in all tests from the Tests/Subtract directory and run them
	tests := ReadSubtractTests("Tests/Subtract/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Subtract(test.v2)
		ourAnswer.x = roundFloat(ourAnswer.x, 2)
		ourAnswer.y = roundFloat(ourAnswer.y, 2)
		ourAnswer.z = roundFloat(ourAnswer.z, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Subtract(%v, %v) = %v, want %v", test.v1, test.v2, ourAnswer, test.result)
		}
	}
}

func ReadSubtractTests(directory string) []SubtractTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]SubtractTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadSubtractTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadVectorFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadSubtractTestFile(file string) SubtractTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var subtractTest SubtractTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			subtractTest.v1 = ReadVector(line)
		} else if index == 1 {
			subtractTest.v2 = ReadVector(line)
		}
		index++
	}
	return subtractTest
}

func TestCross(t *testing.T) {
	// read in all tests from the Tests/Cross directory and run them
	tests := ReadCrossTests("Tests/Cross/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Cross(test.v2)
		ourAnswer.x = roundFloat(ourAnswer.x, 2)
		ourAnswer.y = roundFloat(ourAnswer.y, 2)
		ourAnswer.z = roundFloat(ourAnswer.z, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Cross(%v, %v) = %v, want %v", test.v1, test.v2, ourAnswer, test.result)
		}
	}
}

func ReadCrossTests(directory string) []CrossTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]CrossTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadCrossTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadVectorFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadCrossTestFile(file string) CrossTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var crossTest CrossTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			crossTest.v1 = ReadVector(line)
		} else if index == 1 {
			crossTest.v2 = ReadVector(line)
		}
		index++
	}
	return crossTest
}

func TestDot(t *testing.T) {
	// read in all tests from the Tests/Dot directory and run them
	tests := ReadDotTests("Tests/Dot/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Dot(test.v2)
		ourAnswer = roundFloat(ourAnswer, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Dot(%v, %v) = %v, want %v", test.v1, test.v2, ourAnswer, test.result)
		}
	}
}

func ReadDotTests(directory string) []DotTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]DotTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadDotTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadFloatFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadDotTestFile(file string) DotTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var dotTest DotTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			dotTest.v1 = ReadVector(line)
		} else if index == 1 {
			dotTest.v2 = ReadVector(line)
		}
		index++
	}
	return dotTest
}

func TestScale(t *testing.T) {
	// read in all tests from the Tests/Scale directory and run them
	tests := ReadScaleTests("Tests/Scale/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Scale(test.s)
		ourAnswer.x = roundFloat(ourAnswer.x, 2)
		ourAnswer.y = roundFloat(ourAnswer.y, 2)
		ourAnswer.z = roundFloat(ourAnswer.z, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Scale(%v, %v) = %v, want %v", test.v1, test.s, ourAnswer, test.result)
		}
	}
}

func ReadScaleTests(directory string) []ScaleTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]ScaleTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadScaleTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadVectorFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadScaleTestFile(file string) ScaleTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var scaleTest ScaleTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			scaleTest.v1 = ReadVector(line)
		} else if index == 1 {
			scaleTest.s, _ = strconv.ParseFloat(line, 64)
		}
		index++
	}
	return scaleTest
}

func TestNormalize(t *testing.T) {
	// read in all tests from the Tests/Normalize directory and run them
	tests := ReadNormalizeTests("Tests/Normalize/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Normalize()
		ourAnswer.x = roundFloat(ourAnswer.x, 2)
		ourAnswer.y = roundFloat(ourAnswer.y, 2)
		ourAnswer.z = roundFloat(ourAnswer.z, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Normalize(%v) = %v, want %v", test.v1, ourAnswer, test.result)
		}
	}
}

func ReadNormalizeTests(directory string) []NormalizeTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]NormalizeTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadNormalizeTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadVectorFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadNormalizeTestFile(file string) NormalizeTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var normalizeTest NormalizeTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			normalizeTest.v1 = ReadVector(line)
		}
		index++
	}
	return normalizeTest
}

func TestLength(t *testing.T) {
	// read in all tests from the Tests/Length directory and run them
	tests := ReadLengthTests("Tests/Length/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.Length()
		ourAnswer = roundFloat(ourAnswer, 2)
		// check the result
		if ourAnswer != test.result {
			t.Errorf("Length(%v) = %v, want %v", test.v1, ourAnswer, test.result)
		}
	}
}

func ReadLengthTests(directory string) []LengthTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]LengthTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadLengthTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadFloatFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadLengthTestFile(file string) LengthTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var lengthTest LengthTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			lengthTest.v1 = ReadVector(line)
		}
		index++
	}
	return lengthTest
}

func ReadDirectory(dir string) []fs.DirEntry {
	//read in all files in the given directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}
func ReadFloatFromFile(file string) float64 {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//create a new scanner
	scanner := bufio.NewScanner(f)
	//read in the line
	scanner.Scan()
	line := scanner.Text()
	//convert the line to an int using strconv
	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		panic(err)
	}
	return value
}

func TestEqualsZero(t *testing.T) {
	// read in all tests from the Tests/EqualsZero directory and run them
	tests := ReadEqualsZeroTests("Tests/EqualsZero/")
	for _, test := range tests {
		// run the test
		ourAnswer := test.v1.EqualsZero()
		// check the result
		if ourAnswer != test.result {
			t.Errorf("EqualsZero(%v) = %v, want %v", test.v1, ourAnswer, test.result)
		}
	}
}

func ReadEqualsZeroTests(directory string) []EqualsZeroTest {

	// read in all tests from the directory and run them
	inputFiles := ReadDirectory(directory + "input")
	numFiles := len(inputFiles)

	tests := make([]EqualsZeroTest, numFiles)
	for i, inputFile := range inputFiles {
		// read in the test's map
		tests[i] = ReadEqualsZeroTestFile(directory + "input/" + inputFile.Name())
	}

	// now, read output files
	outputFiles := ReadDirectory(directory + "output")

	// ensure the same number of input and output files
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}

	for i, outputFile := range outputFiles {
		// read in the test's result
		tests[i].result = ReadBoolFromFile(directory + "output/" + outputFile.Name())
	}

	return tests
}

func ReadEqualsZeroTestFile(file string) EqualsZeroTest {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() // this command occurs at the end of the function

	// create a new scanner
	scanner := bufio.NewScanner(f)

	var equalsZeroTest EqualsZeroTest
	// while the scanner still has lines to read,
	// read in the next line
	index := 0
	for scanner.Scan() {
		// read in the line
		line := scanner.Text()
		// split the line into two parts
		if index == 0 {
			equalsZeroTest.v1 = ReadVector(line)
		}
		index++
	}
	return equalsZeroTest
}

func ReadVector(line string) vec3 {
	parts := strings.Split(line, " ")
	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)
	z, _ := strconv.ParseFloat(parts[2], 64)
	return vec3{x, y, z}
}

func ReadVectorFromFile(file string) vec3 {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//create a new scanner
	scanner := bufio.NewScanner(f)
	//read in the line
	scanner.Scan()
	line := scanner.Text()
	//convert the line to an int using strconv
	parts := strings.Split(line, " ")
	x, _ := strconv.ParseFloat(parts[0], 64)
	y, _ := strconv.ParseFloat(parts[1], 64)
	z, _ := strconv.ParseFloat(parts[2], 64)
	return vec3{x, y, z}
}

func ReadBoolFromFile(file string) bool {
	//open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	//create a new scanner
	scanner := bufio.NewScanner(f)
	//read in the line
	scanner.Scan()
	line := scanner.Text()
	//convert the line to an int using strconv
	value, err := strconv.ParseBool(line)
	if err != nil {
		panic(err)
	}
	return value
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
