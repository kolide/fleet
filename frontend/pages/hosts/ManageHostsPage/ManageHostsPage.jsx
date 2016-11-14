import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { filter } from 'lodash';
import ReactCSSTransitionGroup from 'react-addons-css-transition-group'

import entityGetter from 'redux/utilities/entityGetter';
import hostActions from 'redux/nodes/entities/hosts/actions';
import labelActions from 'redux/nodes/entities/labels/actions';
import labelInterface from 'interfaces/label';
import HostDetails from 'components/hosts/HostDetails';
import hostInterface from 'interfaces/host';
import HostSidePanel from 'components/side_panels/HostSidePanel';
import osqueryTableInterface from 'interfaces/osquery_table';
import QueryComposer from 'components/queries/QueryComposer';
import QuerySidePanel from 'components/side_panels/QuerySidePanel';
import { selectOsqueryTable } from 'redux/nodes/components/QueryPages/actions';
import { setSelectedLabel } from 'redux/nodes/components/ManageHostsPage/actions';
import { showRightSidePanel, removeRightSidePanel } from 'redux/nodes/app/actions';

export class ManageHostsPage extends Component {
  static propTypes = {
    allHostLabels: PropTypes.arrayOf(labelInterface),
    dispatch: PropTypes.func,
    hosts: PropTypes.arrayOf(hostInterface),
    hostPlatformLabels: PropTypes.arrayOf(labelInterface),
    hostStatusLabels: PropTypes.arrayOf(labelInterface),
    labels: PropTypes.arrayOf(labelInterface),
    selectedLabel: labelInterface,
    selectedOsqueryTable: osqueryTableInterface,
  };

  constructor (props) {
    super(props);

    this.state = {
      isAddLabel: false,
      labelQuery: '',
    };
  }

  componentWillMount () {
    const {
      allHostLabels,
      dispatch,
      hosts,
      labels,
      selectedLabel,
    } = this.props;

    dispatch(showRightSidePanel);

    if (!hosts.length) {
      dispatch(hostActions.loadAll());
    }

    if (!labels.length) {
      dispatch(labelActions.loadAll());
    }

    if (!selectedLabel) {
      dispatch(setSelectedLabel(allHostLabels[0]));
    }

    return false;
  }

  componentWillReceiveProps (nextProps) {
    const { allHostLabels, dispatch, selectedLabel } = nextProps;
    const allHostLabel = allHostLabels[0];

    if (!selectedLabel && !!allHostLabel) {
      dispatch(setSelectedLabel(allHostLabel));
    }

    return false;
  }

  componentWillUnmount () {
    const { dispatch } = this.props;

    dispatch(removeRightSidePanel);
  }

  onCancelAddLabel = (evt) => {
    evt.preventDefault();

    this.setState({ isAddLabel: false});

    return false;
  }

  onAddLabelClick = (evt) => {
    evt.preventDefault();

    this.setState({
      isAddLabel: true,
    });

    return false;
  }

  onHostDetailActionClick = (type) => {
    return (host) => {
      return (evt) => {
        evt.preventDefault();

        console.log(type, host);
        return false;
      };
    };
  }

  onLabelClick = (selectedLabel) => {
    return (evt) => {
      evt.preventDefault();

      const { dispatch } = this.props;

      dispatch(setSelectedLabel(selectedLabel));

      return false;
    };
  }

  onOsqueryTableSelect = (tableName) => {
    const { dispatch } = this.props;

    dispatch(selectOsqueryTable(tableName));

    return false;
  }

  onSaveAddLabel = (formData) => {
    console.log('Add label form submitted', formData);
    this.setState({ isAddLabel: false });

    return false;
  }

  onTextEditorInputChange = (labelQuery) => {
    this.setState({ labelQuery });

    return false;
  }

  renderHeader = () => {
    const { queryType, selectedLabel } = this.props;
    const { isAddLabel } = this.state;

    if (!selectedLabel || isAddLabel) {
      return false;
    }

    return (
      <div>
        <i className="kolidecon-label" />
        <span>{selectedLabel.title}</span>
        <span>{selectedLabel.hosts_count} Hosts Total</span>
      </div>
    );
  }

  renderBody = () => {
    let Body;
    const { hosts } = this.props;
    const { isAddLabel, labelQuery } = this.state;
    const {
      onCancelAddLabel,
      onHostDetailActionClick,
      onSaveAddLabel,
      onTextEditorInputChange
    } = this;

    if (isAddLabel) {
      return (
        <QueryComposer
          key="query-composer"
          onCancel={onCancelAddLabel}
          onSaveQueryFormSubmit={onSaveAddLabel}
          onTextEditorInputChange={onTextEditorInputChange}
          queryType="label"
          textEditorText={labelQuery}
        />
      );
    }

    return hosts.map((host) => {
      return (
        <HostDetails
          host={host}
          key={host.hostname}
          onDisableClick={onHostDetailActionClick('disable')}
          onQueryClick={onHostDetailActionClick('query')}
        />
      );
    });
  }

  renderSidePanel = () => {
    let Component;
    const { isAddLabel } = this.state;
    const {
      allHostLabels,
      hostPlatformLabels,
      hostStatusLabels,
      selectedLabel,
      selectedOsqueryTable,
    } = this.props;
    const { onAddLabelClick, onLabelClick, onOsqueryTableSelect } = this;

    if (isAddLabel) {
      Component = (
        <QuerySidePanel
          key="query-side-panel"
          onOsqueryTableSelect={onOsqueryTableSelect}
          selectedOsqueryTable={selectedOsqueryTable}
        />
      )
    } else {
      Component = (
        <HostSidePanel
          key="hosts-side-panel"
          allHostGroupItems={allHostLabels}
          hostPlatformGroupItems={hostPlatformLabels}
          hostStatusGroupItems={hostStatusLabels}
          onAddLabelClick={onAddLabelClick}
          onLabelClick={onLabelClick}
          selectedLabel={selectedLabel}
        />
      )
    }

    return (
      <ReactCSSTransitionGroup
        transitionName="hosts-page-side-panel"
        transitionEnterTimeout={500}
        transitionLeaveTimeout={0}
      >
        {Component}
      </ReactCSSTransitionGroup>
    );
  }

  render () {
    const { renderBody, renderHeader, renderSidePanel } = this;

    return (
      <div className="manage-hosts">
        {renderHeader()}
        {renderBody()}
        {renderSidePanel()}
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { entities: hosts } = entityGetter(state).get('hosts');
  const { entities: labels } = entityGetter(state).get('labels');
  const allHostLabels = filter(labels, { type: 'all' });
  const hostStatusLabels = filter(labels, { type: 'status' });
  const hostPlatformLabels = filter(labels, { type: 'platform' });
  const { selectedLabel } = state.components.ManageHostsPage;
  const { selectedOsqueryTable } = state.components.QueryPages;

  return {
    allHostLabels,
    hosts,
    labels,
    hostStatusLabels,
    hostPlatformLabels,
    selectedLabel,
    selectedOsqueryTable,
  };
};

export default connect(mapStateToProps)(ManageHostsPage);
