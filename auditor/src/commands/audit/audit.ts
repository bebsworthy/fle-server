/**
 * Code audit command - analyzes code quality and architecture
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import inquirer from 'inquirer';
import { GlobalOptions, ModulesData } from '../../types';
import { CodeAuditOptions } from '../../types/audit';
import { runAgent } from '../../utils/agent-runner';
import { loadPromptWithPartials } from '../../utils/prompt-loader';
import { displayAuditSummary } from './summary';
import { ReviewEnhancer } from '../../utils/review-enhancer';

/**
 * Run the code audit command
 */
export async function runCodeAudit(
  options: CodeAuditOptions,
  cwdPath: string,
  globalOpts: GlobalOptions
) {
  console.log(chalk.blue('🔍 Codebase Auditor - Code Quality & Architecture Audit'));
  console.log(chalk.gray('─'.repeat(50)));
  
  // Get global options for dry-run
  const dryRun = globalOpts.dryRun || false;
  
  // Store global options for agent runner
  (global as any).globalOptions = globalOpts;
  
  // Resolve paths
  const workingDir = fs.realpathSync(cwdPath);
  const resolvedModulesPath = path.isAbsolute(options.file)
    ? path.resolve(options.file)
    : path.resolve(workingDir, options.file);
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
  
  if (!fs.existsSync(resolvedSourcePath)) {
    console.error(chalk.red(`❌ Source path not found: ${resolvedSourcePath}`));
    process.exit(1);
  }
  
  // Load modules data
  const modules: ModulesData = JSON.parse(fs.readFileSync(resolvedModulesPath, 'utf-8'));
  
  // Filter areas to audit
  let areasToAudit = modules.functional_areas;
  
  // Apply module filter first if specified
  if (options.module) {
    const moduleFilter = options.module.toLowerCase();
    areasToAudit = areasToAudit.filter(area => 
      area.module && area.module.toLowerCase() === moduleFilter
    );
    
    if (areasToAudit.length === 0) {
      console.log(chalk.red(`\n❌ No areas found for module: ${options.module}`));
      
      // Show available modules
      const availableModules = new Set<string>();
      modules.functional_areas.forEach(area => {
        if (area.module) {
          availableModules.add(area.module);
        }
      });
      
      if (availableModules.size > 0) {
        console.log(chalk.gray('Available modules:'));
        availableModules.forEach(mod => {
          console.log(chalk.gray(`  - ${mod}`));
        });
      } else {
        console.log(chalk.gray('No modules defined in this project (not a monorepo)'));
      }
      return;
    }
    
    console.log(chalk.cyan(`\n🎯 Filtering to module: ${options.module}`));
  }
  
  if (options.name && options.name.length > 0) {
    // Filter by specific names
    const nameSet = new Set(options.name.map(n => n.toLowerCase()));
    areasToAudit = areasToAudit.filter(area => 
      nameSet.has(area.name.toLowerCase())
    );
    
    // Check if all requested names were found
    const foundNames = new Set(areasToAudit.map(a => a.name.toLowerCase()));
    const notFound = options.name.filter(n => !foundNames.has(n.toLowerCase()));
    
    if (notFound.length > 0) {
      console.log(chalk.yellow(`\n⚠️  Areas not found: ${notFound.join(', ')}`));
      console.log(chalk.gray('Available areas:'));
      const relevantAreas = options.module 
        ? modules.functional_areas.filter(a => a.module === options.module)
        : modules.functional_areas;
      relevantAreas.forEach(area => {
        console.log(chalk.gray(`  - ${area.name}${area.module ? ` (${area.module})` : ''}`));
      });
    }
    
    console.log(chalk.cyan(`\n📋 Will audit ${areasToAudit.length} specific areas: ${areasToAudit.map(a => a.name).join(', ')}`));
  } else if (options.minComplexity) {
    // Filter by complexity score
    const minComplexity = parseInt(options.minComplexity);
    areasToAudit = areasToAudit.filter(area => 
      (area.complexity_score || 0) >= minComplexity
    );
    
    const moduleMsg = options.module ? ` in module '${options.module}'` : '';
    console.log(chalk.cyan(`\n📋 Will audit ${areasToAudit.length} areas${moduleMsg} with complexity ≥ ${minComplexity}`));
  } else {
    const moduleMsg = options.module ? ` in module '${options.module}'` : ' functional areas';
    console.log(chalk.cyan(`\n📋 Will audit all ${areasToAudit.length}${moduleMsg}`));
  }
  
  if (areasToAudit.length === 0) {
    console.log(chalk.yellow('\n⚠️  No areas match the criteria'));
    return;
  }
  
  // Display focus area if specified
  if (options.focus) {
    console.log(chalk.yellow(`\n🎯 Focus area: ${options.focus}`));
  }
  
  // Create audit directory (unless dry-run)
  const auditPath = path.join(resolvedOutputPath, 'audit');
  if (!dryRun && !fs.existsSync(auditPath)) {
    console.log(chalk.yellow(`📁 Creating audit directory: ${auditPath}`));
    fs.mkdirSync(auditPath, { recursive: true });
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
  
  // Load audit prompt
  const currentDir = path.dirname(new URL(import.meta.url).pathname);
  const promptsDir = path.join(currentDir, '..', '..', '..', 'prompts');
  const auditPromptPath = path.join(promptsDir, '06-code-audit.md');
  
  // Initialize review enhancer for technology-specific profiles
  console.log(chalk.gray('📚 Loading review profiles...'));
  const enhancer = new ReviewEnhancer();
  await enhancer.initialize();
  const totalProfiles = enhancer.getProfiles().length;
  console.log(chalk.gray(`✓ Loaded ${totalProfiles} review profiles\n`));
  
  // Audit each area
  console.log(chalk.blue('🚀 Starting code audit...\n'));
  
  const auditResults = [];
  
  for (let i = 0; i < areasToAudit.length; i++) {
    const area = areasToAudit[i];
    
    console.log(chalk.cyan(`[${i + 1}/${areasToAudit.length}] Auditing: ${area.name}`));
    
    // Get and display matched review profiles
    const matchInfo = await enhancer.getMatchInfo(area);
    if (matchInfo.matched.length > 0) {
      console.log(chalk.blue('📋 Applying Review Profiles:'));
      for (const profile of matchInfo.matched) {
        const typeIcon = profile.type === 'language' ? '🔤' : 
                        profile.type === 'framework' ? '🏗️' : 
                        profile.type === 'library' ? '📚' : '⚙️';
        console.log(chalk.gray(`  ${typeIcon} ${profile.name} (${profile.type}) - ${profile.path}`));
      }
      console.log(chalk.gray(`  → ${matchInfo.matched.length} profiles with ~${matchInfo.guidelineCount} review guidelines\n`));
    } else {
      console.log(chalk.yellow('  ⚠️  No specific review profiles matched (using generic guidelines)\n'));
    }
    
    // Prepare output file name
    const fileName = area.name.toLowerCase().replace(/[^a-z0-9]+/g, '-') + '.md';
    const outputFile = path.join(auditPath, fileName);
    
    // Load prompt with partials
    const promptData = loadPromptWithPartials(auditPromptPath);
    
    // Enhance prompt with review profiles
    let enhancedContent = promptData.content;
    if (matchInfo.matched.length > 0) {
      const profileEnhancement = await enhancer.enhancePrompt(area, '');
      if (profileEnhancement) {
        enhancedContent = promptData.content + '\n\n' + profileEnhancement;
      }
    }
    
    // Execute audit
    const result = await runAgent({
      tool: cliTool as 'claude' | 'gemini',
      prompts: [
        { name: promptData.name, content: enhancedContent, partials: promptData.partials }
      ],
      variables: {
        parameters: {
          source_path: resolvedSourcePath,
          area_name: area.name,
          area_info: JSON.stringify(area, null, 2),
          output_file: outputFile,
          focus_area: options.focus || 'balanced'
        }
      },
      outputPath: outputFile,
      workingDir: workingDir,
      dryRun: dryRun,
      description: `Audit code quality for ${area.name}`
    });
    
    if (!dryRun) {
      if (result.success && fs.existsSync(outputFile)) {
        console.log(chalk.green(`  ✅ Audit complete: ${fileName}`));
        auditResults.push({
          area: area.name,
          file: fileName,
          success: true
        });
      } else {
        console.log(chalk.red(`  ❌ Failed to audit ${area.name}`));
        auditResults.push({
          area: area.name,
          file: fileName,
          success: false
        });
      }
      
      // Small delay between audits
      await new Promise(resolve => setTimeout(resolve, 500));
    }
  }
  
  // Generate audit index
  if (!dryRun && auditResults.length > 0) {
    await generateAuditIndex(auditPath, auditResults, resolvedModulesPath, workingDir, cliTool);
  }
  
  // Display summary
  if (!dryRun) {
    displayAuditSummary(auditResults);
    console.log(chalk.green('\n✨ Code audit complete!'));
    console.log(chalk.gray(`Audit reports saved in: ${auditPath}`));
  }
}

/**
 * Generate an index file for all audit reports
 */
async function generateAuditIndex(
  auditPath: string,
  auditResults: any[],
  modulesPath: string,
  workingDir: string,
  cliTool: string
) {
  console.log(chalk.cyan('\n📚 Generating audit index...'));
  
  const indexContent = `# Code Audit Reports

**Generated**: ${new Date().toISOString().split('T')[0]}
**Total Areas Audited**: ${auditResults.length}
**Success Rate**: ${Math.round(auditResults.filter(r => r.success).length / auditResults.length * 100)}%

## Audit Reports by Area

| Area | Status | Report |
|------|--------|--------|
${auditResults.map(r => 
  `| ${r.area} | ${r.success ? '✅ Complete' : '❌ Failed'} | ${r.success ? `[View](${r.file})` : 'N/A'} |`
).join('\n')}

## Quick Navigation

${auditResults.filter(r => r.success).map(r => 
  `- [${r.area}](${r.file})`
).join('\n')}

## How to Use These Reports

Each audit report contains:
- **Executive Summary**: Quick health score and key metrics
- **Critical Findings**: Issues that need immediate attention
- **Code Quality Assessment**: Detailed quality metrics
- **Architecture Review**: Design and structure analysis
- **Action Plan**: Prioritized improvements

### Priority Levels

- 🔴 **Critical**: Security vulnerabilities, data loss risks, system crashes
- 🟠 **High**: Performance bottlenecks, missing error handling, high complexity
- 🟡 **Medium**: Moderate complexity, inconsistent patterns, missing docs
- 🟢 **Low**: Minor improvements, style issues, optimization opportunities

### Next Steps

1. Review critical issues in each report
2. Create tickets for Priority 1 quick wins
3. Plan refactoring for high-impact improvements
4. Track progress with follow-up audits

---
*Generated by Codebase Auditor*`;
  
  const indexPath = path.join(auditPath, 'README.md');
  fs.writeFileSync(indexPath, indexContent);
  console.log(chalk.green('✅ Audit index created'));
}