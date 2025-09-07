/**
 * Display utilities for formatted console output
 */

import chalk from 'chalk';
import { DocumentationQuality, DocumentationSummary, FunctionalArea, ModulesData, ModuleInfo } from '../types';

/**
 * Calculate summary statistics from functional areas
 */
export function calculateSummaryStats(areas: FunctionalArea[]): DocumentationSummary {
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
    // Sum for averages
    qualitySum += area.doc_quality_score;
    importanceSum += area.doc_importance_score;

    // Count by importance level
    switch(area.doc_importance_score) {
      case 5: summary.mandatory_count++; break;
      case 4: summary.very_important_count++; break;
      case 3: summary.important_count++; break;
      case 2: summary.nice_to_have_count++; break;
      case 1: summary.not_important_count++; break;
    }

    // Count undocumented critical (quality <= 2 AND importance >= 4)
    if (area.doc_quality_score <= 2 && area.doc_importance_score >= 4) {
      summary.undocumented_critical++;
    }
  });

  // Calculate averages
  summary.overall_quality_score = areas.length > 0 ? qualitySum / areas.length : 0;
  summary.overall_importance_score = areas.length > 0 ? importanceSum / areas.length : 0;

  return summary;
}

/**
 * Display documentation quality summary
 */
export function displayDocumentationSummary(qualityData: DocumentationQuality) {
  // Calculate summary if not present
  if (!qualityData.summary) {
    qualityData.summary = calculateSummaryStats(qualityData.functional_areas);
  }
  
  console.log(chalk.blue('\n📊 Documentation Quality Summary'));
  console.log(chalk.blue('═'.repeat(70)));
  
  // Header
  console.log(chalk.white('\nArea Name'.padEnd(35)) + 
              chalk.white('Quality') + '  ' +
              chalk.white('Importance') + '  ' +
              chalk.white('Action'));
  console.log(chalk.gray('─'.repeat(70)));
  
  // Sort by importance score (descending) then by quality score (ascending)
  const sortedAreas = qualityData.functional_areas.sort((a: any, b: any) => {
    if (b.doc_importance_score !== a.doc_importance_score) {
      return b.doc_importance_score - a.doc_importance_score;
    }
    return a.doc_quality_score - b.doc_quality_score;
  });
  
  // Display each area
  sortedAreas.forEach((area: any) => {
    const name = area.name.padEnd(35);
    const quality = '★'.repeat(area.doc_quality_score) + '☆'.repeat(5 - area.doc_quality_score);
    const importance = '★'.repeat(area.doc_importance_score) + '☆'.repeat(5 - area.doc_importance_score);
    
    let action = '';
    switch(area.doc_importance_score) {
      case 5: action = chalk.red('MANDATORY'); break;
      case 4: action = chalk.yellow('VERY IMPORTANT'); break;
      case 3: action = chalk.cyan('IMPORTANT'); break;
      case 2: action = chalk.gray('NICE TO HAVE'); break;
      case 1: action = chalk.gray('NOT IMPORTANT'); break;
    }
    
    console.log(name + '  ' + quality + '  ' + importance + '  ' + action);
  });
  
  // Summary statistics
  console.log(chalk.blue('═'.repeat(70)));
  console.log(chalk.white('\nStatistics:'));
  console.log(chalk.gray(`  Total Areas: ${qualityData.total_areas || qualityData.functional_areas.length}`));
  console.log(chalk.gray(`  Overall Quality: ${qualityData.overall_quality_score?.toFixed(1) || qualityData.summary?.overall_quality_score?.toFixed(1)}/5`));
  console.log(chalk.gray(`  Overall Importance: ${qualityData.overall_importance_score?.toFixed(1) || qualityData.summary?.overall_importance_score?.toFixed(1)}/5`));
  
  if (qualityData.summary) {
    const critical = qualityData.summary.undocumented_critical;
    const needsDocs = qualityData.summary.mandatory_count + qualityData.summary.very_important_count;
    
    console.log(chalk.yellow(`  Undocumented Critical: ${critical} areas (quality ≤2 AND importance ≥4)`));
    console.log(chalk.yellow(`  Areas needing documentation: ${needsDocs} critical`));
  }
  
  console.log('');
}

/**
 * Display unified list of modules with optional documentation quality
 */
export function displayModulesWithQuality(
  modulesData: ModulesData,
  qualityData?: DocumentationQuality,
  detailed?: boolean
) {
  // Create a map of documentation quality by area name
  const qualityMap = new Map<string, FunctionalArea>();
  if (qualityData?.functional_areas) {
    qualityData.functional_areas.forEach(area => {
      qualityMap.set(area.name, area);
    });
    
    // Debug: Show if there's a mismatch in area names
    const moduleAreaNames = new Set(modulesData.functional_areas.map(a => a.name));
    const qualityAreaNames = new Set(qualityData.functional_areas.map(a => a.name));
    const matchingNames = [...moduleAreaNames].filter(name => qualityAreaNames.has(name));
    
    if (matchingNames.length === 0 && qualityData.functional_areas.length > 0) {
      console.log(chalk.yellow('⚠️  Warning: No matching area names between modules.json and documentation_quality.json'));
      console.log(chalk.gray(`   Module areas: ${[...moduleAreaNames].slice(0, 3).join(', ')}...`));
      console.log(chalk.gray(`   Quality areas: ${[...qualityAreaNames].slice(0, 3).join(', ')}...`));
    }
  }

  // Always show doc quality columns to make the feature discoverable
  const hasDocQuality = true; // Always show quality columns
  const hasMatchingQuality = qualityMap.size > 0;
  const hasQualityData = qualityData !== undefined;
  const areas = modulesData.functional_areas;

  console.log(chalk.blue('📋 Functional Areas Overview'));
  console.log(chalk.blue('═'.repeat(100)));
  
  // Display modules info if present
  if (modulesData.modules && modulesData.modules.length > 0) {
    console.log(chalk.cyan('\n📦 Detected Modules:'));
    modulesData.modules.forEach(mod => {
      console.log(chalk.gray(`  • ${mod.name} (${mod.type}): ${mod.main_technology} at ${mod.path}`));
    });
    console.log('');
  }

  // Header
  if (detailed) {
    if (hasDocQuality) {
      console.log(
        chalk.white('No.  ') +
        chalk.white('Module'.padEnd(12)) +
        chalk.white('Name'.padEnd(28)) +
        chalk.white('Type'.padEnd(20)) +
        chalk.white('Complex'.padEnd(10)) +
        chalk.white('Quality'.padEnd(10)) +
        chalk.white('Import.'.padEnd(10)) +
        chalk.white('Files'.padEnd(8)) +
        chalk.white('Tech')
      );
      console.log(chalk.gray('─'.repeat(110)));
    } else {
      console.log(
        chalk.white('No.  ') +
        chalk.white('Module'.padEnd(12)) +
        chalk.white('Name'.padEnd(30)) +
        chalk.white('Type'.padEnd(20)) +
        chalk.white('Complex') + '  ' +
        chalk.white('Files') + '  ' +
        chalk.white('LOC') + '    ' +
        chalk.white('Technologies')
      );
      console.log(chalk.gray('─'.repeat(95)));
    }
  } else {
    // Simple header - always show Type, Complex, Quality, and Priority
    console.log(
      chalk.white('No.  ') +
      chalk.white('Module'.padEnd(12)) +
      chalk.white('Name'.padEnd(30)) +
      chalk.white('Type'.padEnd(20)) +
      chalk.white('Complex'.padEnd(10)) +
      chalk.white('Quality'.padEnd(10)) +
      chalk.white('Priority')
    );
    console.log(chalk.gray('─'.repeat(100)));
  }

  // Sort areas by module, then by importance (if available), then by name
  const sortedAreas = [...areas].sort((a, b) => {
    // First by module
    const modA = a.module || 'main';
    const modB = b.module || 'main';
    if (modA !== modB) return modA.localeCompare(modB);
    
    // Then by importance if available
    if (hasMatchingQuality) {
      const qualA = qualityMap.get(a.name);
      const qualB = qualityMap.get(b.name);
      if (qualA && qualB && qualA.doc_importance_score !== qualB.doc_importance_score) {
        return qualB.doc_importance_score - qualA.doc_importance_score;
      }
    }
    
    // Finally by name
    return a.name.localeCompare(b.name);
  });

  // Display each area
  sortedAreas.forEach((area, index) => {
    const num = `${index + 1}.`.padEnd(5);
    const module = (area.module || 'main').padEnd(12);
    
    // Use new field names with fallback to legacy
    const areaType = (area.area_type || area.type || 'unknown');
    const complexity = area.complexity_score || 0;
    const fileCount = area.metrics?.file_count || area.source_files?.length || area.source_references?.length || 0;
    const loc = area.metrics?.estimated_loc || 0;
    
    const complexityStars = '★'.repeat(complexity) + '☆'.repeat(5 - complexity);
    
    const docQuality = qualityMap.get(area.name);

    if (detailed) {
      const name = area.name.substring(0, 28).padEnd(28);
      const type = areaType.substring(0, 20).padEnd(20);
      
      if (hasDocQuality) {
        // Show quality data if available for this area, otherwise show dashes
        const qualityStars = docQuality 
          ? '★'.repeat(docQuality.doc_quality_score) + '☆'.repeat(5 - docQuality.doc_quality_score)
          : '-----';
        const importStars = docQuality
          ? '★'.repeat(docQuality.doc_importance_score) + '☆'.repeat(5 - docQuality.doc_importance_score)
          : '-----';
        const techList = area.technologies?.slice(0, 3).map(t => t.name).join(', ') || 'N/A';
        
        console.log(
          num + module + name + type +
          complexityStars.padEnd(10) +
          qualityStars.padEnd(10) +
          importStars.padEnd(10) +
          fileCount.toString().padEnd(8) +
          techList
        );
      } else {
        const techList = area.technologies?.slice(0, 3).map(t => t.name).join(', ') || 'N/A';
        console.log(
          num + module + name.padEnd(30) + type +
          complexityStars + '  ' +
          fileCount.toString().padEnd(7) +
          loc.toString().padEnd(7) +
          techList
        );
      }
      
      // Show description in gray
      if (area.description) {
        console.log(chalk.gray(`     └─ ${area.description.substring(0, 80)}`));
      }
    } else {
      // Simple view - always show all columns
      const name = area.name.substring(0, 30).padEnd(30);
      const type = areaType.substring(0, 20).padEnd(20);
      
      // Show quality stars if this area has quality data, otherwise show dashes
      let qualityStars = '-----';
      let priority = chalk.gray('-'.padEnd(14));  // Pad priority to fixed width
      
      if (docQuality) {
        qualityStars = '★'.repeat(docQuality.doc_quality_score) + '☆'.repeat(5 - docQuality.doc_quality_score);
        
        switch(docQuality.doc_importance_score) {
          case 5: priority = chalk.red('MANDATORY'.padEnd(14)); break;
          case 4: priority = chalk.yellow('VERY IMPORTANT'); break;  // 14 chars exactly
          case 3: priority = chalk.cyan('IMPORTANT'.padEnd(14)); break;
          case 2: priority = chalk.gray('NICE TO HAVE'.padEnd(14)); break;
          case 1: priority = chalk.gray('NOT NEEDED'.padEnd(14)); break;
        }
      }
      
      console.log(
        num + module + name + type +
        complexityStars.padEnd(10) +
        qualityStars.padEnd(10) +
        priority
      );
    }
  });

  // Summary statistics
  console.log(chalk.blue('═'.repeat(100)));
  console.log(chalk.cyan('\nSummary:'));
  
  // Count by module
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
  
  // Count by type
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
  if (hasQualityData && qualityData) {
    if (hasMatchingQuality && qualityData.functional_areas.length > 0) {
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