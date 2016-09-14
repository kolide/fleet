import expect from 'expect';
import Kolide from './index';
import { validLoginRequest, validUser } from '../test/mocks';

describe('Kolide - API client', () => {
  describe('defaults', () => {
    it('sets the base URL', () => {
      expect(Kolide.baseURL).toEqual('http://localhost:8080/api');
    });
  });

  describe('#loginUser', () => {
    it('calls the appropriate endpoint with the correct parameters', (done) => {
      const request = validLoginRequest();

      Kolide.loginUser({
        username: 'admin',
        password: 'secret',
      })
        .then((user) => {
          expect(user).toEqual(validUser);
          expect(request.isDone()).toEqual(true);
          done();
        })
        .catch(done);
    });
  });
});
