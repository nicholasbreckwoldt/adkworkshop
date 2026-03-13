# ADK Workflow Agent

This application demonstrates ADK's workflow agents. The agent makes use of 'Sequential' and 'Loop' agent workflows to define an agentic blog writing pipeline.

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

### Running the ADK Agent

To start the agent, run the following:
```bash
go run .
```

### ADK UI Interaction

Interact with the agent at http://localhost:8080/ui/