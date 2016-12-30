import React, { Component, PropTypes } from 'react';
import { isEqual, noop } from 'lodash';

const defaultValidate = () => { return { valid: true, errors: {} }; };

const baseClass = 'form-hoc';

export default (WrappedComponent, { fields, validate = defaultValidate }) => {
  class Form extends Component {
    static propTypes = {
      serverErrors: PropTypes.object, // eslint-disable-line react/forbid-prop-types
      formData: PropTypes.object, // eslint-disable-line react/forbid-prop-types
      handleSubmit: PropTypes.func,
      onChangeFunc: PropTypes.func,
    };

    static defaultProps = {
      errors: {},
      formData: {},
    };

    constructor (props) {
      super(props);

      const { formData } = props;

      this.state = {
        errors: {},
        formData,
      };

      return false;
    }

    componentWillReceiveProps ({ formData, serverErrors }) {
      const {
        formData: oldFormDataProp,
        serverErrors: oldServerErrors,
      } = this.props;

      if (!isEqual(formData, oldFormDataProp)) {
        const { formData: currentFormData } = this.state;

        this.setState({
          formData: {
            ...currentFormData,
            ...formData,
          },
        });
      }

      if (!isEqual(serverErrors, oldServerErrors)) {
        const { errors } = this.state;

        this.setState({
          errors: {
            ...errors,
            ...serverErrors,
          },
        });
      }

      return false;
    }

    onFieldChange = (fieldName) => {
      return (value) => {
        const { errors, formData } = this.state;
        const { onChangeFunc = noop } = this.props;

        onChangeFunc(fieldName, value);

        this.setState({
          errors: { ...errors, [fieldName]: null },
          formData: { ...formData, [fieldName]: value },
        });

        return false;
      };
    }

    onSubmit = (evt) => {
      evt.preventDefault();

      const { handleSubmit } = this.props;
      const { errors, formData } = this.state;
      const { valid, errors: clientErrors } = validate(formData);

      if (valid) {
        return handleSubmit(formData);
      }

      this.setState({
        errors: { ...errors, ...clientErrors },
      });

      return false;
    }

    getError = (fieldName) => {
      const { errors } = this.state;

      return errors[fieldName];
    }

    getFields = () => {
      const { getError, getValue, onFieldChange } = this;
      const fieldProps = {};

      fields.forEach((field) => {
        fieldProps[field] = {
          error: getError(field),
          name: field,
          onChange: onFieldChange(field),
          value: getValue(field),
        };
      });

      return fieldProps;
    }

    getValue = (fieldName) => {
      return this.state.formData[fieldName];
    }

    render () {
      const { getFields, onSubmit, props } = this;
      const { errors } = this.state;

      return (
        <div>
          {errors.base && <div className={`${baseClass}__base-error`}>{errors.base}</div>}
          <WrappedComponent
            {...props}
            fields={getFields()}
            handleSubmit={onSubmit}
          />
        </div>
      );
    }
  }

  return Form;
};
