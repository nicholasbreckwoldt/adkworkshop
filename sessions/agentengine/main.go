package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	aiplatform "google.golang.org/api/aiplatform/v1beta1"
	"google.golang.org/api/option"
)

func init() {
    if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
        log.Fatalf("GOOGLE_CLOUD_PROJECT env not set")
    }
    if os.Getenv("GOOGLE_CLOUD_LOCATION") == "" {
        log.Fatalf("GOOGLE_CLOUD_LOCATION env not set")
    }
}

func main() {

	ctx := context.Background()

	// Initialize the service
	// The endpoint for Vertex AI is region-specific
	service, err := aiplatform.NewService(ctx, option.WithEndpoint(fmt.Sprintf("https://%s-aiplatform.googleapis.com/", os.Getenv("GOOGLE_CLOUD_LOCATION"))))
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	parent := fmt.Sprintf("projects/%s/locations/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), os.Getenv("GOOGLE_CLOUD_LOCATION"))

	engine := &aiplatform.GoogleCloudAiplatformV1beta1ReasoningEngine{
		DisplayName: "Agent Engine",
		ContextSpec: &aiplatform.GoogleCloudAiplatformV1beta1ReasoningEngineContextSpec{},
		Description: "Agent Engine for Session Management",
		Spec: &aiplatform.GoogleCloudAiplatformV1beta1ReasoningEngineSpec{},
	}

	op, err := service.Projects.Locations.ReasoningEngines.Create(parent, engine).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Error creating engine: %v", err)
	}

	for {
		op, err = service.Projects.Locations.ReasoningEngines.Operations.Get(op.Name).Context(ctx).Do()
		if err != nil {
			log.Fatalf("Error polling operation: %v", err)
		}
		if op.Done {
			if op.Error != nil {
				log.Fatalf("Operation failed: %v", op.Error.Message)
			}
			fmt.Println("Reasoning Engine created successfully!")
			fmt.Println(string(op.Response))
			return
		}
		fmt.Println("Operation still in progress...")
		time.Sleep(5*time.Second)
	}
}