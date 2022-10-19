package main

import (
    "log"
    "runtime"
    "fmt"
    "strings"
    "github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
    "github.com/go-gl/glfw/v3.3/glfw"
)

const ( // constants such as width and height
    width  = 500
    height = 500

    // strings with source code for 2 shaders
    vertexShaderSource = `
        #version 410
        in vec3 vp;
        void main() {
            gl_Position = vec4(vp, 1.0);
        }
    ` + "\x00"

    fragmentShaderSource = `
        #version 410
        out vec4 frag_colour;
        void main() {
            frag_colour = vec4(1, 1, 1, 1);
        }
    ` + "\x00"
)

// it's a triangle, duh
var (
    triangle = []float32{ // always use float32 when providing verices to Opengl
        // values represented as xyz
        // ranges from -1 to 1
        // we're not using 3D so z is always 0
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

    vao := makeVao(triangle)
    
    for !window.ShouldClose() {
        draw(vao, window, program) // when the program is closing, draw
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
    
    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }
    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        panic(err)
    }

    prog := gl.CreateProgram()
    gl.AttachShader(prog, vertexShader)
    gl.AttachShader(prog, fragmentShader)
    gl.LinkProgram(prog)
    return prog
}

// Here we call **makeVao** to get our **vao** reference from the **triangle** points we defined before, and pass it as a new argument to the **draw** function:

func draw(vao uint32, window *glfw.Window, prog uint32 ) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // removes anything that was drawn last frame
    gl.UseProgram(prog) // using our program refrence

    gl.BindVertexArray(vao)
    // divides by 3 because there are 3 values so it knows how many vertices to draw
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle) / 3))

    glfw.PollEvents() // checks for events such as keystrokes and mouseclicks
    
    // swaps the buffers wich is important because opengl actually uses 2 canvases, one visible and one invisible
    // when drawing we draw to the invisible one and then later swap them with the bufferswap method
    window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided
// vao = Vertex Array Object
// vbo = vertex buffer object
func makeVao(points []float32) uint32 {
    var vbo uint32 // create a buffer
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    // generate, bind and giv the buffer data

    // 4 x len(points) because a 32 bit float has 4 bytes
    gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

    var vao uint32 // create vertex array
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)
    gl.EnableVertexAttribArray(0)
    // generate, bind and give it attributes
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
    return vao
}


func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csourcers, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csourcers, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }
    return shader, nil
}
