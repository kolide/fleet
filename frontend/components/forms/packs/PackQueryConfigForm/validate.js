import { size } from 'lodash';

const validate = () => {
  const errors = {};

  const valid = !size(errors);

  return { valid, errors };
};

export default validate;
