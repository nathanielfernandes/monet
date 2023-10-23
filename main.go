package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	rl "github.com/nathanielfernandes/rl"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/:payload", view)
	router.POST("/:author/:message/:payload", send)

	fmt.Println("Listening on port 80")
	if err := http.ListenAndServe("0.0.0.0:80", router); err != nil {
		log.Fatal(err)
	}
}

func view(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	payload := ps.ByName("payload")

	im, err := PayloadToImage(payload)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "image/png")
	// cors headers

	// cache in the browser for 1 year
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write(im.Bytes())
}

var RLM = rl.NewRatelimitManager(1, 15*1000)

func send(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	iden := r.Header.Get("X-Forwarded-For")
	if iden == "" {
		iden = r.RemoteAddr
	}

	if RLM.IsRatelimited(iden) {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	author := ps.ByName("author")
	message := ps.ByName("message")

	payload := ps.ByName("payload")

	if author == "" {
		author = "anonymous"
	}
	// truncate author to 32 characters
	if len(author) > 32 {
		author = author[:32]
	}

	fmt.Println("new request from", author)

	// truncate message to 256 characters
	if len(message) > 256 {
		message = message[:256]
	}

	raw_payload := DiscPayload{
		Content: "",
		Embeds: []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Color       int    `json:"color"`
			Image       struct {
				URL string `json:"url"`
			} `json:"image"`
		}{
			{
				Title:       author,
				Description: message,
				URL:         "https://monet.b-cdn.net/" + payload,
				Color:       16748546,
				Image: struct {
					URL string `json:"url"`
				}{
					URL: "https://monet.b-cdn.net/" + payload,
				},
			},
		},
		Username:  "monet",
		AvatarURL: "https://shop.ariustechnology.com/cdn/shop/articles/Claude-Monet-San-Giorgio-Maggiore-at-Dusk_600x.jpg",
	}

	err := SendWebhook(raw_payload)
	if err != nil {
		fmt.Println("Error sending webhook:")
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

type DiscPayload struct {
	Content string `json:"content"`
	Embeds  []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Color       int    `json:"color"`
		Image       struct {
			URL string `json:"url"`
		} `json:"image"`
	} `json:"embeds"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

var WEBHOOK_URL = os.Getenv("WEBHOOK_URL")

func SendWebhook(payload DiscPayload) error {
	// convert payload to json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// send payload to webhook
	resp, err := http.Post(WEBHOOK_URL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("webhook returned status code %d", resp.StatusCode)
	}

	return nil
}
