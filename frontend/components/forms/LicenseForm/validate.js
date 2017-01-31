import { size } from 'lodash';

import validJwtToken from 'components/forms/validators/valid_jwt_token';

export default ({ license }) => {
  const errors = {};

  if (!license) {
    errors.license = 'License must be present';
  }

  if (license && !validJwtToken(license)) {
    errors.license = 'License is not a valid JWT token';
  }

  const valid = !size(errors);

  return { errors, valid };
};
