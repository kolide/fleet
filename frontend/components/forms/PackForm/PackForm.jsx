import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import InputField from 'components/forms/fields/InputField';
import SelectTargetsDropdown from 'components/forms/fields/SelectTargetsDropdown';

const fieldNames = ['title', 'description', 'targets'];
const validate = () => {
  return {
    valid: true,
    errors: {},
  };
};

class PackForm extends Component {
  static propTypes = {
    className: PropTypes.string,
    fields: PropTypes.shape({
      description: formFieldInterface.isRequired,
      targets: formFieldInterface.isRequired,
      title: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func,
    onFetchTargets: PropTypes.func,
    selectedTargetsCount: PropTypes.number,
  };

  render () {
    const {
      className,
      fields,
      handleSubmit,
      onFetchTargets,
      selectedTargetsCount,
    } = this.props;

    return (
      <form className={className} onSubmit={handleSubmit}>
        <InputField
          {...fields.title}
          placeholder="Query Pack Title"
        />
        <InputField
          {...fields.description}
          label="Description"
          placeholder="Add a description of your query"
          type="textarea"
        />
        <div>
          <span>{selectedTargetsCount} total hosts</span>
          <SelectTargetsDropdown
            {...fields.targets}
            label="Select Pack Targets"
            onSelect={fields.targets.onChange}
            onFetchTargets={onFetchTargets}
            selectedTargets={fields.targets.value}
          />
        </div>
        <div>
          <Button
            text="Save Query pack"
            type="submit"
            variant="brand"
          />
        </div>
      </form>
    );
  }
}

export default Form(PackForm, {
  fields: fieldNames,
  validate,
});
