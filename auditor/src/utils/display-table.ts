/**
 * Table-based display utilities using cli-table3
 */

import chalk from 'chalk';
import Table from 'cli-table3';
import { DocumentationQuality, DocumentationSummary, FunctionalArea, ModulesData, ModuleInfo } from '../types';

/**
 * Display modules with quality using proper table formatting
 */
export function displayModulesTable(
  modulesData: ModulesData,
  qualityData?: DocumentationQuality,
  detailed?: boolean,
  showTech?: boolean
) {
  // Create quality map
  const qualityMap = new Map<string, FunctionalArea>();
  if (qualityData?.functional_areas) {
    qualityData.functional_areas.forEach(area => {
      qualityMap.set(area.name, area);
    });
  }

  const areas = modulesData.functional_areas;
  
  console.log(chalk.blue('\n📋 Functional Areas Overview'));
  console.log(chalk.blue('═'.repeat(100)));
  
  // Display modules info if present
  if (modulesData.modules && modulesData.modules.length > 0) {
    console.log(chalk.cyan('\n📦 Detected Modules:'));
    modulesData.modules.forEach(mod => {
      console.log(chalk.gray(`  • ${mod.name} (${mod.type}): ${mod.main_technology} at ${mod.path}`));
    });
    console.log('');
  }

  // Create table with appropriate columns - compact layout without borders
  const table = new Table({
    head: showTech ? [
      chalk.white.bold('No.'),
      chalk.white.bold('Module'),
      chalk.white.bold('Name'),
      chalk.white.bold('Technologies')
    ] : detailed ? [
      chalk.white.bold('No.'),
      chalk.white.bold('Module'),
      chalk.white.bold('Name'),
      chalk.white.bold('Type'),
      chalk.white.bold('Complex'),
      chalk.white.bold('Quality'),
      chalk.white.bold('Priority'),
      chalk.white.bold('Files'),
      chalk.white.bold('Tech')
    ] : [
      chalk.white.bold('No.'),
      chalk.white.bold('Module'),
      chalk.white.bold('Name'),
      chalk.white.bold('Type'),
      chalk.white.bold('Complex'),
      chalk.white.bold('Quality'),
      chalk.white.bold('Priority')
    ],
    colWidths: showTech ?
      [5, 12, 35, 70] :
      detailed ? 
      [5, 12, 30, 20, 10, 10, 16, 8, 30] :
      [5, 12, 32, 22, 10, 10, 16],
    style: {
      'padding-left': 0,
      'padding-right': 0,
      head: [],
      border: []
    },
    chars: {
      'top': '', 'top-mid': '', 'top-left': '', 'top-right': '',
      'bottom': '', 'bottom-mid': '', 'bottom-left': '', 'bottom-right': '',
      'left': '', 'left-mid': '', 'mid': '', 'mid-mid': '',
      'right': '', 'right-mid': '', 'middle': ' '
    }
  });

  // Sort areas
  const sortedAreas = [...areas].sort((a, b) => {
    const modA = a.module || 'main';
    const modB = b.module || 'main';
    if (modA !== modB) return modA.localeCompare(modB);
    
    const qualA = qualityMap.get(a.name);
    const qualB = qualityMap.get(b.name);
    if (qualA && qualB && qualA.doc_importance_score !== qualB.doc_importance_score) {
      return qualB.doc_importance_score - qualA.doc_importance_score;
    }
    
    return a.name.localeCompare(b.name);
  });

  // Add rows to table
  sortedAreas.forEach((area, index) => {
    const num = (index + 1).toString();
    const module = area.module || 'main';
    const name = area.name;
    
    const docQuality = qualityMap.get(area.name);
    let qualityStars = chalk.gray('-----');
    
    if (docQuality) {
      qualityStars = '★'.repeat(docQuality.doc_quality_score) + '☆'.repeat(5 - docQuality.doc_quality_score);
    }
    
    if (showTech) {
      // Technology view - show all technologies without versions
      const techList = area.technologies?.map(t => {
        // Color-code by type
        switch(t.type) {
          case 'language': return chalk.cyan(t.name);
          case 'framework': return chalk.green(t.name);
          case 'library': return chalk.yellow(t.name);
          case 'database': return chalk.magenta(t.name);
          case 'tool': return chalk.gray(t.name);
          default: return t.name;
        }
      }).join(', ') || chalk.gray('No technologies detected');
      
      table.push([
        num,
        module,
        name,
        techList
      ]);
    } else {
      // Standard or detailed view
      const areaType = area.area_type || area.type || 'unknown';
      const complexity = area.complexity_score || 0;
      const complexityStars = '★'.repeat(complexity) + '☆'.repeat(5 - complexity);
      
      let priority = chalk.gray('-');
      if (docQuality) {
        switch(docQuality.doc_importance_score) {
          case 5: priority = chalk.red('MANDATORY'); break;
          case 4: priority = chalk.yellow('VERY IMPORTANT'); break;
          case 3: priority = chalk.cyan('IMPORTANT'); break;
          case 2: priority = chalk.gray('NICE TO HAVE'); break;
          case 1: priority = chalk.gray('NOT NEEDED'); break;
        }
      }
      
      if (detailed) {
        const fileCount = area.metrics?.file_count || area.source_files?.length || area.source_references?.length || 0;
        const techList = area.technologies?.slice(0, 3).map(t => t.name).join(', ') || 'N/A';
        
        table.push([
          num,
          module,
          name,
          areaType,
          complexityStars,
          qualityStars,
          priority,
          fileCount.toString(),
          techList
        ]);
      } else {
        table.push([
          num,
          module,
          name,
          areaType,
          complexityStars,
          qualityStars,
          priority
        ]);
      }
    }
  });

  // Print table with a separator line after header
  const tableStr = table.toString();
  const lines = tableStr.split('\n');
  if (lines.length > 1) {
    // Print header
    console.log(lines[0]);
    // Print separator
    console.log(chalk.gray('─'.repeat(100)));
    // Print rest of table
    console.log(lines.slice(1).join('\n'));
  } else {
    console.log(tableStr);
  }

  // Summary statistics
  console.log(chalk.cyan('\nSummary:'));
  
  // Module counts
  if (modulesData.modules && modulesData.modules.length > 0) {
    const moduleCounts = new Map<string, number>();
    areas.forEach(area => {
      const mod = area.module || 'main';
      moduleCounts.set(mod, (moduleCounts.get(mod) || 0) + 1);
    });
    
    console.log(chalk.gray('  Areas by module:'));
    moduleCounts.forEach((count, mod) => {
      console.log(chalk.gray(`    ${mod}: ${count} areas`));
    });
  }
  
  // Type counts
  const typeCounts = new Map<string, number>();
  areas.forEach(area => {
    const type = area.area_type || area.type || 'unknown';
    typeCounts.set(type, (typeCounts.get(type) || 0) + 1);
  });
  
  console.log(chalk.gray('  Areas by type:'));
  typeCounts.forEach((count, type) => {
    console.log(chalk.gray(`    ${type}: ${count}`));
  });
  
  // Average complexity
  const avgComplexity = areas.reduce((sum, area) => sum + (area.complexity_score || 0), 0) / areas.length;
  console.log(chalk.gray(`  Average Complexity: ${'★'.repeat(Math.round(avgComplexity))}${'☆'.repeat(5 - Math.round(avgComplexity))} (${avgComplexity.toFixed(1)})`));
  
  // Documentation summary
  if (qualityData) {
    if (qualityMap.size > 0 && qualityData.functional_areas.length > 0) {
      const summary = qualityData.summary || calculateSummaryStats(qualityData.functional_areas);
      console.log(chalk.gray(`  Documentation Quality: ${summary.overall_quality_score.toFixed(1)}/5`));
      console.log(chalk.gray(`  Documentation Priority: ${summary.overall_importance_score.toFixed(1)}/5`));
      
      if (summary.undocumented_critical > 0) {
        console.log(chalk.yellow(`  ⚠️  Undocumented Critical: ${summary.undocumented_critical} areas need immediate attention`));
      }
    } else if (qualityData.functional_areas.length === 0) {
      console.log(chalk.gray(`  Documentation Quality: Not assessed`));
      console.log(chalk.gray(`  Documentation Priority: Not assessed`));
    } else {
      console.log(chalk.yellow(`  ⚠️  Documentation quality data exists but no matching areas found`));
      console.log(chalk.gray(`     Run "doc review" to reassess documentation quality`));
    }
  } else {
    console.log(chalk.gray(`  Documentation Quality: Not assessed`));
    console.log(chalk.gray(`  Documentation Priority: Not assessed`));
  }
  
  console.log('');
}

/**
 * Calculate summary statistics (copied from display.ts)
 */
function calculateSummaryStats(areas: FunctionalArea[]): DocumentationSummary {
  const summary: DocumentationSummary = {
    total_areas: areas.length,
    overall_quality_score: 0,
    overall_importance_score: 0,
    mandatory_count: 0,
    very_important_count: 0,
    important_count: 0,
    nice_to_have_count: 0,
    not_important_count: 0,
    undocumented_critical: 0
  };

  let qualitySum = 0;
  let importanceSum = 0;

  areas.forEach(area => {
    qualitySum += area.doc_quality_score;
    importanceSum += area.doc_importance_score;

    switch(area.doc_importance_score) {
      case 5: summary.mandatory_count++; break;
      case 4: summary.very_important_count++; break;
      case 3: summary.important_count++; break;
      case 2: summary.nice_to_have_count++; break;
      case 1: summary.not_important_count++; break;
    }

    if (area.doc_quality_score <= 2 && area.doc_importance_score >= 4) {
      summary.undocumented_critical++;
    }
  });

  summary.overall_quality_score = areas.length > 0 ? qualitySum / areas.length : 0;
  summary.overall_importance_score = areas.length > 0 ? importanceSum / areas.length : 0;

  return summary;
}