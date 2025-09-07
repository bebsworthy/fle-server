/**
 * Review Enhancer - Integrates profile system with review prompts
 */

import * as path from 'path';
import { ProfileLoader } from './profile-loader';
import { ProfileMatcher } from './profile-matcher';
import { ProfileResolver } from './profile-resolver';
import { Profile, ProfileChain } from '../types/profile';
import { ModuleInfo } from '../types';

export class ReviewEnhancer {
  private loader: ProfileLoader;
  private matcher?: ProfileMatcher;
  private resolver?: ProfileResolver;
  private profiles: Profile[] = [];
  private initialized = false;
  
  constructor(profilesDir?: string) {
    const dir = profilesDir || path.join(process.cwd(), 'prompts/review-profiles');
    this.loader = new ProfileLoader(dir);
  }
  
  /**
   * Initialize the review enhancer by loading all profiles
   */
  async initialize(): Promise<void> {
    if (this.initialized) return;
    
    try {
      this.profiles = await this.loader.loadProfiles();
      this.matcher = new ProfileMatcher(this.profiles);
      this.resolver = new ProfileResolver(this.profiles);
      this.initialized = true;
      
      console.log(`✓ Loaded ${this.profiles.length} review profiles`);
    } catch (error) {
      console.warn('Warning: Failed to initialize review profiles:', error);
      this.profiles = [];
      this.initialized = true;
    }
  }
  
  /**
   * Enhance a review prompt with matched profiles
   */
  async enhancePrompt(area: ModuleInfo, basePrompt: string): Promise<string> {
    if (!this.initialized) {
      await this.initialize();
    }
    
    if (!this.matcher || !this.resolver || this.profiles.length === 0) {
      return basePrompt; // No profiles available
    }
    
    // Find matching profiles
    const matched = this.matcher.findMatchingProfiles(area);
    
    if (matched.length === 0) {
      return basePrompt; // No matching profiles
    }
    
    // Resolve inheritance chain
    const chain = this.resolver.resolve(matched);
    
    // Format profile section
    const profileSection = this.formatProfileSection(chain, area);
    
    // Enhance the prompt
    return `${basePrompt}

${profileSection}`;
  }
  
  /**
   * Format the profile section for injection into prompt
   */
  private formatProfileSection(chain: ProfileChain, area: ModuleInfo): string {
    const profiles = chain.profiles;
    
    let section = '## Technology-Specific Review Guidelines\n\n';
    
    // Show which profiles are being applied
    section += '### Applied Profiles\n';
    section += 'The following review profiles have been selected based on the detected technologies:\n\n';
    
    for (const profile of profiles) {
      section += `- **${profile.name}** (${profile.type}): ${profile.keywords.join(', ')}\n`;
    }
    
    section += '\n### Review Guidelines\n\n';
    
    // Add guidelines from each profile
    for (const profile of profiles) {
      if (profile.content) {
        section += `#### ${profile.name} Guidelines\n\n`;
        section += profile.content;
        section += '\n\n';
      }
    }
    
    // Add area-specific context
    section += '### Area Context\n\n';
    section += `- **Area Name**: ${area.name}\n`;
    section += `- **Area Type**: ${area.area_type || 'Not specified'}\n`;
    section += `- **Technologies**: ${area.technologies?.map(t => t.name).join(', ') || 'None detected'}\n`;
    section += `- **Complexity**: ${area.complexity_score || 0}/5\n`;
    
    return section;
  }
  
  /**
   * Get debug information about profile matching
   */
  async debugMatch(area: ModuleInfo): Promise<string> {
    if (!this.initialized) {
      await this.initialize();
    }
    
    if (!this.matcher || !this.resolver) {
      return 'Profile system not initialized';
    }
    
    const matched = this.matcher.findMatchingProfiles(area);
    const chain = this.resolver.resolve(matched);
    
    let debug = `Debug information for area: ${area.name}\n`;
    debug += `Technologies: ${area.technologies?.map(t => t.name).join(', ') || 'None'}\n\n`;
    
    debug += `Matched profiles (${matched.length}):\n`;
    for (const profile of matched) {
      debug += `  - ${profile.name} (${profile.type}, priority: ${profile.priority || 0})\n`;
    }
    
    debug += `\nResolved chain (${chain.profiles.length}):\n`;
    for (const profile of chain.profiles) {
      debug += `  - ${profile.name} (${profile.type})\n`;
    }
    
    return debug;
  }
  
  /**
   * Get all loaded profiles
   */
  getProfiles(): Profile[] {
    return this.profiles;
  }
  
  /**
   * Check if system is initialized
   */
  isInitialized(): boolean {
    return this.initialized;
  }
  
  /**
   * Get match information for an area
   */
  async getMatchInfo(area: ModuleInfo): Promise<{
    matched: Array<{
      name: string;
      type: string;
      path: string;
      priority: number;
      keywords: string[];
    }>;
    chain: string[];
    guidelineCount: number;
  }> {
    if (!this.initialized) {
      await this.initialize();
    }
    
    if (!this.matcher || !this.resolver) {
      return { matched: [], chain: [], guidelineCount: 0 };
    }
    
    const matched = this.matcher.findMatchingProfiles(area);
    const chain = this.resolver.resolve(matched);
    
    // Count approximate guidelines (lines in content)
    let guidelineCount = 0;
    for (const profile of chain.profiles) {
      if (profile.content) {
        // Count checklist items and sections
        const checklistItems = (profile.content.match(/\[[ x]\]/g) || []).length;
        const sections = (profile.content.match(/^###/gm) || []).length;
        guidelineCount += checklistItems + sections;
      }
    }
    
    return {
      matched: matched.map(p => ({
        name: p.name,
        type: p.type,
        path: p.path,
        priority: p.priority || 0,
        keywords: p.keywords
      })),
      chain: chain.getNames(),
      guidelineCount
    };
  }
}