import Dispatcher from '#app/Dispatcher';
import { browserHistory } from 'react-router';

export default {

  fetchInitialState() {
  },

  login(email, password, redirectTo) {
    Dispatcher.dispatch('IS_AUTHENTICATING', true);

    function sleep (time) {
      return new Promise((resolve) => setTimeout(resolve, time));
    }

    var p = sleep(1000);

    p.then(() => {
      var user = {
        id: 'marpaia',
        username: 'marpaia',
        email: 'mike@kolide.co',
        name: 'Mike Arpaia',
        admin: true,
        needs_password_reset: false
      };

      Dispatcher.dispatch('RECEIVE_USER_INFO', user);
      browserHistory.push('/');
    });

    p.then(() => {
      Dispatcher.dispatch('IS_AUTHENTICATING', false);
    });
  },
};