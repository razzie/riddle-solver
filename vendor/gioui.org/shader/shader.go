// SPDX-License-Identifier: Unlicense OR MIT

package shader

type Sources struct {
	Name          string
	GLSL100ES     string
	GLSL300ES     string
	GLSL310ES     string
	GLSL130       string
	GLSL150       string
	DXBC          string
	MetalLib      string
	Uniforms      UniformsReflection
	Inputs        []InputLocation
	Textures      []TextureBinding
	WorkgroupSize [3]int
}

type UniformsReflection struct {
	Blocks    []UniformBlock
	Locations []UniformLocation
	Size      int
}

type TextureBinding struct {
	Name    string
	Binding int
}

type UniformBlock struct {
	Name    string
	Binding int
}

type UniformLocation struct {
	Name   string
	Type   DataType
	Size   int
	Offset int
}

type InputLocation struct {
	// For GLSL.
	Name     string
	Location int
	// For HLSL.
	Semantic      string
	SemanticIndex int

	Type DataType
	Size int
}

type DataType uint8

type DepthFunc uint8

const (
	DataTypeFloat DataType = iota
	DataTypeInt
	DataTypeShort
)
