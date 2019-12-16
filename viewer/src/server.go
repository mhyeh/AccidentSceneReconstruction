package main

import (
	"fmt"
	"log"
	"net/http"
	"context"
	"time"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	
	"github.com/go-gl/mathgl/mgl64"
	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"

	"github.com/nareix/joy4/format/rtmp"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/codec/h264parser"

	"view"
)

func getMouseNDC(w int, h int, mx int, my int, x *float64, y *float64) {
	*x = float64(mx) / float64(w) * 2 - 1
	*y = float64(my) / float64(h) * 2 - 1
}

func int2EG(n int) string {
	s := strconv.FormatInt(int64(n + 1), 2)
	return strings.Repeat("0", len(s) - 1) + s
}

func string2Bytes(b string) []byte {
    var out []byte
    var str string

    for i := len(b); i > 0; i -= 8 {
        if i-8 < 0 {
            str = string(b[0:i])
        } else {
            str = string(b[i-8 : i])
        }
        v, err := strconv.ParseUint(str, 2, 8)
        if err != nil {
            panic(err)
        }
        out = append([]byte{byte(v)}, out...)
    }
    return out
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	
	// go server.Serve()
	// defer server.Close()

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
		// s.SetContext("")
		fmt.Println("connected")
		so.Emit("connected")

		so.On("init", func(w int, h int) {
			windW, windH = w, h
			camera = view.NewCamera(float64(windW) / float64(windH))
			
			dir, _ := os.Getwd()
			model = view.LoadModel(filepath.Join(dir, "../project/", os.Args[1], "/models/model.ply"))
			renderer = &view.Renderer{Camera: camera, Model: model}
			renderer.Init(windW, windH)
			
			conn, _ := rtmp.Dial("rtmp://localhost/model/" + os.Args[1])
			// streams := []av.CodecData{h264parser.CodecData{
			// 	RecordInfo: h264parser.AVCDecoderConfRecord{
			// 		AVCProfileIndication: 77,
			// 		ProfileCompatibility: 0,
			// 		AVCLevelIndication: 41,
			// 		LengthSizeMinusOne: 3,
			// 	},
			// 	SPSInfo: h264parser.SPSInfo{
			// 		ProfileIdc: 77,
			// 		LevelIdc: 41,
			// 		MbWidth: 120,
			// 		MbHeight: 67,
			// 		Width: uint(windW),
			// 		Height: uint(windH),
			// 	},
			// }}
			mbW := w / 16 - 1
			mbH := h / 16 - 1
			binstr := "01100111010011010000000000101001111100010110" + int2EG(mbW) + int2EG(mbH) + "10001"
			if len(binstr) % 8 == 0 {
				binstr += strings.Repeat("0", 8 - len(binstr) % 8)
			}
			// SPS, _ := hex.DecodeString("674D0029F16")
			SPS := string2Bytes(binstr)

			binstr = "0110100011001110001110001"
			if len(binstr) % 8 == 0 {
				binstr += strings.Repeat("0", 8 - len(binstr) % 8)
			}
			PPS := string2Bytes(binstr)
			codecData, _ := h264parser.NewCodecDataFromSPSAndPPS(SPS, PPS)
			streams := []av.CodecData{codecData}
			conn.WriteHeader(streams)
			go renderer.Streaming(w, h, conn)
		})
		so.On("mouseDown", func(button int) {
			if renderer != nil {
				if button == 0 {
					camera.Mode = view.ROTATE
				} else if button == 2 {
					camera.Mode = view.PAN
				}
				var x, y float64
				getMouseNDC(windW, windH, mouseX, mouseY, &x, &y)
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
					var x, y float64
					getMouseNDC(windW, windH, mouseX, mouseY, &x, &y)
					camera.ComputeNow(mgl64.Vec2{x, y})
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
			}
		})
		so.On("disconnection", func() {
			renderer.Stop()
			ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		})
	})
	


	server.On("error", func(so socketio.Socket, e error) {
		fmt.Println("meet error:", e)
	})

	log.Println("Serving at 140.118.127.145:3001...")
	log.Fatal(srv.ListenAndServe())

	
}