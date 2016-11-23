import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import validate from 'components/forms/packs/PackQueryConfigForm/validate';

const fieldNames = ['queries', 'interval', 'platform', 'min_osquery_version', 'logging_type'];
const platformOptions = [
  { label: 'All', value: 'all' },
  { label: 'Windows', value: 'windows' },
  { label: 'Linux', value: 'linux' },
  { label: 'Darwin', value: 'darwin' },
];
const loggingTypeOptions = [
  { label: 'Differential', value: 'differential' },
  { label: 'Differential (with Removed)', value: 'differential_with_removed' },
  { label: 'Snapshot', value: 'snapshot' },
];

class PackQueryConfigForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      interval: formFieldInterface.isRequired,
      logging_type: formFieldInterface.isRequired,
      min_osquery_version: formFieldInterface.isRequired,
      platform: formFieldInterface.isRequired,
      queries: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    onCancel: PropTypes.func,
  };

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  render () {
    const { fields, handleSubmit } = this.props;
    const { onCancel } = this;

    return (
      <form>
        <Button
          onClick={handleSubmit}
          text="Save and Close"
          type="submit"
          variant="brand"
        />
        <Button
          onClick={onCancel}
          text="Cancel"
          variant="inverse"
        />
        <InputField
          {...fields.interval}
          label="Interval"
          placeholder="Interval (seconds)"
        />
        <Dropdown
          {...fields.platform}
          options={platformOptions}
          onSelect={fields.platform.onChange}
          placeholder="Platform"
        />
        <Dropdown
          {...fields.logging_type}
          options={loggingTypeOptions}
          onSelect={fields.logging_type.onChange}
          placeholder="Logging type"
        />
      </form>
    );
  }
}

export default Form(PackQueryConfigForm, {
  fields: fieldNames,
  validate,
});
