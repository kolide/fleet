import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { difference, find, filter, isEqual, noop } from 'lodash';

import Button from 'components/buttons/Button';
import configOptionActions from 'redux/nodes/entities/config_options/actions';
import ConfigOptionsForm from 'components/forms/ConfigOptionsForm';
import configOptionInterface from 'interfaces/config_option';
import entityGetter from 'redux/utilities/entityGetter';
import helpers from 'pages/config/ConfigOptionsPage/helpers';
import replaceArrayItem from 'utilities/replace_array_item';

const baseClass = 'config-options-page';
const DEFAULT_CONFIG_OPTION = { name: '', value: '', read_only: false };

export class ConfigOptionsPage extends Component {
  static propTypes = {
    configOptions: PropTypes.arrayOf(configOptionInterface),
    dispatch: PropTypes.func.isRequired,
  };

  static defaultProps = {
    configOptions: [],
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = {
      configOptions: [],
    };
  }

  componentWillMount () {
    const { configOptions, dispatch } = this.props;

    if (!configOptions.length) {
      dispatch(configOptionActions.loadAll());

      return false;
    }

    this.setState({ configOptions });

    return false;
  }

  componentWillReceiveProps ({ configOptions }) {
    if (!isEqual(configOptions, this.state.configOptions)) {
      this.setState({
        configOptions: [
          ...this.state.configOptions,
          ...configOptions,
        ],
      });
    }

    return false;
  }

  onAddNewOption = (evt) => {
    evt.preventDefault();

    const { configOptions } = this.state;

    if (find(configOptions, DEFAULT_CONFIG_OPTION)) {
      return false;
    }

    this.setState({
      configOptions: [
        ...configOptions,
        DEFAULT_CONFIG_OPTION,
      ],
    });

    return false;
  }

  onOptionUpdate = (option, newOption) => {
    const { configOptions } = this.state;
    const newConfigOptions = replaceArrayItem(configOptions, option, newOption);

    this.setState({
      configOptions: newConfigOptions,
    });

    return false;
  }

  onRemoveOption = (option) => {
    const { configOptions } = this.state;
    const configOptionsWithoutRemovedOption = filter(configOptions, o => !isEqual(o, option));

    if (isEqual(option, DEFAULT_CONFIG_OPTION)) {
      this.setState({ configOptions: configOptionsWithoutRemovedOption });
    } else {
      this.setState({
        configOptions: [
          ...configOptionsWithoutRemovedOption,
          { ...option, value: null },
        ],
      });
    }

    return false;
  }

  onResetConfigOptions = () => {
    const { configOptions } = this.props;

    this.setState({ configOptions });

    return false;
  }

  onSave = () => {
    const changedOptions = this.calculateChangedOptions();

    console.log('changedOptions', changedOptions);

    return false;
  }

  calculateChangedOptions = () => {
    const { configOptions: stateConfigOptions } = this.state;
    const { configOptions: propConfigOptions } = this.props;

    return difference(stateConfigOptions, propConfigOptions);
  }

  render () {
    const { configOptions } = this.state;
    const { onAddNewOption, onOptionUpdate, onRemoveOption, onResetConfigOptions, onSave } = this;
    const availableOptions = filter(configOptions, option => option.value !== null);

    return (
      <div className={`body-wrap ${baseClass}`}>
        <div className={`${baseClass}__header-wrapper`}>
          <div>
            <h1>Manage Additional Osquery Options</h1>
            <p>
              Osquery allows you to set a number of configuration options (Osquery Documentation).
              Since Kolide manages your Osquery configuration, you can set these additional desired
              options on this screen. Some options that Kolide needs to function correctly will be ignored.
            </p>
          </div>
          <div className={`${baseClass}__btn-wrapper`}>
            <Button block className={`${baseClass}__reset-btn`} onClick={onResetConfigOptions} variant="inverse">
              RESET TO DEFAULT
            </Button>
            <Button block className={`${baseClass}__save-btn`} onClick={onSave} variant="brand">
              SAVE OPTIONS
            </Button>
          </div>
        </div>
        <ConfigOptionsForm
          configNameOptions={helpers.configOptionDropdownOptions(configOptions)}
          completedOptions={availableOptions}
          onFormUpdate={onOptionUpdate}
          onRemoveOption={onRemoveOption}
        />
        <Button onClick={onAddNewOption} variant="unstyled">Add New Option</Button>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: configOptions } = entityGetter(state).get('config_options');

  return { configOptions };
};

export default connect(mapStateToProps)(ConfigOptionsPage);
