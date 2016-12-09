import expect from 'expect';
import { mount } from 'enzyme';

import AppSettingsPage from 'pages/admin/AppSettingsPage';
import testHelpers from 'test/helpers';

describe.only('AppSettingsPage - component', () => {
  it('renders', () => {
    const page = mount(testHelpers.connectedComponent(AppSettingsPage));

    expect(page.find('AppSettingsPage').length).toEqual(1);
  });
});
