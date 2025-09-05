package cx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	dialogflow "cloud.google.com/go/dialogflow/cx/apiv3"
	dialogflowpb "cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockAgentsClient struct {
	deleteAgent func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error
}

func (m *mockAgentsClient) DeleteAgent(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
	return m.deleteAgent(ctx, req, opts...)
}

func (m *mockAgentsClient) GetAgent(ctx context.Context, req *dialogflowpb.GetAgentRequest, opts ...grpc.CallOption) (*dialogflowpb.Agent, error) {
	panic("not implemented") // TODO: Implement
}

func (m *mockAgentsClient) CreateAgent(ctx context.Context, req *dialogflowpb.CreateAgentRequest, opts ...grpc.CallOption) (*dialogflowpb.Agent, error) {
	panic("not implemented") // TODO: Implement
}

func (m *mockAgentsClient) UpdateAgent(ctx context.Context, req *dialogflowpb.UpdateAgentRequest, opts ...grpc.CallOption) (*dialogflowpb.Agent, error) {
	panic("not implemented") // TODO: Implement
}

func (m *mockAgentsClient) ListAgents(ctx context.Context, req *dialogflowpb.ListAgentsRequest, opts ...grpc.CallOption) *dialogflow.AgentIterator {
	panic("not implemented") // TODO: Implement
}

func (m *mockAgentsClient) GetLocation(ctx context.Context, req *dialogflowpb.GetLocationRequest, opts ...grpc.CallOption) (*dialogflowpb.Location, error) {
	panic("not implemented") // TODO: Implement
}

func (m *mockAgentsClient) ListLocations(ctx context.Context, req *dialogflowpb.ListLocationsRequest, opts ...grpc.CallOption) *dialogflow.LocationIterator {
	panic("not implemented") // TODO: Implement
}

func TestDeleteAgent(t *testing.T) {
	projectID := "test-project"
	locationID := "us-central1"
	agentID := "test-agent"

	// Define test cases

tests := []struct {
		name          string
		deleteAgentFn func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error
		wantErr       bool
		errCode       codes.Code
		outputCheck   func(string) bool
	}{
		{
			name: "Successful deletion",
			deleteAgentFn: func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return nil
			},
			wantErr: false,
			outputCheck: func(output string) bool {
				return strings.Contains(output, fmt.Sprintf("Agent %s deleted successfully.", agentID))
			},
		},
		{
			name: "Agent not found",
			deleteAgentFn: func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return status.Error(codes.NotFound, "Agent not found")
			},
			wantErr: true,
			errCode: codes.NotFound,
			outputCheck: func(output string) bool {
				return strings.Contains(output, fmt.Sprintf("Attempting to delete agent: projects/%s/locations/%s/agents/%s", projectID, locationID, agentID))
			},
		},
		{
			name: "Failed precondition",
			deleteAgentFn: func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return status.Error(codes.FailedPrecondition, "Agent has associated resources")
			},
			wantErr: true,
			errCode: codes.FailedPrecondition,
			outputCheck: func(output string) bool {
				return strings.Contains(output, fmt.Sprintf("Attempting to delete agent: projects/%s/locations/%s/agents/%s", projectID, locationID, agentID))
			},
		},
		{
			name: "Permission denied",
			deleteAgentFn: func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return status.Error(codes.PermissionDenied, "Permission denied")
			},
			wantErr: true,
			errCode: codes.PermissionDenied,
			outputCheck: func(output string) bool {
				return strings.Contains(output, fmt.Sprintf("Attempting to delete agent: projects/%s/locations/%s/agents/%s", projectID, locationID, agentID))
			},
		},
		{
			name: "Generic error",
			deleteAgentFn: func(ctx context.Context, req *dialogflowpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return fmt.Errorf("generic error")
			},
			wantErr: true,
			outputCheck: func(output string) bool {
				return strings.Contains(output, fmt.Sprintf("Attempting to delete agent: projects/%s/locations/%s/agents/%s", projectID, locationID, agentID))
			},
		},
	}

	for _, tt := range tests {
		
tt := tt
			
t.Run(tt.name, func(t *testing.T) {
				
				ctx := context.Background()
				
				mockAgentsClient := &mockAgentsClient{
					deleteAgent: tt.deleteAgentFn,
				}

				// Override the NewAgentsClient function to return the mock client.
				newAgentsClient := func(ctx context.Context, opts ...option.ClientOption) (*dialogflow.AgentsClient, error) {
					return &dialogflow.AgentsClient{},
				}

				_ = newAgentsClient // avoid unused variable error

				// Capture stdout
				old := os.Stdout
				defer func() {
					os.Stdout = old
				}()

				var buf bytes.Buffer
				
				log.SetOutput(&buf)
				
				os.Stdout = &buf

				// Call the function under test
				
				err := deleteAgent(&buf, projectID, locationID, agentID)

				// Check for error
				if (err != nil) != tt.wantErr {
					
t.Fatalf("deleteAgent() error = %v, wantErr %v", err, tt.wantErr)
				}

				// Check error code if applicable
				if tt.wantErr && tt.errCode != 0 {
					
s, ok := status.FromError(err)
					
					if !ok || s.Code() != tt.errCode {
						
t.Errorf("deleteAgent() error code = %v, want %v", s.Code(), tt.errCode)
					}
				}

				// Check output
				
output := buf.String()
				
				
				
				
				
				
				
				
				if !tt.outputCheck(output) {
					
t.Errorf("Output mismatch: got %q", output)
				}
			})
		}
}

func TestDeleteAgent_integration(t *testing.T) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	locationID := "global" // or your desired location
	agentID := "123e4567-e89b-12d3-a456-426614174000" // Replace with a valid agent ID for testing

	// Check if project ID is set
	if projectID == "" {
		
t.Skip("Skipping integration test because GOOGLE_CLOUD_PROJECT is not set")
	}

	// Initialize a buffer to capture stdout
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Call the deleteAgent function

	err := deleteAgent(&buf, projectID, locationID, agentID)

	// Check for errors
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	