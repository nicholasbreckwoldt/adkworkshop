package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/artifact"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/api"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session/vertexai"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
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

    if os.Getenv("AGENT_ENGINE_LOCATION") == "" {
        log.Fatalf("AGENT_ENGINE_LOCATION env not set")
    }
    if os.Getenv("AGENT_ENGINE_ID") == "" {
        log.Fatalf("AGENT_ENGINE_ID env not set")
    }
}

// This afterModelCallback allows us to 'inject' the tool generated artifact into the agent response.
func afterModelCallback(ctx agent.CallbackContext, llmResponse *model.LLMResponse, llmResponseError error) (*model.LLMResponse, error) {
    artifactFound, _ := ctx.State().Get("tool_generated_artifact")
    if artifactFound != nil {
        if artifactFound.(bool) {
            artifactsRes, err := ctx.Artifacts().Load(ctx, artifactFileName)
            if err != nil {
                return nil, err
            }
            llmResponse.Content.Parts = append(llmResponse.Content.Parts, artifactsRes.Part)
            err = ctx.State().Set("tool_generated_artifact", false)
            if err != nil {
                return nil, err
            }
            return llmResponse, nil
        }
    }
    return nil, nil
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
    	Location:  os.Getenv("AGENT_ENGINE_LOCATION"),
        ReasoningEngine: os.Getenv("AGENT_ENGINE_ID"),
    })
    if err != nil {
        log.Fatalf("Failed to create session service")
    }

     // Define GenerateImage tool
    generateImageTool, err := generateImageTool()
    if err != nil {
        log.Fatalf("Failed to create tool: %v", err)
    }

	editImageTool, err := editImageTool()
    if err != nil {
        log.Fatalf("Failed to create tool: %v", err)
    }

	// Create new agent
    adkAgent, err := llmagent.New(llmagent.Config{
        Name:        "artifactsAgent",
        Model:       model,
        Description: "Assists with image editing",
        Instruction: "You are a helpful assistant that can assist users with generating and editing images. Use the 'generateImage' tool for generation and the 'editImage' tool for image editing.",
        Tools: []tool.Tool{
			generateImageTool,
			editImageTool,
        },
        AfterModelCallbacks: []llmagent.AfterModelCallback{
            afterModelCallback,
        },
    })
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }

    // Configure launcher
    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(adkAgent),
		SessionService: sessionService,
		ArtifactService: artifact.InMemoryService(),
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

