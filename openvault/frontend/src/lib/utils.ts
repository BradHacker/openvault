import { clsx, type ClassValue } from 'clsx';
import { extendTailwindMerge } from 'tailwind-merge';

const twMerge = extendTailwindMerge({
  extend: {
    theme: {
      // We only need to define the custom scale values without the `shadow-` prefix when adding them to the theme object
      'drop-shadow': ['glow'],
      animate: ['shake']
    }
  }
});

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
