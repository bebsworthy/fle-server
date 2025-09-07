/**
 * Profile Resolver - Resolves profile inheritance and creates chains
 */

import { Profile, ProfileChain } from '../types/profile';

export class ProfileChainImpl implements ProfileChain {
  constructor(public profiles: Profile[]) {}
  
  /**
   * Get merged guidelines from all profiles in the chain
   */
  getMergedGuidelines(): string {
    return this.profiles
      .map(profile => {
        const header = `## Guidelines from ${profile.name} (${profile.type})`;
        return `${header}\n\n${profile.content}`;
      })
      .join('\n\n---\n\n');
  }
  
  /**
   * Get profile names in the chain
   */
  getNames(): string[] {
    return this.profiles.map(p => p.name);
  }
}

export class ProfileResolver {
  private profilesByPath: Map<string, Profile>;
  
  constructor(private profiles: Profile[]) {
    this.profilesByPath = new Map();
    profiles.forEach(p => this.profilesByPath.set(p.path, p));
  }
  
  /**
   * Resolve matched profiles into a chain with inheritance
   */
  resolve(matched: Profile[]): ProfileChain {
    const resolved: Profile[] = [];
    const seen = new Set<string>();
    
    // Sort matched profiles by type to process in correct order
    const sorted = this.sortByTypeHierarchy(matched);
    
    // Resolve each profile with its extends
    for (const profile of sorted) {
      this.resolveWithExtends(profile, resolved, seen);
    }
    
    // Auto-include language base profiles
    this.addLanguageBases(matched, resolved, seen);
    
    // Remove duplicates while preserving order
    const uniqueResolved = this.removeDuplicates(resolved);
    
    return new ProfileChainImpl(uniqueResolved);
  }
  
  /**
   * Resolve a profile and its extends recursively
   */
  private resolveWithExtends(
    profile: Profile,
    chain: Profile[],
    seen: Set<string>
  ): void {
    if (seen.has(profile.path)) return;
    seen.add(profile.path);
    
    // First resolve extends (depth-first)
    if (profile.extends) {
      for (const extPath of profile.extends) {
        const extended = this.profilesByPath.get(extPath);
        if (extended) {
          this.resolveWithExtends(extended, chain, seen);
        } else {
          console.warn(`Warning: Profile ${profile.path} extends non-existent profile: ${extPath}`);
        }
      }
    }
    
    // Then add this profile
    chain.push(profile);
  }
  
  /**
   * Add language base profiles for matched frameworks
   */
  private addLanguageBases(
    matched: Profile[],
    resolved: Profile[],
    seen: Set<string>
  ): void {
    // Find all languages mentioned in keywords
    const languages = new Set<string>();
    
    for (const profile of matched) {
      // Check if profile is a framework with language in keywords
      if (profile.type === 'framework' && profile.keywords) {
        // Common language keywords
        const langKeywords = ['python', 'javascript', 'typescript', 'go', 'java', 'ruby', 'rust', 'csharp'];
        
        for (const keyword of profile.keywords) {
          const normalized = keyword.toLowerCase();
          if (langKeywords.includes(normalized)) {
            languages.add(normalized);
          }
        }
      }
    }
    
    // Find and add base profiles for each language
    for (const lang of languages) {
      const basePath = `${lang}/_base.md`;
      const baseProfile = this.profilesByPath.get(basePath);
      
      if (baseProfile && !seen.has(basePath)) {
        // Insert at beginning (base profiles go first)
        resolved.unshift(baseProfile);
        seen.add(basePath);
      }
    }
  }
  
  /**
   * Sort profiles by type hierarchy
   */
  private sortByTypeHierarchy(profiles: Profile[]): Profile[] {
    const typeOrder = { pattern: 1, language: 2, system: 3, framework: 4 };
    
    return [...profiles].sort((a, b) => {
      const orderA = typeOrder[a.type] || 0;
      const orderB = typeOrder[b.type] || 0;
      return orderA - orderB;
    });
  }
  
  /**
   * Remove duplicate profiles while preserving order
   */
  private removeDuplicates(profiles: Profile[]): Profile[] {
    const seen = new Set<string>();
    const unique: Profile[] = [];
    
    for (const profile of profiles) {
      if (!seen.has(profile.path)) {
        seen.add(profile.path);
        unique.push(profile);
      }
    }
    
    return unique;
  }
}