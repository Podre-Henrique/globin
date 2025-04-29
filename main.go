package main

import (
	"crypto/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/podre-Henrique/globin/ts"
)

const (
	URL_LIFETIME = 3 * 60 * 60 // 3 hours
	GC_INTERVAL  = 30 * time.Minute
	CHARSET      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	URL_LEN      = 6

	PORT = "8000"
)

type URLDTO struct {
	Original  string `json:"original,omitempty"`
	Shortened string `json:"shortened"`
}

type URL struct {
	Original string
	start    uint32
}

type DB struct {
	URLs map[string]URL
	sync.RWMutex
}

var db = DB{URLs: map[string]URL{}}

func deleteExpiredURLs() {
	db.Lock()
	defer db.Unlock()
	now := ts.Timestamp()
	var toDelete []*string

	for shortened, url := range db.URLs {
		if now-url.start >= URL_LIFETIME {
			toDelete = append(toDelete, &shortened)
		}
	}
	for _, key := range toDelete {
		delete(db.URLs, *key)
	}
}

func GC() {
	ts.StartTimeStampUpdater()
	go func() {
		tc := time.NewTicker(GC_INTERVAL)
		defer tc.Stop()
		for range tc.C {
			deleteExpiredURLs()
		}
	}()
}

func (u *URLDTO) storeURL(shortened string) bool {
	if _, exists := db.URLs[shortened]; exists {
		return false
	}
	db.URLs[shortened] = URL{u.Original, ts.Timestamp()}
	u.Shortened = shortened
	u.Original = ""
	return true
}

func (u *URLDTO) generateShortURL() {
	result := make([]byte, URL_LEN)
	db.RLock()
	defer db.RUnlock()
	for {
		for i := range result {
			randomByte := make([]byte, 1)
			_, err := rand.Read(randomByte)
			if err != nil {
				panic(err)
			}
			result[i] = CHARSET[int(randomByte[0])%len(CHARSET)]
		}
		if u.storeURL(string(result)) {
			break
		}
	}
}

func validURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func main() {
	GC()
	app := fiber.New(fiber.Config{
		// ServerHeader:                 "MY SHORTENER",
		BodyLimit:                    4096,
		Immutable:                    true,
		CaseSensitive:                true,
		Concurrency:                  333,
		ReadBufferSize:               1024,
		WriteBufferSize:              512,
		DisableKeepalive:             true,
		DisableDefaultDate:           true,
		DisableDefaultContentType:    true,
		DisablePreParseMultipartForm: true,
	})
	app.Use(limiter.New(limiter.Config{
		Max: 34,
	}))

	app.Get("/", func(c fiber.Ctx) error {
		var urlDTO URLDTO
		if err := c.Bind().Body(&urlDTO); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if !validURL(urlDTO.Original) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		urlDTO.generateShortURL()
		return c.Status(fiber.StatusOK).JSON(urlDTO)
	})

	app.Get("/:url", func(c fiber.Ctx) error {
		shortURL := c.Params("url")
		if shortURL == "" || len(shortURL) != URL_LEN {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		db.RLock()
		url, exists := db.URLs[shortURL]
		db.RUnlock()
		if !exists {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Redirect().To(url.Original)
	})
	app.Listen("localhost:" + PORT)
}
