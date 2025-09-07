/**
 * List command - displays functional areas from modules.json
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import { Command } from 'commander';
import { ModulesData, DocumentationQuality } from '../types';
import { displayModulesTable } from '../utils/display-table';

/**
 * Register the list command
 */
export function registerListCommand(program: Command) {
  program
    .command('list')
    .description('List functional areas from a modules.json file')
    .option('--file <path>', 'Path to modules.json file', 'audit_result/modules.json')
    .option('--detailed', 'Show detailed information including area type, complexity, and metrics')
    .option('--tech', 'Show technologies used by each functional area')
    .action(async (options: { file: string; detailed?: boolean; tech?: boolean }) => {
      const globalOpts = program.opts();
      const cwd = globalOpts.cwd ? path.resolve(process.cwd(), globalOpts.cwd) : process.cwd();
      await runList(options.file, cwd, options.detailed, options.tech);
    });
}

/**
 * Run the list command
 */
async function runList(filePath: string, cwdPath: string, detailed?: boolean, showTech?: boolean) {
  // Resolve file path
  const workingDir = fs.realpathSync(cwdPath);
  const resolvedPath = path.isAbsolute(filePath)
    ? path.resolve(filePath)
    : path.resolve(workingDir, filePath);
  
  // Check if file exists
  if (!fs.existsSync(resolvedPath)) {
    console.error(chalk.red(`❌ File not found: ${resolvedPath}`));
    console.error(chalk.gray('Run the explore command first to generate modules.json'));
    process.exit(1);
  }
  
  // Load and parse the modules JSON file
  let modulesData: ModulesData;
  try {
    const content = fs.readFileSync(resolvedPath, 'utf-8');
    modulesData = JSON.parse(content);
  } catch (error) {
    console.error(chalk.red(`❌ Failed to parse JSON file: ${error}`));
    process.exit(1);
  }
  
  // Check for documentation quality file (be smart by default)
  const outputDir = path.dirname(resolvedPath);
  const qualityPath = path.join(outputDir, 'documentation_quality.json');
  let qualityData: DocumentationQuality | undefined;
  
  if (fs.existsSync(qualityPath)) {
    try {
      const qualityContent = fs.readFileSync(qualityPath, 'utf-8');
      qualityData = JSON.parse(qualityContent);
      console.log(chalk.gray(`📚 Documentation quality data found with ${qualityData?.functional_areas?.length || 0} areas`));
    } catch (error) {
      console.log(chalk.yellow('⚠️  Found documentation_quality.json but could not parse it'));
    }
  } else {
    // Also check in the current directory if not found in output dir
    const altQualityPath = path.join(workingDir, 'documentation_quality.json');
    if (fs.existsSync(altQualityPath)) {
      try {
        const qualityContent = fs.readFileSync(altQualityPath, 'utf-8');
        qualityData = JSON.parse(qualityContent);
        console.log(chalk.gray(`📚 Documentation quality data found at ${altQualityPath}`));
      } catch (error) {
        console.log(chalk.yellow('⚠️  Found documentation_quality.json but could not parse it'));
      }
    }
  }
  
  console.log(chalk.gray(`📁 Source: ${resolvedPath}`));
  console.log(chalk.gray(`📊 Total Functional Areas: ${modulesData.functional_areas.length}\n`));
  
  // Use table display function - show tech view if requested
  displayModulesTable(modulesData, qualityData, detailed, showTech);
  
  // Additional tips if no documentation quality
  if (!qualityData) {
    console.log(chalk.gray('\n💡 Tip: Run "doc review" to assess documentation quality for each area'));
  }
}