# Basic ADK Agent with Artifacts

This application demonstrates an ADK Agent that can interact and manage artifacts via an ArtifactsService.

### Setup

#### Step 1: Install dependencies:
```bash
go mod tidy
```

#### Step 2: Configure environment:

If using a Google Cloud Project:
```bash
export GOOGLE_CLOUD_PROJECT="YOUR PROJECT"
export GOOGLE_CLOUD_LOCATION="YOUR LOCATION"
export GOOGLE_GENAI_USE_VERTEXAI=TRUE
```

If using API key:
```bash
export GEMINI_API_KEY="YOUR API KEY"
```

### AgentEngine instance configuration
See details in the Session's agent for setup.

```bash
export AGENT_ENGINE_LOCATION="YOUR LOCATION"
export AGENT_ENGINE_ID="YOUR AGENT ENGINE ID"
```

### Running the ADK Agent

To start the agent, run the following:
```bash
go run .
```

### ADK UI Interaction

Interact with the agent at http://localhost:8080/ui/