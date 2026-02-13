import * as vscode from 'vscode';

async function pollLogs(output: vscode.OutputChannel, runId: string) {
  const addr = vscode.workspace.getConfiguration('openagent').get<string>('backendUrl') || 'http://127.0.0.1:7341';
  let lastLen = 0;
  let done = false;

  while (!done) {
    try {
      const res = await fetch(`${addr}/logs?runId=${encodeURIComponent(runId)}`);
      const json = await res.json();
      const lines: string[] = json.lines || [];
      done = Boolean(json.done);

      if (lines.length > lastLen) {
        for (const line of lines.slice(lastLen)) {
          output.appendLine(line);
        }
        lastLen = lines.length;
      }
    } catch (err: any) {
      output.appendLine(`Log polling failed: ${err?.message || err}`);
      return;
    }
    await new Promise((r) => setTimeout(r, 500));
  }
}

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
    await pollLogs(output, json.runId);
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
