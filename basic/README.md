# Basic ADK Agent

This application demonstrates a simple 'Hello World!' ADK agent.

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

### API Interaction

#### Step 1: List Applications:

```bash
curl -X GET http://localhost:8080/api/list-apps

```

#### Step 2: List available sessions:
```bash
curl -X GET http://localhost:8080/api/apps/simpleAgent/users/user/sessions
```

#### Step 3: Get the current session:

```bash
curl -X GET http://localhost:8080/api/apps/simpleAgent/users/user/sessions/<sessionID>
```

#### Step 4: Create a new session:
```bash
curl -X POST http://localhost:8080/api/apps/simpleAgent/users/user/sessions

```

#### Step 5: Invoke the ADK agent:

```bash
curl -X POST http://localhost:8080/api/run \
    -H "Content-Type: application/json" \
    -d '{
    "appName": "simpleAgent",
    "userId": "user",
    "sessionId": "<sessionID>",
    "newMessage": {
        "role": "user",
        "parts": [{
        "text": "Enter your query here..."
        }]
    }
}'
```