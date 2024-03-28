package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Spotify API Anahtarınız
const token = "YOUR_SPOTIFY_API_KEY"

// API isteği yapma fonksiyonu
func fetchWebAPI(endpoint, method string, body interface{}) map[string]interface{} {
	client := &http.Client{}

	req, err := http.NewRequest(method, fmt.Sprintf("https://api.spotify.com/%s", endpoint), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	return data
}

// En iyi şarkıları almak için API isteği yapma fonksiyonu
func getTopTracks(limit int) map[string]interface{} {
	endpoint := fmt.Sprintf("v1/me/top/tracks?time_range=short_term&limit=%d", limit)
	return fetchWebAPI(endpoint, "GET", nil)
}

// Rastgele şarkıları almak için fonksiyon
func getRandomTracks(limit int) []string {
	rand.Seed(time.Now().UnixNano())

	topTracks := getTopTracks(limit)
	items := topTracks["items"].([]interface{})
	randomIndices := rand.Perm(len(items))

	var randomSongs []string
	for i := 0; i < limit; i++ {
		track := items[randomIndices[i]].(map[string]interface{})
		var artists []string
		for _, artist := range track["artists"].([]interface{}) {
			artists = append(artists, artist.(map[string]interface{})["name"].(string))
		}
		randomSongs = append(randomSongs, fmt.Sprintf("%s by %s", track["name"].(string), artists))
	}

	return randomSongs
}

func main() {
	numTracks := 15 // Almak istediğiniz şarkı sayısı
	randomSongs := getRandomTracks(numTracks)

	fmt.Println("Önerilen Şarkılar:")
	for i, song := range randomSongs {
		fmt.Printf("%d. %s\n", i+1, song)
	}
}
