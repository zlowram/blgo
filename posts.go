package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday"
)

type Post struct {
	Author    string
	Date      time.Time
	Title     string
	Content   string
	Preview   string
	Template  string
	Permalink string
	Comments  bool
}

// Implement the sort.Interface for []Post by Date
type ByDate []Post

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

func loadPost(filename string) Post {
	var post Post
	const dateFormat = "01-02-2006 15:04"

	// Read the post file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Extract metadata and content
	x := strings.Split(string(content), "\n\n")
	metadata := x[0]

	// Convert the markdown to HTML
	post.Content = string(blackfriday.MarkdownCommon([]byte(strings.Trim(strings.Join(x[1:], "\n\n"), "\n"))))

	// Get the preview
	rtags := regexp.MustCompile(`<(pre|code|img).*>(.|\s)*?(</(pre|code|img)>)+`)
	stripped := rtags.ReplaceAllString(post.Content, "[...]")

	words := strings.Split(stripped, " ")
	if len(words) <= site.Config.PreviewLength {
		post.Preview = stripped
	} else {
		post.Preview = strings.Join(words[:site.Config.PreviewLength], " ") + "..."
	}

	// Process metadata
	rauthor := regexp.MustCompile(`Author: (.*)`)
	rdate := regexp.MustCompile(`Date: (.*)`)
	rtitle := regexp.MustCompile(`Title: (.*)`)
	rtemplate := regexp.MustCompile(`Template: (.*)`)
	rcomments := regexp.MustCompile(`Comments: (.*)`)

	if author := rauthor.FindStringSubmatch(metadata); author != nil {
		post.Author = author[1]
	} else {
		log.Fatal("Author not defined for post ", filename)
	}

	if date := rdate.FindStringSubmatch(metadata); date != nil {
		post.Date, _ = time.Parse(dateFormat, date[1])

		// Generate the url
		year := strconv.Itoa(post.Date.Year())
		month := strconv.Itoa(int(post.Date.Month()))
		day := strconv.Itoa(post.Date.Day())
		fname := strings.Split(strings.Split(filename, ".")[0], "/")[1]
		post.Permalink = "/" + year + "/" + month + "/" + day + "/" + fname + "/"
	} else {
		log.Fatal("Date not defined for post ", filename)
	}

	if title := rtitle.FindStringSubmatch(metadata); title != nil {
		post.Title = title[1]
	} else {
		log.Fatal("Title not defined for post ", filename)
	}

	if template := rtemplate.FindStringSubmatch(metadata); template != nil {
		post.Template = template[1]
	} else {
		log.Fatal("Template not defined for post ", filename)
	}

	if comments := rcomments.FindStringSubmatch(metadata); comments != nil {
		if comments[1] == "enabled" {
			post.Comments = true
		}
	}

	return post
}

func (p Post) convertPost() string {
	// Check if enabled comments
	htmlComments := &bytes.Buffer{}
	if p.Comments {
		data := struct {
			DisqusShortname string
			Permalink       string
		}{
			site.Config.DisqusShortname,
			p.Permalink,
		}
		template_file := site.Config.Templates + "/comments.html"
		layout, err := ioutil.ReadFile(template_file)
		if err != nil {
			log.Fatal(err)
		}
		commentsLayout := template.Must(template.New(p.Template).Parse(string(layout)))
		if err := commentsLayout.Execute(htmlComments, data); err != nil {
			log.Fatal(err)
		}
	}

	// Read the post template
	template_file := site.Config.Templates + "/" + p.Template + ".html"
	layout, err := ioutil.ReadFile(template_file)
	if err != nil {
		log.Fatal(err)
	}

	// Run the template
	htmlPost := &bytes.Buffer{}
	data := struct {
		Config   Config
		Post     Post
		Comments string
	}{
		site.Config,
		p,
		htmlComments.String(),
	}
	postLayout := template.Must(template.New(p.Template).Parse(string(layout)))
	if err := postLayout.Execute(htmlPost, data); err != nil {
		log.Fatal(err)
	}

	return htmlPost.String()
}
