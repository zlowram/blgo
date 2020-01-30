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

func (s *site) build() {
	if err := os.RemoveAll(s.Config.Public); err != nil {
		log.Fatalf("error removing site folder: %v\n", err)
	}
	if err := os.Mkdir(s.Config.Public, 0755); err != nil {
		log.Fatalf("error creating site folder: %v\n", err)
	}

	if err := filepath.Walk(s.Config.Posts, readPosts(s)); err != nil {
		log.Fatalf("error reading posts: %v\n", err)
	}

	if err := s.copyTemplateFiles(); err != nil {
		log.Fatalf("error copying tempalte files to site directory: %v\n", err)
	}

	indexPages, err := s.generateIndex()
	if err != nil {
		log.Fatalf("error writing index pages: %v\n", err)
	}

	if err := s.writePosts(); err != nil {
		log.Fatalf("error writing posts pages: %v\n", err)
	}

	if err := s.writeIndex(indexPages); err != nil {
		log.Fatalf("error writing index pages: %v\n", err)
	}
}

func readPosts(s *site) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		p, err := newPost(path)
		if err != nil {
			return err
		}

		s.Posts = append(s.Posts, p)

		return nil
	}
}

func (s *site) writePosts() error {
	for _, p := range s.Posts {
		if err := os.MkdirAll(s.Config.Public+p.Permalink, 0755); err != nil {
			return err
		}
		postHTML, err := p.build(s)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(s.Config.Public+p.Permalink+"index.html", []byte(postHTML), 0755); err != nil {
			return err
		}
	}
	return nil
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
	if err := copyDir(s.Config.Images, s.Config.Public+"/img"); err != nil {
		return err
	}
	return nil
}

func (s *site) generateIndex() ([]string, error) {
	var htmlPages []string

	// Read the index template
	templateFile := s.Config.Templates + "/index.html"
	layout, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}
	indexTemplate := template.Must(template.New("index").Parse(string(layout)))

	// Only posts
	var posts []post
	for _, post := range s.Posts {
		if !post.Page {
			posts = append(posts, post)
		}
	}

	sort.Sort(byDate(posts))

	// Group posts in pages
	pages := int(math.Ceil(float64(len(posts)) / float64(s.Config.PostsPerPage)))
	for i := 0; i < pages; i++ {
		st := 0
		ed := len(posts)
		prev := ""
		next := ""
		curr := ""

		if pages > 1 {
			st = i * s.Config.PostsPerPage
			ed = i*s.Config.PostsPerPage + s.Config.PostsPerPage

			prev = ""
			switch {
			case i == 0:
				curr = "/"
			case i > 1:
				prev = "/p/" + strconv.Itoa(i-1)
				curr = "/p/" + strconv.Itoa(i)
			case i == 1:
				prev = "/"
			}
			next = "/p/" + strconv.Itoa(i+1)
			if i == pages-1 {
				ed = len(posts)
				next = ""
			}
		}

		// Build the index pages
		data := struct {
			Site         *site
			Posts        []post
			CurrentPage  string
			PreviousPage string
			NextPage     string
		}{
			s,
			posts[st:ed],
			curr,
			prev,
			next,
		}

		htmlIndex := &bytes.Buffer{}
		if err := indexTemplate.Execute(htmlIndex, data); err != nil {
			return nil, err
		}

		htmlPages = append(htmlPages, htmlIndex.String())
	}

	return htmlPages, nil
}

func (s *site) writeIndex(indexPages []string) error {
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
