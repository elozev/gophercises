package cyoa

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type StoryHandler struct {
	Stories Story
}

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Text    []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

var styles = `
	<style>
		html {
			font-family: Montserrat, sans-serif;
		}

		.content {
			max-width: 750px;
			text-align: center;
		}

		.text {
			text-transform: italic;
		}

		body {
			display: flex;
			justify-content: center;
		}

		a {
			display: inline-block;
			margin: 10px;
			padding: 5px 15px;
			background: lightblue;
			color: white;
			border-radius: 4px;
		}
	</style>
`

func (sh *StoryHandler) StoryHandler(w http.ResponseWriter, r *http.Request) {
	chapter := r.PathValue("chapter")

	var tpl = fmt.Sprintf(`
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
		<div class="content">
			<h1>{{ .Title }}</h1>
			<br/>

			<p class="text">{{ .Text }}</p>

			{{ range $opt := .Options }}
				<p class="text">{{ $opt.Text }}</p>
				<a href="{{ $opt.Chapter}}">Go to {{ $opt.Chapter}}!</a>
			{{ else }}
				<h5>Oh, it seems you have to head back to the start</h5>
				<a href="/">Start over!</a>
			{{ end }}
		</div>
	</body>
	%s
</html>
`, styles)

	var story Chapter
	var ok bool

	if story, ok = sh.Stories[chapter]; !ok {
		log.Printf("story %s not found; redirecting to home page\n", story)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	t, err := template.New("story").Parse(tpl)
	Check(err)

	err = t.Execute(w, story)
	Check(err)
}

func (sh *StoryHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = fmt.Sprintf(`
<html>
	<head>
		<meta charset="UTF-8">
		<title>Choose Your Own Story</title>
	</head>
	<body>
		<div class="content">
			<h1>{{.Title}}</h1>
			<a href="/story/intro">Start your adventure</a>
		</div>
	</body>
	%s
</html>
`, styles)

	t, err := template.New("index").Parse(tpl)
	Check(err)
	data := struct {
		Title string
	}{
		Title: "My Story Time",
	}

	w.WriteHeader(http.StatusOK)
	err = t.Execute(w, data)
	Check(err)
}
