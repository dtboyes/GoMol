package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	imageWidth  = 800
	imageHeight = 600
)

var (
	pixels []uint8
)

func init() {
	runtime.LockOSThread()
}

func main() {

	// download a pdb file from the RCSB PDB database
	pdbURL := "https://files.rcsb.org/download/" + os.Args[1] + ".pdb"
	localPath := "pdbfiles/" + os.Args[1] + ".pdb"

	err := downloadPDB(pdbURL, localPath)
	if err != nil {
		fmt.Println("Error downloading PDB file:", err)
	}

	// Initialize GLFW and create a window
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(imageWidth, imageHeight, "Ray Generation", nil, nil)
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
	atoms := ParsePDB("pdbfiles/" + os.Args[1] + ".pdb")
	fmt.Println(atoms[0])

	// Generate rays and render them
	RenderScene()

	// Main loop
	for !window.ShouldClose() {
		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
		window.SwapBuffers()
		glfw.PollEvents()

	}
}
