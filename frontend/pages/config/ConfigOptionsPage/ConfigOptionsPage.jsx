import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter, noop } from 'lodash';

import Button from 'components/buttons/Button';
import configOptionActions from 'redux/nodes/entities/config_options/actions';
import ConfigOptionsForm from 'components/forms/ConfigOptionsForm';
import configOptionInterface from 'interfaces/config_option';
import entityGetter from 'redux/utilities/entityGetter';
import helpers from 'pages/config/ConfigOptionsPage/helpers';

const baseClass = 'config-options-page';

export class ConfigOptionsPage extends Component {
  static propTypes = {
    configOptions: PropTypes.arrayOf(configOptionInterface),
    dispatch: PropTypes.func.isRequired,
  };

  static defaultProps = {
    configOptions: [],
    dispatch: noop,
  };

  componentWillMount () {
    const { configOptions, dispatch } = this.props;

    if (!configOptions.length) {
      dispatch(configOptionActions.loadAll());
    }

    return false;
  }

  onRemoveOption = (option) => {
    console.log('option removed', option);

    return false;
  }

  render () {
    const { configOptions } = this.props;
    const { onRemoveOption } = this;
    const completedOptions = filter(configOptions, option => option.value);

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
            <Button block className={`${baseClass}__reset-btn`} variant="inverse">
              RESET TO DEFAULT
            </Button>
            <Button block className={`${baseClass}__save-btn`} variant="brand">
              SAVE OPTIONS
            </Button>
          </div>
        </div>
        <ConfigOptionsForm
          configNameOptions={helpers.configOptionDropdownOptions(configOptions)}
          completedOptions={completedOptions}
          onRemoveOption={onRemoveOption}
        />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: configOptions } = entityGetter(state).get('config_options');

  return { configOptions };
};

export default connect(mapStateToProps)(ConfigOptionsPage);
