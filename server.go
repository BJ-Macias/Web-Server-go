package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var materias map[string]map[string]float32

func server() {
	materias = make(map[string]map[string]float32)

	http.HandleFunc("/", root)
	http.HandleFunc("/Califica", Califica)
	http.HandleFunc("/Promedio", Promedio)
	http.HandleFunc("/PromedioAlumno", PromedioAlumno)
	http.HandleFunc("/PromedioMateria", PromedioMateria)

	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}

func root(res http.ResponseWriter, req *http.Request) {
	cargarHtml("form.html", &res, "")
}

func cargarHtml(path string, res *http.ResponseWriter, msgs ...interface{}) {
	(*res).Header().Set(
		"Content-Type",
		"text/html",
	)

	html, _ := ioutil.ReadFile(path)

	fmt.Fprintf(*res, string(html), msgs...)
}

func Califica(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}

		var errors []string
		materia, ok := req.PostForm["materia"]
		if !ok || materia[0] == "" {
			errors = append(errors, "Materia es requerida")
		}

		alumno, ok := req.PostForm["alumno"]
		if !ok || alumno[0] == "" {
			errors = append(errors, "Materia es requerida")
		}

		califica, ok := req.PostForm["califica"]
		if !ok || califica[0] == "" {
			errors = append(errors, "Materia es requerida")
		}
		//fmt.Println(req.PostForm) 

		califfloat, err := strconv.ParseFloat(califica[0], 32)
		if err != nil {
			errors = append(errors, "Calificacion invalida")
		} else if califfloat < 0 || califfloat > 100 {
			errors = append(errors, "Fuera de rango")
		}

		mat, ok := materias[materia[0]]
		if !ok {
			mat = make(map[string]float32)
			materias[materia[0]] = mat
		}

		cal, existe := mat[alumno[0]]
		if existe {
			errors = append(errors, alumno[0]+" ya tiene una calificacion en "+materia[0]+": "+strconv.FormatFloat(float64(cal), 'f', 2, 32))
		} else {
			mat[alumno[0]] = float32(califfloat)
		}
		//fmt.Println(req.PostForm)

		if len(errors) > 0 {

			alert := ""
			for _, err := range errors {
				alert += "<p>" + err + "</p>\n"
			}
			alert += "</div>"

			cargarHtml("form.html", &res, alert)
		} else {
			cargarHtml("registrado.html", &res)
		}

	} else {
		cargarHtml("form.html", &res, "")
	}
}

func Promedio(res http.ResponseWriter, req *http.Request) {
	matcount := 0
	var total float32 = 0
	for _, v := range materias {
		alumnos := 0
		var local float32 = 0
		for _, cal := range v {
			alumnos++
			local += cal
		}
		if alumnos > 0 {
			matcount++
			total += local / float32(alumnos)
		}
	}
	var msg string
	if matcount > 0 && total > 0 {
		total = total / float32(matcount)
		msg = "El promedio general es: " + strconv.FormatFloat(float64(total), 'f', 2, 32)
	} else {
		msg = "No hay materias para promediar."
	}

	cargarHtml("promedio.html", &res, msg)
}

func PromedioAlumno(res http.ResponseWriter, req *http.Request) {
	//fmt.Println(req.PostForm)
	if req.FormValue("alumnoP") == "" {
		cargarHtml("promedioAlumno.html", &res, "Nombre invalido")
		return
	}

	var suma float32 = 0
	matcount := 0
	alumnoP := req.FormValue("alumnoP")
	for _, v := range materias {
		cal, ok := v[alumnoP]
		if ok {
			matcount++
			suma += cal
		}
	}

	var msg string
	if matcount == 0 {
		cargarHtml("promedioAlumno.html", &res, alumnoP+" no se encuentra en ninguna clase.")
	} else {
		total := suma / float32(matcount)
		msg = "El promedio de " + alumnoP + " es: " + strconv.FormatFloat(float64(total), 'f', 2, 32)
		cargarHtml("promedioAlumno.html", &res, msg)
	}
}

func PromedioMateria(res http.ResponseWriter, req *http.Request) {
	if req.FormValue("materiaP") == "" {
		cargarHtml("promedioMateria.html", &res, "Invalido")
		return
	}

	materia := req.FormValue("materiaP")
	mat, existe := materias[materia]
	if !existe {
		cargarHtml("promedioMateria.html", &res, "Materia no existe")
		return
	}

	alumnos := 0
	var total float32 = 0
	for _, cal := range mat {
		alumnos++
		total += cal
	}

	//var msg string
	if alumnos == 0 {
		cargarHtml("promedioMateria.html", &res, "Clase sin alumnos")
	} else {
		total := total / float32(alumnos)
		cargarHtml("promedioMateria.html", &res, "El promedio de la materia es: "+strconv.FormatFloat(float64(total), 'f', 2, 32))
	}
}

func main() {
	go server()

	var input string
	fmt.Scanln(&input)
}
