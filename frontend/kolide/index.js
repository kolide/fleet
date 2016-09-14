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
}

export default new Kolide();
