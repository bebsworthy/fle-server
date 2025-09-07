/**
 * Audit command registration
 */

import { Command } from 'commander';
import * as path from 'path';
import { GlobalOptions } from '../../types';
import { CodeAuditOptions } from '../../types/audit';
import { runCodeAudit } from './audit';

/**
 * Register the audit command
 */
export function registerAuditCommand(program: Command) {
  program
    .command('audit')
    .description('Audit code quality and architecture for functional areas')
    .option('--file <path>', 'Path to modules.json file', 'audit_result/modules.json')
    .option('--source <path>', 'Path to source code directory', '.')
    .option('--output <path>', 'Path to output directory', 'audit_result')
    .option('--min-complexity <number>', 'Minimum complexity score to audit (1-5)')
    .option('--module <name>', 'Filter areas by module (for monorepos)')
    .option('--name <names...>', 'Audit specific areas by name')
    .option('--focus <aspect>', 'Focus on specific aspect (quality|architecture|security|performance)')
    .action(async (options: CodeAuditOptions) => {
      const globalOpts = program.opts() as GlobalOptions;
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      await runCodeAudit(options, cwd, globalOpts);
    });
}