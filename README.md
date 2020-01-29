blgo
====

A minimal and static blog engine written in Go.

Overview
--------

Blgo is a static blog engine written in Go, aimed at minimalism and flexibility.

Blgo takes a directory with content written in Markdown and another one with templates
and renders a full static website.

Getting started
---------------

### Install

```
go get github.com/zlowram/blgo
```

### Creating your site

Blgo needs to know where the directories with the content and templates are, as well as some parameters, and this
can be done via the configuration file, which is in [toml](https://github.com/toml-lang/toml) format.

An example of Blgo configuration file is the following:

``` Markdown
# Site parameters
Title = "blgo" # Title of your site
Description = "A static blog engine written in Go" # Description for your site
Keywords = "blog, Go" # Keywords for your site
PreviewLength = 30 # Length of the post preview in the index (number of words)
PostsPerPage = 1 # Index posts pagination

# Disqus parameters
DisqusShortname = "blgo"

# GoogleAnalytics
GoogleAnalyticsID = "UA-XXXXXXXXX-X"

# Directories
Posts = "posts" # Directory containing the posts
Templates = "templates" # Directory containing the templates
Public = "static" # Directory where the generated static site will be stored
```

### Writing posts

In blgo each post is a plaintext file with the following structure:

``` Markdown
Author: zlowram
Date: 01-08-2015 00:45
Title: Hello, world!
Template: post
Comments: enabled
Page: false
Pinned : false


Here starts the content of the post.
```

Posts content is created using the [Markdown syntax](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet), and file must be placed in the posts directory.

It is recommended to name the file with the title of the post replacing spaces with hyphens, because it will be included in the post's permalink. The permalink of the example post above, if the file is named "test-post.md" would be:

    http://example.com/2015/01/08/test-post/

The values `Page` and `Pinned` from the metadata are optional, and their meaning is:

* Page: Defines whether the file will be treated as a page or not. A page is different from a post because are not indexed.
* Pinned: Defines whether the post or page will be pinned or not. This can be then leveraged in the template to create navigation bars, starred post section, etc. (check _example_site)

## Writing templates

Blgo has support for templates so, if you have web-design skills, you may want to write your own templates. Then, good news for you: writing themes for blgo is stupidly simple! You need to create a file template for the index and for the posts, and place it in the templates directory.

The template directory structure is the following:

```
templates/
├── css/
├── fonts/
├── images/
├── index.html
├── js/
└── post.html
```

A template is an HTML file that contains [Go template syntax](http://golang.org/pkg/text/template/) at the concrete places where you want the content and information to appear.

The data struct passed to the template is different for the index and the posts:

#### Post

``` Go
 data := struct {
		Site *site
		Post post
 }{
		s,
		*p,
 }

 type site struct {
	Config config
	Posts  []post
 }

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
 }

 type Config struct {
 	Title         string
 	Description   string
 	Posts         string
 	Templates     string
 	Public        string
 	PreviewLength int
 }
```

 Examples:
  * Print the title of the post: {{.Post.Title}}
  * Print the Date: {{.Post.Date.Month}} {{.Post.Date.Day}}, {{.Post.Date.Year}}
  * Print the title of the site: {{.Config.Title}}


#### Index

``` Go
 type Site struct {
 	Config Config
 	Posts  []Post
 }

 data := struct {
     Config       Config
     Posts        []Post
     PreviousPage string
     NextPage     string
 }{
     s.Config,
     posts[st:ed],
     prev,
     next,
 }
```

 Examples:
  * Iterate over the posts and print the title of the post: {{range .Posts}} {{.Title}} {{end}}
  * Check if next page exists and print the path to it: {{if .NextPage}} {{.NextPage}} {{end}}

Whether you don't want to make a template, or you created one and would like to share it, check the [official templates repository](http://github.com/zlowram/blgo-templates)!

### Using Google Analytics

Blgo currently has built-in support for Google Analytics. However, the Template that you use must also support it. If the template does not support Google Analytics,
it is really easy to do so, just add the following script tag within the <head> </head> tags in each template page:

``` Html
<script>
	(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
	 (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
	 m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
	 })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

	ga('create', '{{.Config.GoogleAnalyticsID}}', 'auto');
	ga('send', 'pageview');
</script>

```

### Deploying your site

The easiest way to deploy your site generated with blgo, is to copy the contents of the public directory in the root folder of your favorite web server.

Blgo is a command-line tool, so the deploy process can be automated in different ways (Makefiles, Git Hooks, etc.).If you have a cool deploying method, let us know!


Current features
----------------

* Index pagination
* Comments via disqus
* Support for GoogleAnalytics
* Support for pages.
* Support for pinned posts / pages.

Future features
---------------

* Tags for posts.
* Permalink customization.
* Etc.

Sites using blgo
----------------

If you are using blgo for your site, we would like to know! You can either contact me via [twitter](http://twitter.com/zlowram_) or you can send a pull-request adding your URL in the following list:

* [zlowr.am](http://zlowr.am) (zlowram's site)
* your site here

Contributing
------------
If you would like to contribute to this project, you can do it in several ways:

* Opening an Issue to report a bug, suggest a feature, whatever.
* Sending a pull-request with your changes.

In both cases it will be reviewed and studied if it should be merged or not in the project.
