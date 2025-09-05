const assert = require('node:assert');
const {describe, it, mock} = require('node:test');
const {AgentsClient} = require('@google-cloud/dialogflow-cx');

const projectId = process.env.GOOGLE_CLOUD_PROJECT || 'test-project';
const location = 'global';
const agentDisplayName = 'test-agent';


describe('createAgent-async', () => {
  describe('mocked unit tests', () => {
    it('should write to stdout a string including the resource name', async t => {
      const mockAgent = {
        name: `projects/${projectId}/locations/${location}/agents/${agentDisplayName}`,
        displayName: agentDisplayName,
        defaultLanguageCode: 'en',
        timeZone: 'America/Los_Angeles',
      };

      const mockCreateAgent = t.mock.method(
        AgentsClient.prototype, // Corrected this line
        'createAgent',
        async () => [
          mockAgent,
        ]
      );

      const mockedLog = t.mock.method(console, 'log');
      const sample = require('../agents-client-agent-create-async.js');
      await sample.createAgent(projectId, location, agentDisplayName);

      const callCount = mockCreateAgent.mock.callCount();
      const logs = mockedLog.mock.calls
        .map(log => log.arguments.join('\n'))
        .join('\n');
      t.mock.restoreAll();

      assert.strictEqual(callCount, 1, 'should have called createAgent once');
      const expected = `Agent created: ${agentDisplayName} (${mockAgent.name})`;
      assert.ok(
        logs.includes(expected),
        `Failed: Sample output does not include '${expected}', instead shows: ${logs}`
      );
    });

    it('should handle ALREADY_EXISTS error', async t => {
      const mockCreateAgent = t.mock.method(
        AgentsClient.prototype,
        'createAgent',
        async () => {
          const error = new Error('Agent already exists');
          error.code = 6; // status.ALREADY_EXISTS
          throw error;
        }
      );

      const mockedWarn = t.mock.method(console, 'warn');
      const sample = require('../agents-client-agent-create-async.js');
      await sample.createAgent(projectId, location, agentDisplayName);

      const callCount = mockCreateAgent.mock.callCount();
      const warns = mockedWarn.mock.calls
        .map(warn => warn.arguments.join('\n'))
        .join('\n');
      t.mock.restoreAll();

      assert.strictEqual(callCount, 1, 'should have called createAgent once');
      const expected = `Agent with display name "${agentDisplayName}" already exists in location "${location}" of project "${projectId}".`;
      assert.ok(
        warns.includes(expected),
        `Failed: Sample output does not include '${expected}', instead shows: ${warns}`
      );
    });

    it('should re-throw unexpected errors', async t => {
      const mockCreateAgent = t.mock.method(
        AgentsClient.prototype,
        'createAgent',
        async () => {
          const error = new Error('Unexpected error');
          error.code = 13; // Example: status.INTERNAL
          throw error;
        }
      );

      const mockedError = t.mock.method(console, 'error');
      const sample = require('../agents-client-agent-create-async.js');

      await assert.rejects(
        sample.createAgent(projectId, location, agentDisplayName),
        /Unexpected error/
      );

      const callCount = mockCreateAgent.mock.callCount();
      const errors = mockedError.mock.calls
        .map(err => err.arguments.join('\n'))
        .join('\n');
      t.mock.restoreAll();

      assert.strictEqual(callCount, 1, 'should have called createAgent once');
      assert.ok(
          errors.includes(`Error creating agent "${agentDisplayName}":`),
          `Error message not found in console.error output: ${errors}`
      );

    });
  });
});
