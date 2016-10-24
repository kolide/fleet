import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import 'brace/ext/linking';
import radium from 'radium';
import Select from 'react-select';
import 'react-select/dist/react-select.css';

import './mode';
import './theme';
import componentStyles from './styles';
import debounce from '../../../utilities/debounce';
import SaveQueryForm from '../../forms/queries/SaveQueryForm';
import SaveQuerySection from './SaveQuerySection';
import targetInterface from '../../../interfaces/target';
import ThemeDropdown from './ThemeDropdown';
import { validateQuery } from './helpers';

class NewQuery extends Component {
  static propTypes = {
    isLoadingTargets: PropTypes.bool,
    onInvalidQuerySubmit: PropTypes.func,
    onNewQueryFormSubmit: PropTypes.func,
    onOsqueryTableSelect: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    targets: PropTypes.arrayOf(targetInterface),
    textEditorText: PropTypes.string,
  };

  constructor (props) {
    super(props);

    this.state = {
      saveQuery: false,
      selectedTargets: [],
      theme: 'kolide',
    };
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

  onSaveQueryFormSubmit = debounce((formData) => {
    const {
      onInvalidQuerySubmit,
      onNewQueryFormSubmit,
      textEditorText,
    } = this.props;
    const { selectedTargets } = this.state;

    const { error } = validateQuery(textEditorText);

    if (error) {
      return onInvalidQuerySubmit(error);
    }

    return onNewQueryFormSubmit({
      ...formData,
      query: textEditorText,
      selectedTargets,
    });
  })

  onTargetSelect = (selectedTargets) => {
    this.setState({ selectedTargets });
    return false;
  }

  onThemeSelect = (evt) => {
    evt.preventDefault();

    this.setState({
      theme: evt.target.value,
    });

    return false;
  }

  onToggleSaveQuery = () => {
    const { saveQuery } = this.state;

    this.setState({
      saveQuery: !saveQuery,
    });

    return false;
  }

  render () {
    const {
      containerStyles,
      selectTargetsHeaderStyles,
      titleStyles,
    } = componentStyles;
    const { isLoadingTargets, onTextEditorInputChange, targets, textEditorText } = this.props;
    const { saveQuery, selectedTargets, theme } = this.state;
    const {
      onLoad,
      onSaveQueryFormSubmit,
      onTargetSelect,
      onThemeSelect,
      onToggleSaveQuery,
    } = this;

    return (
      <div style={containerStyles}>
        <p style={titleStyles}>
          New Query Page
        </p>
        <ThemeDropdown onSelectChange={onThemeSelect} theme={theme} />
        <div style={{ marginTop: '20px' }}>
          <AceEditor
            enableBasicAutocompletion
            enableLiveAutocompletion
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={4}
            maxLines={4}
            name="query-editor"
            onLoad={onLoad}
            onChange={onTextEditorInputChange}
            setOptions={{ enableLinking: true }}
            showGutter
            showPrintMargin={false}
            theme={theme}
            value={textEditorText}
            width="100%"
          />
        </div>
        <div>
          <p style={selectTargetsHeaderStyles}>Select Targets</p>
          <Select
            className="target-select"
            isLoading={isLoadingTargets}
            multi
            name="targets"
            options={targets}
            onChange={onTargetSelect}
            placeholder="Type to search"
            resetValue={[]}
            value={selectedTargets}
            valueKey="name"
          />
        </div>
        <SaveQuerySection onToggleSaveQuery={onToggleSaveQuery} saveQuery={saveQuery} />
        <SaveQueryForm onSubmit={onSaveQueryFormSubmit} saveQuery={saveQuery} />
      </div>
    );
  }
}

export default radium(NewQuery);
