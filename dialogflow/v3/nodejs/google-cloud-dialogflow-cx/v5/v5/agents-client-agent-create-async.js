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

const process = require('process');

// [START dialogflow_v3_agents_agent_create_async]
const {AgentsClient} = require('@google-cloud/dialogflow-cx');
const {status} = require('@grpc/grpc-js');

const client = new AgentsClient();

/**
 * Creates a new Dialogflow CX agent.
 *
 * An agent is a virtual agent that handles conversations with your end-users.
 * It is a core component of Dialogflow CX.
 *
 * @param {string} projectId The Google Cloud project ID. Example: 'my-project-id'
 * @param {string} locationId The Google Cloud location (region) of the agent.
 *     Example: 'global', 'us-central1'
 * @param {string} displayName The human-readable name of the agent. Example: 'my-new-cx-agent'
 */
async function createAgent(
  projectId,
  locationId = 'global',
  displayName = 'my-new-cx-agent'
) {
  // Construct the parent path for the agent.
  // Format: projects/<Project ID>/locations/<Location ID>
  const parent = client.locationPath(projectId, locationId);

  // Define the agent to be created.
  // Hardcoded API parameters as per instructions.
  const agent = {
    displayName: displayName,
    defaultLanguageCode: 'en', // Example: English
    timeZone: 'America/New_York', // Example: New York time zone
  };

  const request = {
    parent: parent,
    agent: agent,
  };

  try {
    const [response] = await client.createAgent(request);

    // Print the details of the created agent.
    console.log(`Agent created: ${response.name}`);
    console.log(`Display Name: ${response.displayName}`);
    console.log(`Default Language: ${response.defaultLanguageCode}`);
    console.log(`Time Zone: ${response.timeZone}`);
  } catch (err) {
    // Handle specific API errors.
    if (err.code === status.ALREADY_EXISTS) {
      console.warn(`Agent with display name '${displayName}' already exists in location '${locationId}'.`);
      console.warn(`Please try a different display name or delete the existing agent.`);
    } else {
      console.error('Error creating agent:', err);
      // Re-throw the error for the command-line runner to handle.
      throw err;
    }
  }
}
// [END dialogflow_v3_agents_agent_create_async]

async function main(args) {
  if (args.length < 3) {
    console.error(`Usage: node ${process.argv[1]} <projectId> <locationId> <displayName>`);
    process.exit(1);
  }

  const projectId = args[0];
  const locationId = args[1];
  const displayName = args[2];

  await createAgent(projectId, locationId, displayName);
}

if (require.main === module) {
  main(process.argv.slice(2)).catch(err => {
    console.error(`Error running sample: ${err.message}`);
    process.exitCode = 1;
  });
}

module.exports = {
  createAgent,
};
