package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday"
)

type post struct {
	Author          string
	Date            time.Time
	Title           string
	Content         string
	Preview         string
	Page            bool
	Pinned          bool
	Template        string
	Permalink       string
	Comments        string
	commentsEnabled bool
}

type byDate []post

func (a byDate) Len() int           { return len(a) }
func (a byDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDate) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

func newPost(filename string) (post, error) {
	var p post

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	contentSplit := strings.Split(string(content), "\n\n")
	if len(contentSplit) < 2 {
		return p, errors.New("either metadata or post content is missing")
	}
	metadata := contentSplit[0]
	p.Content = string(blackfriday.MarkdownCommon([]byte(strings.Trim(strings.Join(contentSplit[1:], "\n\n"), "\n"))))
	p.Preview = getPostPreview(p.Content)

	if err := p.processMetadata(metadata, filename); err != nil {
		return p, err
	}

	return p, nil
}

func getPostPreview(content string) string {
	var preview string

	rtags := regexp.MustCompile(`<(pre|code|img).*>(.|\s)*?(</(pre|code|img)>)+`)
	stripped := rtags.ReplaceAllString(content, "[...]")

	words := strings.Split(stripped, " ")
	if len(words) <= blog.Config.PreviewLength {
		preview = stripped
	} else {
		preview = strings.Join(words[:blog.Config.PreviewLength], " ") + "..."
	}

	return preview
}

func (p *post) processMetadata(metadata string, postName string) error {
	const dateFormat = "01-02-2006 15:04"

	rauthor := regexp.MustCompile(`Author: (.*)`)
	rdate := regexp.MustCompile(`Date: (.*)`)
	rtitle := regexp.MustCompile(`Title: (.*)`)
	rpage := regexp.MustCompile(`Page: (.*)`)
	rpinned := regexp.MustCompile(`Pinned: (.*)`)
	rtemplate := regexp.MustCompile(`Template: (.*)`)
	rcommentsEnabled := regexp.MustCompile(`Comments: (.*)`)

	if author := rauthor.FindStringSubmatch(metadata); author != nil {
		p.Author = author[1]
	} else {
		return errors.New("author not defined for post " + postName)
	}

	if page := rpage.FindStringSubmatch(metadata); page != nil {
		if strings.ToLower(page[1]) == "true" {
			p.Page = true
		} else {
			p.Page = false
		}
	} else {
		p.Page = false
	}

	if date := rdate.FindStringSubmatch(metadata); date != nil {
		p.Date, _ = time.Parse(dateFormat, date[1])

		year := strconv.Itoa(p.Date.Year())
		month := strconv.Itoa(int(p.Date.Month()))
		day := strconv.Itoa(p.Date.Day())
		fname := strings.Split(strings.Split(postName, ".")[0], "/")[1]
		if p.Page {
			p.Permalink = "/etc/" + fname + "/"
		} else {
			p.Permalink = "/" + year + "/" + month + "/" + day + "/" + fname + "/"
		}
	} else {
		return errors.New("date not defined for post " + postName)
	}

	if title := rtitle.FindStringSubmatch(metadata); title != nil {
		p.Title = title[1]
	} else {
		return errors.New("title not defined for post " + postName)
	}

	if template := rtemplate.FindStringSubmatch(metadata); template != nil {
		p.Template = template[1]
	} else {
		return errors.New("template not defined for post " + postName)
	}

	if pinned := rpinned.FindStringSubmatch(metadata); pinned != nil {
		if strings.ToLower(pinned[1]) == "true" {
			p.Pinned = true
		} else {
			p.Pinned = false
		}
	} else {
		p.Pinned = false
	}

	if commentsEnabled := rcommentsEnabled.FindStringSubmatch(metadata); commentsEnabled != nil {
		if strings.ToLower(commentsEnabled[1]) == "enabled" {
			p.commentsEnabled = true
		}
	}

	return nil
}

func (p *post) build(s *site) (string, error) {
	htmlComments, err := p.buildComments()
	if err != nil {
		return "", err
	}
	p.Comments = htmlComments

	templateFile := blog.Config.Templates + "/" + p.Template + ".html"
	layout, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
	}

	htmlPost := &bytes.Buffer{}
	data := struct {
		Site *site
		Post post
	}{
		s,
		*p,
	}
	postLayout := template.Must(template.New(p.Template).Parse(string(layout)))
	if err := postLayout.Execute(htmlPost, data); err != nil {
		return "", err
	}

	return htmlPost.String(), nil
}

func (p *post) buildComments() (string, error) {
	htmlComments := &bytes.Buffer{}

	if p.commentsEnabled {
		data := struct {
			DisqusShortname string
			Permalink       string
		}{
			blog.Config.DisqusShortname,
			p.Permalink,
		}
		templateFile := blog.Config.Templates + "/comments.html"
		layout, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return "", err
		}
		commentsLayout := template.Must(template.New(p.Template).Parse(string(layout)))
		if err := commentsLayout.Execute(htmlComments, data); err != nil {
			return "", err
		}
	}

	return htmlComments.String(), nil
}
