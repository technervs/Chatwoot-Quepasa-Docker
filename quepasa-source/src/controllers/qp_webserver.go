package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	models "github.com/sufficit/sufficit-quepasa/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	//docs "github.com/sufficit/sufficit-quepasa/docs"
	//ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	// swagger embed files
	httpSwagger "github.com/swaggo/http-swagger"
)

func QPWebServerStart() {
	r := newRouter()
	webAPIPort := os.Getenv("WEBAPIPORT")
	webAPIHost := os.Getenv("WEBAPIHOST")
	if len(webAPIPort) == 0 {
		webAPIPort = "31000"
	}

	var timeout = 30 * time.Second
	server := http.Server{
		Addr:         webAPIHost + ":" + webAPIPort,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Handler:      r,
	}

	log.Infof("Starting Web Server on Port: %s", webAPIPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func NormalizePathsToLower(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" {
			r.URL.Path = strings.ToLower(r.URL.Path)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(NormalizePathsToLower)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	shouldLog, _ := models.GetEnvBool("HTTPLOGS", false)
	if shouldLog {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// web routes
	// authenticated web routes
	r.Group(RegisterFormAuthenticatedControllers)

	// unauthenticated web routes
	r.Group(RegisterFormControllers)

	// api routes
	addAPIRoutes(r)

	// static files
	workDir, _ := os.Getwd()
	assetsDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(assetsDir))

	// Swagger Ui
	ServeSwaggerUi(r)

	// Metrics
	ServeMetrics(r)
	return r
}

func addAPIRoutes(r chi.Router) {
	r.Group(RegisterAPIControllers)
	r.Group(RegisterAPIV2Controllers)
	r.Group(RegisterAPIV3Controllers)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func ServeSwaggerUi(r chi.Router) {
	log.Debugln("Starting SwaggerUi Service")
	r.Mount("/swagger", httpSwagger.WrapHandler)
}

func ServeMetrics(r chi.Router) {
	log.Debugln("Starting Metrics Service")
	r.Handle("/metrics", promhttp.Handler())
}
