package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

var client *genai.Client

const (
	artifactFileName = "artifact.png"
)

func init() {
	var err error
	ctx := context.Background()
	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
		Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
		Backend:  genai.BackendVertexAI,
		// APIKey: "", Remove above arguments if using api key
	})
	if err != nil {
		log.Fatalf("failed to initialised genai.Client()")
	}
}

// Request message for the GenerateImage tool
type GenerateImageRequest struct {
    // Image description
    Description string `json:"description"`
}

// Response message for the GenerateImage tool
type GenerateImageResponse struct {
    Result string `json:"result"`
}

func generateImageHandler(ctx tool.Context, args GenerateImageRequest) (GenerateImageResponse, error) {
	// Use Gemini 2.5 flash model to generate an image
	prompt := fmt.Sprintf("Generate an interesting visual representation based on the following provided description: %s", args.Description)
	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash-image",
		genai.Text(prompt), &genai.GenerateContentConfig{
			ImageConfig: &genai.ImageConfig{
				AspectRatio:    "1:1",
				ImageSize:      "1K",
				OutputMIMEType: "image/png",
			},
	})
	if err != nil {
		return GenerateImageResponse{}, err
	}

	var imageBytes []byte
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			return GenerateImageResponse{}, nil
		} else if part.InlineData != nil {
			imageBytes = part.InlineData.Data
			break
		}
	}

	// Save to artifacts service for later
	_, err = ctx.Artifacts().Save(ctx, artifactFileName, genai.NewPartFromBytes(imageBytes, "image/png"))
	if err != nil {
		return GenerateImageResponse{}, err
	}

	// Save to state to communicate to rest of app that artifact is available
	err = ctx.State().Set("tool_generated_artifact", true)
	if err != nil {
		return GenerateImageResponse{}, err
	}

	return GenerateImageResponse{
		Result: "success",
	}, nil
}

func generateImageTool() (tool.Tool, error) {
	return functiontool.New(
        functiontool.Config{
            Name:        "generate_image",
            Description: "Generates an image based on provided description",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"description": {
						Type: "string",
						Description: "Detailed description of image to generate",
					},
				},
				Required: []string{"description"},
			},
        },
        generateImageHandler,
    )
}

// Request message for the EditImage tool
type EditImageRequest struct {
    // Editing instruciton
    Instruction string `json:"instruction"`
}

// Response message for the EditImage tool
type EditImageResponse struct {
    Result string `json:"result"`
}

func editImageHandler(ctx tool.Context, args EditImageRequest) (EditImageResponse, error) {
	// Load the artifact
	loadRes, err := ctx.Artifacts().Load(ctx, artifactFileName)
	if err != nil {
		return EditImageResponse{}, err
	}

	// Apply editing using gemini image model
	prompt := fmt.Sprintf("Edit the below image according to the following instructions: %s", args.Instruction)
	resp, err := client.Models.GenerateContent(ctx, 
		"gemini-2.5-flash-image",
		[]*genai.Content{
			{
				Role:  genai.RoleUser,
				Parts: []*genai.Part{
					genai.NewPartFromText(prompt),
					genai.NewPartFromBytes(loadRes.Part.InlineData.Data, loadRes.Part.InlineData.MIMEType),
				},
				
			},
		},
		&genai.GenerateContentConfig{
			ImageConfig: &genai.ImageConfig{
				AspectRatio:    "1:1",
				ImageSize:      "1K",
				OutputMIMEType: "image/png",
			},
	})
	if err != nil {
		return EditImageResponse{}, err
	}

	var imageBytes []byte
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.InlineData != nil {
			imageBytes = part.InlineData.Data
			break
		}
	}

	// Save to artifacts service for later
	_, err = ctx.Artifacts().Save(ctx, artifactFileName, genai.NewPartFromBytes(imageBytes, "image/png"))
	if err != nil {
		return EditImageResponse{}, err
	}

	// Save to state to communicate to rest of app that artifact is available
	err = ctx.State().Set("tool_generated_artifact", true)
	if err != nil {
		return EditImageResponse{}, err
	}

	return EditImageResponse{
		Result: "success",
	}, nil
}

func editImageTool() (tool.Tool, error) {
	return functiontool.New(
        functiontool.Config{
            Name:        "edit_image",
            Description: "Edit an existing image based on provided instruction",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"instruction": {
						Type: "string",
						Description: "Detailed instructions for editing the image",
					},
				},
				Required: []string{"instruction"},
			},
        },
        editImageHandler,
    )
}
