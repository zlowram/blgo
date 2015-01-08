package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

	// Generate the index
	if err := ioutil.WriteFile(s.Config.Public+"/index.html", []byte(s.generateIndex()), 0755); err != nil {
		log.Fatal(err)
	}
}

func (s Site) generateIndex() string {
	// Read the index template used
	template_file := s.Config.Templates + "/index.html"
	layout, err := ioutil.ReadFile(template_file)
	if err != nil {
		log.Fatal(err)
	}

	// Run the template
	htmlIndex := &bytes.Buffer{}
	indexTemplate := template.Must(template.New("index").Parse(string(layout)))
	if err := indexTemplate.Execute(htmlIndex, s); err != nil {
		log.Fatal(err)
	}

	return htmlIndex.String()
}
