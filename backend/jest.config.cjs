module.exports = {
    testEnvironment: 'node',
    // Map the ESM-only cuid2 package used by dependencies to a simple CommonJS mock
    moduleNameMapper: {
        '^@paralleldrive/cuid2$': '<rootDir>/tests/__mocks__/cuid2.js'
    }
};
