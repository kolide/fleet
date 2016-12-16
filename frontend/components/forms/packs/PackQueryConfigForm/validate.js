import { size } from 'lodash';

import validateNumericality from 'components/forms/validators/validate_numericality';

const validate = (formData) => {
  const errors = {};

  if (!formData.query_id) {
    errors.query_id = 'A query must be selected';
  }

  if (!formData.interval) {
    errors.interval = 'Interval must be present';
  }

  if (formData.interval && !validateNumericality(formData.interval)) {
    errors.interval = 'Interval must be a number';
  }

  if (!formData.platform) {
    errors.query_id = 'A platform must be selected';
  }

  if (!formData.logging_type) {
    errors.query_id = 'A Logging Type must be selected';
  }

  const valid = !size(errors);

  return { valid, errors };
};

export default validate;
