import expect from 'expect';
import reducer from './reducer';
import {
  resetOptionsStart,
  resetOptionsSuccess,
  resetOptionsFailure,
} from './actions';

const options = [
  {id: 1, name: 'option1', type: 'int', value: 10},
  {id: 2, name: 'option2', type: 'string', value: 'changed'},
];

const resetOptions = [
  {id: 1, name: 'option1', type: 'int', value: 10},
  {id: 2, name: 'option2', type: 'string', value: 'original'},
];

describe('Options - reducer', () => {
  describe('reset', () => {

    it('should return options on success', () => {
      const initState = {
        loading: true,
        errors: {},
        data: {},
      };
      const newState = reducer(initState, resetOptionsSuccess(resetOptions));
      expect(newState).toEqual({
        ...initState,
        loading: false,
        data: resetOptions,
      });

    });
  });
});
