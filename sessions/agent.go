package main

import (
	"context"
	"log"
	"os"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/cmd/launcher/web/api"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
    "google.golang.org/adk/session/vertexai"
    "time"
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

// GetCurrentRequest defines the arguments for the getCurrentTime tool.
type GetCurrentTimeRequest struct {
    // The timezone for which to retrieve the current time
    TimeZone string `json:"time_zone"`
}

//  GetCurrentTimeResponset defines the output of the getCurrentTime tool.
type GetCurrentTimeResponse struct {
    // Current time
    Time string `json:"current_time"`
}

// Tool handler
func getCurrentTimeHandler(ctx tool.Context, args GetCurrentTimeRequest) (GetCurrentTimeResponse, error) {
    loc, err := time.LoadLocation(args.TimeZone)
	if err != nil {
		return GetCurrentTimeResponse{}, err
	}
    return  GetCurrentTimeResponse{
        Time: time.Now().In(loc).String(),
    }, nil
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

    // Define getCurrentTime tool
    getCurrentTimeTool, err := functiontool.New(
        functiontool.Config{
            Name:        "getCurrentTime",
            Description: "Returns the current time given the timezone. Example: 'America/New_York', 'Europe/London', 'Asia/Tokyo'",
        },
        getCurrentTimeHandler,
    )
    if err != nil {
        log.Fatalf("Failed to create tool: %v", err)
    }

	// Create new agent
    adkAgent, err := llmagent.New(llmagent.Config{
        Name:        "sessionsAgent",
        Model:       model,
        Description: "Assists with user queries",
        Instruction: "You are a helpful assistant that can assist users with their queries",
        Tools: []tool.Tool{
            getCurrentTimeTool,
        },
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

