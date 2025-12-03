const { defineConfig } = require('vitest/config');

module.exports = defineConfig({
    test: {
        globals: true,
        environment: 'node',
        include: ['tests/**/*.test.js', 'tests/**/*.spec.js'],
        setupFiles: './tests/setup-vitest.js',
        // mongodb-memory-server does not always play well with worker threads
        threads: false,
        // Increase default timeouts to allow mongodb-memory-server to download/start
        timeout: 120000,
        hookTimeout: 120000,
    },
});
