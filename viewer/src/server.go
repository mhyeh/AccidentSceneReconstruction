package main

import (
	"fmt"
	"log"
	"net/http"
	"context"
	"time"
	"os"
	"path/filepath"
	"runtime"
	
	"github.com/go-gl/mathgl/mgl64"
	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"

	"view"
)

func init() {
	runtime.LockOSThread()
}

func getMouseNDC(w int, h int, mx int, my int) (x float64, y float64) {
	x = float64(mx) / float64(w) * 2 - 1
	y = float64(my) / float64(h) * 2 - 1
	return
}


func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/" + os.Args[1] + "/", server)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	baseURL := "140.118.127.145:3001"

	srv := &http.Server{Addr: baseURL, Handler: handler}

	var (
		camera *view.Camera
		model *view.ModelData
		renderer *view.Renderer
		mouseX, mouseY int
		windW, windH int
	)

	server.On("connection", func(so socketio.Socket) {
		fmt.Println("connected")
		so.Emit("connected")

		so.On("init", func(w int, h int) {
			// so.Emit("frame", renderer.Frame2Base64())
		})
		so.On("mouseDown", func(button int) {
			if renderer != nil {
				if button == 0 {
					camera.Mode = view.ROTATE
				} else if button == 2 {
					camera.Mode = view.PAN
				}
				x, y := getMouseNDC(windW, windH, mouseX, mouseY)
				camera.MouseDown(mgl64.Vec2{x, y})
			}
		})
		so.On("mouseUp", func(button int) {
			if renderer != nil {
				camera.Mode = view.NONE
			}
		})
		so.On("mouseMove", func(x int, y int) {
			if renderer != nil {
				mouseX = x
				mouseY = y
				if camera.Mode != view.NONE {
					x, y := getMouseNDC(windW, windH, mouseX, mouseY)
					camera.ComputeNow(mgl64.Vec2{x, y})
					if renderer.Count % 10 == 0 {
						so.Emit("frame", renderer.Frame2Base64())
					}
				}
			}
		})
		so.On("wheel", func(delta float64) {
			if renderer != nil {
				var zamt float64
				if delta < 0 {
					zamt = 1.1
				} else {
					zamt = 1 / 1.1
				}
				camera.Position[2] *= zamt

				if renderer.Count % 10 == 0 {
					so.Emit("frame", renderer.Frame2Base64())
				}
			}
		})
		so.On("disconnection", func() {
			if renderer != nil {
				renderer.Stop()
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		})
	})
	
	server.On("error", func(so socketio.Socket, e error) {
		fmt.Println("meet error:", e)
	})
	log.Println("Serving at 140.118.127.145:3001...")
	go srv.ListenAndServe()

	windW, windH = 1600, 900
	camera = view.NewCamera(float64(windW) / float64(windH))
	
	dir, _ := os.Getwd()
	model = view.LoadModel(filepath.Join(dir, "../project/", os.Args[1], "/models/model.ply"))
	renderer = &view.Renderer{Camera: camera, Model: model}
	renderer.Init(windW, windH)
	renderer.Render(windW, windH)
}