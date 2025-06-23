import antfu from '@antfu/eslint-config'

export default antfu({
  ignores: [
    'dist/*',
    'node_modules/*',
    'components/ui/*',
  ],
  rules: {
    'unused-imports/no-unused-vars': 'warn',
    'no-console': 'warn',
    'symbol-description': 'off',
    'ts/no-empty-object-type': 'off',
  },
})
