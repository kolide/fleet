import fetch from 'isomorphic-fetch';
import config from '../../config';
import endpoints from './endpoints';

class APIClient {
  constructor () {
    this.baseURL = this.setBaseURL();
  }

  setBaseURL () {
    const {
      settings: { env },
      environments: { development },
    } = config;

    if (env === development) {
      return 'http://localhost:8080/api';
    }

    throw new Error(`API base URL is not configured for environment: ${env}`);
  }

  setBearerToken (bearerToken) {
    this.bearerToken = bearerToken;
  }

  loginUser ({ username, password }) {
    const { LOGIN } = endpoints;
    const endpoint = this.baseURL + LOGIN;

    return fetch(endpoint, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    })
      .then(response => {
        return response.json()
          .then(user => {
            return this.setBearerToken(user.token);
          });
      });
  }
}

export default new APIClient();
