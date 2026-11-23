import { defineConfig } from 'eslint/config'
import globals from 'globals'
import js from '@eslint/js'
import stylistic from '@stylistic/eslint-plugin'
import unicorn from 'eslint-plugin-unicorn'

export default defineConfig([
  {
    extends: ['js/all', unicorn.configs.all, stylistic.configs.all],
    files: ['**/*.{js,mjs,cjs}'],
    languageOptions: { globals: globals.browser },
    linterOptions: {
      reportUnusedDisableDirectives: 'error',
      reportUnusedInlineConfigs: 'error'
    },
    plugins: { js },
    rules: {
      // eslint-disable-next-line no-magic-numbers
      '@stylistic/indent': ['error', 2],
      '@stylistic/lines-between-class-members': 'off',
      '@stylistic/no-multiple-empty-lines': ['error', { max: 1, maxBOF: 0, maxEOF: 0 }],
      '@stylistic/padded-blocks': ['error', 'never'],
      '@stylistic/quote-props': ['error', 'as-needed'],
      '@stylistic/quotes': ['error', 'single', { avoidEscape: true }],
      '@stylistic/semi': ['error', 'never'],
      '@stylistic/space-before-function-paren': ['error', 'never'],
      'init-declarations': 'off',
      'no-magic-numbers': ['error', { ignore: [0, 1] }],
      'no-ternary': 'off',
      'unicorn/name-replacements': ['error', { allowList: { params: true } }],
      'unicorn/prefer-string-slice': 'off',
      'unicorn/prefer-temporal': 'off'
    }
  }
])
