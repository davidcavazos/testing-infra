import pytest
import unittest
from unittest.mock import patch, MagicMock
import os

from google.api_core.exceptions import NotFound, PermissionDenied, InvalidArgument

# Import the code under test
import sys

sys.path.append('..')
from agents_client_agents_list import list_dialogflow_cx_agents


@patch("google.cloud.dialogflowcx_v3.AgentsClient", autospec=True)
def test_list_dialogflow_cx_agents_success(mock_agents_client, capsys):
    # Configure the mock client to return a list of agents
    mock_client = mock_agents_client.return_value
    mock_agent = MagicMock()
    mock_agent.name = "agent1"
    mock_agent.display_name = "Agent One"
    mock_agent.default_language_code = "en"
    mock_agent.time_zone = "UTC"
    mock_client.list_agents.return_value = [mock_agent]

    project_id = "test-project"
    location = "global"

    list_dialogflow_cx_agents(project_id, location)

    mock_agents_client.assert_called_once()
    mock_client.common_location_path.assert_called_with(project_id, location)
    mock_client.list_agents.assert_called_with(parent=mock_client.common_location_path.return_value)

    captured = capsys.readouterr()
    assert f"Attempting to list agents for project: {project_id}, location: {location}" in captured.out
    assert "Agent Name: agent1" in captured.out
    assert "Display Name: Agent One" in captured.out
    assert "Default Language Code: en" in captured.out
    assert "Time Zone: UTC" in captured.out


@patch("google.cloud.dialogflowcx_v3.AgentsClient", autospec=True)
def test_list_dialogflow_cx_agents_no_agents(mock_agents_client, capsys):
    # Configure the mock client to return an empty list of agents
    mock_client = mock_agents_client.return_value
    mock_client.list_agents.return_value = []

    project_id = "test-project"
    location = "global"

    list_dialogflow_cx_agents(project_id, location)

    captured = capsys.readouterr()
    assert f"No agents found in project '{project_id}' at location '{location}'." in captured.out


@patch("google.cloud.dialogflowcx_v3.AgentsClient", autospec=True)
def test_list_dialogflow_cx_agents_not_found(mock_agents_client, capsys):
    # Configure the mock client to raise a NotFound exception
    mock_client = mock_agents_client.return_value
    mock_client.list_agents.side_effect = NotFound("Location not found")

    project_id = "test-project"
    location = "global"

    list_dialogflow_cx_agents(project_id, location)

    captured = capsys.readouterr()
    assert f"Error: The specified location '{location}' does not exist for project '{project_id}'." in captured.out


@patch("google.cloud.dialogflowcx_v3.AgentsClient", autospec=True)
def test_list_dialogflow_cx_agents_permission_denied(mock_agents_client, capsys):
    # Configure the mock client to raise a PermissionDenied exception
    mock_client = mock_agents_client.return_value
    mock_client.list_agents.side_effect = PermissionDenied("Permission denied")

    project_id = "test-project"
    location = "global"

    list_dialogflow_cx_agents(project_id, location)

    captured = capsys.readouterr()
    assert f"Error: Permission denied to list agents in project '{project_id}' at location '{location}'." in captured.out


@patch("google.cloud.dialogflowcx_v3.AgentsClient", autospec=True)
def test_list_dialogflow_cx_agents_invalid_argument(mock_agents_client, capsys):
    # Configure the mock client to raise an InvalidArgument exception
    mock_client = mock_agents_client.return_value
    mock_client.list_agents.side_effect = InvalidArgument("Invalid argument")

    project_id = "test-project"
    location = "global"

    list_dialogflow_cx_agents(project_id, location)

    captured = capsys.readouterr()
    assert "Error: Invalid argument provided." in captured.out


@patch("google.cloud.dialogflowcx_v3.AgentsClient", autospec=True)
def test_list_dialogflow_cx_agents_unexpected_error(mock_agents_client, capsys):
    # Configure the mock client to raise an unexpected exception
    mock_client = mock_agents_client.return_value
    mock_client.list_agents.side_effect = Exception("Unexpected error")

    project_id = "test-project"
    location = "global"

    list_dialogflow_cx_agents(project_id, location)

    captured = capsys.readouterr()
    assert "An unexpected error occurred: Unexpected error" in captured.out
