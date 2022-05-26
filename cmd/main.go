package main

import (
	"log"
	"runtime"

	"github.com/inkyblackness/imgui-go/v4"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  = 1280
	WINDOW_HEIGHT = 768
)

type SDL struct {
	imguiIO imgui.IO

	window     *sdl.Window
	shouldStop bool

	time uint64
}

func (platform *SDL) Dispose() {
	if platform.window != nil {
		platform.window.Destroy()
		platform.window = nil
	}
	sdl.Quit()
}

func newSDL(io imgui.IO) (*SDL, error) {
	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, errors.Wrap(err, "failed to initialize SDL")
	}

	window, err := sdl.CreateWindow(
		"oh no",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		WINDOW_WIDTH,
		WINDOW_HEIGHT,
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		sdl.Quit()
		return nil, errors.Wrap(err, "failed to create window")
	}

	platform := &SDL{
		imguiIO: io,
		window:  window,
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)

	glContext, err := window.GLCreateContext()
	if err != nil {
		platform.Dispose()
		return nil, errors.Wrap(err, "failed to create OpenGL context")
	}
	if err := window.GLMakeCurrent(glContext); err != nil {
		platform.Dispose()
		return nil, errors.Wrap(err, "failed to set current OpenGL context")
	}

	sdl.GLSetSwapInterval(1)

	return platform, nil
}

func main() {
	context := imgui.CreateContext(nil)
	defer context.Destroy()

	io := imgui.CurrentIO()
	platform, err := newSDL(io)
	if err != nil {
		log.Fatalf("failed to initialize SDL: %s", err)
	}
	defer platform.Dispose()
}
