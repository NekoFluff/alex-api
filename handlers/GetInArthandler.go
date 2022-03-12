package handlers

import (
	"addi/data"
	"addi/models"
	"addi/restapi/operations"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
)

func GetInArtHandler(params operations.GetInArtParams) middleware.Responder {

	page := params.Page
	skip := (page - 1) * 50
	var limit int64 = 50
	twitterMediaList, err := data.GetTwitterMedia(&skip, &limit)
	var inArt []*models.InArt

	if err != nil {
		fmt.Errorf("%v", err)
	} else {
		for _, twitterMedia := range twitterMediaList {
			fmt.Println(twitterMedia)

			author := string(twitterMedia.Author.UserName)
			height := int64(twitterMedia.Height)
			width := int64(twitterMedia.Width)
			sensitive := bool(twitterMedia.PossiblySensitive)
			tweetId := string(twitterMedia.TweetId)
			url := string(twitterMedia.Url)
			m := models.InArt{
				Author:            &author,
				Height:            &height,
				PossiblySensitive: &sensitive,
				TweetID:           &tweetId,
				URL:               &url,
				Width:             &width,
			}
			inArt = append(inArt, &m)
		}
	}

	return operations.NewGetInArtOK().WithPayload(inArt)
}
