package main

import (
	"strings"
	"html/template"
	"log"
	"net/http"
	"github.com/arulthayalan/cyoa/cyoa"
	"os"
	"flag"
	"fmt"
)


func main() {
	port :=  flag.Int("port", 3000, "the port to start the CYOA web application")
	filename := flag.String("fileName", "../resource/gophoer.json", "JSON file with CYOA story")
	flag.Parse()
	fmt.Printf("Using the story file %s\n", *filename)

	file, err := os.Open(*filename)

	if (err != nil) {
		fmt.Errorf("%s: %v", *filename, err)
	}

	story, err := cyoa.JsonStory(file)
	if err != nil {
		fmt.Errorf("%s: %v", "Unable to parse story json", err)
	}

	tpl := template.Must(template.New("").Parse(storyTemplate))

	h := cyoa.NewHandler(story,
		cyoa.WithTemplate(tpl),
		cyoa.WithPathFunc(pathFn))

	fmt.Printf("Starting the server on port %d\n", *port)

	mux := http.NewServeMux()
	mux.Handle("/story/", h)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

	//%+v prints out fieldname as well
	//fmt.Printf("%+v\n", story)]
}

//setting context path like
func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	// "/story/intro" ==> "intro"
	path = path[len("/story/"):]
	return path
}

var storyTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`