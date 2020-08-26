'use strict';

import calculator from '../lib/index.js';

test('adds 1 + 2 to equal 3', () => {
    expect(calculator(1, 2)).toBe(3);
});