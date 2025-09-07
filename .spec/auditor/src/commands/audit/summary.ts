/**
 * Audit summary display utilities
 */

import chalk from 'chalk';

/**
 * Display audit summary after completion
 */
export function displayAuditSummary(auditResults: any[]) {
  console.log(chalk.blue('\n📊 Audit Summary'));
  console.log(chalk.blue('═'.repeat(70)));
  
  const successful = auditResults.filter(r => r.success);
  const failed = auditResults.filter(r => !r.success);
  
  console.log(chalk.white('\nResults:'));
  console.log(chalk.green(`  ✅ Successfully audited: ${successful.length} areas`));
  
  if (failed.length > 0) {
    console.log(chalk.red(`  ❌ Failed to audit: ${failed.length} areas`));
    console.log(chalk.gray('\nFailed areas:'));
    failed.forEach(r => {
      console.log(chalk.gray(`    - ${r.area}`));
    });
  }
  
  if (successful.length > 0) {
    console.log(chalk.white('\n📝 Reports Generated:'));
    successful.forEach(r => {
      console.log(chalk.gray(`  - ${r.file}`));
    });
  }
  
  console.log(chalk.blue('═'.repeat(70)));
}