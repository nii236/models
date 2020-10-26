package main

import (
	. "github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
)

const (
	scale  = 1    // optional supersampling
	width  = 1920 // output width in pixels
	height = 1080 // output height in pixels
	fovy   = 30   // vertical field of view in degrees
	near   = 1    // near clipping plane
	far    = 10   // far clipping plane
)
const distance = 5

var isometric_eye = V(-3, -3, 3)                     // camera position
var isometric_center = V(0, 0, 0)                    // view center position
var isometric_up = V(0, 0, 1)                        // up vector
var isometric_light = V(-0.75, -1, 0.25).Normalize() // light direction
var isometric_color = HexColor("#468966")            // object color

func RenderToPNG(in string, out string) {
	// load a mesh
	mesh, err := LoadSTL(in)
	if err != nil {
		panic(err)
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(Radians(30))

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("#424242"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(isometric_eye, isometric_center, isometric_up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := NewPhongShader(matrix, isometric_light, isometric_eye)
	shader.ObjectColor = isometric_color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	// save image
	SavePNG(out, image)
}
