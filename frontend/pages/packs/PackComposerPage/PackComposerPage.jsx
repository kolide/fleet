import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { noop } from 'lodash';

import PackForm from 'components/forms/PackForm';
import queryActions from 'redux/nodes/entities/queries/actions';
import queryInterface from 'interfaces/query';
import QueriesListWrapper from 'components/queries/QueriesListWrapper';
import stateEntityGetter from 'redux/utilities/entityGetter';

export class PackComposerPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    queries: PropTypes.arrayOf(queryInterface),
  };

  static defaultProps = {
    dispatch: noop,
  };

  constructor (props) {
    super(props);

    this.state = { selectedTargetsCount: 0 };
  }

  componentDidMount () {
    const { dispatch } = this.props;

    dispatch(queryActions.loadAll());
  }

  onFetchTargets = (query, targetsResponse) => {
    const {
      selected_targets_count: selectedTargetsCount,
    } = targetsResponse;

    this.setState({ selectedTargetsCount });

    return false;
  }

  handleSubmit = (formData) => {
    console.log(formData);

    return false;
  }

  render () {
    const { handleSubmit, onFetchTargets } = this;
    const { selectedTargetsCount } = this.state;
    const { queries } = this.props;

    return (
      <div>
        <PackForm
          handleSubmit={handleSubmit}
          onFetchTargets={onFetchTargets}
          selectedTargetsCount={selectedTargetsCount}
        />
        <QueriesListWrapper queries={queries} />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: queries } = stateEntityGetter(state).get('queries');

  return { queries };
};

export default connect(mapStateToProps)(PackComposerPage);
