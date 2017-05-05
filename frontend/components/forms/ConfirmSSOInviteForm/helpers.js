import { size } from 'lodash';
import validateEquality from 'components/forms/validators/validate_equality';

const validate = (formData) => {
  const errors = {};
  const {
    name,
    username,
  } = formData;

  if (!name) {
    errors.name = 'Full name must be present';
  }

  if (!username) {
    errors.username = 'Username must be present';
  }

  const valid = !size(errors);

  return { valid, errors };
};

export default { validate };
