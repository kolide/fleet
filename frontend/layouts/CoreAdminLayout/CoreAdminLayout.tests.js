import expect from 'expect';
import { mount } from 'enzyme';
import ConnectedCoreAdminLayout from './CoreAdminLayout';
import { connectedComponent, reduxMockStore } from '../../test/helpers';

describe('CoreAdminLayout - layout', () => {
  const redirectToHomeAction = {
    type: '@@router/CALL_HISTORY_METHOD',
    payload: {
      method: 'push',
      args: ['/'],
    },
  };

  it('redirects to the homepage if the user is not an admin', () => {
    const user = { id: 1, admin: false };
    const storeWithoutAdminUser = { auth: { user } };
    const mockStore = reduxMockStore(storeWithoutAdminUser);
    mount(
      connectedComponent(ConnectedCoreAdminLayout, { mockStore })
    );

    expect(mockStore.getActions()).toInclude(redirectToHomeAction);
  });

  it('does not redirect if the user is an admin', () => {
    const user = { id: 1, admin: true };
    const storeWithAdminUser = { auth: { user } };
    const mockStore = reduxMockStore(storeWithAdminUser);
    mount(
      connectedComponent(ConnectedCoreAdminLayout, { mockStore })
    );

    expect(mockStore.getActions()).toNotInclude(redirectToHomeAction);
  });
});
