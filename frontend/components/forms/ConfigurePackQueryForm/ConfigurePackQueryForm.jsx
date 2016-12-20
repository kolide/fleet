import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import validate from 'components/forms/ConfigurePackQueryForm/validate';

const baseClass = 'configure-pack-query-form';
const fieldNames = ['query_id', 'interval', 'platform', 'min_osquery_version', 'logging_type'];
const platformOptions = [
  { label: 'All', value: 'all' },
  { label: 'Windows', value: 'windows' },
  { label: 'Linux', value: 'linux' },
  { label: 'macOS', value: 'darwin' },
];
const loggingTypeOptions = [
  { label: 'Differential', value: 'differential' },
  { label: 'Differential (Ignore Removals)', value: 'differential_ignore_removals' },
  { label: 'Snapshot', value: 'snapshot' },
];

class ConfigurePackQueryForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      interval: formFieldInterface.isRequired,
      logging_type: formFieldInterface.isRequired,
      min_osquery_version: formFieldInterface.isRequired,
      platform: formFieldInterface.isRequired,
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
      <form className={`${baseClass}__wrapper`} onSubmit={handleSubmit}>
        <div className={`${baseClass}__body-section`}>
          <InputField
            {...fields.interval}
            inputWrapperClass={`${baseClass}__form-field ${baseClass}__form-field--interval`}
            placeholder="Interval (seconds)"
          />
          <Dropdown
            {...fields.platform}
            options={platformOptions}
            placeholder="Platform"
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--platform`}
          />
          <Dropdown
            {...fields.logging_type}
            options={loggingTypeOptions}
            placeholder="Logging type"
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--logging`}
          />
          <div className={`${baseClass}__btn-wrapper`}>
            <Button
              className={`${baseClass}__cancel-btn`}
              onClick={onCancel}
              text="Cancel"
              variant="inverse"
            />
            <Button
              className={`${baseClass}__submit-btn`}
              text="Save"
              type="submit"
              variant="brand"
            />
          </div>
        </div>
      </form>
    );
  }
}

export default Form(ConfigurePackQueryForm, {
  fields: fieldNames,
  validate,
});
