/**
 * Doc update command - updates existing documentation with incremental changes
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
 * Run the doc update command
 */
export async function runDocUpdate(
  modulesPath: string,
  qualityPath: string,
  sourcePath: string,
  outputPath: string,
  cwdPath: string,
  minScore: number | undefined,
  areaNames: string[] | undefined,
  globalOpts: GlobalOptions
) {
  console.log(chalk.blue('📝 Codebase Auditor - Documentation Update (Incremental)'));
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
  
  // Load prompt templates - use the new update prompt
  const promptsDir = path.join(path.dirname(path.dirname(path.dirname(__dirname))), 'prompts');
  const updatePromptPath = path.join(promptsDir, '04-doc-update-feature.md');
  
  // Update documentation for each area
  console.log(chalk.blue('\n🚀 Starting incremental documentation update...\n'));
  
  for (let i = 0; i < areasToDocument.length; i++) {
    const area = areasToDocument[i];
    const moduleInfo = modules.functional_areas.find(m => m.name === area.name);
    
    // Prepare output file name
    const fileName = area.name.toLowerCase().replace(/[^a-z0-9]+/g, '-') + '.md';
    const outputFile = path.join(docPath, fileName);
    
    // Check if file exists for update
    const fileExists = fs.existsSync(outputFile);
    if (fileExists) {
      console.log(chalk.cyan(`[${i + 1}/${areasToDocument.length}] Updating existing documentation for: ${area.name}`));
    } else {
      console.log(chalk.cyan(`[${i + 1}/${areasToDocument.length}] Creating new documentation for: ${area.name}`));
    }
    
    // Load prompt with partials
    const promptData = loadPromptWithPartials(updatePromptPath);
    
    // Execute update using runAgent with parameters as variables
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
      description: fileExists ? `Update documentation for ${area.name}` : `Create documentation for ${area.name}`
    });
    
    if (!dryRun) {
      if (result.success && fs.existsSync(outputFile)) {
        if (fileExists) {
          console.log(chalk.green(`  ✅ Updated: ${fileName}`));
        } else {
          console.log(chalk.green(`  ✅ Created: ${fileName}`));
        }
      } else {
        console.log(chalk.red(`  ❌ Failed to update/create documentation for ${area.name}`));
      }
      
      // Small delay between updates
      await new Promise(resolve => setTimeout(resolve, 500));
    }
  }
  
  // Update index
  await updateDocumentationIndex(
    resolvedOutputPath,
    resolvedModulesPath,
    workingDir,
    cliTool,
    dryRun
  );
  
  if (!dryRun) {
    console.log(chalk.green('\n✨ Documentation update complete!'));
    console.log(chalk.gray(`Documentation updated in: ${path.join(resolvedOutputPath, 'doc')}`));
  }
}