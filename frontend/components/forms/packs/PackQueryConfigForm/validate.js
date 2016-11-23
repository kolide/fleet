import { size } from 'lodash';

const validate = (formData) => {
  const errors = {};

  if (!size(formData.queries)) {
    errors.queries = 'Queries must be selected';
  }

  const valid = !size(errors);

  return { valid, errors };
};

export default validate;
