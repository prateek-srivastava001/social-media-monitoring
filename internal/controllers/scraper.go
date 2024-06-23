package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TwitterScraper(ctx echo.Context) error {
	url := "https://x.com/i/api/graphql/TQmyZ_haUqANuyBcFBLkUw/SearchTimeline"

	payload := []byte(`{"variables":{"rawQuery":"narendra modi until:2024-06-22 since:2024-06-18","count":20,"querySource":"typed_query","product":"Top"},"features":{"rweb_tipjar_consumption_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":true,"creator_subscriptions_tweet_preview_api_enabled":true,"responsive_web_graphql_timeline_navigation_enabled":true,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"communities_web_enable_tweet_community_results_fetch":true,"c9s_tweet_anatomy_moderator_badge_enabled":true,"articles_preview_enabled":true,"tweetypie_unmention_optimization_enabled":true,"responsive_web_edit_tweet_api_enabled":true,"graphql_is_translatable_rweb_tweet_is_translatable_enabled":true,"view_counts_everywhere_api_enabled":true,"longform_notetweets_consumption_enabled":true,"responsive_web_twitter_article_tweet_consumption_enabled":true,"tweet_awards_web_tipping_enabled":false,"creator_subscriptions_quote_tweet_preview_enabled":false,"freedom_of_speech_not_reach_fetch_enabled":true,"standardized_nudges_misinfo":true,"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled":true,"rweb_video_timestamps_enabled":true,"longform_notetweets_rich_text_read_enabled":true,"longform_notetweets_inline_media_enabled":true,"responsive_web_enhance_cards_enabled":false}}`)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(payload))
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

	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, jsonResponse)
}
