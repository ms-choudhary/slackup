package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ms-choudhary/slackup/pkg/store"
)

type Server struct {
	Store *store.Store
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not found: %v", r)
}

func (s *Server) error(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "Internal Error: %#v", err)
}

func (s *Server) write(statusCode int, object interface{}, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	output, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		s.error(err, w)
		return
	}
	w.Write(output)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.RequestURI)
	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		server.error(err, w)
		return
	}

	requestParts := strings.Split(u.Path, "/")[1:]
	//log.Printf("url path: %s", u.Path)
	//log.Printf("len req parts: %v %d", requestParts, len(requestParts))

	if len(requestParts) != 2 {
		server.notFound(w, req)
		return
	}

	project, channel := requestParts[0], requestParts[1]
	channelId, err := server.Store.GetChannel(project, channel)
	if err != nil {
		server.error(err, w)
	}

	switch req.Method {
	case "GET":
		q, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			server.error(err, w)
		}

		messages, err := server.Store.Query(channelId,
			store.Filter{
				User: q["user"][0],
				Text: q["text"][0]})

		if err != nil {
			server.error(err, w)
		}

		log.Printf("got messages: %v", messages)
		server.write(200, messages, w)
	default:
		server.notFound(w, req)
	}
}
