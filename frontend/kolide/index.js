import fetch from 'isomorphic-fetch';
import Base from './base';
import endpoints from './endpoints';
import local from '../utilities/local';

class Kolide extends Base {
  loginUser ({ username, password }) {
    const { LOGIN } = endpoints;
    const loginEndpoint = this.baseURL + LOGIN;

    return this.post(loginEndpoint, JSON.stringify({ username, password }));
  }

  forgotPassword ({ email }) {
    const { FORGOT_PASSWORD } = endpoints;
    const forgotPasswordEndpoint = this.baseURL + FORGOT_PASSWORD;

    return this.post(forgotPasswordEndpoint, JSON.stringify({ email }));
  }

  resetPassword ({ newPassword, passwordResetToken }) {
    const { RESET_PASSWORD } = endpoints;
    const resetPasswordEndpoint = this.baseURL + RESET_PASSWORD;

    return this.post(resetPasswordEndpoint, JSON.stringify({
      new_password: newPassword,
      password_reset_token: passwordResetToken,
    }));
  }
}

export default new Kolide();
