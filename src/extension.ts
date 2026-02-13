import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
  const output = vscode.window.createOutputChannel('OpenAgent');
  context.subscriptions.push(output);

  const runCmd = vscode.commands.registerCommand('openagent.runWorkflow', async () => {
    output.show(true);
    output.appendLine('OpenAgent: Run workflow (stub)');
    vscode.window.showInformationMessage('OpenAgent: run workflow (stub).');
  });

  const configCmd = vscode.commands.registerCommand('openagent.openConfig', async () => {
    const doc = await vscode.workspace.openTextDocument({
      language: 'yaml',
      content: '# OpenAgent config (stub)\nframework: langgraph\n',
    });
    await vscode.window.showTextDocument(doc, { preview: false });
  });

  context.subscriptions.push(runCmd, configCmd);
}

export function deactivate() {}
