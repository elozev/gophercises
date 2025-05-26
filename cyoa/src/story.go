package cyoa

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
			font-style: italic;
		}

		body {
			display: flex;
			justify-content: center;
		}

		a {
			display: inline-block;
			margin: 10px;
			padding: 10px 15px;
			background: coral;
			color: white;
			border-radius: 4px;

			transition: border-radius 0.2s;
		}

		a:hover {
			border-radius: 16px;
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
			{{ range $paragraph := .Text }}
			<p class="text">{{ $paragraph }}</p>
			{{ end }}

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

func printChapter(chapter Chapter) {
	fmt.Printf("-----\n%s\n", chapter.Title)

	for _, v := range chapter.Text {
		fmt.Println(v)
	}

	if len(chapter.Options) > 0 {
		fmt.Printf("\n--\nContinue by selecting:\n")
	} else {
		fmt.Println("Game over")
	}

	for _, opt := range chapter.Options {
		fmt.Printf("[%s] %s\n", opt.Chapter, opt.Text)
	}
}

func (sh *StoryHandler) CliHandler() {

	var story = sh.Stories["intro"]

main:
	for true {
		printChapter(story)
		var userInput string
		valid := false

	out:
		for !valid {
			if len(story.Options) == 0 {
				break main
			}

			fmt.Printf("Enter your choice: ")
			fmt.Scanf("%s\n", &userInput)

			for _, opt := range story.Options {
				if opt.Chapter == strings.TrimSpace(userInput) {
					valid = true
					break out
				}
			}
			fmt.Printf("%s is not a valid option; try again\n", userInput)
		}

		nextStory := sh.Stories[userInput]
		story = nextStory
	}
}

func JsonStory(r io.Reader) (Story, error) {
	var storiesHolder Story
	d := json.NewDecoder(r)

	err := d.Decode(&storiesHolder)
	if err != nil {
		return nil, err
	}

	return storiesHolder, nil
}
