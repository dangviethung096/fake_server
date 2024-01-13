package route

import (
	"html/template"
	"net/http"
)

type PageVariables struct {
	Title   string
	Message string
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	HomePageVars := PageVariables{ //store the page variables in a struct
		Title:   "Homepage",
		Message: "Welcome to the Homepage!",
	}

	t, err := template.ParseFiles("./render_html/index.html") //parse the html file homepage.html
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
