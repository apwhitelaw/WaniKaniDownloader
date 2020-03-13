package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var client *http.Client = &http.Client{
	Timeout: time.Second * 10,
}
var reqURL string = "https://api.wanikani.com/v2/subjects?types=vocabulary" // 2467 - example vocabulary
var apiKey string

func main() {

	apiKey = getAPIKey()

	var result autoGenerated = makeRequest(reqURL)

	f, err := os.Create("sentences.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// OUTPUTS SENTENCES TO FILE IN 'JA|EN' FORMAT.
	// WORKS FOR 99% OF THE SENTENCES, BUT DOES NOT ACCOUNT FOR
	// THE VERY SMALL PERCENTAGE OF SENTENCES THAT HAVE A NEWLINE
	var itemString string
	for {
		var nextPage *string = &result.Pages.NextURL
		for _, subject := range result.Data {
			if subject.Object == "vocabulary" {
				for _, item := range subject.Data.ContextSentences {
					itemString = item.Ja + "|" + item.En
					fmt.Fprintln(f, itemString)
				}
			}
		}
		fmt.Println("Done")
		if *nextPage == "" {
			break
		} else {
			result = makeRequest(*nextPage)
		}
	}
}

func getAPIKey() string {
	f, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		fmt.Println(err)
	}

	return string(f)
}

func makeRequest(url string) autoGenerated {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Wanikani-Revision", "20170710")
	req.Header.Add("Authorization", "Bearer 6e40c19d-e418-45be-a4e4-be5558272408")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var result autoGenerated
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

type autoGenerated struct {
	Object string `json:"object"`
	URL    string `json:"url"`
	Pages  struct {
		PerPage     int    `json:"per_page"`
		NextURL     string `json:"next_url"`
		PreviousURL string `json:"previous_url"`
	} `json:"pages"`
	TotalCount    int       `json:"total_count"`
	DataUpdatedAt time.Time `json:"data_updated_at"`
	Data          []struct {
		ID            int       `json:"id"`
		Object        string    `json:"object"`
		URL           string    `json:"url"`
		DataUpdatedAt time.Time `json:"data_updated_at"`
		Data          struct {
			CreatedAt   time.Time `json:"created_at"`
			Level       int       `json:"level"`
			Slug        string    `json:"slug"`
			HiddenAt    time.Time `json:"hidden_at"`
			DocumentURL string    `json:"document_url"`
			Characters  string    `json:"characters"`
			Meanings    []struct {
				Meaning        string `json:"meaning"`
				Primary        bool   `json:"primary"`
				AcceptedAnswer bool   `json:"accepted_answer"`
			} `json:"meanings"`
			Readings []struct {
				Type           string `json:"type"`
				Primary        bool   `json:"primary"`
				AcceptedAnswer bool   `json:"accepted_answer"`
				Reading        string `json:"reading"`
			} `json:"readings"`
			ComponentSubjectIds       []int  `json:"component_subject_ids"`
			AmalgamationSubjectIds    []int  `json:"amalgamation_subject_ids"`
			VisuallySimilarSubjectIds []int  `json:"visually_similar_subject_ids"`
			MeaningMnemonic           string `json:"meaning_mnemonic"`
			MeaningHint               string `json:"meaning_hint"`
			ReadingMnemonic           string `json:"reading_mnemonic"`
			ReadingHint               string `json:"reading_hint"`
			LessonPosition            int    `json:"lesson_position"`
			CharacterImages           []struct {
				URL      string `json:"url"`
				Metadata struct {
					InlineStyles bool `json:"inline_styles"`
				} `json:"metadata"`
				ContentType string `json:"content_type"`
			} `json:"character_images"`
			ContextSentences []struct {
				En string `json:"en"`
				Ja string `json:"ja"`
			} `json:"context_sentences"`
			PartsOfSpeech       []string `json:"parts_of_speech"`
			PronunciationAudios []struct {
				URL      string `json:"url"`
				Metadata struct {
					Gender           string `json:"gender"`
					SourceID         int    `json:"source_id"`
					Pronunciation    string `json:"pronunciation"`
					VoiceActorID     int    `json:"voice_actor_id"`
					VoiceActorName   string `json:"voice_actor_name"`
					VoiceDescription string `json:"voice_description"`
				} `json:"metadata"`
				ContentType string `json:"content_type"`
			} `json:"pronunciation_audios"`
		} `json:"data"`
	} `json:"data"`
}
