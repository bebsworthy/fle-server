/**
 * Doc review command - reviews and scores documentation quality
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import inquirer from 'inquirer';
import { GlobalOptions, DocumentationQuality, ModulesData } from '../../types';
import { runAgent } from '../../utils/agent-runner';
import { loadPromptWithPartials } from '../../utils/prompt-loader';
import { displayDocumentationSummary, calculateSummaryStats } from '../../utils/display';
import { ReviewEnhancer } from '../../utils/review-enhancer';

/**
 * Run the doc review command
 */
export async function runDocReview(
  filePath: string, 
  outputPath: string, 
  cwdPath: string,
  globalOpts: GlobalOptions
) {
  console.log(chalk.blue('📝 Codebase Auditor - Documentation Review'));
  console.log(chalk.gray('─'.repeat(50)));
  
  // Get global options for dry-run
  const dryRun = globalOpts.dryRun || false;
  
  // Store global options for agent runner
  (global as any).globalOptions = globalOpts;
  
  // Resolve paths
  const workingDir = fs.realpathSync(cwdPath);
  const resolvedFilePath = path.isAbsolute(filePath)
    ? path.resolve(filePath)
    : path.resolve(workingDir, filePath);
  const resolvedOutputPath = path.isAbsolute(outputPath)
    ? path.resolve(outputPath)
    : path.resolve(workingDir, outputPath);
  
  // Validate modules file exists
  if (!fs.existsSync(resolvedFilePath)) {
    console.error(chalk.red(`❌ Modules file not found: ${resolvedFilePath}`));
    console.error(chalk.gray('Run the explore command first to generate modules.json'));
    process.exit(1);
  }
  
  // Create output directory if it doesn't exist (unless dry-run)
  if (!dryRun && !fs.existsSync(resolvedOutputPath)) {
    console.log(chalk.yellow(`📁 Creating output directory: ${resolvedOutputPath}`));
    fs.mkdirSync(resolvedOutputPath, { recursive: true });
  }
  
  // Remove existing quality file to ensure fresh review (unless dry-run)
  const qualityFilePath = path.join(resolvedOutputPath, 'documentation_quality.json');
  if (!dryRun && fs.existsSync(qualityFilePath)) {
    console.log(chalk.yellow(`🗑️  Removing existing quality assessment: ${qualityFilePath}`));
    fs.unlinkSync(qualityFilePath);
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
  
  console.log(chalk.cyan(`\n🤖 Using ${cliTool === 'claude' ? 'Claude' : 'Gemini'} CLI for review\n`));
  
  // Load prompt templates with partials
  const promptsDir = path.join(path.dirname(path.dirname(path.dirname(__dirname))), 'prompts');
  const reviewPromptPath = path.join(promptsDir, '03-doc-review.md');
  const promptData = loadPromptWithPartials(reviewPromptPath);
  
  // Initialize review enhancer for technology-specific profiles
  const enhancer = new ReviewEnhancer();
  await enhancer.initialize();
  
  // Load modules data to enhance prompts per area
  let modulesData: ModulesData | null = null;
  let enhancedContent = promptData.content;
  
  try {
    const modulesContent = fs.readFileSync(resolvedFilePath, 'utf-8');
    modulesData = JSON.parse(modulesContent);
    
    // If we have functional areas, enhance the prompt with profile guidelines
    if (modulesData && modulesData.functional_areas && modulesData.functional_areas.length > 0) {
      // For now, we'll enhance the base prompt with general context
      // In a more advanced version, we could run separate reviews per area
      const profileSections: string[] = [];
      
      for (const area of modulesData.functional_areas) {
        const enhanced = await enhancer.enhancePrompt(area, '');
        if (enhanced) {
          profileSections.push(`\n### Profile for ${area.name}\n${enhanced}`);
        }
      }
      
      if (profileSections.length > 0) {
        enhancedContent = promptData.content + '\n\n## Technology-Specific Profiles Applied\n' + 
                         'The following technology-specific review guidelines have been loaded:\n' +
                         profileSections.slice(0, 3).join('\n'); // Limit to avoid prompt overflow
      }
    }
  } catch (error) {
    console.warn(chalk.yellow('⚠️  Could not enhance prompts with technology profiles:', error));
  }
  
  // Run the review
  console.log(chalk.blue('🚀 Starting documentation review...\n'));
  
  const result = await runAgent({
    tool: cliTool as 'claude' | 'gemini',
    prompts: [
      { name: promptData.name, content: enhancedContent, partials: promptData.partials }
    ],
    variables: {
      parameters: {
        modules_file: resolvedFilePath,
        output_path: resolvedOutputPath
      }
    },
    outputPath: qualityFilePath,
    workingDir: workingDir,
    dryRun: dryRun,
    description: 'Review documentation quality for all functional areas'
  });
  
  if (!dryRun) {
    if (result.success && fs.existsSync(qualityFilePath)) {
      console.log(chalk.green('\n✅ Documentation review complete!'));
      
      // Load and display summary
      try {
        const content = fs.readFileSync(qualityFilePath, 'utf-8');
        const qualityData: DocumentationQuality = JSON.parse(content);
        
        // Add summary stats if not present
        if (!qualityData.summary && qualityData.functional_areas) {
          qualityData.summary = calculateSummaryStats(qualityData.functional_areas);
          // Save back with summary
          fs.writeFileSync(qualityFilePath, JSON.stringify(qualityData, null, 2));
        }
        
        displayDocumentationSummary(qualityData);
      } catch (error) {
        console.error(chalk.red(`\n⚠️  Could not display summary: ${error}`));
      }
      
      console.log(chalk.gray(`\nResults saved to: ${qualityFilePath}`));
    } else {
      console.error(chalk.red('\n❌ Documentation review failed'));
      process.exit(1);
    }
  }
}