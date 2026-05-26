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
		--set-env-vars=AGENT_ENGINE_ID=$(AGENT_ENGINE_ID) \
		--no-allow-unauthenticated
	@echo "Agent deployment complete!"