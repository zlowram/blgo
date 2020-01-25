package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"
)

type site struct {
	Config config
	Posts  []post
}

var blog site

func newSite(cfg config) site {
	blog = site{Config: cfg}
	return blog
}

func (s *site) generateSite() {
	if err := os.RemoveAll(s.Config.Public); err != nil {
		log.Fatalf("error removing site folder: %v\n", err)
	}
	if err := os.Mkdir(s.Config.Public, 0755); err != nil {
		log.Fatalf("error creating site folder: %v\n", err)
	}

	if err := filepath.Walk(s.Config.Posts, generatePost(s)); err != nil {
		log.Fatalf("error creating posts html files: %v\n", err)
	}

	if err := s.copyTemplateFiles(); err != nil {
		log.Fatalf("error copying tempalte files to site directory: %v\n", err)
	}

	indexPages, err := s.generateIndex()
	if err != nil {
		log.Fatalf("error writing index pages: %v\n", err)
	}

	if err := s.writeIndexFiles(indexPages); err != nil {
		log.Fatalf("error writing index pages: %v\n", err)
	}
}

func generatePost(s *site) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		p, err := loadPost(path)
		if err != nil {
			return err
		}
		s.Posts = append(s.Posts, p)

		if err := os.MkdirAll(s.Config.Public+p.Permalink, 0755); err != nil {
			return err
		}
		postHTML, err := p.convertPost()
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(s.Config.Public+p.Permalink+"index.html", []byte(postHTML), 0755); err != nil {
			return err
		}
		return nil
	}
}

func (s *site) copyTemplateFiles() error {
	if err := copyDir(s.Config.Templates+"/css", s.Config.Public+"/css"); err != nil {
		return err
	}
	if err := copyDir(s.Config.Templates+"/js", s.Config.Public+"/js"); err != nil {
		return err
	}
	if err := copyDir(s.Config.Templates+"/fonts", s.Config.Public+"/fonts"); err != nil {
		return err
	}
	if err := copyDir(s.Config.Templates+"/images", s.Config.Public+"/images"); err != nil {
		return err
	}
	return nil
}

func (s *site) generateIndex() ([]string, error) {
	var htmlPages []string

	templateFile := s.Config.Templates + "/index.html"
	layout, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}
	indexTemplate := template.Must(template.New("index").Parse(string(layout)))

	sort.Sort(byDate(s.Posts))

	pages := int(math.Ceil(float64(len(s.Posts)) / float64(s.Config.PostsPerPage)))
	for i := 0; i < pages; i++ {
		st := 0
		ed := len(s.Posts)
		prev := ""
		next := ""

		if pages > 1 {
			st = i * s.Config.PostsPerPage
			ed = i*s.Config.PostsPerPage + s.Config.PostsPerPage

			prev = ""
			switch {
			case i > 1:
				prev = "/p/" + strconv.Itoa(i-1)
			case i == 1:
				prev = "/"
			}
			next = "/p/" + strconv.Itoa(i+1)
			if i == pages-1 {
				ed = len(s.Posts) - 1
				next = ""
			}
		}

		htmlIndex, err := generateIndexPageHTML(indexTemplate, s.Config, s.Posts[st:ed], prev, next)
		if err != nil {
			return nil, err
		}

		htmlPages = append(htmlPages, htmlIndex)
	}

	return htmlPages, nil
}

func generateIndexPageHTML(indexTemplate *template.Template, cfg config, posts []post, prev string, next string) (string, error) {
	data := struct {
		Config       config
		Posts        []post
		PreviousPage string
		NextPage     string
	}{
		cfg,
		posts,
		prev,
		next,
	}

	htmlIndex := &bytes.Buffer{}
	if err := indexTemplate.Execute(htmlIndex, data); err != nil {
		return "", err
	}

	return htmlIndex.String(), nil
}

func (s *site) writeIndexFiles(indexPages []string) error {
	for i, page := range indexPages {
		path := s.Config.Public
		if i > 0 {
			path = s.Config.Public + "/p/" + strconv.Itoa(i)
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		}

		if err := ioutil.WriteFile(path+"/index.html", []byte(page), 0755); err != nil {
			return err
		}
	}
	return nil
}
