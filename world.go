package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	vertexShaderSource = `
		#version 400

		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 400

		uniform vec3 un_color;

		out vec4 frag_colour;
		void main() {
  			frag_colour = vec4(un_color, 1.0);
		}
	` + "\x00"
)

const (
	minWidth = 200
	minHeight = 200
)

type t_action func(uint32) error

var (
	handlers []t_action
)

func run (action t_action, duration int64, width int, height int) error {
	if duration < 1 {
		errors.New("uiDuration")
	}

	if width < minWidth {
		errors.New("uiWidth")
	}

	if height < minHeight {
		errors.New("uiHeight")
	}

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Focused, glfw.True)

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Test window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.Maximize()

	prog := initGL()

	t := time.Now()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		if err := action(prog); err != nil {
			panic(err)
		}

		glfw.PollEvents()
		window.SwapBuffers()

		time.Sleep(time.Second/time.Duration(duration) - time.Since(t))
		t = time.Now()
	}

	return nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
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

func initGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	log.Printf("OpenGL version: %v", gl.GoStr(gl.GetString(gl.VERSION)))

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	pointer := gl.CreateProgram()
	gl.AttachShader(pointer, vertexShader)
	gl.AttachShader(pointer, fragmentShader)
	gl.LinkProgram(pointer)

	return pointer
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func draw(prog uint32, data []float32, color t_color)  {
	loc := gl.GetUniformLocation(prog, gl.Str("un_color" + "\x00"))
	gl.Uniform3f(loc, float32(color.R), float32(color.G), float32(color.B))

	gl.BindVertexArray(makeVao(data))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data) / 3))
}