#ifdef GL_ES
precision mediump float;
#endif

varying vec2 vTexcoord;

uniform sampler2D sTexture;
uniform bool uTextured;

void main()
{
        if(uTextured) {
                gl_FragColor = texture2D(sTexture, vTexcoord);
        } else {
                gl_FragColor = vec4(1.0, 1.0, 1.0, 1.0);
        }
}

