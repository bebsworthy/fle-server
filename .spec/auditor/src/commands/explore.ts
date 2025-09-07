/**
 * Explore command - analyzes codebase to identify functional areas
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import inquirer from 'inquirer';
import { Command } from 'commander';
import { GlobalOptions } from '../types';
import { runAgent } from '../utils/agent-runner';
import { loadPromptWithPartials } from '../utils/prompt-loader';

/**
 * Register the explore command
 */
export function registerExploreCommand(program: Command) {
  program
    .command('explore')
    .description('Explore and analyze the codebase to identify functional areas')
    .option('--doc <path>', 'Path to documentation directory (optional, default: none)')
    .option('--output <path>', 'Path to output directory (default: <cwd>/audit_result)')
    .argument('<source>', 'Path to source code directory')
    .action(async (source: string, options: { doc?: string; output?: string }) => {
      const globalOpts = program.opts() as GlobalOptions;
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      const outputPath = options.output || path.join(cwd, 'audit_result');
      await runExplore(source, options.doc, outputPath, cwd, globalOpts);
    });
}

/**
 * Run the explore command
 */
async function runExplore(
  sourcePath: string, 
  docPath: string | undefined, 
  outputPath: string, 
  cwdPath: string,
  globalOpts: GlobalOptions
) {
  console.log(chalk.blue('🔍 Codebase Auditor - Explore Mode'));
  console.log(chalk.gray('─'.repeat(50)));
  
  // Get global options for dry-run
  const dryRun = globalOpts.dryRun || false;
  
  // Store global options for agent runner
  (global as any).globalOptions = globalOpts;
  
  // Convert paths to canonical absolute paths
  const canonicalWorkingDir = fs.realpathSync(cwdPath);
  
  // Resolve source path
  sourcePath = path.isAbsolute(sourcePath)
    ? path.resolve(sourcePath)
    : path.resolve(canonicalWorkingDir, sourcePath);
    
  if (!fs.existsSync(sourcePath)) {
    console.error(chalk.red(`❌ Source path does not exist: ${sourcePath}`));
    console.error(chalk.gray(`  Working directory: ${canonicalWorkingDir}`));
    console.error(chalk.gray(`  Original source argument: ${sourcePath}`));
    process.exit(1);
  }
  const canonicalSourcePath = fs.realpathSync(sourcePath);
  
  // Resolve documentation path if provided
  let canonicalDocPath: string | undefined;
  if (docPath) {
    docPath = path.isAbsolute(docPath)
      ? path.resolve(docPath)
      : path.resolve(canonicalWorkingDir, docPath);
      
    if (!fs.existsSync(docPath)) {
      console.error(chalk.red(`❌ Documentation path does not exist: ${docPath}`));
      process.exit(1);
    }
    canonicalDocPath = fs.realpathSync(docPath);
  }
  
  // Output path might not exist yet, so just resolve it
  const canonicalOutputPath = outputPath;
  
  // Display resolved canonical paths
  console.log(chalk.gray(`📂 Working directory: ${canonicalWorkingDir}`));
  console.log(chalk.gray(`📄 Source path: ${canonicalSourcePath}`));
  console.log(chalk.gray(`📚 Documentation: ${canonicalDocPath || 'Not provided'}`));
  console.log(chalk.gray(`💾 Output path: ${canonicalOutputPath}`));
  
  // Create output directory if it doesn't exist (unless dry-run)
  if (!dryRun && !fs.existsSync(canonicalOutputPath)) {
    console.log(chalk.yellow(`📁 Creating output directory: ${canonicalOutputPath}`));
    fs.mkdirSync(canonicalOutputPath, { recursive: true });
  }
  
  // Select CLI tool
  let cliTool = 'claude'; // default for dry-run
  if (!dryRun) {
    const response = await inquirer.prompt([
      {
        type: 'list',
        name: 'cliTool',
        message: 'Select the AI CLI tool to use:',
        choices: [
          { name: 'Claude CLI', value: 'claude' },
          { name: 'Gemini CLI', value: 'gemini' }
        ]
      }
    ]);
    cliTool = response.cliTool;
  }
  
  console.log(chalk.cyan(`\n🤖 Using ${cliTool === 'claude' ? 'Claude' : 'Gemini'} CLI for analysis\n`));
  
  // Check if modules.json already exists for automatic update mode
  const modulesJsonPath = path.join(canonicalOutputPath, 'modules.json');
  const isUpdateMode = fs.existsSync(modulesJsonPath);
  let existingModulesData = '';
  
  if (isUpdateMode) {
    console.log(chalk.yellow('📝 Existing modules.json detected - running in UPDATE mode'));
    console.log(chalk.gray('   Will preserve existing data and merge changes'));
    existingModulesData = fs.readFileSync(modulesJsonPath, 'utf-8');
  } else {
    console.log(chalk.cyan('🆕 No existing modules.json - running in CREATE mode'));
  }
  
  // Load prompt template with partials
  const currentDir = path.dirname(new URL(import.meta.url).pathname);
  const promptsDir = path.join(currentDir, '..', '..', 'prompts');
  const explorePromptPath = path.join(promptsDir, '01-explore.md');
  const promptData = loadPromptWithPartials(explorePromptPath);
  
  // Run the exploration using the agent runner
  console.log(chalk.blue('🚀 Starting codebase exploration...\n'));
  
  const result = await runAgent({
    tool: cliTool as 'claude' | 'gemini',
    prompts: [
      { name: promptData.name, content: promptData.content, partials: promptData.partials }
    ],
    variables: {
      parameters: {
        code_repository: canonicalSourcePath,
        documentation_repository: canonicalDocPath || 'None',
        output_path: modulesJsonPath,
        update_mode: isUpdateMode ? 'true' : 'false',
        existing_modules: isUpdateMode ? existingModulesData : ''
      }
    },
    outputPath: modulesJsonPath,
    workingDir: canonicalWorkingDir,
    dryRun: dryRun,
    description: isUpdateMode ? 'Update codebase analysis' : 'Explore and analyze the codebase'
  });
  
  if (!dryRun) {
    if (result.success) {
      console.log(chalk.green('\n✅ Exploration complete!'));
      console.log(chalk.gray(`Results saved to: ${path.join(canonicalOutputPath, 'modules.json')}`));
    } else {
      console.error(chalk.red('\n❌ Exploration failed'));
      process.exit(1);
    }
  }
}