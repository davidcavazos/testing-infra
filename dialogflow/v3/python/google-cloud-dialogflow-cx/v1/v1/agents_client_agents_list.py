# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import argparse

# [START dialogflow_v3_agents_list_agents]
from google.cloud import dialogflowcx_v3
from google.api_core.exceptions import NotFound, PermissionDenied, InvalidArgument

def list_dialogflow_cx_agents(
    project_id: str,
    location: str,
) -> None:
    """Lists all Dialogflow CX agents in the specified Google Cloud project and location.

    Dialogflow CX agents are virtual agents that handle conversations with end-users.
    Listing them allows you to manage your deployed conversational AI models and
    verify their presence and basic configuration.

    Args:
        project_id: The Google Cloud project ID.
        location: The Google Cloud location (region) of the agents.
    """
    client = dialogflowcx_v3.AgentsClient()

    # The parent resource name, for example,
    # "projects/PROJECT_ID/locations/LOCATION_ID"
    # The common_location_path helper constructs the correct resource name.
    parent = client.common_location_path(project_id, location)

    print(f"Attempting to list agents for project: {project_id}, location: {location}")

    try:
        # Make the request to list agents. This method returns an iterable
        # that handles pagination automatically.
        agents_found = False
        for agent in client.list_agents(parent=parent):
            agents_found = True
            print(f"Agent Name: {agent.name}")
            print(f"Display Name: {agent.display_name}")
            print(f"Default Language Code: {agent.default_language_code}")
            print(f"Time Zone: {agent.time_zone}")
            print("---------------------------------")
        
        if not agents_found:
            print(f"No agents found in project '{project_id}' at location '{location}'.")
            print("Ensure that there are agents deployed in this location or that the project ID and location are correct.")

    except NotFound:
        print(f"Error: The specified location '{location}' does not exist for project '{project_id}'.")
        print("Please verify the location ID and ensure it is valid for your project.")
    except PermissionDenied:
        print(f"Error: Permission denied to list agents in project '{project_id}' "
              f"at location '{location}'.")
        print("Please ensure your account has the 'Dialogflow CX Agent Reader' "
              "or 'Dialogflow CX API Admin' role for this project and location.")
    except InvalidArgument as e:
        print(f"Error: Invalid argument provided. This may indicate an improperly formatted "
              f"project ID or location. Details: {e}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

# [END dialogflow_v3_agents_list_agents]

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Lists Dialogflow CX agents in a given project and location."
    )
    parser.add_argument(
        "--project_id",
        type=str,
        required=True, # Make project_id a required argument
        help="The Google Cloud project ID.",
    )
    parser.add_argument(
        "--location",
        type=str,
        required=True, # Make location a required argument
        help="The Google Cloud location (region) of the agents (e.g., 'global', 'us-central1').",
    )

    args = parser.parse_args()

    list_dialogflow_cx_agents(args.project_id, args.location)
