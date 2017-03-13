import expect from 'expect';
import { mount } from 'enzyme';

import { connectedComponent, reduxMockStore } from 'test/helpers';
import local from 'utilities/local';
import LoginPage from 'pages/LoginPage';
import LoginPageObject from 'test/PageObjects/LoginPage';

describe('LoginPage - component', () => {
  describe('acceptance', () => {
    after(() => global.webDriver.quit());

    it('logs a user in', () => {
      const driver = global.webDriver;
      const page = LoginPageObject(driver);

      page.navigate();
      page.enterUsername('admin');
      page.enterPassword('p@ssw0rd');
      page.submit();

      return page.getSuccessText()
        .then((text) => {
          expect(text).toInclude('Taking you to the Kolide application...');
        });
    }).timeout(5000);
  });

  context('when the user is not logged in', () => {
    const mockStore = reduxMockStore({ auth: {} });

    it('renders the LoginForm', () => {
      const page = mount(connectedComponent(LoginPage, { mockStore }));

      expect(page.find('LoginForm').length).toEqual(1);
    });
  });

  context('when the users session is not recognized', () => {
    const mockStore = reduxMockStore({
      auth: {
        errors: { base: 'Unable to authenticate the current user' },
      },
    });

    it('renders the LoginForm base errors', () => {
      const page = mount(connectedComponent(LoginPage, { mockStore }));
      const loginForm = page.find('LoginForm');

      expect(loginForm.length).toEqual(1);
      expect(loginForm.prop('serverErrors')).toEqual({
        base: 'Unable to authenticate the current user',
      });
    });
  });

  context('when the user is logged in', () => {
    beforeEach(() => {
      local.setItem('auth_token', 'fake-auth-token');
    });

    const user = { id: 1, firstName: 'Bill', lastName: 'Shakespeare' };

    it('redirects to the home page', () => {
      const mockStore = reduxMockStore({ auth: { user } });
      const props = { pathname: '/login' };
      const redirectAction = {
        type: '@@router/CALL_HISTORY_METHOD',
        payload: {
          method: 'push',
          args: ['/'],
        },
      };

      mount(connectedComponent(LoginPage, { props, mockStore }));
      expect(mockStore.getActions()).toInclude(redirectAction);
    });
  });
});
