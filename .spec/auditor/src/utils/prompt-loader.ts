/**
 * Prompt loading utilities with partial support
 */

import * as fs from 'fs';
import * as path from 'path';
import matter from 'gray-matter';
import { PromptWithMetadata } from '../types';

/**
 * Load a prompt file with its partials
 * Partials declared in frontmatter are concatenated to the main content
 */
export function loadPromptWithPartials(promptPath: string): PromptWithMetadata {
  const promptContent = fs.readFileSync(promptPath, 'utf-8');
  const { data: frontmatter, content } = matter(promptContent);
  
  let finalContent = content;
  const includedPartials: string[] = [];
  
  // If there are partials declared in frontmatter, concatenate them
  if (frontmatter.partials && Array.isArray(frontmatter.partials)) {
    const promptDir = path.dirname(promptPath);
    
    for (const partial of frontmatter.partials) {
      const partialPath = path.join(promptDir, partial);
      
      if (!fs.existsSync(partialPath)) {
        console.warn(`Warning: Partial file not found: ${partialPath}`);
        continue;
      }
      
      const partialContent = fs.readFileSync(partialPath, 'utf-8');
      includedPartials.push(partial);
      
      // Simply append the partial content to the end
      finalContent += `\n\n---\n\n${partialContent}`;
    }
  }
  
  return {
    content: finalContent,
    name: path.basename(promptPath),
    partials: includedPartials
  };
}