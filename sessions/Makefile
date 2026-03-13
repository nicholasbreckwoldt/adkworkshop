GOOGLE_CLOUD_PROJECT=YOUR_GOOGLE_CLOUD_PROJECT
GOOGLE_CLOUD_LOCATION=YOUR_GOOGLE_CLOUD_LOCATION
GOOGLE_GENAI_USE_VERTEXAI=TRUE
AGENT_ENGINE_LOCATION=YOUR_AGENT_ENGINE_LOCATION
AGENT_ENGINE_ID=YOUR_AGENT_ENGINE_ID

# Build and deploy the agent to Cloud Run
deploy:
	@echo "Deploying agent to Cloud Run..."
	gcloud run deploy my-agent \
		--source=. \
		--region $(GOOGLE_CLOUD_LOCATION) \
		--project $(GOOGLE_CLOUD_PROJECT) \
		--set-env-vars=GOOGLE_CLOUD_PROJECT=$(GOOGLE_CLOUD_PROJECT) \
		--set-env-vars=GOOGLE_GENAI_USE_VERTEXAI=$(GOOGLE_CLOUD_PROJECT) \
		--set-env-vars=GOOGLE_CLOUD_LOCATION=$(GOOGLE_CLOUD_LOCATION) \
		--set-env-vars=AGENT_ENGINE_LOCATION=$(AGENT_ENGINE_LOCATION) \
		--set-env-vars=AGENT_ENGINE_ID=$(AGENT_ENGINE_ID) \
		--no-allow-unauthenticated
	@echo "Agent deployment complete!"