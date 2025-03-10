package main

import (
	"fmt"
	"net"
	"net/http"
	"text/template"
)

const BUCKET_A_NAME = "bucket-a"
const BUCKET_B_NAME = "bucket-b"

const TEMPLATE = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>AB test example</title>
  </head>
  <body>
    <h1>AB test example</h1>
    <p>Bucket: {{.bucket}}</p>
  </body>
</html>
`

var tmpl = template.Must(template.New("index").Parse(TEMPLATE))

func handler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "userip: %q is not IP:port\n", r.RemoteAddr)
		return
	}
	lastChar := ip[len(ip)-1]

	// This sort of modulo thing is just an example of how users could be assigned.
	// One benefit is that users in the same household will likely have a shared experience
	// in case it comes up in conversation.
	switch bucket := lastChar % 2; bucket {
	case 0:
		tmpl.Execute(w, map[string]string{"bucket": BUCKET_A_NAME})
	case 1:
		tmpl.Execute(w, map[string]string{"bucket": BUCKET_B_NAME})
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
