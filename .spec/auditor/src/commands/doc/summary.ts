/**
 * Doc summary command - displays documentation quality summary
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import { DocumentationQuality, ModulesData } from '../../types';
import { displayDocumentationSummary } from '../../utils/display';
import { displayModulesTable } from '../../utils/display-table';

/**
 * Run the doc summary command
 */
export async function runDocSummary(cwdPath: string) {
  const workingDir = fs.realpathSync(cwdPath);
  const qualityFile = path.join(workingDir, 'audit_result', 'documentation_quality.json');
  const modulesFile = path.join(workingDir, 'audit_result', 'modules.json');
  
  // Check if documentation quality exists
  if (fs.existsSync(qualityFile)) {
    // Load and display the quality data
    try {
      const content = fs.readFileSync(qualityFile, 'utf-8');
      const qualityData: DocumentationQuality = JSON.parse(content);
      
      console.log(chalk.blue('📚 Documentation Status'));
      console.log(chalk.blue('═'.repeat(70)));
      displayDocumentationSummary(qualityData);
      
      console.log(chalk.gray('\nNext steps:'));
      console.log(chalk.cyan('  • Run "doc generate" to create documentation for important areas'));
      console.log(chalk.cyan('  • Run "doc update" to update existing documentation'));
      console.log(chalk.cyan('  • Run "list --detailed" to see all areas with quality scores'));
      return;
    } catch (error) {
      console.error(chalk.red(`❌ Failed to load documentation quality: ${error}`));
    }
  }
  
  // Fall back to modules.json if quality doesn't exist
  if (fs.existsSync(modulesFile)) {
    try {
      const content = fs.readFileSync(modulesFile, 'utf-8');
      const modulesData: ModulesData = JSON.parse(content);
      
      console.log(chalk.blue('📚 Documentation Status'));
      console.log(chalk.blue('═'.repeat(70)));
      console.log(chalk.yellow('\n⚠️  Documentation quality not yet assessed\n'));
      
      // Show basic module info
      displayModulesTable(modulesData, undefined, false);
      
      console.log(chalk.gray('\nTo assess documentation quality:'));
      console.log(chalk.cyan('  npx tsx audit.ts doc review'));
      console.log(chalk.gray('\nThis will analyze each functional area and score:'));
      console.log(chalk.gray('  - Current documentation quality (1-5)'));
      console.log(chalk.gray('  - Documentation importance/priority (1-5)'));
      console.log(chalk.gray('  - Identify which areas need documentation most urgently'));
      return;
    } catch (error) {
      console.error(chalk.red(`❌ Failed to load modules data: ${error}`));
    }
  }
  
  // No data found at all
  console.log(chalk.yellow('⚠️  No codebase analysis found'));
  console.log(chalk.gray('\nTo get started:'));
  console.log(chalk.cyan('  1. Run "explore <source>" to analyze your codebase'));
  console.log(chalk.cyan('  2. Run "doc review" to assess documentation needs'));
  console.log(chalk.cyan('  3. Run "doc generate" to create documentation'));
}