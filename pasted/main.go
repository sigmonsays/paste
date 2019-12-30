package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	opts := &PasteOptions{
		BindAddr: ":3555",
		DataDir:  "/tmp/paste",
	}

	flag.StringVar(&opts.BindAddr, "bindaddr", opts.BindAddr, "address to listen")
	flag.StringVar(&opts.DataDir, "data", opts.DataDir, "location of paste data")
	flag.Parse()

	os.MkdirAll(opts.DataDir, 0755)

	h := NewPasteHandler(opts)

	log.Printf("serving at %s\n", opts.BindAddr)

	err := http.ListenAndServe(opts.BindAddr, h)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}

}

func NewPasteHandler(opts *PasteOptions) *Paste {
	h := &Paste{
		PasteOptions: opts,
		rdx:          NewRadix(Digits + strings.ToLower(Alpha) + Alpha),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/paste", h.Paste)
	mux.HandleFunc("/id/", h.PasteId)
	h.mux = mux
	return h
}

type PasteOptions struct {
	BindAddr string
	DataDir  string
}

type Paste struct {
	*PasteOptions
	mux *http.ServeMux
	rdx *radix
}

func (me *Paste) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	me.mux.ServeHTTP(w, r)
}

var IndexPage = `
<html>
<head>
	<title>paste</title>
</head>
<body>
<h2>paste</h2>
client
<pre>
#!/bin/bash
paste() { curl --data-binary @- http://{{.Server}}/paste ;  }
</pre>

Usage:
<pre>
% echo whatever | paste
</pre>
<p/>

<small><a href="http://github.com/sigmonsays/paste">github</a></small>

</body>
</html>
`

type IndexContext struct {
	Server string
}

func (me *Paste) Index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index").Parse(IndexPage))

	ctx := &IndexContext{
		Server: r.Host,
	}

	err := t.Execute(w, ctx)
	if err != nil {
		log.Printf("error %s\n", err)
	}

}

func (me *Paste) Error(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) {
	w.WriteHeader(400)
	log.Printf("Error %s\n", fmt.Sprintf(s, args...))
	fmt.Fprintf(w, s, args...)
}

func (me *Paste) Paste(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		me.Error(w, r, "invalid request")
		return
	}

	idnum := int64(rand.Int31())
	id := me.rdx.Encode(idnum)

	path := filepath.Join(me.DataDir, id)
	f, err := os.Create(path)
	if err != nil {
		me.Error(w, r, "invalid request: %s", err)
		return
	}
	defer f.Close()

	written, err := io.Copy(f, r.Body)
	if err != nil {
		me.Error(w, r, "read request: %s", err)
		return
	}

	location := r.Host + fmt.Sprintf("/id/%s", id)

	w.WriteHeader(302)
	fmt.Fprintf(w, "see %s\n", location)
	h := w.Header()
	h.Set("Location", location)

	log.Printf("%s %s %s written=%d location=%s",
		r.RemoteAddr, r.Method, r.URL, written, location)

}

func (me *Paste) PasteId(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path, "/")
	if len(tmp) < 3 {
		me.Error(w, r, "no such paste id")
		return
	}
	id := tmp[2]

	path := filepath.Join(me.DataDir, id)

	f, err := os.Open(path)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "not found: %s\n", id)
		return
	}
	defer f.Close()

	// TODO: deal with content types
	w.Header().Set("Content-Type", "text/plain")

	_, err = io.Copy(w, f)
	if err != nil {
		w.WriteHeader(403)
		me.Error(w, r, "error: %s: %s\n", id, err)
		return
	}
}

var (
	Alpha  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits = "01234567890"
	Other  = "-_"
)

func NewRadix(codeset string) *radix {
	return &radix{codeset, int64(len(codeset))}
}

type radix struct {
	codeset string
	base    int64
}

func (me *radix) Encode(v int64) string {
	t := make([]byte, 0)
	id := uint64(v)
	for id > 0 {
		t = append(t, me.codeset[int64(id)%me.base])
		id = id / uint64(me.base)
	}
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
	return string(t)
}

func (me *radix) Decode(v string) int64 {
	id := int64(0)
	p := float64(len(v)) - 1
	for i := 0; i < len(v); i++ {
		id += int64(strings.Index(me.codeset, string(v[i])) * int(math.Pow(float64(len(me.codeset)), p)))
		p--
	}
	return id
}
