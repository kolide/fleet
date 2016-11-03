import expect from 'expect';
import { mount } from 'enzyme';
import { connectedComponent, reduxMockStore } from 'test/helpers';
import { validGetQueryRequest, invalidGetQueryRequest } from 'test/mocks';

import QueryPageWrapper from './QueryPageWrapper';

const bearerToken = 'abc123';
const storeWithoutQuery = {
  entities: {
    queries: {
      data: {},
    },
  },
};

describe.only('QueryPageWrapper - component', () => {
  beforeEach(() => {
    global.localStorage.setItem('KOLIDE::auth_token', bearerToken);
  });

  describe('/queries/:id', () => {
    const queryID = '10';
    const locationProp = { params: { id: queryID } };

    it('dispatches an action to get the query when there is no query', (done) => {
      validGetQueryRequest(bearerToken, queryID);

      const mockStore = reduxMockStore(storeWithoutQuery);

      mount(connectedComponent(QueryPageWrapper, { mockStore, props: locationProp }));

      setTimeout(() => {
        const dispatchedActions = mockStore.getActions().map((action) => { return action.type; });
        expect(dispatchedActions).toInclude('queries_LOAD_SUCCESS');
        done();
      }, 1500);
    });

    it('dispatches an action to transition route and display flash when API call errors', (done) => {
      invalidGetQueryRequest(bearerToken, queryID);

      const mockStore = reduxMockStore(storeWithoutQuery);

      mount(connectedComponent(QueryPageWrapper, { mockStore, props: locationProp }));

      setTimeout(() => {
        const dispatchedActions = mockStore.getActions().map((action) => { return action.type; });
        expect(dispatchedActions).toInclude('@@router/CALL_HISTORY_METHOD');
        expect(dispatchedActions).toInclude('RENDER_FLASH');
        done();
      }, 1500);
    });
  });
});
