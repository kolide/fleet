import React, { Component, PropTypes } from 'react';

import Button from 'components/buttons/Button';
import QueryForm from 'components/forms/queries/QueryForm';
import queryInterface from 'interfaces/query';
import SelectTargetsDropdown from 'components/forms/fields/SelectTargetsDropdown';
import targetInterface from 'interfaces/target';

const baseClass = 'query-composer';

class QueryComposer extends Component {
  static propTypes = {
    onFetchTargets: PropTypes.func,
    onFormCancel: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    onRunQuery: PropTypes.func,
    onSave: PropTypes.func,
    onStopQuery: PropTypes.func,
    onTargetSelect: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    onUpdate: PropTypes.func,
    query: queryInterface,
    queryIsRunning: PropTypes.bool,
    queryType: PropTypes.string,
    selectedTargets: PropTypes.arrayOf(targetInterface),
    targetsCount: PropTypes.number,
    queryText: PropTypes.string,
  };

  static defaultProps = {
    queryType: 'query',
    targetsCount: 0,
  };

  renderForm = () => {
    const {
      onFormCancel,
      onRunQuery,
      onSave,
      onStopQuery,
      onTextEditorInputChange,
      onOsqueryTableSelect,
      onUpdate,
      query,
      queryIsRunning,
      queryText,
      queryType,
    } = this.props;

    return (
      <QueryForm
        onCancel={onFormCancel}
        onRunQuery={onRunQuery}
        onSave={onSave}
        onStopQuery={onStopQuery}
        onTextEditorInputChange={onTextEditorInputChange}
        onOsqueryTableSelect={onOsqueryTableSelect}
        onUpdate={onUpdate}
        query={query}
        queryIsRunning={queryIsRunning}
        queryType={queryType}
        queryText={queryText}
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
      <div className={`${baseClass}__target-select`}>
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

  renderRunQueryButton = () => {
    const {
      onRunQuery,
      onStopQuery,
      queryIsRunning,
    } = this.props;
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

    return runQueryButton;
  };

  render () {
    const { queryType } = this.props;
    const { renderForm, renderTargetsInput, renderRunQueryButton } = this;

    return (
      <div className={`${baseClass}__wrapper`}>
        <div className={`${baseClass}__query body-wrap`}>
          <h1>{queryType === 'label' ? 'New Label Query' : 'New Query'}</h1>
          {renderForm()}
        </div>

        <div className={`${baseClass}__targets body-wrap`}>
          <h2>No Results Until Query Run</h2>
          {renderRunQueryButton()}
          {renderTargetsInput()}
        </div>
      </div>
    );
  }
}

export default QueryComposer;
