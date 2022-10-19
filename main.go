package main

import (
    "log"
    "runtime"

    "github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
    "github.com/go-gl/glfw/v3.3/glfw"
)

const ( // constants such as width and height
    width  = 500
    height = 500
)

// it's a triangle, duh
var (
    triangle = []float32{
        0, 0.5, 0, // top
        -0.5, -0.5, 0, // left
        0.5, -0.5, 0, // right
    }
)

func main() {
    runtime.LockOSThread() // makes sure all code is executed on the thread it was initialized on

    window := initGlfw()
    defer glfw.Terminate() // this code runs after main is done because defer

    program := initOpengl()

    for !window.ShouldClose() {
        draw(window, program) // when the program is closing, draw
    }
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
    if err := glfw.Init(); err != nil {
            panic(err)
    }
    
    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
    if err != nil {
            panic(err)
    }
    window.MakeContextCurrent() // saves all the contexts

    return window
}

// initGlfw initializes Opengl and returns an initialized program

// returns a program wich is important because it gives us a reference to store shaders
func initOpengl() uint32 { 
    if err := gl.Init();
    err != nil {
        panic(err)
    }
    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("Opengl version", version)

    prog := gl.CreateProgram()
    gl.LinkProgram(prog)
    return prog
}

func draw(window *glfw.Window, prog uint32 ) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // removes anything that was drawn last frame
    gl.UseProgram(prog) // using our program refrence
    
    glfw.PollEvents() // checks for events such as keystrokes and mouseclicks
    
    // swaps the buffers wich is important because opengl actually uses 2 canvases, one visible and one invisible
    // when drawing we draw to the invisible one and then later swap them with the bufferswap method
    window.SwapBuffers()
}