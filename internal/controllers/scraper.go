package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/labstack/echo/v4"
)

func TwitterScraper(ctx echo.Context) error {
	hateKeywords := []string{"Hindutva", "Cow vigilante", "Babri Masjid", "CAA NRC", "Bharat tere tukde honge", "Azaadi", "Azad Kashmir", "Khalistan", "Dalit", "Triple Talaq", "Saffron terror", "Jihadi", "Jihad", "Gazwa e Hind", "Godhra", "Hinduphobia", "Islamophobia", "Dictator", "Love Jihad", "Jai Bhim", "Sickular", "Cow Piss", "black lives matter", "Kafir", "anti national", "sanghi", "libtard", "woke", "rohingya", "genocide", "lynch", "kill", "tan se juda", "kashmir", "bhakt", "nazi", "fascist"}

	var allResults []map[string]interface{}
	for _, keyword := range hateKeywords {
		results := searchTweets(keyword)
		if results != nil {
			allResults = append(allResults, results...)
		}
	}

	var hateSpeechTweets []map[string]interface{}
	for _, tweet := range allResults {
		text := tweet["full_text"].(string)
		if isHateSpeech(text) {
			hateSpeechTweets = append(hateSpeechTweets, tweet)
		}
	}

	err := generatePDFReport(hateSpeechTweets)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "PDF report generated successfully"})
}

func isHateSpeech(text string) bool {
	payload := map[string]string{"text": text}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return false
	}

	resp, err := http.Post("http://localhost:8000/predict", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error making prediction request:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading prediction response body:", err)
		return false
	}

	var prediction map[string]float64
	if err := json.Unmarshal(body, &prediction); err != nil {
		fmt.Println("Error parsing prediction response:", err)
		return false
	}

	if probability, ok := prediction["hate_speech_probability"]; ok && probability > 0.70 {
		return true
	}

	return false
}

func generatePDFReport(tweets []map[string]interface{}) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hate Speech Tweets Report")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	for i, tweet := range tweets {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, fmt.Sprintf("Tweet %d:", i+1))
		pdf.Ln(8)
		for key, value := range tweet {
			pdf.SetFont("Arial", "B", 12)
			pdf.CellFormat(40, 10, fmt.Sprintf("%s:", key), "", 0, "", false, 0, "")
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(0, 10, fmt.Sprintf("%v", value), "", "", false)
		}
		pdf.Ln(10)
	}
	err := pdf.OutputFileAndClose("hate_speech_report.pdf")
	if err != nil {
		return err
	}
	return nil
}

func searchTweets(keyword string) []map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	url := "https://x.com/i/api/graphql/TQmyZ_haUqANuyBcFBLkUw/SearchTimeline"
	today := time.Now()
	sinceDate := today.AddDate(0, 0, -1)

	until := today.Format("2006-01-02")
	since := sinceDate.Format("2006-01-02")

	payload := fmt.Sprintf(`{
        "variables": {
            "rawQuery": "%s until:%s since:%s",
            "count": 20,
            "querySource": "typed_query",
            "product": "Top"
        },
        "features": {
            "rweb_tipjar_consumption_enabled": true,
            "responsive_web_graphql_exclude_directive_enabled": true,
            "verified_phone_label_enabled": true,
            "creator_subscriptions_tweet_preview_api_enabled": true,
            "responsive_web_graphql_timeline_navigation_enabled": true,
            "responsive_web_graphql_skip_user_profile_image_extensions_enabled": false,
            "communities_web_enable_tweet_community_results_fetch": true,
            "c9s_tweet_anatomy_moderator_badge_enabled": true,
            "articles_preview_enabled": true,
            "tweetypie_unmention_optimization_enabled": true,
            "responsive_web_edit_tweet_api_enabled": true,
            "graphql_is_translatable_rweb_tweet_is_translatable_enabled": true,
            "view_counts_everywhere_api_enabled": true,
            "longform_notetweets_consumption_enabled": true,
            "responsive_web_twitter_article_tweet_consumption_enabled": true,
            "tweet_awards_web_tipping_enabled": false,
            "creator_subscriptions_quote_tweet_preview_enabled": false,
            "freedom_of_speech_not_reach_fetch_enabled": true,
            "standardized_nudges_misinfo": true,
            "tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
            "rweb_video_timestamps_enabled": true,
            "longform_notetweets_rich_text_read_enabled": true,
            "longform_notetweets_inline_media_enabled": true,
            "responsive_web_enhance_cards_enabled": false
        }
    }`, keyword, until, since)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "x.com")
	req.Header.Set("Cookie", "auth_token=<auth_token>; ct0=<x_csrf_token>")
	req.Header.Set("X-Csrf-Token", "<x_csrf_token>")
	req.Header.Set("Authorization", "Bearer <bearer_token>")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return nil
	}

	entries := jsonResponse["data"].(map[string]interface{})["search_by_raw_query"].(map[string]interface{})["search_timeline"].(map[string]interface{})["timeline"].(map[string]interface{})["instructions"].([]interface{})[0].(map[string]interface{})["entries"].([]interface{})

	var results []map[string]interface{}
	for _, entry := range entries {
		content := entry.(map[string]interface{})["content"].(map[string]interface{})
		itemContent, exists := content["itemContent"]
		if exists {
			tweetResults := itemContent.(map[string]interface{})["tweet_results"].(map[string]interface{})["result"].(map[string]interface{})
			legacy := tweetResults["legacy"].(map[string]interface{})

			if fullText, ok := legacy["full_text"].(string); ok {
				user := tweetResults["core"].(map[string]interface{})["user_results"].(map[string]interface{})["result"].(map[string]interface{})
				userLegacy := user["legacy"].(map[string]interface{})

				result := map[string]interface{}{
					"full_text": fullText,
				}

				if followersCount, ok := userLegacy["followers_count"].(float64); ok {
					result["followers_count"] = followersCount
				}
				if location, ok := userLegacy["location"].(string); ok {
					result["location"] = location
				}
				if name, ok := userLegacy["name"].(string); ok {
					result["name"] = name
				}
				if possiblySensitive, ok := legacy["possibly_sensitive"].(bool); ok {
					result["possibly_sensitive"] = possiblySensitive
				}
				if verifiedPhoneStatus, ok := user["verified_phone_status"].(bool); ok {
					result["verified_phone_status"] = verifiedPhoneStatus
				}
				if conversationIdStr, ok := legacy["conversation_id_str"].(string); ok {
					result["conversation_id_str"] = conversationIdStr
				}
				if createdAt, ok := legacy["created_at"].(string); ok {
					result["created_at"] = createdAt
				}
				if mediaArray, ok := legacy["entities"].(map[string]interface{})["media"].([]interface{}); ok && len(mediaArray) > 0 {
					if expandedUrl, ok := mediaArray[0].(map[string]interface{})["expanded_url"].(string); ok {
						result["expanded_url"] = expandedUrl
					}
				}
				if userIdStr, ok := user["rest_id"].(string); ok {
					result["user_id_str"] = userIdStr
				}
				if idStr, ok := legacy["id_str"].(string); ok {
					result["id_str"] = idStr
				}
				if isBlueVerified, ok := user["is_blue_verified"].(bool); ok {
					result["is_blue_verified"] = isBlueVerified
				}
				if description, ok := userLegacy["description"].(string); ok {
					result["description"] = description
				}
				if screenName, ok := userLegacy["screen_name"].(string); ok {
					result["screen_name"] = screenName
				}

				results = append(results, result)
			}
		}
	}

	return results
}
