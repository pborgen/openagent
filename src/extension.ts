import * as vscode from 'vscode';

async function callBackendRun(output: vscode.OutputChannel) {
  const addr = vscode.workspace.getConfiguration('openagent').get<string>('backendUrl') || 'http://127.0.0.1:7341';
  const body = JSON.stringify({ workflow: 'default', params: {} });

  try {
    const res = await fetch(`${addr}/run`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body,
    });
    const json = await res.json();
    output.appendLine(`Run started: ${JSON.stringify(json)}`);
    vscode.window.showInformationMessage(`OpenAgent: run started (${json.runId})`);
  } catch (err: any) {
    output.appendLine(`Run failed: ${err?.message || err}`);
    vscode.window.showErrorMessage('OpenAgent: backend not reachable. Start the Go server.');
  }
}

export function activate(context: vscode.ExtensionContext) {
  const output = vscode.window.createOutputChannel('OpenAgent');
  context.subscriptions.push(output);

  const runCmd = vscode.commands.registerCommand('openagent.runWorkflow', async () => {
    output.show(true);
    output.appendLine('OpenAgent: Run workflow');
    await callBackendRun(output);
  });

  const configCmd = vscode.commands.registerCommand('openagent.openConfig', async () => {
    const doc = await vscode.workspace.openTextDocument({
      language: 'yaml',
      content: '# OpenAgent config\nframework: langgraph\n',
    });
    await vscode.window.showTextDocument(doc, { preview: false });
  });

  context.subscriptions.push(runCmd, configCmd);
}

export function deactivate() {}
