package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var input string
	title := `
	===============================
		      GoMol
	===============================
	A Protein Analysis and Molecular
	   Visualization Tool
	===============================
	`
	fmt.Println(title)
	fmt.Print("Enter PDB ID #1: ")
	fmt.Scanln(&input)
	os.Args = append(os.Args, input)
	fmt.Print("Enter PDB ID #2: ")
	fmt.Scanln(&input)
	os.Args = append(os.Args, input)

	// initialize number of processors
	numProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(numProcs)

	// DOWNLOAD PDB FILE
	pdbURL := "https://files.rcsb.org/download/" + os.Args[1] + ".pdb"
	localPath := "pdbfiles/" + os.Args[1] + ".pdb"
	err := downloadPDB(pdbURL, localPath)
	if err != nil {
		fmt.Println("Error downloading PDB file:", err)
	}
	localPath = "pdbfiles/" + os.Args[2] + ".pdb"
	err = downloadPDB(pdbURL, localPath)
	if err != nil {
		fmt.Println("Error downloading PDB file:", err)
	}
	// Initialize GLFW and create a window
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(imageWidth, imageHeight, "GoMol: "+strings.ToUpper(os.Args[1]), nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.MakeContextCurrent()
	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	gl.Viewport(0, 0, imageWidth, imageHeight)

	// parse pdb file to get list of atom objects
	atoms1 = ParsePDB("pdbfiles/" + os.Args[1] + ".pdb")

	atoms2 = ParsePDB("pdbfiles/" + os.Args[2] + ".pdb")

	menu := `
	==============================
		    Options
	==============================
	1. Press "1" - Color protein by different side chain
	2. Press "2" - Color protein by different atom
	3. Press "3" - Color protein by differing regions from sequence alignment
	==============================
	`

	fmt.Println(menu)

	atoms1_sequence := GetQuerySequence(atoms1)
	atoms2_sequence := GetQuerySequence(atoms2)

	for _, val := range atoms1 {
		val.amino = ConvertAminoAcidToSingleChar(val.amino)
	}
	alignedSeq1, alignedSeq2, matchLine, percentSimilarity = NeedlemanWunsch(atoms1_sequence, atoms2_sequence)

	// initialize camera and light
	camera = InitializeCamera(atoms1)
	light = ParseLight("input/light.txt")
	light = InitializeLight(atoms1)

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetKeyCallback(keyCallback)
	window.SetScrollCallback(scrollCallback)

	// main loop to render scene
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		RotateAtoms(atoms1, rotationX, rotationY)
		pixels := make([]uint8, 4*imageWidth*imageHeight)
		RenderMultiProc(pixels, numProcs, atoms1, atoms2)
		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
		gl.LoadIdentity()
		window.SwapBuffers()
		glfw.PollEvents()

	}
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			leftMouseButtonPressed = true
			lastX, lastY = w.GetCursorPos()
		} else if action == glfw.Release {
			leftMouseButtonPressed = false
		}
	}
}

func cursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	if leftMouseButtonPressed {
		dx := xpos - lastX
		dy := ypos - lastY

		// Adjust rotation angles based on mouse movement
		rotationY += dx * 0.005
		rotationX += dy * 0.005

		// Limit rotation around X-axis to prevent flipping
		rotationX = math.Max(-math.Pi/2, math.Min(math.Pi/2, rotationX))
	} else {
		rotationX = 0.0
		rotationY = 0.0
	}

	lastX, lastY = xpos, ypos
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		} else if key == glfw.KeyW {
			camera.position.y += 0.1
		} else if key == glfw.KeyS {
			camera.position.y -= 0.1
		} else if key == glfw.KeyA {
			camera.position.x += 0.1
		} else if key == glfw.KeyD {
			camera.position.x -= 0.1
		} else if key == glfw.Key1 {
			colorByChain = true
			colorByAtom = false
			colorByDifferingRegions = false
		} else if key == glfw.Key2 {
			colorByChain = false
			colorByAtom = true
			colorByDifferingRegions = false
		} else if key == glfw.Key3 {
			colorByChain = false
			colorByAtom = false
			colorByDifferingRegions = true
		} else if key == glfw.Key4 {
			colorByChain = false
			colorByAtom = false
			colorByDifferingRegions = false
		}
	}
}

func scrollCallback(window *glfw.Window, xoff, yoff float64) {
	camera.position.z += yoff * 0.1
	camera.position.z -= xoff * 0.1
}
