// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains samples for Google Cloud Dialogflow CX API.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	dialogflowcx "cloud.google.com/go/dialogflow/cx/apiv3"
	cxpb "cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// [START dialogflow_cx_delete_agent]

// deleteAgent deletes a specified Dialogflow CX agent.
//
// An agent is a virtual agent that handles conversations with end-users.
// Deleting an agent permanently removes it and all its associated data,
// including flows, intents, entity types, and webhooks.
func deleteAgent(w io.Writer, projectID, locationID, agentID string) error {
	// projectID := "your-project-id"
	// locationID := "global" // or "us-central1", etc.
	// agentID := "your-agent-id"

	ctx := context.Background()

	// Create a new AgentsClient.
	// The client will automatically authenticate using Application Default Credentials (ADC).
	client, err := dialogflowcx.NewAgentsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewAgentsClient: %w", err)
	}
	// It's a good practice to close the client when it is no longer needed.
	// This ensures that resources are released and connections are properly terminated.
	defer client.Close()

	// Construct the full resource name of the agent to delete.
	// Format: `projects/<ProjectID>/locations/<LocationID>/agents/<AgentID>`
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s", projectID, locationID, agentID)

	// Create the DeleteAgentRequest.
	req := &cxpb.DeleteAgentRequest{
		Name: name,
	}

	// Call the DeleteAgent method.
	err = client.DeleteAgent(ctx, req)
	if err != nil {
		// Handle specific errors, such as the agent not being found.
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.NotFound {
				return fmt.Errorf("agent %q not found. Please ensure the agent ID and location are correct: %w", name, err)
			}
		}
		return fmt.Errorf("DeleteAgent failed for agent %q: %w", name, err)
	}

	fmt.Fprintf(w, "Agent %s deleted successfully.\n", agentID)

	return nil
}

// [END dialogflow_cx_delete_agent]

func main() {
	projectID := flag.String("project-id", "", "Your Google Cloud Project ID")
	locationID := flag.String("location-id", "global", "The GCP location (e.g., global, us-central1)")
	agentID := flag.String("agent-id", "", "The ID of the Dialogflow CX agent to delete")

	flag.Parse()

	if *projectID == "" {
		fmt.Fprintln(os.Stderr, "Error: --project-id must be provided.")
		flag.Usage()
		os.Exit(1)
	}
	if *agentID == "" {
		fmt.Fprintln(os.Stderr, "Error: --agent-id must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	if err := deleteAgent(os.Stdout, *projectID, *locationID, *agentID); err != nil {
		log.Fatalf("Failed to delete agent: %v", err)
	}
}
