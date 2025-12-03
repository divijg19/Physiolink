// simple CommonJS mock for cuid2 used in tests
module.exports = {
	createId: () => `cuidtest_${Math.random().toString(36).slice(2, 10)}`,
	init: () => {},
	getConstants: () => ({}),
	isCuid: () => false,
};
