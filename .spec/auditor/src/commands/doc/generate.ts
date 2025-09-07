/**
 * Doc generate command - generates documentation for functional areas
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import { GlobalOptions } from '../../types';
import { displayDocumentationSummary } from '../../utils/display';
import { runAgent } from '../../utils/agent-runner';
import { loadPromptWithPartials } from '../../utils/prompt-loader';
import {
  DocGenerateOptions,
  setupDocCommand,
  filterAreasToDocument,
  selectCLITool,
  updateDocumentationIndex
} from './base';

/**
 * Run the doc generate command
 */
export async function runDocGenerate(
  modulesPath: string,
  qualityPath: string,
  sourcePath: string,
  outputPath: string,
  cwdPath: string,
  minScore: number | undefined,
  areaNames: string[] | undefined,
  globalOpts: GlobalOptions
) {
  console.log(chalk.blue('📝 Codebase Auditor - Documentation Generation'));
  console.log(chalk.gray('─'.repeat(50)));
  
  // Setup command
  const setup = await setupDocCommand(
    { file: modulesPath, quality: qualityPath, source: sourcePath, output: outputPath },
    cwdPath,
    globalOpts
  );
  
  const {
    dryRun,
    workingDir,
    resolvedModulesPath,
    resolvedSourcePath,
    resolvedOutputPath,
    modules,
    quality
  } = setup;
  
  // Display summary
  displayDocumentationSummary(quality);
  
  // Filter areas to document
  const areasToDocument = filterAreasToDocument(quality, minScore, areaNames);
  
  if (areasToDocument.length === 0) {
    console.log(chalk.yellow('\n⚠️  No areas match the criteria'));
    return;
  }
  
  // Create doc directory (unless dry-run)
  const docPath = path.join(resolvedOutputPath, 'doc', 'features');
  if (!dryRun && !fs.existsSync(docPath)) {
    console.log(chalk.yellow(`📁 Creating documentation directory: ${docPath}`));
    fs.mkdirSync(docPath, { recursive: true });
  }
  
  // Select CLI tool
  const cliTool = await selectCLITool(dryRun);
  
  // Load prompt templates
  const promptsDir = path.join(path.dirname(path.dirname(path.dirname(__dirname))), 'prompts');
  const featurePromptPath = path.join(promptsDir, '04-doc-generate-feature.md');
  
  // Generate documentation for each area
  console.log(chalk.blue('\n🚀 Starting documentation generation...\n'));
  
  for (let i = 0; i < areasToDocument.length; i++) {
    const area = areasToDocument[i];
    const moduleInfo = modules.functional_areas.find(m => m.name === area.name);
    
    console.log(chalk.cyan(`[${i + 1}/${areasToDocument.length}] Generating documentation for: ${area.name}`));
    
    // Prepare output file name
    const fileName = area.name.toLowerCase().replace(/[^a-z0-9]+/g, '-') + '.md';
    const outputFile = path.join(docPath, fileName);
    
    // Load prompt with partials
    const promptData = loadPromptWithPartials(featurePromptPath);
    
    // Execute generation using runAgent with parameters as variables
    const result = await runAgent({
      tool: cliTool as 'claude' | 'gemini',
      prompts: [
        { name: promptData.name, content: promptData.content, partials: promptData.partials }
      ],
      variables: {
        parameters: {
          source_path: resolvedSourcePath,
          area_name: area.name,
          area_info: JSON.stringify(moduleInfo, null, 2),
          output_file: outputFile
        }
      },
      outputPath: outputFile,
      workingDir: workingDir,
      dryRun: dryRun,
      description: `Generate documentation for ${area.name}`
    });
    
    if (!dryRun) {
      if (result.success && fs.existsSync(outputFile)) {
        console.log(chalk.green(`  ✅ Created: ${fileName}`));
      } else {
        console.log(chalk.red(`  ❌ Failed to generate documentation for ${area.name}`));
      }
      
      // Small delay between generations
      await new Promise(resolve => setTimeout(resolve, 500));
    }
  }
  
  // Generate index
  await updateDocumentationIndex(
    resolvedOutputPath,
    resolvedModulesPath,
    workingDir,
    cliTool,
    dryRun
  );
  
  if (!dryRun) {
    console.log(chalk.green('\n✨ Documentation generation complete!'));
    console.log(chalk.gray(`Documentation created in: ${path.join(resolvedOutputPath, 'doc')}`));
  }
}