import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Dropdown from 'components/forms/fields/Dropdown';
import dropdownOptionInterface from 'interfaces/dropdownOption';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import validate from 'components/forms/packs/PackQueryConfigForm/validate';

const baseClass = 'pack-query-config-form';
const fieldNames = ['query_id', 'interval', 'platform', 'min_osquery_version', 'logging_type'];
const platformOptions = [
  { label: 'All', value: 'all' },
  { label: 'Windows', value: 'windows' },
  { label: 'Linux', value: 'linux' },
  { label: 'Darwin', value: 'darwin' },
];
const loggingTypeOptions = [
  { label: 'Differential', value: 'differential' },
  { label: 'Differential (Ignore Removals)', value: 'differential_ignore_removals' },
  { label: 'Snapshot', value: 'snapshot' },
];

class PackQueryConfigForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      interval: formFieldInterface.isRequired,
      logging_type: formFieldInterface.isRequired,
      min_osquery_version: formFieldInterface.isRequired,
      platform: formFieldInterface.isRequired,
      query_id: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    onCancel: PropTypes.func,
    queryOptions: PropTypes.arrayOf(dropdownOptionInterface),
  };

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  render () {
    const { fields, handleSubmit, queryOptions } = this.props;
    const { onCancel } = this;

    return (
      <form className={`${baseClass}__wrapper`}>
        <div className={`${baseClass}__body-section`}>
          <Dropdown
            options={queryOptions}
            onSelect={fields.query_id.onChange}
            placeholder="Select Query"
            wrapperClassName={`${baseClass}__form-field`}
          />
          <InputField
            {...fields.interval}
            inputWrapperClass={`${baseClass}__form-field`}
            placeholder="Interval (seconds)"
          />
          <Dropdown
            {...fields.platform}
            options={platformOptions}
            onSelect={fields.platform.onChange}
            placeholder="Platform"
            wrapperClassName={`${baseClass}__form-field`}
          />
          <Dropdown
            {...fields.logging_type}
            options={loggingTypeOptions}
            onSelect={fields.logging_type.onChange}
            placeholder="Logging type"
            wrapperClassName={`${baseClass}__form-field`}
          />
          <div className={`${baseClass}__btn-wrapper`}>
            <Button
              className={`${baseClass}__cancel-btn`}
              onClick={onCancel}
              text="cancel"
              variant="unstyled"
            />
            <Button
              className={`${baseClass}__submit-btn`}
              onClick={handleSubmit}
              text="add to pack"
              type="submit"
              variant="unstyled"
            />
          </div>
        </div>
      </form>
    );
  }
}

export default Form(PackQueryConfigForm, {
  fields: fieldNames,
  validate,
});
