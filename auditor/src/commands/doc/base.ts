/**
 * Shared utilities for doc generate and update commands
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import inquirer from 'inquirer';
import { GlobalOptions, DocumentationQuality, ModulesData, FunctionalArea } from '../../types';
import { displayDocumentationSummary } from '../../utils/display';
import { runAgent } from '../../utils/agent-runner';
import { loadPromptWithPartials } from '../../utils/prompt-loader';

export interface DocGenerateOptions {
  file: string;
  quality: string;
  source: string;
  output: string;
  minScore?: string;
  name?: string[];
}

/**
 * Common setup for doc generate/update commands
 */
export async function setupDocCommand(
  options: DocGenerateOptions,
  cwdPath: string,
  globalOpts: GlobalOptions
) {
  // Get global options for dry-run
  const dryRun = globalOpts.dryRun || false;
  
  // Store global options for agent runner
  (global as any).globalOptions = globalOpts;
  
  // Resolve paths
  const workingDir = fs.realpathSync(cwdPath);
  const resolvedModulesPath = path.isAbsolute(options.file)
    ? path.resolve(options.file)
    : path.resolve(workingDir, options.file);
  const resolvedQualityPath = path.isAbsolute(options.quality)
    ? path.resolve(options.quality)
    : path.resolve(workingDir, options.quality);
  const resolvedSourcePath = path.isAbsolute(options.source)
    ? path.resolve(options.source)
    : path.resolve(workingDir, options.source);
  const resolvedOutputPath = path.isAbsolute(options.output)
    ? path.resolve(options.output)
    : path.resolve(workingDir, options.output);
  
  // Validate files exist
  if (!fs.existsSync(resolvedModulesPath)) {
    console.error(chalk.red(`❌ Modules file not found: ${resolvedModulesPath}`));
    console.error(chalk.gray('Run `explore` command first.'));
    process.exit(1);
  }
  
  if (!fs.existsSync(resolvedQualityPath)) {
    console.error(chalk.red(`❌ Quality assessment not found: ${resolvedQualityPath}`));
    console.error(chalk.gray('Run `doc review` command first.'));
    process.exit(1);
  }
  
  if (!fs.existsSync(resolvedSourcePath)) {
    console.error(chalk.red(`❌ Source path not found: ${resolvedSourcePath}`));
    process.exit(1);
  }
  
  // Load data
  const modules: ModulesData = JSON.parse(fs.readFileSync(resolvedModulesPath, 'utf-8'));
  const quality: DocumentationQuality = JSON.parse(fs.readFileSync(resolvedQualityPath, 'utf-8'));
  
  // Validate data consistency
  const moduleNames = new Set(modules.functional_areas.map(m => m.name));
  const qualityNames = new Set(quality.functional_areas.map(q => q.name));
  
  for (const name of qualityNames) {
    if (!moduleNames.has(name)) {
      console.error(chalk.red(`❌ Inconsistency: "${name}" in quality file but not in modules`));
      console.error(chalk.gray('Regenerate documentation quality with `doc review`.'));
      process.exit(1);
    }
  }
  
  return {
    dryRun,
    workingDir,
    resolvedModulesPath,
    resolvedQualityPath,
    resolvedSourcePath,
    resolvedOutputPath,
    modules,
    quality
  };
}

/**
 * Filter areas to document based on options
 */
export function filterAreasToDocument(
  quality: DocumentationQuality,
  minScore?: number,
  areaNames?: string[]
): FunctionalArea[] {
  let areasToDocument: FunctionalArea[];
  
  if (areaNames && areaNames.length > 0) {
    // Filter by specific names
    const nameSet = new Set(areaNames.map(n => n.toLowerCase()));
    areasToDocument = quality.functional_areas.filter(area => 
      nameSet.has(area.name.toLowerCase())
    );
    
    // Check if all requested names were found
    const foundNames = new Set(areasToDocument.map(a => a.name.toLowerCase()));
    const notFound = areaNames.filter(n => !foundNames.has(n.toLowerCase()));
    
    if (notFound.length > 0) {
      console.log(chalk.yellow(`\n⚠️  Areas not found: ${notFound.join(', ')}`));
      console.log(chalk.gray('Available areas:'));
      quality.functional_areas.forEach(area => {
        console.log(chalk.gray(`  - ${area.name}`));
      });
    }
    
    console.log(chalk.cyan(`\n📋 Will process ${areasToDocument.length} specific areas: ${areasToDocument.map(a => a.name).join(', ')}`));
  } else {
    // Filter by importance score
    if (!minScore) {
      minScore = 3;
      console.log(chalk.gray(`Using default minimum score: ${minScore}`));
    }
    
    areasToDocument = quality.functional_areas.filter(area => 
      area.doc_importance_score >= minScore!
    );
    
    console.log(chalk.cyan(`\n📋 Will process ${areasToDocument.length} areas with importance ≥ ${minScore}`));
  }
  
  return areasToDocument;
}

/**
 * Select CLI tool for execution
 */
export async function selectCLITool(dryRun: boolean | string): Promise<string> {
  if (dryRun) {
    return 'claude'; // default for dry-run
  }
  
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
  
  return response.cliTool;
}

/**
 * Generate/update documentation index
 */
export async function updateDocumentationIndex(
  outputPath: string,
  modulesPath: string,
  workingDir: string,
  cliTool: string,
  dryRun: boolean | string
) {
  if (dryRun) return;
  
  console.log(chalk.cyan('\n📚 Updating documentation index...'));
  
  const promptsDir = path.join(path.dirname(path.dirname(path.dirname(__dirname))), 'prompts');
  const indexPrompt = fs.readFileSync(path.join(promptsDir, '05-doc-update-index.md'), 'utf-8');
  
  const indexResult = await runAgent({
    tool: cliTool as 'claude' | 'gemini',
    prompts: [
      { name: '05-doc-update-index.md', content: indexPrompt }
    ],
    variables: {
      parameters: {
        doc_path: path.join(outputPath, 'doc'),
        project_name: 'Project',
        modules_file: modulesPath
      }
    },
    outputPath: path.join(outputPath, 'doc', 'README.md'),
    workingDir: workingDir,
    dryRun: false, // Never dry-run the index generation
    description: 'Update documentation index'
  });
  
  if (indexResult.success) {
    console.log(chalk.green('✅ Documentation index updated'));
  } else {
    console.log(chalk.red('❌ Failed to update documentation index'));
  }
}