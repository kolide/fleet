import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import queryInterface from 'interfaces/query';
import validate from 'components/forms/packs/PackQueryConfigForm/validate';

const baseClass = 'pack-query-config-form';
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
    formData: PropTypes.shape({
      queries: PropTypes.arrayOf(queryInterface),
    }),
    handleSubmit: PropTypes.func,
    onCancel: PropTypes.func,
  };

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  render () {
    const { fields, formData, handleSubmit } = this.props;
    const { onCancel } = this;
    const queryCount = formData.queries.length;

    return (
      <form className={`${baseClass}__wrapper`}>
        <div className={`${baseClass}__header-section`}>
          <span>Configure {queryCount} Selected {queryCount === 1 ? 'Query' : 'Queries'}</span>
          <div className={`${baseClass}__btn-wrapper`}>
            <Button
              onClick={onCancel}
              text="Cancel"
              variant="unstyled"
            />
            <span> | </span>
            <Button
              onClick={handleSubmit}
              text="Save and Close"
              type="submit"
              variant="unstyled"
            />
          </div>
        </div>
        <div className={`${baseClass}__body-section`}>
          <InputField
            {...fields.interval}
            inputWrapperClass={`${baseClass}__form-field`}
            label="Interval"
            placeholder="Interval (seconds)"
          />
          <Dropdown
            {...fields.platform}
            className={`${baseClass}__form-field`}
            options={platformOptions}
            onSelect={fields.platform.onChange}
            placeholder="Platform"
          />
          <Dropdown
            {...fields.logging_type}
            className={`${baseClass}__form-field`}
            options={loggingTypeOptions}
            onSelect={fields.logging_type.onChange}
            placeholder="Logging type"
          />
        </div>
      </form>
    );
  }
}

export default Form(PackQueryConfigForm, {
  fields: fieldNames,
  validate,
});
