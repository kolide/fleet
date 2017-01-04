import expect, { spyOn, restoreSpies } from 'expect';
import { mount } from 'enzyme';

import CoreLayout from './CoreLayout';
import helpers from '../../test/helpers';
import { userStub } from 'test/stubs';

const {
  connectedComponent,
  reduxMockStore,
} = helpers;

describe('CoreLayout - layouts', () => {
  const store = { app: { config: {} }, auth: { user: userStub }, notifications: {} };
  const mockStore = reduxMockStore(store);
  const component = mount(
    connectedComponent(CoreLayout, { mockStore })
  );

  it('renders the FlashMessage component when notifications are present', () => {
    const storeWithNotifications = {
      app: { config: {} },
      auth: {
        user: userStub,
      },
      notifications: {
        alertType: 'success',
        isVisible: true,
        message: 'nice jerb!',
      },
    };
    const mockStoreWithNotifications = reduxMockStore(storeWithNotifications);
    const componentWithFlash = connectedComponent(CoreLayout, {
      mockStore: mockStoreWithNotifications,
    });
    const componentWithoutFlash = connectedComponent(CoreLayout, {
      mockStore,
    });

    const appWithFlash = mount(componentWithFlash);
    const appWithoutFlash = mount(componentWithoutFlash);

    expect(appWithFlash.length).toEqual(1);
    expect(appWithoutFlash.length).toEqual(1);

    expect(appWithFlash.find('FlashMessage').html()).toExist();
    expect(appWithoutFlash.find('FlashMessage').html()).toNotExist();
  });
});
