/**
 * CLI execution utilities
 */

import { spawn } from 'child_process';
import * as fs from 'fs';
import chalk from 'chalk';
import { filterClaudeOutput } from './cli-filter';

/**
 * Execute CLI command and capture output
 */
export async function executeCLI(
  command: string, 
  args: string[], 
  prompt: string, 
  cwd: string,
  filterOutput: boolean = false,
  agentOutputFile?: string
): Promise<{ success: boolean; output: string }> {
  return new Promise((resolve) => {
    let output = '';
    let lineBuffer = '';
    let agentOutputStream: fs.WriteStream | undefined;
    
    // Create write stream for agent output if file specified
    if (agentOutputFile) {
      agentOutputStream = fs.createWriteStream(agentOutputFile, { flags: 'w' });
      console.log(chalk.gray(`📝 Writing raw agent output to: ${agentOutputFile}`));
    }
    
    const child = spawn(command, args, {
      stdio: ['pipe', 'pipe', 'pipe'],
      shell: true,
      cwd: cwd
    });
    
    // Handle stdout - apply filtering if needed
    child.stdout?.on('data', (data) => {
      const text = data.toString();
      output += text;
      
      // Write to agent output file if specified
      if (agentOutputStream) {
        agentOutputStream.write(data);
      }
      
      if (filterOutput) {
        // Process line by line for JSON stream
        lineBuffer += text;
        const lines = lineBuffer.split('\n');
        lineBuffer = lines.pop() || ''; // Keep incomplete line in buffer
        
        for (const line of lines) {
          if (line.trim()) {
            const filtered = filterClaudeOutput(line);
            if (filtered !== null) {
              console.log(filtered);
            }
          }
        }
      } else {
        // Raw output
        process.stdout.write(data);
      }
    });
    
    // Handle any remaining buffered content when stream ends
    child.stdout?.on('end', () => {
      if (filterOutput && lineBuffer.trim()) {
        const filtered = filterClaudeOutput(lineBuffer);
        if (filtered !== null) {
          console.log(filtered);
        }
      }
    });
    
    // Always pass stderr through unfiltered
    child.stderr?.on('data', (data) => {
      output += data.toString();
      // Write to agent output file if specified
      if (agentOutputStream) {
        agentOutputStream.write(data);
      }
      process.stderr.write(data);
    });
    
    child.stdin?.write(prompt);
    child.stdin?.end();
    
    child.on('close', (code) => {
      // Close the output stream if it exists
      if (agentOutputStream) {
        agentOutputStream.end();
      }
      resolve({
        success: code === 0,
        output: output
      });
    });
    
    child.on('error', (err) => {
      console.error(chalk.red(`Failed to execute ${command}:`), err.message);
      resolve({
        success: false,
        output: err.message
      });
    });
  });
}

/**
 * Get CLI arguments based on tool type
 */
export function getCLIArgs(tool: string, outputPath: string): string[] {
  if (tool === 'claude') {
    return [
      '-p',
      '--permission-mode', 'acceptEdits',
      '--verbose',
      '--output-format', 'stream-json'
    ];
  } else {
    // Gemini
    return [
      '--outputPath', outputPath
    ];
  }
}