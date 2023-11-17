package main

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"os"
// 	"runtime"
// 	"strings"

// 	"github.com/go-gl/gl/v2.1/gl"
// 	"github.com/go-gl/glfw/v3.3/glfw"
// )

// const (
// 	imageWidth     = 800
// 	aspectRatio    = 16.0 / 9.0
// 	imageHeight    = imageWidth / aspectRatio
// 	viewportHeight = 2.0
// 	viewportWidth  = viewportHeight * float64(imageWidth) / float64(imageHeight)
// )

// var (
// 	camera1 *Camera
// 	camera2 *Camera
// 	light   *Light
// 	atoms1  []*Atom
// 	atoms2  []*Atom
// )

// var (
// 	pixels1 []uint8
// 	pixels2 []uint8
// )

// var (
// 	leftMouseButtonPressed bool
// 	lastX, lastY           float64
// )

// func init() {
// 	runtime.LockOSThread()
// }

// func main() {
// 	// initialize number of processors
// 	numProcs := runtime.NumCPU()
// 	runtime.GOMAXPROCS(numProcs)

// 	// DOWNLOAD PDB FILE
// 	pdbURL := "https://files.rcsb.org/download/" + os.Args[1] + ".pdb"
// 	localPath := "pdbfiles/" + os.Args[1] + ".pdb"
// 	err := downloadPDB(pdbURL, localPath)
// 	if err != nil {
// 		fmt.Println("Error downloading PDB file:", err)
// 	}

// 	// Initialize GLFW and create a window
// 	if err := glfw.Init(); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer glfw.Terminate()

// 	glfw.WindowHint(glfw.Resizable, glfw.False)

// 	window1, err := glfw.CreateWindow(imageWidth, imageHeight, "GoMol: "+strings.ToUpper(os.Args[1]), nil, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer window1.Destroy()

// 	window2, err := glfw.CreateWindow(imageWidth, imageHeight, "GoMol: "+strings.ToUpper(os.Args[2]), nil, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer window2.Destroy()
// 	// Initialize OpenGL
// 	if err := gl.Init(); err != nil {
// 		log.Fatal(err)
// 	}
// 	gl.Viewport(0, 0, imageWidth, imageHeight)

// 	// parse pdb file to get list of atom objects
// 	atoms1 = ParsePDB("pdbfiles/" + os.Args[1] + ".pdb")

// 	// initialize camera and light
// 	camera1 = InitializeCamera(atoms1)
// 	camera2 = InitializeCamera(atoms2)
// 	light = ParseLight("input/light.txt")

// 	window1.SetMouseButtonCallback(mouseButtonCallback)
// 	window1.SetCursorPosCallback(cursorPosCallback)
// 	window1.SetKeyCallback(keyCallback)
// 	window1.SetScrollCallback(scrollCallback)
// 	window2.SetMouseButtonCallback(mouseButtonCallback)
// 	window2.SetCursorPosCallback(cursorPosCallback)
// 	window2.SetKeyCallback(keyCallback)
// 	window2.SetScrollCallback(scrollCallback)
// 	// RenderScene(camera, light, atoms)
// 	// pixels = make([]uint8, 4*imageWidth*imageHeight)
// 	// RenderScene(camera, light, atoms, 0, imageHeight, pixels)

// 	// main loop to render scene
// 	for !window1.ShouldClose() && !window2.ShouldClose() {
// 		window1.MakeContextCurrent()
// 		pixels1 = make([]uint8, 4*imageWidth*imageHeight)
// 		RenderMultiProc(camera1, atoms1, pixels1, numProcs)
// 		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels1))
// 		window1.SwapBuffers()

// 		window2.MakeContextCurrent()
// 		pixels2 = make([]uint8, 4*imageWidth*imageHeight)
// 		RenderMultiProc(camera2, atoms2, pixels2, numProcs)
// 		gl.DrawPixels(imageWidth, imageHeight, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels2))
// 		window2.SwapBuffers()

// 		// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
// 		glfw.PollEvents()

// 	}
// }

// func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
// 	if button == glfw.MouseButtonLeft {
// 		if action == glfw.Press {
// 			leftMouseButtonPressed = true
// 			lastX, lastY = w.GetCursorPos()
// 		} else if action == glfw.Release {
// 			leftMouseButtonPressed = false
// 		}
// 	}
// }

// func cursorPosCallback(camerawindow *glfw.Window, xpos, ypos float64) {
// 	if leftMouseButtonPressed {
// 		camera1.yaw += camera1.speed * (xpos - lastX)
// 		camera1.pitch += camera1.speed * (ypos - lastY)
// 		if camera1.pitch > 89.0 {
// 			camera1.pitch = 89.0
// 		}
// 		if camera1.pitch < -89.0 {
// 			camera1.pitch = -89.0
// 		}
// 		camera1.position.x += camera1.radius * math.Cos(degToRad(camera1.yaw)) * math.Cos(degToRad(camera1.pitch))
// 		camera1.position.y = camera1.radius * math.Sin(degToRad(camera1.pitch))
// 		camera1.position.z += camera1.radius * math.Sin(degToRad(camera1.yaw)) * math.Cos(degToRad(camera1.pitch))
// 	}
// }

// func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
// 	if action == glfw.Press || action == glfw.Repeat {
// 		if key == glfw.KeyEscape {
// 			window.SetShouldClose(true)
// 		} else if key == glfw.KeyW {
// 			light.position.y += 0.1
// 		} else if key == glfw.KeyS {
// 			light.position.y -= 0.1
// 		} else if key == glfw.KeyA {
// 			light.position.x -= 0.1
// 		} else if key == glfw.KeyD {
// 			light.position.x += 0.1
// 		}
// 	}
// }

// func scrollCallback(window *glfw.Window, xoff, yoff float64) {
// 	camera1.position.z += yoff * 0.1
// 	camera1.position.z -= xoff * 0.1
// }

// func degToRad(degrees float64) float64 {
// 	return degrees * math.Pi / 180.0
// }
