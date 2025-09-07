/**
 * Agent runner with dry-run support
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import { AgentConfig, AgentResult, GlobalOptions } from '../types';
import { executeCLI, getCLIArgs } from './cli-executor';
import { loadPromptWithPartials } from './prompt-loader';

/**
 * Run an AI agent with the specified configuration
 */
export async function runAgent(config: AgentConfig): Promise<AgentResult> {
  const { tool, prompts, variables = {}, outputPath, workingDir, dryRun, description } = config;
  
  // Build final prompt by concatenating all prompts
  let finalPrompt = prompts.map(p => p.content).join('\n\n');
  
  // Replace variables in the prompt - both [VAR] style and <parameters> blocks
  Object.entries(variables).forEach(([key, value]) => {
    // Handle [VAR] style replacements
    const regex = new RegExp(`\\[${key}\\]`, 'g');
    finalPrompt = finalPrompt.replace(regex, value);
    
    // Handle <parameter_name>value</parameter_name> style
    const paramRegex = new RegExp(`<${key}>[^<]*</${key}>`, 'g');
    finalPrompt = finalPrompt.replace(paramRegex, `<${key}>${value}</${key}>`);
  });
  
  // Handle <parameters> block replacement if parameters variable exists
  if (variables.parameters) {
    const parametersBlock = Object.entries(variables.parameters)
      .map(([key, value]) => `  <${key}>${value}</${key}>`)
      .join('\n');
    finalPrompt = finalPrompt.replace(
      /<parameters>[\s\S]*?<\/parameters>/,
      `<parameters>\n${parametersBlock}\n</parameters>`
    );
  }
  
  if (dryRun) {
    // Determine verbosity level
    let verbosity = 0; // Default: minimal
    if (dryRun === true || dryRun === 'true') {
      verbosity = 0; // Default when just --dry-run is used
    } else if (typeof dryRun === 'string') {
      const level = parseInt(dryRun);
      if (!isNaN(level) && level >= 0 && level <= 2) {
        verbosity = level;
      }
    }
    
    console.log(chalk.blue('\n🔍 DRY RUN MODE - No agent will be executed'));
    console.log(chalk.gray('═'.repeat(70)));
    
    if (description) {
      console.log(chalk.cyan(`\nTask: ${description}`));
    }
    
    console.log(chalk.yellow('\n📋 Configuration:'));
    console.log(chalk.gray(`  Tool: ${tool === 'claude' ? 'Claude CLI' : 'Gemini CLI'}`));
    console.log(chalk.gray(`  Working Directory: ${workingDir}`));
    console.log(chalk.gray(`  Output Path: ${outputPath}`));
    
    if (verbosity >= 1) {
      console.log(chalk.yellow('\n📝 Prompts to be concatenated:'));
      prompts.forEach((p, i) => {
        console.log(chalk.gray(`  ${i + 1}. ${p.name}`));
        if (p.partials && p.partials.length > 0) {
          console.log(chalk.gray(`     └─ Includes partials:`));
          p.partials.forEach(partial => {
            console.log(chalk.gray(`        • ${partial}`));
          });
        }
      });
      
      if (Object.keys(variables).length > 0) {
        console.log(chalk.yellow('\n🔧 Variables/Parameters to be injected:'));
        Object.entries(variables).forEach(([key, value]) => {
          if (key === 'parameters' && typeof value === 'object') {
            console.log(chalk.gray(`  Parameters:`));
            Object.entries(value).forEach(([pKey, pValue]) => {
              const displayValue = typeof pValue === 'string' && pValue.length > 100 
                ? pValue.substring(0, 100) + '...' 
                : pValue;
              console.log(chalk.gray(`    <${pKey}>: ${displayValue}`));
            });
          } else {
            const displayValue = typeof value === 'string' && value.length > 100 
              ? value.substring(0, 100) + '...' 
              : value;
            console.log(chalk.gray(`  ${key}: ${displayValue}`));
          }
        });
      }
    }
    
    if (verbosity >= 2) {
      console.log(chalk.yellow('\n📄 Prompt Preview (first 1000 chars):'));
      console.log(chalk.gray('─'.repeat(70)));
      console.log(chalk.gray(finalPrompt.substring(0, 1000)));
      if (finalPrompt.length > 1000) {
        console.log(chalk.gray('... [truncated]'));
      }
      console.log(chalk.gray('─'.repeat(70)));
    }
    
    console.log(chalk.yellow('\n🚀 Command that would be executed:'));
    const args = getCLIArgs(tool, outputPath);
    console.log(chalk.gray(`  ${tool} ${args.join(' ')}`));
    
    console.log(chalk.green('\n✅ Dry run complete - no changes made'));
    
    return { success: true };
  }
  
  // Execute the agent
  const args = getCLIArgs(tool, outputPath);
  
  // Get global options to check for raw output and agent output file
  const globalOpts = (global as any).globalOptions as GlobalOptions || {};
  const filterOutput = tool === 'claude' && !globalOpts.raw;
  
  const result = await executeCLI(
    tool, 
    args, 
    finalPrompt, 
    workingDir,
    filterOutput,
    globalOpts.agentOutput
  );
  
  return result;
}