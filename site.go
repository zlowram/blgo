package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"text/template"
)

type Site struct {
	Config Config
	Posts  []Post
}

var site Site

func newSite(cfg Config) Site {
	site = Site{Config: cfg}
	return site
}

func (s Site) generateSite() {
	// Delete the site folder and create it again
	if err := os.RemoveAll(s.Config.Public); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir(s.Config.Public, 0755); err != nil {
		log.Fatal(err)
	}

	// Get all the posts and convert them
	dirlist, err := ioutil.ReadDir(s.Config.Posts)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dirlist {
		// load the post
		p := loadPost(s.Config.Posts + "/" + f.Name())

		// Add the post to post slice
		s.Posts = append(s.Posts, p)

		// create the path
		if err := os.MkdirAll(s.Config.Public+p.Permalink, 0755); err != nil {
			log.Fatal(err)
		}

		// convert & write
		if err := ioutil.WriteFile(s.Config.Public+p.Permalink+"index.html", []byte(p.convertPost()), 0755); err != nil {
			log.Fatal(err)
		}
	}

	// Sort the posts by date
	sort.Sort(ByDate(s.Posts))

	// Copy all necessary to the site folder (css, js and assets)
	if err := copyDir(s.Config.Templates+"/css", s.Config.Public+"/css"); err != nil {
		log.Fatal(err)
	}
	if err := copyDir(s.Config.Templates+"/js", s.Config.Public+"/js"); err != nil {
		log.Fatal(err)
	}
	if err := copyDir(s.Config.Templates+"/fonts", s.Config.Public+"/fonts"); err != nil {
		log.Fatal(err)
	}
	if err := copyDir(s.Config.Templates+"/images", s.Config.Public+"/images"); err != nil {
		log.Fatal(err)
	}

	// Generate the index main page
	index := s.generateIndex()
	if err := ioutil.WriteFile(s.Config.Public+"/index.html", []byte(index[0]), 0755); err != nil {
		log.Fatal(err)
	}

	index = index[1:len(index)]

	// Now the rest of the pages (if existent)
	if len(index) <= 0 {
		return
	}

	for i, page := range index {
		path := s.Config.Public + "/p/" + strconv.Itoa(i+1)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(path+"/index.html", []byte(page), 0755); err != nil {
			log.Fatal(err)
		}
	}
}

func (s Site) generateIndex() []string {
	var htmlPages []string

	// Read the index template used and parse it
	template_file := s.Config.Templates + "/index.html"
	layout, err := ioutil.ReadFile(template_file)
	if err != nil {
		log.Fatal(err)
	}
	indexTemplate := template.Must(template.New("index").Parse(string(layout)))

	// Pagination
	pages := int(math.Ceil(float64(len(s.Posts)) / float64(s.Config.PostsPerPage)))

	for i := 0; i < pages; i++ {
		var st, ed int
		var next, prev string

		posts := i*s.Config.PostsPerPage + s.Config.PostsPerPage
		if posts >= len(s.Posts) {
			st = i * s.Config.PostsPerPage
			ed = len(s.Posts)
		} else {
			st = i * s.Config.PostsPerPage
			ed = i*s.Config.PostsPerPage + s.Config.PostsPerPage
		}

		next = "/p/" + strconv.Itoa(i+1)
		prev = "/p/" + strconv.Itoa(i-1)

		switch {
		// First page
		case i <= 0:
			prev = ""
		// Second page when = 2 pages
		case i-1 == 0 && i >= pages-1:
			prev = "/"
			next = ""
		// Second page
		case i-1 == 0:
			prev = "/"
		// Last page
		case i >= pages-1:
			next = ""
		}

		data := struct {
			Config       Config
			Posts        []Post
			PreviousPage string
			NextPage     string
		}{
			s.Config,
			s.Posts[st:ed],
			prev,
			next,
		}

		// Run the template
		htmlIndex := &bytes.Buffer{}
		if err := indexTemplate.Execute(htmlIndex, data); err != nil {
			log.Fatal(err)
		}

		htmlPages = append(htmlPages, htmlIndex.String())

	}

	return htmlPages
}
