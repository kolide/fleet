import expect from 'expect';
import { mount } from 'enzyme';

import ConfirmInvitePage from 'pages/ConfirmInvitePage';
import { connectedComponent, reduxMockStore } from 'test/helpers';

describe('ConfirmInvitePage - component', () => {
  const inviteToken = 'abc123';
  const params = { invite_token: inviteToken };
  const component = connectedComponent(ConfirmInvitePage, {
    props: { params },
    mockStore: reduxMockStore(),
  });
  const page = mount(component);

  it('renders', () => {
    expect(page.length).toEqual(1);
    expect(
      page.find('ConfirmInvitePage').prop('inviteFormData')
    ).toEqual(params);
  });

  it('renders a ConfirmInviteForm', () => {
    expect(page.find('ConfirmInviteForm').length).toEqual(1);
  });
});
