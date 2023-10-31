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
	// PrintPixelsArray()
}

func RenderScene() {
	// Initialize pixel data
	pixels = make([]uint8, 4*imageWidth*imageHeight)

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			index := (y*imageWidth + x) * 4

			// Calculate smooth gradient colors
			r := uint8(float32(x) / float32(imageWidth) * 255)  // Red component
			g := uint8(float32(y) / float32(imageHeight) * 255) // Green component
			b := uint8(0)                                       // Blue component
			a := uint8(255)                                     // Alpha (fully opaque)

			// Set pixel color in RGBA format
			pixels[index] = r
			pixels[index+1] = g
			pixels[index+2] = b
			pixels[index+3] = a
		}
	}

	// Use glDrawPixels to render the scene
}
func PrintPixelsArray() {
	// Print a simplified representation of the pixels array
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			index := (y*imageWidth + x) * 4
			fmt.Printf("(%3d,%3d,%3d,%3d) ", pixels[index], pixels[index+1], pixels[index+2], pixels[index+3])
		}
		fmt.Println()
	}
}
