package web

import (
	"fmt"
	"github.com/aja-video/contra/src/configuration"
	"html/template"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// webDevice holds information about configured devices we want displayed on the web UI
type webDevice struct {
	Name  string
	Model string
}

// Web holds webUI config and device information
type Web struct {
	webroot   string
	workspace string
	devices   []webDevice
	tls       bool
	cert      string
	key       string
	username  string
	password  string
	auth      auth
}

// rootHandler renders the main web template
func (w *Web) rootHandler(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/index.html"))

	if err := tmpl.Execute(res, w.devices); err != nil {
		http.Error(res, err.Error(), http.StatusNotImplemented)
		fmt.Println("Error executing template:", err)
		return
	}

}

// viewHandler displays raw file data
func (w *Web) viewHandler(res http.ResponseWriter, req *http.Request) {
	var file string
	if !strings.HasPrefix(req.URL.Path, "/view/") {
		http.Error(res, "Invalid Request", http.StatusBadRequest)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/view/") {
		// strip leading /view/
		file = req.URL.Path[6:]
	} else {
		// diffs are in their own folder
		file = req.URL.Path
	}
	data, err := ioutil.ReadFile(path.Join(w.workspace, file))
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		fmt.Println("Error loading config")
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.Write(data)
}

// downloadHandler serves config files for download
func (w *Web) downloadHandler(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, "/dl/") {
		http.Error(res, "Invalid Request", http.StatusBadRequest)
		return
	}
	// strip leading /dl/
	file := req.URL.Path[4:]
	if _, err := os.Stat(path.Join(w.workspace, file)); err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		fmt.Println("Error loading config:", file)
	}
	// Get the MIME type
	contentType := mime.TypeByExtension(filepath.Ext(file))
	// default to octet-stream if we can't find a MIME
	if contentType == "" {
		contentType = "application/octet-stream" // Default for unknown types
	}
	// Name file for download
	res.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(file))
	// Set Content-Type
	res.Header().Set("Content-Type", contentType)
	http.ServeFile(res, req, path.Join(w.workspace, file))
}

// diffHandler displays the most recent diffs for a device
func (w *Web) diffHandler(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, "/diffs/") {
		http.Error(res, "Invalid Request", http.StatusBadRequest)
		return
	}
	file := req.URL.Path
	data, err := ioutil.ReadFile(path.Join(w.workspace, file))
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		fmt.Println("Error loading diffs for config", file)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.Write(data)
}

// cssHandler serves css files
func (w *Web) cssHandler(res http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, "/css/") {
		http.Error(res, "Invalid Request", http.StatusBadRequest)
		fmt.Println("Error loading css")
		return
	}
	http.ServeFile(res, req, w.webroot+req.URL.Path)
}

// initWeb initializes the web configuration
func initWeb(c *configuration.Config) Web {
	var devices []webDevice
	// we don't need (or want) every detail of the device config, just bring over the name and type
	for _, d := range c.Devices {
		devices = append(devices, webDevice{d.Name, d.Type})
	}
	w := Web{
		webroot:   c.WebRoot,
		workspace: c.Workspace,
		devices:   devices,
	}
	w.auth = w.buildAuth(c)
	if c.WebTLS {
		w.cert = c.WebCert
		w.key = c.WebKey
		w.tls = true
	}
	return w

}

// ServeWeb runs the web server and handles client requests
func ServeWeb(c *configuration.Config) {
	w := initWeb(c)
	if w.auth == nil {
		log.Println("Refusing to launch web UI without authentication")
		return
	}

	http.HandleFunc("/", w.auth.authenticate(w.rootHandler))
	http.HandleFunc("/view/", w.auth.authenticate(w.viewHandler))
	http.HandleFunc("/diffs/", w.auth.authenticate(w.diffHandler))
	http.HandleFunc("/dl/", w.auth.authenticate(w.downloadHandler))
	http.HandleFunc("/css/", w.auth.authenticate(w.cssHandler))

	log.Printf("Starting web server at :%s", c.WebPort)
	if w.tls {
		log.Fatal(http.ListenAndServeTLS(":"+c.WebPort, w.cert, w.key, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+c.WebPort, nil))
	}
}
