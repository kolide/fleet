import expect from 'expect';
import APIClient from './index';
import { validLoginRequest } from '../../test/mocks';

describe('API client - utility', () => {
  describe('defaults', () => {
    it('sets the base URL', () => {
      expect(APIClient.baseURL).toEqual('http://localhost:8080/api');
    });
  });

  describe('#loginUser', () => {
    it('sets the bearer token', (done) => {
      const request = validLoginRequest();

      APIClient.loginUser({
        username: 'admin',
        password: 'secret',
      })
        .then(() => {
          expect(request.isDone()).toEqual(true);
          expect(APIClient.bearerToken).toEqual('auth_token');
          done();
        })
        .catch(done);
    });
  });
});
