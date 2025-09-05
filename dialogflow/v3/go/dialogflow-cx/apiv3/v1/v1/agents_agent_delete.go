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

// Package cx contains samples for Google Cloud Dialogflow CX API.
package cx

import (
	"flag"
	"log"
	"os"
)

// [START dialogflow_cx_delete_agent]
import (
	"context"
	"fmt"
	"io"

	dialogflow "cloud.google.com/go/dialogflow/cx/apiv3"
	dialogflowpb "cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// deleteAgent deletes a Dialogflow CX agent.
// This sample demonstrates how to delete an existing Dialogflow CX agent.
// Deleting an agent permanently removes all its associated resources, such as flows, intents, and entity types.
// It's important to note that an agent can only be deleted if it has no active environments or versions associated with it,
// as these dependencies must be removed first.
func deleteAgent(w io.Writer, projectID, locationID, agentID string) error {
	ctx := context.Background()
	agentsClient, err := dialogflow.NewAgentsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewAgentsClient: %w", err)
	}
	defer agentsClient.Close()

	// The agent name has the format: `projects/<ProjectID>/locations/<LocationID>/agents/<AgentID>`
	agentName := fmt.Sprintf("projects/%s/locations/%s/agents/%s", projectID, locationID, agentID)

	req := &dialogflowpb.DeleteAgentRequest{
		Name: agentName,
	}

	fmt.Fprintf(w, "Attempting to delete agent: %s\n", agentName)

	err = agentsClient.DeleteAgent(ctx, req)
	if err != nil {
		// Check for specific error types and provide actionable messages.
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.NotFound:
				return fmt.Errorf("Agent %s not found. Please ensure the project ID, location ID, and agent ID are correct and the agent exists: %w", agentName, err)
			case codes.FailedPrecondition:
				// This error often means the agent has associated environments or versions.
				// The error message from the API usually provides details.
				return fmt.Errorf("Failed to delete agent %s. It might have associated environments or versions that need to be deleted first: %w", agentName, err)
			case codes.PermissionDenied:
				return fmt.Errorf("Permission denied to delete agent %s. Ensure your service account has the 'Dialogflow Editor' or 'Dialogflow Admin' role: %w", agentName, err)
			default:
				return fmt.Errorf("DeleteAgent failed for %s: %w", agentName, err)
			}
		}
		return fmt.Errorf("DeleteAgent failed for %s: %w", agentName, err)
	}

	fmt.Fprintf(w, "Agent %s deleted successfully.\n", agentID)
	return nil
}

// [END dialogflow_cx_delete_agent]

func main() {
	projectID := flag.String("project-id", os.Getenv("GOOGLE_CLOUD_PROJECT"), "Your Google Cloud Project ID")
	locationID := flag.String("location-id", "global", "The Dialogflow CX location ID (e.g., 'global', 'us-central1')")
	agentID := flag.String("agent-id", "", "The ID of the agent to delete")

	flag.Parse()

	if *projectID == "" {
		log.Fatal("project-id must be set. Use --project-id or set GOOGLE_CLOUD_PROJECT environment variable.")
	}
	if *agentID == "" {
		log.Fatal("agent-id must be set. Please provide the ID of the agent to delete.")
	}

	// Before running this sample, ensure the agent exists and has no associated environments or versions.
	// You can create a temporary agent using the `createAgent` sample if needed.

	err := deleteAgent(os.Stdout, *projectID, *locationID, *agentID)
	if err != nil {
		log.Fatalf("Failed to delete agent: %v", err)
	}
}
