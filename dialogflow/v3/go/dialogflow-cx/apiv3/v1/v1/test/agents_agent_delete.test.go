package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	dialogflowcx "cloud.google.com/go/dialogflow/cx/apiv3"
	cxpb "cloud.google.com/go/dialogflow/cx/apiv3/cxpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"../agents_agent_delete"
)

type mockAgentsClient struct {
	deleteAgent func(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error
}

func (m *mockAgentsClient) DeleteAgent(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
	return m.deleteAgent(ctx, req, opts...)
}

func (m *mockAgentsClient) GetAgent(ctx context.Context, req *cxpb.GetAgentRequest, opts ...grpc.CallOption) (*cxpb.Agent, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) ListAgents(ctx context.Context, req *cxpb.ListAgentsRequest, opts ...grpc.CallOption) (*cxpb.ListAgentsResponse, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) UpdateAgent(ctx context.Context, req *cxpb.UpdateAgentRequest, opts ...grpc.CallOption) (*cxpb.Agent, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) CreateAgent(ctx context.Context, req *cxpb.CreateAgentRequest, opts ...grpc.CallOption) (*cxpb.Agent, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) RestoreAgent(ctx context.Context, req *cxpb.RestoreAgentRequest, opts ...grpc.CallOption) (*cxpb.Agent, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) ExportAgent(ctx context.Context, req *cxpb.ExportAgentRequest, opts ...grpc.CallOption) (*cxpb.ExportAgentOperationMetadata, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) GetAgentValidationResult(ctx context.Context, req *cxpb.GetAgentValidationResultRequest, opts ...grpc.CallOption) (*cxpb.AgentValidationResult, error) {
	panic("not implemented")
}
func (m *mockAgentsClient) ValidateAgent(ctx context.Context, req *cxpb.ValidateAgentRequest, opts ...grpc.CallOption) (*cxpb.AgentValidationResult, error) {
	panic("not implemented")
}

func TestDeleteAgent(t *testing.T) {
	projectID := "test-project"
	locationID := "us-central1"
	agentID := "test-agent"
	agentName := fmt.Sprintf("projects/%s/locations/%s/agents/%s", projectID, locationID, agentID)

	tests := []struct {
		name         string
		deleteAgentFunc func(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error
		wantErr      bool
		err          error
		outputContains string
	}{
		{
			name: "Successful deletion",
			deleteAgentFunc: func(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				if req.Name != agentName {
					return fmt.Errorf("incorrect agent name in request: got %q, want %q", req.Name, agentName)
				}
				return nil
			},
			wantErr:      false,
			outputContains: agentID,
		},
		{
			name: "Agent not found",
			deleteAgentFunc: func(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return status.Error(codes.NotFound, "agent not found")
			},
			wantErr:      true,
			err:          status.Error(codes.NotFound, "agent not found"),
			outputContains: "not found",
		},
		{
			name: "Deletion failed",
			deleteAgentFunc: func(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
				return fmt.Errorf("deletion failed")
			},
			wantErr:      true,
			err:          fmt.Errorf("deletion failed"),
			outputContains: "failed",
		},
	}

	for _, tt := range tests {
		clientOption := option.WithGRPCConn(&grpc.ClientConn{})
		ctx := context.Background()

		buf := &bytes.Buffer{}

		t.Run(tt.name, func(t *testing.T) {
			agentsClient := &mockAgentsClient{
				deleteAgent: tt.deleteAgentFunc,
			}

			// Override the global variable with the mock client.
			deleteAgentOverride = func(ctx context.Context, opts ...option.ClientOption) (*dialogflowcx.AgentsClient, error) {
				return agentsClient, nil
			}
			defer func() {
				deleteAgentOverride = nil // Reset the override after the test
			}()

			err := deleteAgent(buf, projectID, locationID, agentID)

			if (err != nil) != tt.wantErr {
				t.Errorf("deleteAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				// Compare errors if both are not nil and an error is expected
				if !strings.Contains(err.Error(), tt.err.Error()) {
					// Check if the returned error contains the expected error message
					t.Errorf("deleteAgent() error = %v, want error to contain %v", err, tt.err)
					return
				}
			}

			output := buf.String()
			if !strings.Contains(output, tt.outputContains) {
				// Check if the output contains the expected string
				fmt.Println(output)
				t.Errorf("deleteAgent() output = %q, want to contain %q", output, tt.outputContains)
			}
		})
	}
}

// deleteAgentOverride is a variable that can be overridden in tests.
// This is necessary because the AgentsClient is created using the dialogflowcx.NewAgentsClient function.
// In tests, we want to use a mock AgentsClient instead of the real one.
var deleteAgentOverride func(ctx context.Context, opts ...option.ClientOption) (*dialogflowcx.AgentsClient, error)

// Override the NewAgentsClient function to use the mock client in tests.
func init() {
	newAgentsClient = func(ctx context.Context, opts ...option.ClientOption) (*dialogflowcx.AgentsClient, error) {
		if deleteAgentOverride != nil {
			return deleteAgentOverride(ctx, opts...)
		}
		return dialogflowcx.NewAgentsClient(ctx, opts...)
	}
}

// newAgentsClient is a variable that holds the function for creating a new AgentsClient.
// This is necessary because we want to override it in tests.
var newAgentsClient func(ctx context.Context, opts ...option.ClientOption) (*dialogflowcx.AgentsClient, error)

// Override the dialogflowcx.NewAgentsClient function to use the newAgentsClient variable.
func TestMain(m *testing.M) {
	// Save the original NewAgentsClient function.
	originalNewAgentsClient := dialogflowcx.NewAgentsClient

	// Override the NewAgentsClient function with our variable.
	// This allows us to override it in tests.
dialogflowcx.NewAgentsClient = func(ctx context.Context, opts ...option.ClientOption) (*dialogflowcx.AgentsClient, error) {
		return newAgentsClient(ctx, opts...)
	}

	// Run the tests.
	code := m.Run()

	// Restore the original NewAgentsClient function.
dialogflowcx.NewAgentsClient = originalNewAgentsClient

	// Exit.
	os.Exit(code)
}

func TestDeleteAgent_NoCrash(t *testing.T) {
	projectID := "test-project"
	locationID := "us-central1"
	agentID := "test-agent"

	// Create a mock AgentsClient that does nothing.
	agentsClient := &mockAgentsClient{
		deleteAgent: func(ctx context.Context, req *cxpb.DeleteAgentRequest, opts ...grpc.CallOption) error {
			return nil
		},
	}

	// Override the global variable with the mock client.
	deleteAgentOverride = func(ctx context.Context, opts ...option.ClientOption) (*dialogflowcx.AgentsClient, error) {
		return agentsClient, nil
	}
	defer func() {
		deleteAgentOverride = nil // Reset the override after the test
	}()

	// Call the deleteAgent function with the mock client.
	err := deleteAgent(io.Discard, projectID, locationID, agentID)

	// Check that the function did not crash.
	if err != nil {
		t.Errorf("deleteAgent() should not have crashed, but returned error: %v", err)
	}
}
