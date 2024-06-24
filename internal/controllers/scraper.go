package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func TwitterScraper(ctx echo.Context) error {
	query := ctx.Param("query")
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
    }`, query, until, since)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "x.com")
	req.Header.Set("Cookie", "auth_token=33636fa184dc3cc6519ca1d2d54c9f30a66d609f; ct0=6acd47438d54d800e85851c5b58c24dd5b402dfc58dbc3f5ca017b9c5a2d0085c07eba33d5742af1a54f8edf88b6c20d4b2fbaf5f73b803ddeb9cc25d87f2bcf829a52a1c97ce89a0bcb2fc3508823d4")
	req.Header.Set("X-Csrf-Token", "6acd47438d54d800e85851c5b58c24dd5b402dfc58dbc3f5ca017b9c5a2d0085c07eba33d5742af1a54f8edf88b6c20d4b2fbaf5f73b803ddeb9cc25d87f2bcf829a52a1c97ce89a0bcb2fc3508823d4")
	req.Header.Set("Authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
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

	return ctx.JSON(http.StatusOK, results)
}
