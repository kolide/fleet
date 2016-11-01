import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';

import Kolide from '../../../kolide';
import NewQuery from '../../../components/queries/NewQuery';
import osqueryTableInterface from '../../../interfaces/osquery_table';
import queryActions from '../../../redux/nodes/entities/queries/actions';
import QuerySidePanel from '../../../components/side_panels/QuerySidePanel';
import { removeRightSidePanel, selectOsqueryTable, showRightSidePanel } from '../../../redux/nodes/app/actions';
import { renderFlash } from '../../../redux/nodes/notifications/actions';

class NewQueryPage extends Component {
  static propTypes = {
    dispatch: PropTypes.func,
    selectedOsqueryTable: osqueryTableInterface,
  };

  componentWillMount () {
    const { dispatch } = this.props;

    this.state = {
      isLoadingTargets: false,
      moreInfoTarget: null,
      selectedTargetsCount: 0,
      targets: [],
      textEditorText: 'SELECT * FROM users u JOIN groups g WHERE u.gid = g.gid',
    };

    dispatch(showRightSidePanel);
    this.fetchTargets();

    return false;
  }

  componentWillUnmount () {
    const { dispatch } = this.props;

    dispatch(removeRightSidePanel);

    return false;
  }

  onNewQueryFormSubmit = (formData) => {
    console.log('New Query Form submitted', formData);
    const { dispatch } = this.props;

    dispatch(queryActions.create(formData));
  }

  onInvalidQuerySubmit = (errorMessage) => {
    const { dispatch } = this.props;

    dispatch(renderFlash('error', errorMessage));

    return false;
  }

  onOsqueryTableSelect = (tableName) => {
    const { dispatch } = this.props;

    dispatch(selectOsqueryTable(tableName));

    return false;
  }

  onRemoveMoreInfoTarget = (evt) => {
    evt.preventDefault();

    this.setState({ moreInfoTarget: null });

    return false;
  }

  onTargetSelectMoreInfo = (moreInfoTarget) => {
    return (evt) => {
      evt.preventDefault();

      const { target_type: targetType } = moreInfoTarget;

      if (targetType.toLowerCase() === 'labels') {
        return Kolide.getLabelHosts(moreInfoTarget.id)
          .then((hosts) => {
            console.log('hosts', hosts);
            this.setState({
              moreInfoTarget: {
                ...moreInfoTarget,
                hosts,
              },
            });

            return false;
          });
      }


      this.setState({ moreInfoTarget });

      return false;
    };
  }

  onTextEditorInputChange = (textEditorText) => {
    this.setState({ textEditorText });

    return false;
  }

  fetchTargets = (search) => {
    this.setState({ isLoadingTargets: true });

    return Kolide.getTargets({ search })
      .then((response) => {
        const {
          selected_targets_count: selectedTargetsCount,
          targets,
        } = response;

        this.setState({
          isLoadingTargets: false,
          selectedTargetsCount,
          targets,
        });

        return search;
      })
      .catch((error) => {
        this.setState({ isLoadingTargets: false });

        throw error;
      });
  }

  render () {
    const {
      fetchTargets,
      onNewQueryFormSubmit,
      onInvalidQuerySubmit,
      onOsqueryTableSelect,
      onRemoveMoreInfoTarget,
      onTargetSelectMoreInfo,
      onTextEditorInputChange,
    } = this;
    const {
      isLoadingTargets,
      moreInfoTarget,
      selectedTargetsCount,
      targets,
      textEditorText,
    } = this.state;
    const { selectedOsqueryTable } = this.props;

    return (
      <div>
        <NewQuery
          isLoadingTargets={isLoadingTargets}
          moreInfoTarget={moreInfoTarget}
          onNewQueryFormSubmit={onNewQueryFormSubmit}
          onInvalidQuerySubmit={onInvalidQuerySubmit}
          onOsqueryTableSelect={onOsqueryTableSelect}
          onRemoveMoreInfoTarget={onRemoveMoreInfoTarget}
          onTargetSelectInputChange={fetchTargets}
          onTargetSelectMoreInfo={onTargetSelectMoreInfo}
          onTextEditorInputChange={onTextEditorInputChange}
          selectedTargetsCount={selectedTargetsCount}
          selectedOsqueryTable={selectedOsqueryTable}
          targets={targets}
          textEditorText={textEditorText}
        />
        <QuerySidePanel
          onOsqueryTableSelect={onOsqueryTableSelect}
          onTextEditorInputChange={onTextEditorInputChange}
          selectedOsqueryTable={selectedOsqueryTable}
        />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const { selectedOsqueryTable } = state.app;

  return { selectedOsqueryTable };
};

export default connect(mapStateToProps)(NewQueryPage);
