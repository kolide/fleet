import React, { Component, PropTypes } from 'react';
import radium from 'radium';
import componentStyles from './styles';
import GradientButton from '../../../buttons/GradientButton';
import InputField from '../../fields/InputField';
import validatePresence from '../../validators/validate_presence';

class SaveQueryForm extends Component {
  static propTypes = {
    onSubmit: PropTypes.func,
  };

  constructor (props) {
    super(props);

    this.state = {
      errors: {
        queryName: null,
        queryDescription: null,
        queryDuration: null,
        queryPlatform: null,
        queryHosts: null,
        queryHostsPercentage: null,
        scanInterval: null,
      },
      formData: {
        queryName: null,
        queryDescription: null,
        queryDuration: 'short',
        queryPlatform: 'all',
        queryHosts: 'all',
        queryHostsPercentage: null,
        scanInterval: 0,
      },
      runType: 'run',
      showMoreOptions: false,
    };
  }

  onFieldChange = (fieldName) => {
    return ({ target }) => {
      const { errors, formData } = this.state;

      this.setState({
        errors: {
          ...errors,
          [fieldName]: null,
        },
        formData: {
          ...formData,
          [fieldName]: target.value,
        },
      });
    };
  }

  onFormSubmit = (evt) => {
    evt.preventDefault();

    const { formData, runType } = this.state;
    const { onSubmit } = this.props;
    const { validate } = this;

    if (validate()) return onSubmit({ formData, runType });

    return false;
  }

  validate = () => {
    const {
      errors,
      formData: {
        queryName,
      },
    } = this.state;

    if (!validatePresence(queryName)) {
      this.setState({
        errors: {
          ...errors,
          queryName: 'Query Name field must be completed',
        },
      });

      return false;
    }

    return true;
  }

  toggleShowMoreOptions = () => {
    const { showMoreOptions } = this.state;

    this.setState({
      showMoreOptions: !showMoreOptions,
    });

    return false;
  };

  runAndSaveQuery = (evt) => {
    evt.preventDefault();

    this.setState({
      runType: 'runAndSave',
    });

    return this.onFormSubmit(evt);
  }

  saveQuery = (evt) => {
    evt.preventDefault();

    this.setState({
      runType: 'save',
    });

    return this.onFormSubmit(evt);
  }

  renderMoreOptionsCtaSection = () => {
    const { moreOptionsIconStyles, moreOptionsCtaSectionStyles, moreOptionsTextStyles } = componentStyles;
    const { showMoreOptions } = this.state;
    const { toggleShowMoreOptions } = this;

    if (showMoreOptions) {
      return (
        <div style={moreOptionsCtaSectionStyles}>
          <span onClick={toggleShowMoreOptions} style={moreOptionsTextStyles}>
            Fewer Options
            <i className="kolidecon-upcarat" style={moreOptionsIconStyles} />
          </span>
        </div>
      );
    }

    return (
      <div style={moreOptionsCtaSectionStyles}>
        <span onClick={toggleShowMoreOptions} style={moreOptionsTextStyles}>
          More Options
          <i className="kolidecon-downcarat" style={moreOptionsIconStyles} />
        </span>
      </div>
    );
  }

  renderMoreOptionsFormFields = () => {
    const {
      errors,
      formData: {
        queryDuration,
        queryPlatform,
        queryHosts,
      },
      showMoreOptions,
    } = this.state;
    const {
      dropdownInputStyles,
      formSectionStyles,
      helpTextStyles,
      labelStyles,
      queryDescriptionInputStyles,
      queryHostsPercentageStyles,
      queryNameInputStyles,
    } = componentStyles;
    const { onFieldChange } = this;

    if (!showMoreOptions) return false;

    return (
      <div>
        <div style={formSectionStyles}>
          <InputField
            error={errors.queryDescription}
            label="Query Description"
            labelStyles={labelStyles}
            name="queryDescription"
            onChange={onFieldChange('queryDescription')}
            placeholder="e.g. This query does x, y, & z because n"
            style={queryDescriptionInputStyles}
            type="textarea"
          />
          <small style={helpTextStyles}>
            If your query is really complex and/or it is not clear why you wrote this query, you should write a description so others can reuse this query for the correct reason.
          </small>
        </div>
        <div style={formSectionStyles}>
          <div>
            <label htmlFor="queryDuration" style={labelStyles}>Query Duration</label>
            <select
              key="queryDuration"
              name="queryDuration"
              value={queryDuration}
              onChange={onFieldChange('queryDuration')}
              style={dropdownInputStyles}
            >
              <option value="short">Short</option>
              <option value="long">Long</option>
            </select>
          </div>
          <small style={helpTextStyles}>
            Individual hosts are not always online. A longer duration will return more complete results. You can view results of any in-progress query at any time.
          </small>
        </div>
        <div style={formSectionStyles}>
          <div>
            <label htmlFor="queryPlatform" style={labelStyles}>Query Platform</label>
            <select
              key="queryPlatform"
              name="queryPlatform"
              value={queryPlatform}
              onChange={onFieldChange('queryPlatform')}
              style={dropdownInputStyles}
            >
              <option value="all">ALL PLATFORMS</option>
              <option value="none">NO PLATFORMS</option>
            </select>
          </div>
          <small style={helpTextStyles}>
            Specifying a platform allows you to restrict the query from running on a certain platform (even on hosts specifically targeted that do not match).
          </small>
        </div>
        <div style={formSectionStyles}>
          <div>
            <label htmlFor="queryHosts" style={labelStyles}>Run On All Hosts?</label>
            <div>
              <input
                checked={queryHosts === 'all'}
                onChange={onFieldChange('queryHosts')}
                type="radio"
                value="all"
              /> Run Query On All Hosts
              <br />
              <input
                checked={queryHosts === 'percentage'}
                onChange={onFieldChange('queryHosts')}
                type="radio"
                value="percentage"
              /> Run Query On
              <InputField
                inputWrapperStyles={{ display: 'inline-block' }}
                inputOptions={{ maxLength: 3 }}
                onChange={onFieldChange('queryHostsPercentage')}
                style={queryHostsPercentageStyles}
                type="tel"
              />% Of All Hosts
            </div>
          </div>
          <small style={helpTextStyles}>
            Specifying a platform allows you to restrict the query from running on a certain platform (even on hosts specifically targeted that do not match).
          </small>
        </div>
        <div style={formSectionStyles}>
          <InputField
            error={errors.scanInterval}
            label="Scan Interval (seconds)"
            labelStyles={labelStyles}
            name="scanInterval"
            onChange={onFieldChange('scanInterval')}
            placeholder="e.g. 300"
            style={queryNameInputStyles}
            type="tel"
          />
          <small style={helpTextStyles}>
            You can use queries you write in "scans". The interval can be used to control how frequently the query runs when it is running continuously.
          </small>
        </div>
      </div>
    );
  };

  render () {
    const {
      buttonInvertStyles,
      buttonStyles,
      buttonWrapperStyles,
      labelStyles,
      helpTextStyles,
      queryNameInputStyles,
      queryNameWrapperStyles,
    } = componentStyles;
    const { errors } = this.state;
    const {
      onFieldChange,
      onFormSubmit,
      renderMoreOptionsFormFields,
      renderMoreOptionsCtaSection,
      runAndSaveQuery,
      saveQuery,
    } = this;

    return (
      <form onSubmit={onFormSubmit}>
        <div style={queryNameWrapperStyles}>
          <InputField
            error={errors.queryName}
            label="Query Name"
            labelStyles={labelStyles}
            name="queryName"
            onChange={onFieldChange('queryName')}
            placeholder="e.g. Interesting Query Name"
            style={queryNameInputStyles}
          />
          <small style={helpTextStyles}>
            Write a name that describes the query and its intent. Pick a name that others will find useful.
          </small>
        </div>
        {renderMoreOptionsCtaSection()}
        {renderMoreOptionsFormFields()}
        <div style={buttonWrapperStyles}>
          <GradientButton
            onClick={saveQuery}
            style={buttonInvertStyles}
            text="Save Query Only"
          />

          <GradientButton
            onClick={runAndSaveQuery}
            style={buttonStyles}
            text="Run & Save Query"
          />
        </div>
      </form>
    );
  }
}

export default radium(SaveQueryForm);
