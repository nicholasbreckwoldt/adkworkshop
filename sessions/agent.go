package main

import (
	"context"
	"log"
	"os"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/cmd/launcher/web/api"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
    "google.golang.org/adk/session/vertexai"
)

func init() {
    if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
        log.Fatalf("GOOGLE_CLOUD_PROJECT env not set")
    }
    if os.Getenv("GOOGLE_CLOUD_LOCATION") == "" {
        log.Fatalf("GOOGLE_CLOUD_LOCATION env not set")
    }

    // Uncomment if using API key instead
	// if os.Getenv("GEMINI_API_KEY") == "" {
	// 	log.Fatalf("GEMINI_API_KEY env not set")
	// }
}

func main() {
	ctx := context.Background()

    // Initialise new model
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
        Backend: genai.BackendVertexAI,
        Project: os.Getenv("GOOGLE_CLOUD_PROJECT"),
        Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
        // APIKey: "", Remove above arguments if using api key
    })
    if err != nil {
        log.Fatalf("Failed to create model: %v", err)
    }

    // Initialise Vertex AI Engine Session Service
    sessionService, err := vertexai.NewSessionService(ctx, vertexai.VertexAIServiceConfig{
    	ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
    	Location: "REASONING_ENGINE_LOCATION",  // TODO: Populate
        ReasoningEngine: "REASONING_ENGINE_ID", // TODO: Populate
    })
    if err != nil {
        log.Fatalf("Failed to create session service")
    }

	// Create new agent
    adkAgent, err := llmagent.New(llmagent.Config{
        Name:        "sessionsAgent",
        Model:       model,
        Description: "Assists with user queries",
        Instruction: "You are a helpful assistant that can assist users with their queries",
    })
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }  

    // Configure launcher
    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(adkAgent),
		SessionService: sessionService,
    }

	// Set web launcher to interact via the ADK webui
	webLauncher := web.NewLauncher(
		webui.NewLauncher(),
		api.NewLauncher(),
	)
	_, err = webLauncher.Parse([]string{"--port", "8080", "webui", "api"})
	if err != nil {
		log.Fatalf("webLauncher.Parse() error = %v", err)
	}

    if err := webLauncher.Run(ctx, config); err != nil {
		log.Fatalf("webLauncher.Run() error = %v", err)
	}
}

