#!/usr/bin/env node

/**
 * Codebase Auditor - AI-driven codebase analysis tool
 * Main entry point
 */

import { Command } from 'commander';
import { registerExploreCommand } from './commands/explore';
import { registerListCommand } from './commands/list';
import { registerDocCommands } from './commands/doc';
import { registerAuditCommand } from './commands/audit';

// Create the main program
const program = new Command();

program
  .name('audit')
  .description('AI-driven codebase analysis tool')
  .version('1.0.0')
  .option('--cwd <path>', 'Working directory to launch CLI from (default: current directory)')
  .option('--dry-run [level]', 'Show what would be executed without actually running the agent (0=minimal, 1=normal, 2=verbose with prompt)')
  .option('--raw', 'Show raw output from Claude without filtering (Claude only)')
  .option('--agent-output <file>', 'Write raw agent output (stdout+stderr) to file');

// Register all commands
registerExploreCommand(program);
registerListCommand(program);
registerDocCommands(program);
registerAuditCommand(program);

// Main function to run the CLI
export function main() {
  // Parse command line arguments
  program.parse(process.argv);

  // Show help if no command provided
  if (!process.argv.slice(2).length) {
    program.outputHelp();
  }
}

// Only run CLI if this is the main module
if (import.meta.url === `file://${process.argv[1]}`) {
  main();
}