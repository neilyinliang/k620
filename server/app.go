package server

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/neilyinliang/k620/global"
	"github.com/neilyinliang/k620/public"
)

type App struct {
	cfg            *global.Config
	uid            string
	svr            *http.Server
	exitSignal     chan os.Signal
	bufferPool     *sync.Pool
	upGrader       *websocket.Upgrader
	dialer         *net.Dialer
	specialDomains []string
}

func (app *App) httpSvr() {
	mux := http.NewServeMux()
	wsPath := "/" + app.uid
	mux.HandleFunc(wsPath, app.WsVLESS)

	mime.AddExtensionType(".js", "application/javascript")
	content, _ := fs.Sub(public.Public, "dist")
	fileServer := http.FileServer(http.FS(content))
	mux.Handle("/", http.StripPrefix("/", fileServer))

	server := &http.Server{
		Addr:         app.cfg.ListenAddr(),
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app.svr = server
}

func NewApp(c *global.Config, sig chan os.Signal) *App {
	bufferSize := c.GetBufferSize()
	app := &App{
		cfg:        c,
		uid:        c.AllowUsers,
		exitSignal: sig,
		svr:        nil,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, bufferSize)
			},
		},
		upGrader: &websocket.Upgrader{
			HandshakeTimeout: 2 * time.Second,
			ReadBufferSize:   bufferSize,
			WriteBufferSize:  bufferSize,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all connections by default
				return true
			},
		},
		dialer: &net.Dialer{
			Timeout:   5 * time.Second,
			DualStack: true,
			Resolver: &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: time.Millisecond * time.Duration(5000),
					}
					return d.DialContext(ctx, "udp", "8.8.8.8:53")
				},
			},
		},
		specialDomains: func() []string {
			if domainsStr := os.Getenv("DOMAINS"); domainsStr != "" {
				log.Printf("DOMAINS: %s", domainsStr)
				return strings.Split(domainsStr, ",")
			}
			return []string{}
		}(),
		/**
		specialDomains: []string{
			"youtube",
			"facebook",
			"telegram",
		},
		*/
	}

	app.httpSvr()
	go app.loopPush()
	return app
}

func (app *App) Run() {
	log.Println("server starting on http://", app.cfg.ListenAddr())
	if err := app.svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not listen on %s: %v\n", app.cfg.ListenAddr(), err)
	}
}

func (app *App) Shutdown(ctx context.Context) {
	log.Println("Shutting down the server...")
	if err := app.svr.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

func (app *App) loopPush() {
	tk := time.NewTicker(app.cfg.PushInterval())
	defer tk.Stop()
	for {
		select {
		case sig := <-app.exitSignal:
			app.exitSignal <- sig
			//app.PushNode() //last push
			return
		case <-tk.C:
			//app.PushNode()
		}
	}
}

func (app *App) IsUserNotAllowed(uuid string) (isNotAllowed bool) {
	return uuid != app.uid
}
