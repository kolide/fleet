import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import { isEqual, size } from 'lodash';

import Button from 'components/buttons/Button';
import DropdownButton from 'components/buttons/DropdownButton';
import Dropdown from 'components/forms/fields/Dropdown';
import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import helpers from 'components/forms/queries/QueryForm/helpers';
import InputField from 'components/forms/fields/InputField';
import queryInterface from 'interfaces/query';
import SelectTargetsDropdown from 'components/forms/fields/SelectTargetsDropdown';
import targetInterface from 'interfaces/target';
import validatePresence from 'components/forms/validators/validate_presence';

const baseClass = 'query-form';

const validate = (formData) => {
  const errors = {};

  if (!formData.name) {
    errors.name = 'Query title must be present';
  }

  const valid = !size(errors);

  return { valid, errors };
};

class QueryForm extends Component {
  static propTypes = {
    fields: PropTypes.shape({
      description: formFieldInterface.isRequired,
      name: formFieldInterface.isRequired,
      query: formFieldInterface.isRequired,
    }).isRequired,
    onCancel: PropTypes.func,
    onFetchTargets: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    onRunQuery: PropTypes.func,
    onSave: PropTypes.func,
    onStopQuery: PropTypes.func,
    onTargetSelect: PropTypes.func,
    onUpdate: PropTypes.func,
    query: queryInterface,
    queryIsRunning: PropTypes.bool,
    queryType: PropTypes.string,
    selectedTargets: PropTypes.arrayOf(targetInterface),
    targetsCount: PropTypes.number,
  };

  constructor (props) {
    super(props);

    this.state = { errors: {} };
  }

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  onLoad = (editor) => {
    editor.setOptions({
      enableLinking: true,
    });

    editor.on('linkClick', (data) => {
      const { type, value } = data.token;
      const { onOsqueryTableSelect } = this.props;

      if (type === 'osquery-token') {
        return onOsqueryTableSelect(value);
      }

      return false;
    });
  }

  onUpdate = (evt) => {
    evt.preventDefault();

    const { fields } = this.props;
    const { valid } = this;
    const { onUpdate: handleUpdate } = this.props;

    if (valid()) {
      handleUpdate(fields);
    }

    return false;
  }

  valid = () => {
    const { errors } = this.state;
    const { fields, queryType } = this.props;

    const namePresent = validatePresence(fields.name.value);

    if (!namePresent) {
      this.setState({
        errors: {
          ...errors,
          name: `${queryType === 'label' ? 'Label title' : 'Query title'} must be present`,
        },
      });

      return false;
    }

    // TODO: validate queryText

    return true;
  }

  renderButtons = () => {
    const { canSaveAsNew, canSaveChanges } = helpers;
    const {
      fields,
      onRunQuery,
      onStopQuery,
      query,
      queryIsRunning,
      queryType,
    } = this.props;
    const { onCancel, onSave, onUpdate } = this;

    const dropdownBtnOptions = [{
      disabled: !canSaveChanges(formData, query),
      label: 'Save Changes',
      onClick: onUpdate,
    }, {
      disabled: !canSaveAsNew(formData, query),
      label: 'Save As New...',
      onClick: onSave,
    }];

    let runQueryButton;

    if (queryIsRunning) {
      runQueryButton = (
        <Button
          className={`${baseClass}__stop-query-btn`}
          onClick={onStopQuery}
          variant="alert"
        >
          Stop Query
        </Button>
      );
    } else {
      runQueryButton = (
        <Button
          className={`${baseClass}__run-query-btn`}
          onClick={onRunQuery}
          variant="brand"
        >
          Run Query
        </Button>
      );
    }

    if (queryType === 'label') {
      return (
        <div className={`${baseClass}__button-wrap`}>
          <Button
            className={`${baseClass}__save-changes-btn`}
            onClick={onCancel}
            variant="inverse"
          >
            Cancel
          </Button>
          <Button
            className={`${baseClass}__save-as-new-btn`}
            disabled={!canSaveAsNew(fields, query)}
            type="submit"
            variant="brand"
          >
            Save Label
          </Button>
        </div>
      );
    }

    return (
      <div className={`${baseClass}__button-wrap`}>
        <DropdownButton
          className={`${baseClass}__save`}
          options={dropdownBtnOptions}
          variant="success"
        >
          Save
        </DropdownButton>

        {runQueryButton}
      </div>
    );
  }

  renderPlatformDropdown = () => {
    const { fields, queryType } = this.props;

    if (queryType !== 'label') {
      return false;
    }

    const { platformOptions } = helpers;

    return (
      <Dropdown
        {...fields.platform}
        options={platformOptions}
      />
    );
  }

  renderTargetsInput = () => {
    const {
      onFetchTargets,
      onTargetSelect,
      queryType,
      selectedTargets,
      targetsCount,
    } = this.props;

    if (queryType === 'label') {
      return false;
    }


    return (
      <div>
        <SelectTargetsDropdown
          onFetchTargets={onFetchTargets}
          onSelect={onTargetSelect}
          selectedTargets={selectedTargets}
          targetsCount={targetsCount}
          label="Select Targets"
        />
      </div>
    );
  }

  render () {
    const { errors } = this.state;
    const { fields, handleSubmit, queryIsRunning, queryType } = this.props;
    const { onLoad, renderPlatformDropdown, renderButtons, renderTargetsInput } = this;

    return (
      <form className={baseClass} onSubmit={handleSubmit}>
        <h1>{queryType === 'label' ? 'New Label Query' : 'New Query'}</h1>
        <div className="query-composer__text-editor-wrapper">
          <AceEditor
            {...fields.query}
            enableBasicAutocompletion
            enableLiveAutocompletion
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={2}
            maxLines={20}
            name="query-editor"
            onLoad={onLoad}
            readOnly={queryIsRunning}
            setOptions={{ enableLinking: true }}
            showGutter
            showPrintMargin={false}
            theme="kolide"
            width="100%"
            fontSize={14}
          />
        </div>
        {renderTargetsInput()}
        <InputField
          {...fields.name}
          error={fields.name.error || errors.name}
          inputClassName={`${baseClass}__query-title`}
          label={queryType === 'label' ? 'Label title' : 'Query Title'}
        />
        <InputField
          {...fields.description}
          inputClassName={`${baseClass}__query-description`}
          label="Description"
          type="textarea"
        />
        {renderPlatformDropdown()}
        {renderButtons()}
      </form>
    );
  }
}

export default Form(QueryForm, {
  fields: ['description', 'name', 'platform', 'query'],
  validate,
});
