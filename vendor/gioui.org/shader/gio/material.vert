#version 310 es

// SPDX-License-Identifier: Unlicense OR MIT

#extension GL_GOOGLE_include_directive : enable

precision highp float;

#include "common.h"

layout(binding = 0) uniform Block {
	vec2 scale;
	vec2 pos;
} _block;

layout(location = 0) in vec2 pos;
layout(location = 1) in vec2 uv;

layout(location = 0) out vec2 vUV;

void main() {
	vUV = uv;
	vec3 p = vec3(pos*_block.scale + _block.pos, 1.0);
	gl_Position = vec4(transform3x2(fboTransform, p), 1.0);
}
