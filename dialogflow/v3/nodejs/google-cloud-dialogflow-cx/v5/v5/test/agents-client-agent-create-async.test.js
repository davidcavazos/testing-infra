const assert = require('node:assert');
const {describe, it} = require('node:test');
const {AgentsClient} = require('@google-cloud/dialogflow-cx');
const {status} = require('@grpc/grpc-js');

const {createAgent} = require('../agents-client-agent-create-async.js');

describe('CreateAgent-async', () => {
  const projectId = 'your-project-id';
  const locationId = 'us-central1';
  const displayName = 'my-new-cx-agent';
  const agentName = `projects/${projectId}/locations/${locationId}/agents/${displayName}`;

  describe('mocked unit tests', () => {
    it('should write to stdout a string including the resource name', async t => {
      const mockCreateAgent = t.mock.method(
        AgentsClient.prototype, // Corrected this line
        'createAgent',
        async () => [
          {
            name: agentName,
            displayName: displayName,
            defaultLanguageCode: 'en',
            timeZone: 'America/New_York',
          },
        ]
      );

      const mockedLog = t.mock.method(console, 'log');

      await createAgent(projectId, locationId, displayName);

      const callCount = mockCreateAgent.mock.callCount();
      const logs = mockedLog.mock.calls
        .map(log => log.arguments.join('\n'))
        .join('\n');
      t.mock.restoreAll();

      assert.ok(
        callCount > 0,
        'Failed: Mocked console.log method not called at least once'
      );
      const expected = `Agent created: ${agentName}`;
      assert.ok(
        logs.includes(expected),
        `Failed: Sample output does not include '${expected}', instead shows: ${logs}`
      );
    });

    it('should handle ALREADY_EXISTS error and log a warning', async t => {
      const mockCreateAgent = t.mock.method(
        AgentsClient.prototype, // Corrected this line
        'createAgent',
        async () => {
          const error = new Error('Agent already exists');
          error.code = status.ALREADY_EXISTS;
          throw error;
        }
      );

      const mockedLog = t.mock.method(console, 'warn');

      await createAgent(projectId, locationId, displayName);

      const callCount = mockCreateAgent.mock.callCount();
      const logs = mockedLog.mock.calls
        .map(log => log.arguments.join('\n'))
        .join('\n');
      t.mock.restoreAll();

      assert.ok(
        callCount > 0,
        'Failed: Mocked createAgent method not called at least once'
      );
      const expected = `Agent with display name '${displayName}' already exists in location '${locationId}'.`;
      assert.ok(
        logs.includes(expected),
        `Failed: Sample output does not include '${expected}', instead shows: ${logs}`
      );
    });

    it('should handle other errors and re-throw', async t => {
      const mockCreateAgent = t.mock.method(
        AgentsClient.prototype,
        'createAgent',
        async () => {
          const error = new Error('Some other error');
          error.code = status.PERMISSION_DENIED;
          throw error;
        }
      );

      const mockedError = t.mock.method(console, 'error');

      await assert.rejects(
        createAgent(projectId, locationId, displayName),
        /Some other error/
      );

      const callCount = mockCreateAgent.mock.callCount();
      t.mock.restoreAll();

      assert.ok(
        callCount > 0,
        'Failed: Mocked createAgent method not called at least once'
      );
    });
  });
});
