import expect from 'expect';

import reducer from './reducer';
import {
  requirePasswordResetRequest,
  requirePasswordResetFailure,
  requirePasswordResetSuccess
} from './actions';

describe('Users - reducer', () => {
  const initialState = {
    loading: false,
    errors: {},
    data: {},
  };

  it('updates state when request is dispatched', () => {
    const newState = reducer(initialState, requirePasswordResetRequest);

    expect(newState).toEqual({
      loading: true,
      errors: {},
      data: {},
    });
  });

  it('updates state when request is successful', () => {
    let initState = {
      loading: true,
      errors: {},
      data: {},
    };
    const newState = reducer(initState, requirePasswordResetSuccess);

    expect(newState).toEqual({
      loading: false,
      errors: {},
      data: {},
    });
  });

  it('updates state when request fails', () => {
    const errors = { base: 'Unable to require password reset' };
    const newState = reducer(initialState, requirePasswordResetFailure(errors));

    expect(newState).toEqual({
      loading: false,
      errors: errors,
      data: {},
    });
  });
});
