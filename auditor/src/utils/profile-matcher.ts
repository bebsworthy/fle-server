/**
 * Profile Matcher - Matches profiles to functional areas based on technologies
 */

import { Profile } from '../types/profile';
import { ModuleInfo } from '../types';

export class ProfileMatcher {
  private profiles: Profile[];
  private profilesByName: Map<string, Profile>;
  
  constructor(profiles: Profile[]) {
    this.profiles = profiles;
    this.profilesByName = new Map();
    
    // Build index for fast lookup
    profiles.forEach(profile => {
      this.profilesByName.set(profile.name, profile);
    });
  }
  
  /**
   * Find all profiles that match a functional area
   */
  findMatchingProfiles(area: ModuleInfo): Profile[] {
    const matches: Profile[] = [];
    const areaTechnologies = this.extractTechnologies(area);
    
    for (const profile of this.profiles) {
      if (this.profileMatches(profile, area, areaTechnologies)) {
        matches.push(profile);
      }
    }
    
    // Sort by priority (higher first) and type specificity
    return matches.sort((a, b) => {
      // First by priority
      if (a.priority !== b.priority) {
        return (b.priority || 0) - (a.priority || 0);
      }
      
      // Then by type specificity: framework > system > language > pattern
      const typeOrder = { framework: 4, system: 3, language: 2, pattern: 1 };
      return (typeOrder[b.type] || 0) - (typeOrder[a.type] || 0);
    });
  }
  
  /**
   * Get a profile by its path
   */
  getProfileByPath(path: string): Profile | undefined {
    return this.profiles.find(p => p.path === path);
  }
  
  /**
   * Extract normalized technology names from area
   */
  private extractTechnologies(area: ModuleInfo): Set<string> {
    const techs = new Set<string>();
    
    if (area.technologies) {
      area.technologies.forEach(tech => {
        // Add original name
        techs.add(this.normalizeString(tech.name));
        
        // Add without version numbers
        const nameWithoutVersion = tech.name.replace(/[0-9.]+$/, '').trim();
        if (nameWithoutVersion !== tech.name) {
          techs.add(this.normalizeString(nameWithoutVersion));
        }
      });
    }
    
    return techs;
  }
  
  /**
   * Check if a profile matches an area
   */
  private profileMatches(
    profile: Profile,
    area: ModuleInfo,
    areaTechnologies: Set<string>
  ): boolean {
    // 1. Match by profile name against technologies
    if (areaTechnologies.has(this.normalizeString(profile.name))) {
      return true;
    }
    
    // 2. Match by keywords against technologies
    if (profile.keywords) {
      for (const keyword of profile.keywords) {
        if (areaTechnologies.has(this.normalizeString(keyword))) {
          return true;
        }
      }
    }
    
    // 3. Match by explicit matching rules
    if (profile.matches) {
      return this.matchesByRules(profile, area, areaTechnologies);
    }
    
    return false;
  }
  
  /**
   * Check if profile matches by explicit rules
   */
  private matchesByRules(
    profile: Profile,
    area: ModuleInfo,
    areaTechnologies: Set<string>
  ): boolean {
    const matches = profile.matches!;
    
    // Match area_type
    if (matches.area_type && area.area_type) {
      const areaType = area.area_type.toLowerCase();
      if (matches.area_type.some(type => type.toLowerCase() === areaType)) {
        return true;
      }
    }
    
    // Match name_contains
    if (matches.name_contains) {
      const areaName = area.name.toLowerCase();
      if (matches.name_contains.some(term => areaName.includes(term.toLowerCase()))) {
        return true;
      }
    }
    
    // Match technologies
    if (matches.technologies) {
      for (const tech of matches.technologies) {
        if (areaTechnologies.has(this.normalizeString(tech))) {
          return true;
        }
      }
    }
    
    return false;
  }
  
  /**
   * Normalize string for comparison (lowercase, remove special chars)
   */
  private normalizeString(str: string): string {
    return str.toLowerCase()
      .replace(/[^a-z0-9]/g, '')
      .replace(/js$/, 'javascript'); // Common normalization
  }
}