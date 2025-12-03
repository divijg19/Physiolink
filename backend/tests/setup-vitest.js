// Minimal setup for vitest tests.
// Ensure test env is 'test' and stabilize timers if needed.
process.env.NODE_ENV = process.env.NODE_ENV || 'test';

// If any globals need shimming for older Jest-based tests, add them here.
// Example: make sure globalThis.fetch exists when using cross-fetch in frontend-adjacent code.

// Leave intentionally minimal — most tests perform their own setup/teardown.
