#version 310 es

// SPDX-License-Identifier: Unlicense OR MIT

precision mediump float;

layout(location=0) in vec2 vUV;

{{.Header}}

layout(location = 0) out vec4 fragColor;

void main() {
	fragColor = {{.FetchColorExpr}};
}
