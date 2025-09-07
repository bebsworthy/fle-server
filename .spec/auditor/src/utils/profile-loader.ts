/**
 * Profile Loader - Loads and parses review profile markdown files
 */

import * as fs from 'fs/promises';
import * as path from 'path';
import matter from 'gray-matter';
import { glob } from 'glob';
import { Profile } from '../types/profile';

export class ProfileLoader {
  private profilesDir: string;
  
  constructor(profilesDir: string = 'prompts/review-profiles') {
    this.profilesDir = profilesDir;
  }
  
  /**
   * Load all profile markdown files
   */
  async loadProfiles(): Promise<Profile[]> {
    const profiles: Profile[] = [];
    
    // Find all markdown files
    const pattern = path.join(this.profilesDir, '**/*.md');
    const files = await glob(pattern);
    
    for (const file of files) {
      try {
        const content = await fs.readFile(file, 'utf-8');
        const profile = this.parseProfile(file, content);
        this.validateProfile(profile);
        profiles.push(profile);
      } catch (error) {
        console.warn(`Warning: Failed to load profile ${file}:`, error);
      }
    }
    
    return profiles;
  }
  
  /**
   * Parse a profile markdown file
   */
  private parseProfile(filePath: string, content: string): Profile {
    const { data: frontmatter, content: body } = matter(content);
    
    // Get relative path from profiles directory
    const relativePath = path.relative(this.profilesDir, filePath);
    
    return {
      path: relativePath,
      type: frontmatter.type || this.inferType(relativePath),
      name: frontmatter.name || this.inferName(filePath),
      keywords: frontmatter.keywords || [],
      extends: frontmatter.extends,
      priority: frontmatter.priority || 0,
      matches: frontmatter.matches,
      content: body.trim()
    };
  }
  
  /**
   * Validate a profile has required fields
   */
  private validateProfile(profile: Profile): void {
    if (!profile.type) {
      throw new Error(`Profile ${profile.path} missing required field: type`);
    }
    
    if (!['framework', 'pattern', 'language', 'system'].includes(profile.type)) {
      throw new Error(`Profile ${profile.path} has invalid type: ${profile.type}`);
    }
    
    if (!profile.name) {
      throw new Error(`Profile ${profile.path} missing required field: name`);
    }
    
    if (!profile.keywords || profile.keywords.length === 0) {
      throw new Error(`Profile ${profile.path} missing required field: keywords`);
    }
  }
  
  /**
   * Infer type from file path
   */
  private inferType(relativePath: string): 'framework' | 'pattern' | 'language' | 'system' {
    if (relativePath.startsWith('patterns/')) return 'pattern';
    if (relativePath.startsWith('systems/')) return 'system';
    if (relativePath.includes('/_base.md')) return 'language';
    return 'framework';
  }
  
  /**
   * Infer name from file path
   */
  private inferName(filePath: string): string {
    const basename = path.basename(filePath, '.md');
    return basename === '_base' ? path.basename(path.dirname(filePath)) : basename;
  }
}