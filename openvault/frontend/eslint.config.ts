import js from '@eslint/js';
import globals from 'globals';
import tseslint from 'typescript-eslint';
import pluginReact from 'eslint-plugin-react';
import { defineConfig } from 'eslint/config';
import pluginRouter from '@tanstack/eslint-plugin-router';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import unusedImports from 'eslint-plugin-unused-imports';

export default defineConfig([
  {
    ignores: ['dist/**', '**/routeTree.gen.ts', '**/wailsjs/**']
  },
  {
    files: ['**/*.{js,mjs,cjs,jsx}'],
    plugins: { js },
    extends: ['js/recommended'],
    languageOptions: { globals: globals.browser }
  },
  {
    files: ['**/*.{ts,mts,cts,tsx}'],
    extends: [
      // tseslint.configs.recommendedTypeChecked,
      tseslint.configs.recommended,
      pluginReact.configs.flat.recommended,
      ...pluginRouter.configs['flat/recommended'],
      {
        plugins: { 'unused-imports': unusedImports },
        rules: {
          'react/react-in-jsx-scope': 'off',
          'react/jsx-uses-react': 'off',
          'react/no-children-prop': 'off',
          '@typescript-eslint/only-throw-error': [
            'error',
            {
              allow: [
                {
                  from: 'package',
                  package: '@tanstack/router-core',
                  name: 'Redirect'
                }
              ]
            }
          ],
          '@typescript-eslint/no-unused-vars': 'off', // Or "@typescript-eslint/no-unused-vars": "off",
          'unused-imports/no-unused-imports': 'error',
          'unused-imports/no-unused-vars': [
            'warn',
            {
              vars: 'all',
              varsIgnorePattern: '^_',
              args: 'after-used',
              argsIgnorePattern: '^_'
            }
          ]
        }
      }
    ],
    languageOptions: {
      parserOptions: {
        projectService: true
      }
    },
    settings: {
      react: {
        version: 'detect'
      }
    }
  },
  // {
  //   files: ['**/*.css'],
  //   // @ts-expect-error: Weird typing issue with css plugin
  //   plugins: { css },
  //   language: 'css/css',
  //   languageOptions: {
  //     customSyntax: tailwind4,  //     tolerant: true,
  //   },
  //   extends: ['css/recommended']
  // },
  eslintPluginPrettierRecommended
]);
