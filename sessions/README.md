# Basic ADK Agent with VertexAI SessionService

This application demonstrates an ADK Agent with a persistent SessionService (using VertexAI Agent Engine).

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

### Creating an AgentEngine instance (Vertex AI)
Run the following script:
```bash
cd ./agentengine && go run main.go
```

Then grab the reasoning engine and set the following envs:
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


### Deploying to CloudRun

To deploy to CloudRun, run the following:
```bash
make deploy
```
This makes use of the gcloud CLI to deploy to Google Cloud.

### Remote API Interaction

#### Step 1: List Applications:

```bash
curl -X GET -H "Authorization: Bearer $(gcloud auth print-identity-token)" <cloud_run_url>/api/list-apps

```

#### Step 2: Create a new session:
```bash
curl -X POST -H "Authorization: Bearer $(gcloud auth print-identity-token)" <cloud_run_url>/api/apps/sessionsAgent/users/<user>/sessions

```

#### Step 3: Invoke the ADK agent:

```bash
curl -X POST -H "Authorization: Bearer $(gcloud auth print-identity-token)" <cloud_run_url>/api/run \
    -H "Content-Type: application/json" \
    -d '{
    "appName": "sessionsAgent",
    "userId": "<user>",
    "sessionId": "<sessionID>",
    "newMessage": {
        "role": "user",
        "parts": [{
        "text": "Enter your query here..."
        }]
    }
}'
```