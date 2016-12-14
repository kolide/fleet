import { size } from 'lodash';

const validate = (formData) => {
  const errors = {};

  const valid = !size(errors);

  return { valid, errors };
};

export default validate;
