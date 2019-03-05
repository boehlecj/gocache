//
// cache.go
//
// A basic cache written in Go.
// This is an exercise where we will implement a simple key value store written in Go.
// We will use a simple readline interface and two commands: PUT and GET.
//
// Requirements:
//
// 1. PUT key value     Set a value in the cache.
// 2. GET key           Get a value stored in the cache.
// 3. EXIT/QUIT         Exits the interactive prompt (can also be done with Ctrl-d thanks to the readline pkg).
// 4. Use only packages from the stdlib (except for the readline package already imported below).
//
package main

import (
	"io"
	"log"
	"sync"
	"strings"
	"time"
	"github.com/chzyer/readline"
)

const MaxCached int = 256

type Action int

const (
	GET  Action = iota
	PUT
	EXIT
	NONE
)

type Cache struct {
	mu sync.Mutex
	entries map[string]*Entry
}

type Entry struct {
	data []byte
	size int
	modified time.Time
}

func getAction(action string) Action {
	a := NONE
	switch {
	case action == "get":
		a = GET
	case action == "put":
		a = PUT
	case action == "quit" || action == "exit":
		a = EXIT
	}

	return a
}

func NewCache() *Cache {
	return &Cache{entries: make(map[string]*Entry)}
}

func (c *Cache) removeOldestEntry() {
	time := time.Now()
	key := ""

	for name, ent := range c.entries {
		if key == "" {
			time = ent.modified
			key = name
		} else if ent.modified.Before(time) {
			time = ent.modified
			key = name
		}
	}
	if key != "" {
		c.delete(key)
	}
}

func (c *Cache) delete(k string) {
	_, found := c.get(k)
	if found {
		c.mu.Lock()
		delete(c.entries, k)
		c.mu.Unlock()
	}
}

func (c *Cache) put(k string, v string) {
	if(len(c.entries) == MaxCached) {
		c.removeOldestEntry()
	}

	entry := &Entry{
		data: []byte(v),
		size: len(v),
		modified: time.Now(),
	}

	c.mu.Lock()
	c.entries[k] = entry
	c.mu.Unlock()
}

func (c *Cache) get(k string) (entry *Entry, found bool) {
	c.mu.Lock()
	entry, ok := c.entries[k]

	if !ok {
		c.mu.Unlock()
		log.Println("No cached entry for", k)
		return nil, false
	}

	c.mu.Unlock()
	return entry, true
}

func main() {
	cache := NewCache()
	prompt, err := readline.New("> ")
	if err != nil {
		log.Fatal(err)
	}
	defer prompt.Close()

	for {
		line, err := prompt.Readline()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)
		args := strings.Fields(line)
		action := getAction(strings.ToLower(args[0]))

		if action == EXIT {
			log.Println("Bye!")
			break
		}

		switch {
		case len(args) == 2 && action == GET:
			response, ok := cache.get(args[1])
			if(ok) {
				log.Println(string(response.data[:]))
			}
		case len(args) == 3 && action == PUT:
			cache.put(args[1],args[2])
		default:
			log.Println("Why did you type", line, "?")
		}
		prompt.SetPrompt("> ")
	}
}
