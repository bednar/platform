module.exports = {
  projects: [
    {
      displayName: 'test',
      testPathIgnorePatterns: [
        'build',
        '<rootDir>/node_modules/(?!(jest-test))',
      ],
      modulePaths: ['<rootDir>', '<rootDir>/node_modules/'],
      moduleDirectories: ['src'],
      setupFiles: ['<rootDir>/testSetup.js'],
      transform: {
        '^.+\\.tsx?$': 'ts-jest',
        '^.+\\.js$': 'babel-jest',
      },
      testRegex: '(/__tests__/.*|(\\.|/)(test))\\.(jsx?|tsx?)$',
      moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
      transformIgnorePatterns: ['/node_modules/(?!dygraphs)'],
      snapshotSerializers: ['enzyme-to-json/serializer'],
      moduleNameMapper: {
        '\\.(css|scss)$': 'identity-obj-proxy',
      },
    },
    {
      runner: 'jest-runner-eslint',
      displayName: 'eslint',
      testMatch: ['<rootDir>/**/*.test.js'],
    },
    {
      runner: 'jest-runner-tslint',
      displayName: 'tslint',
      moduleFileExtensions: ['ts', 'tsx'],
      testMatch: ['<rootDir>/**/*.test.ts', '<rootDir>/**/*.test.tsx'],
    },
  ],
}
