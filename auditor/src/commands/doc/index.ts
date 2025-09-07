/**
 * Doc command orchestrator - registers all documentation subcommands
 */

import { Command } from 'commander';
import * as path from 'path';
import { GlobalOptions } from '../../types';
import { runDocSummary } from './summary';
import { runDocReview } from './review';
import { runDocGenerate } from './generate';
import { runDocUpdate } from './update';

/**
 * Register all doc commands
 */
export function registerDocCommands(program: Command) {
  // Doc command (parent command - shows summary when no subcommand)
  const docCommand = program
    .command('doc')
    .description('Documentation analysis and generation')
    .action(async () => {
      // Get the global cwd option
      const globalOpts = program.opts() as GlobalOptions;
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      await runDocSummary(cwd);
    });

  // Doc review subcommand
  docCommand
    .command('review')
    .description('Review and score documentation quality for each functional area')
    .option('--file <path>', 'Path to modules.json file', 'audit_result/modules.json')
    .option('--output <path>', 'Path to output directory', 'audit_result')
    .action(async (options: { file: string; output: string }) => {
      const globalOpts = program.opts() as GlobalOptions;
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      await runDocReview(options.file, options.output, cwd, globalOpts);
    });

  // Doc generate subcommand
  docCommand
    .command('generate')
    .description('Generate documentation for functional areas')
    .option('--file <path>', 'Path to modules.json file', 'audit_result/modules.json')
    .option('--quality <path>', 'Path to documentation_quality.json', 'audit_result/documentation_quality.json')
    .option('--source <path>', 'Path to source code directory', '.')
    .option('--output <path>', 'Path to output directory', 'audit_result')
    .option('--min-score <number>', 'Minimum importance score to document (1-5)')
    .option('--name <names...>', 'Generate documentation for specific areas by name')
    .action(async (options: { file: string; quality: string; source: string; output: string; minScore?: string; name?: string[] }) => {
      const globalOpts = program.opts() as GlobalOptions;
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      const minScore = options.minScore ? parseInt(options.minScore) : undefined;
      await runDocGenerate(options.file, options.quality, options.source, options.output, cwd, minScore, options.name, globalOpts);
    });

  // Doc update subcommand (incremental updates)
  docCommand
    .command('update')
    .description('Update existing documentation or create new documentation with incremental update support')
    .option('--file <path>', 'Path to modules.json file', 'audit_result/modules.json')
    .option('--quality <path>', 'Path to documentation_quality.json', 'audit_result/documentation_quality.json')
    .option('--source <path>', 'Path to source code directory', '.')
    .option('--output <path>', 'Path to output directory', 'audit_result')
    .option('--min-score <number>', 'Minimum importance score to document (1-5)')
    .option('--name <names...>', 'Update documentation for specific areas by name')
    .action(async (options: { file: string; quality: string; source: string; output: string; minScore?: string; name?: string[] }) => {
      const globalOpts = program.opts() as GlobalOptions;
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      const minScore = options.minScore ? parseInt(options.minScore) : undefined;
      await runDocUpdate(options.file, options.quality, options.source, options.output, cwd, minScore, options.name, globalOpts);
    });

  return docCommand;
}