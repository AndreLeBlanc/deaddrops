import expect from 'expect';

function randomFunc() {
	return "Random func";
}

describe('Simple Math', () => {
	it('2 should be equal to 2', () => {
		expect(2).toEqual(2);
	})
})

describe('Random Function', () => {
	it('Should return "Random func"', () => {
		expect(randomFunc()).toEqual("Random func");
	})
})