import expect from 'expect';

import helpers from 'redux/middlewares/nag_message/helpers';
import { licenseStub } from 'test/stubs';

describe('Nag message middleware - helpers', () => {
  describe('#shouldNagUser', () => {
    const { shouldNagUser } = helpers;

    it('returns true when there are more hosts than allowed hosts', () => {
      const overusedLicense = {
        ...licenseStub(),
        allowed_hosts: 2,
        hosts: 3,
      };
      const validLicense = {
        ...licenseStub(),
        allowed_hosts: 2,
        hosts: 2,
      };

      expect(shouldNagUser({ license: overusedLicense })).toEqual(true);
      expect(shouldNagUser({ license: validLicense })).toEqual(false);
    });

    it('returns true when the license is expired', () => {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);

      const license = licenseStub();
      const expiredLicense = { ...license, expiry: yesterday.toISOString() };

      expect(shouldNagUser({ license: expiredLicense })).toEqual(true, 'Expected the expired license to return true');
      expect(shouldNagUser({ license })).toEqual(false, 'Expected the valid license to return false');
    });
  });
});
