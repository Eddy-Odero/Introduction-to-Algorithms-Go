// Package handlers renders every page. Data is hardcoded here for now —
// there's no database until Phase 3, no auth until Phase 4. Once those
// exist, these functions swap the mock data for real queries/session
// lookups without the templates needing to change.
package handlers

import (
	"net/http"

	"forum/internal/view"
)

// Post mirrors what a template needs to render one feed card or a full
// post page. Field names match what the templates range/access.
type Post struct {
	ID       int
	Title    string
	Excerpt  string
	Body     []string // paragraphs
	Author   string
	Initials string
	Avatar   string // CSS avatar color class
	TimeAgo  string
	Tags     []string
	Likes    int
	Dislikes int
	Comments int
	HasImage bool
	LoggedIn bool // stamped from PageData.LoggedIn just before render — see stampLogin
}

type Comment struct {
	ID       int
	Author   string
	Initials string
	Avatar   string
	TimeAgo  string
	Body     string
	Likes    int
	Dislikes int
}

// PageData is the one struct every template renders from. Not every page
// uses every field — that's fine, unused fields are just zero-valued.
type PageData struct {
	Title    string
	PageCSS  string
	Minimal  bool // true for login/register: bare nav, no search/icons
	LoggedIn bool
	Username string
	Initials string
	Avatar   string
	Active   string // which nav icon is "current"

	Posts        []Post
	Post         *Post
	Comments     []Comment
	CategoryName string
	Error        string
}

func mockPosts() []Post {
	return []Post{
		{
			ID: 1, Author: "kanyi_dev", Initials: "KD", Avatar: "avatar-purple", TimeAgo: "2 hours ago",
			Title:   "Edmonds-Karp got me further than BFS ever did on lem-in",
			Excerpt: "Spent three days stuck on suboptimal paths before switching to max-flow with room-splitting. Sharing the writeup in case it saves someone else the detour…",
			Body: []string{
				"Spent three days stuck on suboptimal paths before switching to max-flow with room-splitting. BFS alone finds a path, but it doesn't guarantee the minimum number of turns when there are multiple valid routes through the colony.",
				"Once I modeled each room as two split nodes (in/out) with capacity 1, Edmonds-Karp on the resulting flow network gave me the actual optimal path count. Sharing the writeup in case it saves someone else the detour.",
			},
			Tags: []string{"technology"}, Likes: 14, Dislikes: 2, Comments: 6,
		},
		{
			ID: 2, Author: "achieng_m", Initials: "AM", Avatar: "avatar-teal", TimeAgo: "5 hours ago",
			Title:   "Anyone else's SQLite locking up under concurrent writes?",
			Excerpt: `Getting "database is locked" errors on my forum project the moment two requests hit at once. Is WAL mode enough or do I need to serialize writes myself?`,
			Body: []string{
				`Getting "database is locked" errors on my forum project the moment two requests hit at once. Is WAL mode enough or do I need to serialize writes myself?`,
			},
			Tags: []string{"technology"}, Likes: 3, Dislikes: 3, Comments: 3,
		},
		{
			ID: 3, Author: "mo_otieno", Initials: "MO", Avatar: "avatar-orange", TimeAgo: "1 day ago",
			Title:   "Unpopular opinion: tabs are objectively better than spaces",
			Excerpt: "Come at me. Tabs let everyone set their own indent width, spaces force yours on everyone else's editor.",
			Body: []string{
				"Come at me. Tabs let everyone set their own indent width, spaces force yours on everyone else's editor.",
			},
			Tags: []string{"random"}, Likes: 1, Dislikes: 3, Comments: 11,
		},
		{
			ID: 4, Author: "benja_k", Initials: "BK", Avatar: "avatar-pink", TimeAgo: "2 days ago",
			Title:   "Kisumu Bitcoin Bootcamp Hackathon — recap + lessons learned",
			Excerpt: "Our team built a Lightning-powered contact form in 36 hours. Here's what worked, what didn't, and what we'd do differently next time.",
			Body: []string{
				"Our team built a Lightning-powered contact form in 36 hours. Here's what worked, what didn't, and what we'd do differently next time.",
			},
			Tags: []string{"general", "technology"}, Likes: 4, Dislikes: 0, Comments: 2, HasImage: true,
		},
	}
}

func mockComments() []Comment {
	return []Comment{
		{ID: 1, Author: "achieng_m", Initials: "AM", Avatar: "avatar-teal", TimeAgo: "1 hour ago",
			Body: "This is exactly the trap I fell into. Did you cap the intermediate room splits at capacity 1 across the board or only for rooms with multiple neighbors?",
			Likes: 3, Dislikes: 0},
		{ID: 2, Author: "mo_otieno", Initials: "MO", Avatar: "avatar-orange", TimeAgo: "40 minutes ago",
			Body: "Would love to see the 3D visualizer you mentioned in standup, sounds like a great way to debug the ant paths visually.",
			Likes: 1, Dislikes: 1},
	}
}

// mockAuthed is a temporary, query-string-driven stand-in for real session
// checks. Until Phase 4 exists, every request is genuinely a guest — this
// just lets us preview the logged-in layout with e.g. "/?preview=1" while
// building. Delete this the moment real auth lands.
func mockAuthed(r *http.Request) bool {
	return r.URL.Query().Get("preview") == "1"
}

func withNav(d PageData, r *http.Request) PageData {
	if mockAuthed(r) {
		d.LoggedIn = true
		d.Username = "eddy_k"
		d.Initials = "EK"
		d.Avatar = "avatar-blue"
	}
	return d
}

// stampLogin copies PageData.LoggedIn onto every Post so the shared
// postcard partial (which only ever sees one Post at a time inside a
// {{range}}) knows whether to render reaction forms or read-only counts,
// without needing Go template "dict" helpers to pass extra context.
func stampLogin(d PageData) PageData {
	for i := range d.Posts {
		d.Posts[i].LoggedIn = d.LoggedIn
	}
	if d.Post != nil {
		d.Post.LoggedIn = d.LoggedIn
	}
	return d
}

// Home is "/" — public, and genuinely view-only until auth exists: with no
// session system yet, every visitor here really is a guest.
func Home(w http.ResponseWriter, r *http.Request) {
	d := withNav(PageData{Title: "forum", PageCSS: "index.css", Active: "feed", Posts: mockPosts()}, r)
	view.Render(w, http.StatusOK, "index", stampLogin(d))
}

func Login(w http.ResponseWriter, r *http.Request) {
	d := PageData{Title: "Log in — forum", PageCSS: "login.css", Minimal: true}
	view.Render(w, http.StatusOK, "login", d)
}

func Register(w http.ResponseWriter, r *http.Request) {
	d := PageData{Title: "Register — forum", PageCSS: "register.css", Minimal: true}
	view.Render(w, http.StatusOK, "register", d)
}

// NewPost assumes a logged-in user for now — Phase 4 will add the real
// redirect-to-login-if-guest middleware in front of this handler.
func NewPost(w http.ResponseWriter, r *http.Request) {
	d := withNav(PageData{Title: "New post — forum", PageCSS: "new-post.css", Active: "new-post"}, r)
	d.LoggedIn = true
	if d.Username == "" {
		d.Username, d.Initials, d.Avatar = "eddy_k", "EK", "avatar-blue"
	}
	view.Render(w, http.StatusOK, "new-post", d)
}

// Post is public — guests can read a post and its comments, just not
// react or comment (enforced by the template's LoggedIn check for now;
// real enforcement happens server-side once POST /reaction and
// POST /comment exist in later phases).
func PostDetail(w http.ResponseWriter, r *http.Request) {
	posts := mockPosts()
	p := posts[0]
	d := withNav(PageData{Title: p.Title + " — forum", PageCSS: "post.css", Post: &p, Comments: mockComments()}, r)
	view.Render(w, http.StatusOK, "post", stampLogin(d))
}

func Category(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "technology"
	}
	d := withNav(PageData{Title: "#" + name + " — forum", PageCSS: "category.css", CategoryName: name, Posts: mockPosts()}, r)
	view.Render(w, http.StatusOK, "category", stampLogin(d))
}

// Profile is inherently a logged-in-only page — Phase 4 will add the real
// redirect-to-login-if-guest check in front of this handler.
func Profile(w http.ResponseWriter, r *http.Request) {
	d := withNav(PageData{Title: "eddy_k — forum", PageCSS: "profile.css", Active: "profile", Posts: mockPosts()[3:4]}, r)
	d.LoggedIn = true
	if d.Username == "" {
		d.Username, d.Initials, d.Avatar = "eddy_k", "EK", "avatar-blue"
	}
	view.Render(w, http.StatusOK, "profile", stampLogin(d))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	view.Render(w, http.StatusNotFound, "404", PageData{Title: "404 — forum"})
}