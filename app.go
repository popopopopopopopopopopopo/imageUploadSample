package main

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/show", showHandler)
	http.HandleFunc("/", indexHandler)
	_ = http.ListenAndServe(":8888", nil)
}

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/show.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "index"}
	renderTemplate(w, "index", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", data); err != nil {
		log.Fatalln("Enable to execute template" + err.Error())
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Allowed POST method only", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("upload")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	baseDirectory := "tmp"
	if f, err := os.Stat(baseDirectory); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(baseDirectory, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	baseFileName := uuid.New().String() + ".jpeg"
	f, err := os.Create(filepath.Join(baseDirectory, baseFileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	targetFileName = baseFileName
	http.Redirect(w, r, "/show", http.StatusFound)
}

var targetFileName = ""

func showHandler(w http.ResponseWriter, r *http.Request) {
	showHandle(w, r, targetFileName)
}

func showHandle(w http.ResponseWriter, r *http.Request, fName string) {
	baseDirectory := "tmp"
	file, err := os.Open(filepath.Join(baseDirectory, fName))
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeImageWithTemplate(w, "show", &img)
}

func writeImageWithTemplate(w http.ResponseWriter, tmpl string, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("Unable to encode image.")
	}
	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	data := map[string]interface{}{"Title": targetFileName, "Image": str}
	renderTemplate(w, tmpl, data)
}
