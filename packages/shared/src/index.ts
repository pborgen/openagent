export type RunRequest = {
  workflow: string;
  params: Record<string, string>;
};

export type RunResponse = {
  runId: string;
};

export type LogResponse = {
  lines: string[];
  done: boolean;
};

export type OpenAgentConfig = {
  framework: string;
  workflow: string;
};
