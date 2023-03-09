package gpt3

import "context"

type ImageSizeType string

const (
	IST256  ImageSizeType = "256x256"
	IST512  ImageSizeType = "512x512"
	IST1024 ImageSizeType = "1024x1024"
)

type CreateImageReq struct {
	// A text description of the desired image(s). The maximum length is 1000 characters.
	Prompt string `json:"prompt"`
	// The number of images to generate. Must be between 1 and 10.
	N int `json:"n"`
	// The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024.
	Size ImageSizeType `json:"size"`
}

type urls struct {
	URL string `json:"url"`
}

type CreateImageResp struct {
	Data    []urls `json:"data"`
	Created int64  `json:"created"`
}

func (c *client) CreateImage(ctx context.Context, request CreateImageReq) (*CreateImageResp, error) {
	req, err := c.newRequest(ctx, "POST", "/images/generations", request)
	if err != nil {
		return nil, err
	}
	resp, err := c.performRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(CreateImageResp)
	if err := getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}
