package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"github.com/willfaught/ebnf"
)

// Future work: Use GRPC

func main() {
	const name = "ebnf-live"
	log.SetFlags(0)
	log.SetPrefix(name + ": ")
	r := chi.NewRouter()
	r.Use(httprate.Limit(
		100, // Future work: Determine best number
		time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusTooManyRequests)
			render.JSON(w, r, response{
				Error: map[string]any{
					"message": "too many requests",
				},
			})
		}),
	))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.BasicAuth(name, map[string]string{ // Future work: Replace with mTLS
		"nic":  "apple",
		"will": "banana",
	}))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.ContentCharset("utf-8"))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "not found",
			},
		})
	}))
	r.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusMethodNotAllowed)
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "method not allowed",
			},
		})
	}))
	r.Route("/v1", func(r chi.Router) {
		r.Post("/conflict", conflict)
		r.Post("/first", first)
		r.Post("/follow", follow)
		r.Post("/validate", validate)
	})
	s := &http.Server{
		Addr:    ":80",
		Handler: r,
		// Future work: Add TLS
	}
	if err := s.ListenAndServe(); err != nil { // Future work: Add TLS
		log.Fatalf("error: %v", err)
	}
}

type request struct {
	Data struct {
		Grammar string `json:"grammar"`
	} `json:"data"`
}

type response struct {
	Data  map[string]any `json:"data,omitempty"`
	Error map[string]any `json:"error,omitempty"`
}

func grammar(w http.ResponseWriter, r *http.Request) *ebnf.Grammar {
	var req request
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		if err == io.EOF {
			render.JSON(w, r, response{
				Error: map[string]any{
					"message": "cannot read request body",
				},
			})
			return nil
		}
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "cannot decode json",
				"details": err.Error(),
			},
		})
		return nil
	}
	g, err := ebnf.Parse(req.Data.Grammar)
	if err != nil {
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "cannot parse grammar",
				"details": err.Error(),
			},
		})
		return nil
	}
	return g
}

func validGrammar(w http.ResponseWriter, r *http.Request) *ebnf.Grammar {
	g := grammar(w, r)
	if g == nil {
		return nil
	}
	if err := g.Validate(); err != nil {
		render.JSON(w, r, response{
			Error: map[string]any{
				"message": "grammar is invalid",
				"details": err.Error(),
			},
		})
		return nil
	}
	return g
}

func conflict(w http.ResponseWriter, r *http.Request) {
	g := validGrammar(w, r)
	if g == nil {
		return
	}
	if err := g.Conflict(); err != nil {
		render.JSON(w, r, response{
			Data: map[string]any{
				"message": "grammar has a conflict",
				"details": err.Error(),
			},
		})
		return
	}
	render.JSON(w, r, response{
		Data: map[string]any{
			"message": "grammar has no conflicts",
		},
	})
}

func first(w http.ResponseWriter, r *http.Request) {
	g := validGrammar(w, r)
	if g == nil {
		return
	}
	anyFirstSets := g.FirstNonterminals()
	stringFirstSets := map[string][]string{}
	for nonterm, termVals := range anyFirstSets {
		termStrings := make([]string, 0, len(termVals))
		for termVal := range termVals {
			termStrings = append(termStrings, termVal.(fmt.Stringer).String())
		}
		sort.Strings(termStrings)
		stringFirstSets[nonterm] = termStrings
	}
	render.JSON(w, r, response{
		Data: map[string]any{
			"first": stringFirstSets,
		},
	})
}

func follow(w http.ResponseWriter, r *http.Request) {
	g := validGrammar(w, r)
	if g == nil {
		return
	}
	anyFollowSets := g.Follow()
	stringFollowSets := map[string][]string{}
	for nonterm, termVals := range anyFollowSets {
		termStrings := make([]string, 0, len(termVals))
		for termVal := range termVals {
			termStrings = append(termStrings, termVal.(fmt.Stringer).String())
		}
		sort.Strings(termStrings)
		stringFollowSets[nonterm] = termStrings
	}
	render.JSON(w, r, response{
		Data: map[string]any{
			"follow": stringFollowSets,
		},
	})
}

func validate(w http.ResponseWriter, r *http.Request) {
	g := grammar(w, r)
	if g == nil {
		return
	}
	if err := g.Validate(); err != nil {
		render.JSON(w, r, response{
			Data: map[string]any{
				"message": "grammar is invalid",
				"details": err.Error(),
			},
		})
		return
	}
	render.JSON(w, r, response{
		Data: map[string]any{
			"message": "grammar is valid",
		},
	})
}
