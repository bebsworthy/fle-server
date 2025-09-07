/**
 * CLI output filtering for Claude's JSON stream
 */

import chalk from 'chalk';

/**
 * Filter Claude's JSON stream output to show only essential information
 */
export function filterClaudeOutput(line: string): string | null {
  try {
    const json = JSON.parse(line);
    
    // System/Init message
    if (json.type === 'system' && json.subtype === 'init') {
      return chalk.gray(`Starting ${json.model} in ${json.cwd}`);
    }
    
    // Assistant message with content array
    if (json.type === 'assistant' && json.message?.type === 'message' && Array.isArray(json.message?.content)) {
      const outputs: string[] = [];
      
      // Process each content item in the array
      for (const contentItem of json.message.content) {
        if (contentItem.type === 'text') {
          // Text content
          outputs.push(contentItem.text);
        } else if (contentItem.type === 'tool_use') {
          // Tool use content
          const toolName = contentItem.name;
          const input = contentItem.input;
          
          // Special handling for TodoWrite
          if (toolName === 'TodoWrite' && input.todos) {
            const todoOutput: string[] = [chalk.cyan(`[${toolName}] Task list updated:`)];
            for (const todo of input.todos) {
              let statusIcon = '⏳'; // pending
              if (todo.status === 'completed') {
                statusIcon = '✅';
              } else if (todo.status === 'in_progress') {
                statusIcon = '🔄';
              }
              todoOutput.push(chalk.gray(`  ${todo.id}. ${statusIcon} ${todo.content}`));
            }
            outputs.push(todoOutput.join('\n'));
          } else {
            // Extract relevant input based on tool type
            let inputDesc = '';
            if (input.file_path) {
              inputDesc = input.file_path;
            } else if (input.path) {
              inputDesc = input.path;
            } else if (input.command) {
              inputDesc = input.command;
            } else if (input.pattern) {
              inputDesc = input.pattern;
            } else if (input.url) {
              inputDesc = input.url;
            } else {
              // For other tools, show first meaningful value
              const values = Object.values(input).filter(v => v && typeof v === 'string');
              if (values.length > 0) {
                inputDesc = values[0] as string;
              }
            }
            
            outputs.push(chalk.cyan(`[${toolName}] ${inputDesc}`));
          }
        }
        // Could handle other content types here if needed
      }
      
      // Return all content items joined with newlines
      return outputs.length > 0 ? outputs.join('\n') : null;
    }
    
    // Result message
    if (json.type === 'result') {
      const status = json.is_error ? 'error' : 'success';
      const message = json.result || '';
      // Truncate long result messages
      const truncatedMessage = message.length > 200 ? message.substring(0, 200) + '...' : message;
      return chalk.green(`[${status}] ${truncatedMessage}`);
    }
    
    // Skip other message types
    return null;
  } catch (e) {
    // If JSON parsing fails, return the line as-is (might be non-JSON output)
    return line;
  }
}