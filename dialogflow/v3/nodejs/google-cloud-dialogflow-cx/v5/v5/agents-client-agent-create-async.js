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

'use strict';

// Imports for the command-line runner
const process = require('process');

// [START dialogflow_v3_agents_agent_create_async]
const {AgentsClient} = require('@google-cloud/dialogflow-cx');
const {status} = require('@grpc/grpc-js');

// Global client instance (best practice for production)
const client = new AgentsClient();

/**
 * Creates a new Dialogflow CX agent in the specified Google Cloud project and location.
 *
 * An agent acts as a virtual agent that handles conversations with end-users.
 * It processes natural language input, understands user intent, and responds accordingly.
 * Creating an agent is the first step to building a conversational AI experience.
 *
 * @param {string} projectId The Google Cloud project ID. Example: 'my-project-id-123'
 * @param {string} location The Google Cloud location (e.g., 'global', 'us-central1'). Example: 'global'
 * @param {string} agentDisplayName The display name for the new agent. This must be unique within the project and location. Example: 'my-new-cx-agent-nodejs-example'
 */
async function createAgent(
  projectId: string = 'my-project-id-123',
  location: string = 'global',
  agentDisplayName: string = 'my-new-cx-agent-nodejs-example'
) {
  // Construct the parent path for the agent.
  // This specifies where the agent will be created.
  const parent = `projects/${projectId}/locations/${location}`;

  // Define the agent configuration.
  // displayName, defaultLanguageCode, and timeZone are required fields.
  const agent = {
    displayName: agentDisplayName,
    defaultLanguageCode: 'en', // The default language for the agent. 'en' (English) is commonly used.
    timeZone: 'America/Los_Angeles', // The time zone for the agent, e.g., 'America/Los_Angeles'.
  };

  // Prepare the request object for the createAgent API call.
  const request = {
    parent: parent,
    agent: agent,
  };

  try {
    // Act: Make the API call to create the agent.
    const [newAgent] = await client.createAgent(request);
    // Assert: Print the success message with relevant details.
    console.log(`Agent created: ${newAgent.displayName} (${newAgent.name})`);
  } catch (err) {
    // Handle specific error for existing agent.
    // If an agent with the same display name already exists, Dialogflow CX returns an ALREADY_EXISTS error.
    if (err.code === status.ALREADY_EXISTS) {
      console.warn(
        `Agent with display name "${agentDisplayName}" already exists in location "${location}" of project "${projectId}".`
      );
      // In a real application, you might want to retrieve the existing agent or prompt the user for a different name.
    } else {
      // Log and re-throw other unexpected errors.
      console.error(`Error creating agent "${agentDisplayName}":`, err);
      throw err;
    }
  }
}
// [END dialogflow_v3_agents_agent_create_async]

// Command-line execution boilerplate
async function main() {
  const args = process.argv.slice(2);

  // If command-line arguments are provided, use them. Otherwise, the createAgent
  // function will use its hardcoded default example values.
  if (args.length >= 3) {
    const projectId = args[0];
    const location = args[1];
    const agentDisplayName = args[2];
    await createAgent(projectId, location, agentDisplayName);
  } else {
    console.log('No command-line arguments provided. Using hardcoded example values.');
    await createAgent(); // Calls with default example values
  }
}

if (require.main === module) {
  process.on('uncaughtException', err => {
    console.error(`Error running sample: ${err.message}`);
    console.error(`
To run this sample from the command-line, specify three arguments:
 - Google Cloud Project ID (e.g., 'my-project-123')
 - Google Cloud Location (e.g., 'global' or 'us-central1')
 - Agent Display Name (e.g., 'my-new-agent')

Usage:
  node createAgent.js <projectId> <location> <agentDisplayName>
`);
    process.exitCode = 1;
  });

  main();
}

module.exports = {
  createAgent,
};
