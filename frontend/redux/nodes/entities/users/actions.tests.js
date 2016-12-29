import expect, { createSpy, restoreSpies, spyOn } from 'expect';

import reduxConfig from 'redux/nodes/entities/base/reduxConfig';
import { reduxMockStore } from '../../../../test/helpers';
import * as Kolide from 'kolide';

import {
  requirePasswordReset,
  REQUIRE_PASSWORD_RESET_REQUEST,
  REQUIRE_PASSWORD_RESET_FAILURE,
  REQUIRE_PASSWORD_RESET_SUCCESS
} from './actions';

const store = { entities: { invites: {}, users: {} } };
const user = { id: 1, email: 'zwass@kolide.co' };

describe('Users - actions', () => {
  afterEach(restoreSpies);

  describe('dispatching the require password reset action', () => {
    describe('successful request', () => {
      const mockStore = reduxMockStore(store);

      it('calls the resetFunc', () => {
        const resetFunc = spyOn(Kolide.default, "requirePasswordReset").andCall(() => {
          return Promise.resolve();
        });

        mockStore.dispatch(requirePasswordReset(user, true));

        expect(resetFunc).toHaveBeenCalledWith(user, true);
      });

      it('dispatches the correct actions', () => {
        const resetFunc = spyOn(Kolide.default, "requirePasswordReset").andCall(() => {
          return Promise.resolve();
        });

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

      it('calls the resetFunc', () => {
        const resetFunc = spyOn(Kolide.default, "requirePasswordReset").andCall(() => {
          return Promise.reject({ errors });
        });

        mockStore.dispatch(requirePasswordReset(user, true));

        expect(resetFunc).toHaveBeenCalledWith(user, true);
      });

      it('dispatches the correct actions', () => {
        const resetFunc = spyOn(Kolide.default, "requirePasswordReset").andCall(() => {
          return Promise.reject({ errors });
        });

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
