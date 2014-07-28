attribute vec4 aPosition;
attribute vec2 aTexcoord;

varying vec2 vTexcoord;

uniform vec2 uCamera;
uniform vec2 uPosition;
uniform float uRight;
uniform float uTop;

void main()
{
        vTexcoord = aTexcoord;

        // This sets far == 1, near == -1 and symmetry across top/bottom and left/right
        mat4 window_scale = mat4(vec4(1.0 / uRight, 0.0, 0.0, 0.0),
                vec4(0.0, 1.0 / uTop, 0.0, 0.0),
                vec4(0.0, 0.0, 1.0, 0.0),
                vec4(0.0, 0.0, 0.0, 1.0));
        gl_Position = window_scale * (aPosition - vec4(uCamera, 0.0, 0.0) + vec4(uPosition, 0.0, 0.0));
}
