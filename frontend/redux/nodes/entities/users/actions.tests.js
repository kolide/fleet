import expect, { restoreSpies, spyOn } from 'expect';

import * as Kolide from 'kolide';

import { reduxMockStore } from '../../../../test/helpers';

import {
  requirePasswordReset,
  REQUIRE_PASSWORD_RESET_REQUEST,
  REQUIRE_PASSWORD_RESET_FAILURE,
  REQUIRE_PASSWORD_RESET_SUCCESS,
} from './actions';

const store = { entities: { invites: {}, users: {} } };
const user = { id: 1, email: 'zwass@kolide.co' };

describe('Users - actions', () => {
  describe('dispatching the require password reset action', () => {
    describe('successful request', () => {
      const mockStore = reduxMockStore(store);

      beforeEach(() => {
        spyOn(Kolide.default, 'requirePasswordReset').andCall(() => {
          return Promise.resolve();
        });
      });

      afterEach(restoreSpies);

      it('calls the resetFunc', () => {
        mockStore.dispatch(requirePasswordReset(user, true));

        expect(Kolide.default.requirePasswordReset).toHaveBeenCalledWith(user, true);
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(requirePasswordReset());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude(REQUIRE_PASSWORD_RESET_REQUEST);
        expect(dispatchedActionTypes).toInclude(REQUIRE_PASSWORD_RESET_SUCCESS);
        expect(dispatchedActionTypes).toNotInclude(REQUIRE_PASSWORD_RESET_FAILURE);
      });
    });

    describe('unsuccessful request', () => {
      const mockStore = reduxMockStore(store);
      const errors = { base: 'Unable to require password reset' };

      beforeEach(() => {
        spyOn(Kolide.default, 'requirePasswordReset').andCall(() => {
          return Promise.reject({ errors });
        });
      });

      afterEach(restoreSpies);

      it('calls the resetFunc', () => {
        mockStore.dispatch(requirePasswordReset(user, true));

        expect(Kolide.default.requirePasswordReset).toHaveBeenCalledWith(user, true);
      });

      it('dispatches the correct actions', () => {
        mockStore.dispatch(requirePasswordReset());

        const dispatchedActions = mockStore.getActions();
        const dispatchedActionTypes = dispatchedActions.map((action) => { return action.type; });

        expect(dispatchedActionTypes).toInclude(REQUIRE_PASSWORD_RESET_REQUEST);
        expect(dispatchedActionTypes).toNotInclude(REQUIRE_PASSWORD_RESET_SUCCESS);
        expect(dispatchedActionTypes).toInclude(REQUIRE_PASSWORD_RESET_FAILURE);
      });
    });
  });
});
