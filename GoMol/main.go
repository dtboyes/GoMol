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
	fmt.Print("Only render chain A? (y/n): ")
	fmt.Scanln(&input)
	if input == "y" {
		onlyChainA = true
	}

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
	pdbURL = "https://files.rcsb.org/download/" + os.Args[2] + ".pdb"
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

	atoms1_sequence = GetQuerySequence(atoms1)
	atoms2_sequence = GetQuerySequence(atoms2)

	alignedSeq1, alignedSeq2, matchLine, percentSimilarity = NeedlemanWunsch(atoms1_sequence, atoms2_sequence)
	fmt.Println(alignedSeq1)
	fmt.Println(matchLine)
	fmt.Println(alignedSeq2)
	fmt.Printf("The percent identity of the two sequences using Needleman-Wunsch is %.2f%%\n\n", percentSimilarity)

	// initialize camera and light
	// camera = InitializeCamera(atoms1)
	light = ParseLight("input/light.txt")
	light = InitializeLight(atoms1)

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetKeyCallback(keyCallback)
	window.SetScrollCallback(scrollCallback)

	// main loop to render scene
	var results1, results2, resultsFinal []*Atom
	var rmsd float64
	if len(atoms1_sequence) == len(atoms2_sequence) {
		results1, results2, rmsd = RunKabsch(atoms1, atoms2)
		resultsFinal = append(results1, results2...)
		atoms1_len := len(atoms1)
		for i := atoms1_len; i < len(resultsFinal); i++ {
			resultsFinal[i].x -= 50
		}
	}
	tempAtoms1 := atoms1
	tempAtoms2 := atoms2
	tempAtomsResult := resultsFinal
	fmt.Println("RMSD from Kabsch algorithm: ", rmsd)
	camera = InitializeCamera(atoms1)
	for !window.ShouldClose() {
		if renderProtein1 {
			atoms1 = tempAtoms1
		} else if renderProtein2 {
			atoms1 = tempAtoms2
		} else if renderBoth {
			atoms1 = tempAtomsResult
		}
		gl.Clear(gl.COLOR_BUFFER_BIT)
		RotateAtoms(atoms1, rotationX, rotationY)
		pixels := make([]uint8, 4*imageWidth*imageHeight)
		RenderMultiProc(pixels, numProcs, true)
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
		} else if key == glfw.KeyW && action == glfw.Press {
			camera.position.y += 0.1
		} else if key == glfw.KeyS && action == glfw.Press {
			camera.position.y -= 0.1
		} else if key == glfw.KeyA && action == glfw.Press {
			camera.position.x += 0.1
		} else if key == glfw.KeyD && action == glfw.Press {
			camera.position.x -= 0.1
		} else if key == glfw.Key1 && action == glfw.Press {
			colorByChain = true
			colorByAtom = false
			colorByDifferingRegions = false
		} else if key == glfw.Key2 && action == glfw.Press {
			colorByChain = false
			colorByAtom = true
			colorByDifferingRegions = false
		} else if key == glfw.Key3 && action == glfw.Press {
			colorByChain = false
			colorByAtom = false
			colorByDifferingRegions = true
		} else if key == glfw.Key4 && action == glfw.Press {
			colorByChain = false
			colorByAtom = false
			colorByDifferingRegions = false
		} else if key == glfw.KeyF1 && action == glfw.Press {
			renderProtein1 = true
			renderProtein2 = false
			renderBoth = false
		} else if key == glfw.KeyF2 && action == glfw.Press {
			renderProtein1 = false
			renderProtein2 = true
			renderBoth = false
		} else if key == glfw.KeyF3 && action == glfw.Press {
			renderProtein1 = false
			renderProtein2 = false
			renderBoth = true
		}
	}
}

func scrollCallback(window *glfw.Window, xoff, yoff float64) {
	camera.position.z += yoff * 0.1
	camera.position.z -= xoff * 0.1
}
