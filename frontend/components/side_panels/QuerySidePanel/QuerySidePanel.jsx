import React, { Component, PropTypes } from 'react';

import osqueryTableInterface from 'interfaces/osquery_table';
import { osqueryTableNames } from 'utilities/osquery_tables';
import iconClassForLabel from 'utilities/icon_class_for_label';
import Dropdown from 'components/forms/fields/Dropdown';
import Icon from 'components/Icon';
import SecondarySidePanelContainer from '../SecondarySidePanelContainer';
import Button from '../../buttons/Button';
import {
  availability,
  columnsToRender,
  displayTypeForDataType,
  numAdditionalColumns,
  shouldShowAllColumns,
} from './helpers';

const baseClass = 'query-side-panel';

class QuerySidePanel extends Component {
  static propTypes = {
    onOsqueryTableSelect: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    selectedOsqueryTable: osqueryTableInterface,
  };

  componentWillMount () {
    const { selectedOsqueryTable } = this.props;
    const showAllColumns = shouldShowAllColumns(selectedOsqueryTable);

    this.setState({ showAllColumns });
  }

  componentWillReceiveProps (nextProps) {
    const { selectedOsqueryTable } = nextProps;

    if (this.props.selectedOsqueryTable !== selectedOsqueryTable) {
      const showAllColumns = shouldShowAllColumns(selectedOsqueryTable);

      this.setState({ showAllColumns });
    }

    return false;
  }

  onSelectTable = ({ value }) => {
    const { onOsqueryTableSelect } = this.props;

    onOsqueryTableSelect(value);

    return false;
  }

  onShowAllColumns = () => {
    this.setState({ showAllColumns: true });
  }

  onSuggestedQueryClick = (query) => {
    return (evt) => {
      evt.preventDefault();

      const { onTextEditorInputChange } = this.props;

      return onTextEditorInputChange(query);
    };
  };

  renderColumns = () => {
    const { selectedOsqueryTable } = this.props;
    const { showAllColumns } = this.state;
    const columns = columnsToRender(selectedOsqueryTable, showAllColumns);
    const columnBaseClass = 'query-column-list';

    return columns.map((column) => {
      return (
        <li key={column.name} className={`${columnBaseClass}__item`}>
          <span className={`${columnBaseClass}__name`}>{column.name}</span>
          <div className={`${columnBaseClass}__description`}>
            <span className={`${columnBaseClass}__type`}>{displayTypeForDataType(column.type)}</span>
            <Icon name="help-solid" className={`${columnBaseClass}__help`} title={column.description} />
          </div>
        </li>
      );
    });
  }

  renderMoreColumns = () => {
    const { selectedOsqueryTable } = this.props;
    const { showAllColumns } = this.state;
    const { onShowAllColumns } = this;

    if (showAllColumns) {
      return false;
    }

    return (
      <div className={`${baseClass}__column-wrapper`}>
        <span className={`${baseClass}__more-columns`}>{numAdditionalColumns(selectedOsqueryTable)} MORE COLUMNS</span>
        <button className={`button--unstyled ${baseClass}__show-columns`} onClick={onShowAllColumns}>SHOW</button>
      </div>
    );
  }

  renderSuggestedQueries = () => {
    const { onSuggestedQueryClick } = this;
    const { selectedOsqueryTable } = this.props;

    return selectedOsqueryTable.examples.map((example) => {
      return (
        <div key={example} className={`${baseClass}__column-wrapper`}>
          <span className={`${baseClass}__suggestion`}>{example}</span>
          <Button
            onClick={onSuggestedQueryClick(example)}
            className={`${baseClass}__load-suggestion`}
            text="LOAD"
          />
        </div>
      );
    });
  }

  renderTableSelect = () => {
    const { onSelectTable } = this;
    const { selectedOsqueryTable } = this.props;

    const tableNames = osqueryTableNames.map((name) => {
      return { label: name, value: name };
    });

    return (
      <Dropdown
        options={tableNames}
        value={selectedOsqueryTable.name}
        onSelect={onSelectTable}
        placeholder="Choose Table..."
      />
    );
  }

  render () {
    const {
      renderColumns,
      renderMoreColumns,
      renderTableSelect,
      renderSuggestedQueries,
    } = this;
    const { selectedOsqueryTable: { description, platform } } = this.props;
    const platformArr = availability(platform);

    return (
      <SecondarySidePanelContainer className={baseClass}>
        <div className={`${baseClass}__choose-table`}>
          <h2 className={`${baseClass}__header`}>Choose a Table</h2>
          {renderTableSelect()}
          <p className={`${baseClass}__description`}>{description}</p>
        </div>

        <div className={`${baseClass}__os-availability`}>
          <h2 className={`${baseClass}__header`}>OS Availability</h2>
          <ul className={`${baseClass}__platforms`}>
            {platformArr.map((os, idx) => {
              return <li key={idx}><Icon name={iconClassForLabel(os)} /> {os.display_text}</li>;
            })}
          </ul>
        </div>

        <div className={`${baseClass}__columns`}>
          <h2 className={`${baseClass}__header`}>Columns</h2>
          <ul className={`${baseClass}__column-list`}>
            {renderColumns()}
          </ul>
          {renderMoreColumns()}
        </div>

        <div className={`${baseClass}__joins`}>
          <h2 className={`${baseClass}__header`}>Joins</h2>
        </div>

        <div className={`${baseClass}__suggested-queries`}>
          <h2 className={`${baseClass}__header`}>Suggested Queries</h2>
          {renderSuggestedQueries()}
        </div>
      </SecondarySidePanelContainer>
    );
  }
}

export default QuerySidePanel;
