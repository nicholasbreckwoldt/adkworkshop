package sequential

import (
	"context"
	"log"
	"os"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
	"google.golang.org/adk/agent/workflowagents/loopagent"
)

func LoopAgent(ctx context.Context) (agent.Agent, error) {
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
        Backend: genai.BackendVertexAI,
        Project: os.Getenv("GOOGLE_CLOUD_PROJECT"),
        Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
    })
    if err != nil {
        log.Fatalf("Failed to create model: %v", err)
    }

    // Critic agent
    // Reviews existing blog draft and prepares constructive critique.
	critic, err := llmagent.New(llmagent.Config{
        Name:        "critic",
        Model:       model,
        Description: "Critiques the blog draft and suggest improvements",
        Instruction: `
			Your role is to critique the provided DRAFT blog post and suggest improvements. 
			Provide only actionable suggestion for improvement. Provide critique in a short concise paragraph only.

			# DRAFT:
			{blog}
		`,
		OutputKey: "critique",
    })
    if err != nil {
        return nil, err
    }

    // Refiner agent
    // Updates blog draft based on critique generated in previous step.
	refiner, err := llmagent.New(llmagent.Config{
        Name:  "refiner",
        Model: model,
        Instruction: `Refine the blog post DRAFT based on the provided CRITIQUE.
            # DRAFT:
            {blog}

            # CRITIQUE
            {critique}
		`,
        OutputKey: "blog",
    })
	 if err != nil {
        return nil, err
    }

    return loopagent.New(loopagent.Config{
        AgentConfig: agent.Config{
            Name:      "refiner_loop",
            SubAgents: []agent.Agent{
				critic, refiner,
			},
        },
        MaxIterations: 2,
    })
}

func SequentialAgent(ctx context.Context) (agent.Agent, error) {

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
        Backend: genai.BackendVertexAI,
        Project: os.Getenv("GOOGLE_CLOUD_PROJECT"),
        Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
    })
    if err != nil {
        log.Fatalf("Failed to create model: %v", err)
    }

    // Planner agent
    // Plans blog before preparing draft.
	planner, err := llmagent.New(llmagent.Config{
        Name:        "planner",
        Model:       model,
        Description: "Generates a plan for preparing a blog post",
        Instruction: "Prepare an outline plan of the blog post. Keep the plan to one or two paragraphs maximum",
		OutputKey: "plan",
    })
	if err != nil {
		return nil, err
	}

    // Drafter agent
    // Prepares first draft according to 'plan' generated in previous step.
	drafter, err := llmagent.New(llmagent.Config{
        Name:        "drafter",
        Model:       model,
        Description: "Drafts a blog post based on prepared plan",
        Instruction: `
            Using the below PLAN, write a short blog post. 
            Write concisely and to the point, keeping the blog post length to a maxium of 2-3 paragraphs.

            # PLAN:
            {plan}
        `,
		OutputKey: "blog",
    })

    // Refine agent
    // Refines blog through iterative feedback and redrafting.
	refiner, err := LoopAgent(ctx)
	if err != nil {
		return nil, err
	}

	return sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "sequential_writer",
			Description: "Generates a refined blog post",
			SubAgents:   []agent.Agent{
				planner,
				drafter,
			 	refiner,
			},
		},
	})
}