import config from '../config';

class Base {
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

  post(endpoint, body = {}, overrideHeaders = {}) {
    return this._request('POST', endpoint, body, overrideHeaders);
  }

  _request (method, endpoint, body, overrideHeaders) {
    const headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    };

    return fetch(endpoint, {
      method,
      headers: {
        ...headers,
        ...overrideHeaders
      },
      body,
    })
      .then(response => {
        if (response.ok) {
          return response.json();
        }

        const error = new Error(response.statusText);
        error.response = response;

        throw error;
      });
  }
}

export default Base;

