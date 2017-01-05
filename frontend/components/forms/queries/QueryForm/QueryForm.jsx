import React, { Component, PropTypes } from 'react';
import { isEqual } from 'lodash';
import AceEditor from 'react-ace';
import 'brace/mode/sql';
import 'brace/ext/linking';

import Button from 'components/buttons/Button';
import DropdownButton from 'components/buttons/DropdownButton';
import Dropdown from 'components/forms/fields/Dropdown';
import helpers from 'components/forms/queries/QueryForm/helpers';
import InputField from 'components/forms/fields/InputField';
import queryInterface from 'interfaces/query';
import validatePresence from 'components/forms/validators/validate_presence';
import './mode';
import './theme';

const baseClass = 'query-form';

class QueryForm extends Component {
  static propTypes = {
    onCancel: PropTypes.func,
    onSave: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    onUpdate: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    query: queryInterface,
    queryIsRunning: PropTypes.bool,
    queryText: PropTypes.string.isRequired,
    queryType: PropTypes.string,
  };

  static defaultProps = {
    query: {},
  };

  constructor (props) {
    super(props);

    const {
      query: { description, name },
      queryText,
      queryType,
    } = this.props;
    const errors = { description: null, name: null };
    const formData = { description, name, query: queryText };

    if (queryType === 'label') {
      const { allPlatforms: platform } = helpers;

      this.state = {
        errors: { ...errors, platform: null },
        formData: { ...formData, platform: platform.value },
      };
    } else {
      this.state = { errors, formData };
    }
  }

  componentDidMount = () => {
    const { query, queryText } = this.props;
    const { description, name } = query;
    const { formData } = this.state;

    this.setState({
      formData: {
        ...formData,
        description,
        name,
        query: queryText,
      },
    });
  }

  componentWillReceiveProps = (nextProps) => {
    const { query, queryText } = nextProps;
    const { query: staleQuery, queryText: staleQueryText } = this.props;

    if (!isEqual(query, staleQuery) || !isEqual(queryText, staleQueryText)) {
      const { formData } = this.state;

      this.setState({
        formData: {
          ...formData,
          description: query.description || formData.description,
          name: query.name || formData.name,
          query: queryText,
        },
      });
    }

    return false;
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

  onCancel = (evt) => {
    evt.preventDefault();

    const { onCancel: handleCancel } = this.props;

    return handleCancel();
  }

  onFieldChange = (name) => {
    return (value) => {
      const { errors, formData } = this.state;

      this.setState({
        errors: {
          ...errors,
          [name]: null,
        },
        formData: {
          ...formData,
          [name]: value,
        },
      });

      return false;
    };
  }

  onSave = (evt) => {
    evt.preventDefault();

    const { formData } = this.state;
    const { valid } = this;
    const { onSave: handleSave } = this.props;

    if (valid()) {
      handleSave(formData);
    }

    return false;
  }

  onUpdate = (evt) => {
    evt.preventDefault();

    const { formData } = this.state;
    const { valid } = this;
    const { onUpdate: handleUpdate } = this.props;

    if (valid()) {
      handleUpdate(formData);
    }

    return false;
  }

  valid = () => {
    const { errors, formData: { name } } = this.state;
    const { queryType } = this.props;

    const namePresent = validatePresence(name);

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
    const { formData } = this.state;
    const {
      query,
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
            disabled={!canSaveAsNew(formData, query)}
            onClick={onSave}
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
      </div>
    );
  }

  renderPlatformDropdown = () => {
    const { queryType } = this.props;

    if (queryType !== 'label') {
      return false;
    }

    const { formData: { platform } } = this.state;
    const { onFieldChange } = this;
    const { platformOptions } = helpers;

    return (
      <Dropdown
        options={platformOptions}
        onChange={onFieldChange('platform')}
        value={platform}
      />
    );
  }

  render () {
    const {
      errors,
      formData: {
        description,
        name,
      },
    } = this.state;
    const { onLoad, onFieldChange, renderPlatformDropdown, renderButtons } = this;
    const { queryType, onTextEditorInputChange, queryIsRunning, queryText } = this.props;

    return (
      <form className={baseClass}>
        {renderButtons()}
        <InputField
          error={errors.name}
          label={queryType === 'label' ? 'Label title' : 'Query Title'}
          name="name"
          onChange={onFieldChange('name')}
          value={name}
          inputClassName={`${baseClass}__query-title`}
        />

        <div className={`${baseClass}__text-editor-wrapper`}>
          <label className="form-field__label" htmlFor="query-editor">SQL</label>
          <AceEditor
            enableBasicAutocompletion
            enableLiveAutocompletion
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={2}
            maxLines={20}
            name="query-editor"
            onLoad={onLoad}
            onChange={onTextEditorInputChange}
            readOnly={queryIsRunning}
            setOptions={{ enableLinking: true }}
            showGutter
            showPrintMargin={false}
            theme="kolide"
            value={queryText}
            width="100%"
            fontSize={14}
          />
        </div>

        <InputField
          error={errors.description}
          label="Description"
          name="description"
          onChange={onFieldChange('description')}
          value={description}
          type="textarea"
          inputClassName={`${baseClass}__query-description`}
        />
        {renderPlatformDropdown()}
      </form>
    );
  }
}

export default QueryForm;
