import React, { Component, PropTypes } from 'react';
import AceEditor from 'react-ace';
import 'brace/mode/sql';
import 'brace/theme/dreamweaver';
import 'brace/theme/cobalt';
import 'brace/theme/eclipse';
import 'brace/theme/github';
import 'brace/theme/idle_fingers';
import 'brace/theme/iplastic';
import 'brace/theme/katzenmilch';
import 'brace/theme/kr_theme';
import 'brace/theme/kuroir';
import 'brace/theme/merbivore';
import 'brace/theme/merbivore_soft';
import 'brace/theme/mono_industrial';
import 'brace/theme/monokai';
import 'brace/theme/solarized_light';
import 'brace/theme/sqlserver';
import 'brace/theme/tomorrow';
import 'brace/ext/linking';
import radium from 'radium';
import './mode';
import './theme';
import componentStyles from './styles';
import Slider from '../../buttons/Slider';
import GradientButton from '../../buttons/GradientButton';

class NewQuery extends Component {
  static propTypes = {
    onOsqueryTableSelect: PropTypes.func,
    onTextEditorInputChange: PropTypes.func,
    textEditorText: PropTypes.string,
  };

  constructor (props) {
    super(props);

    this.state = {
      saveQuery: false,
      theme: 'kolide',
    };
  }

  componentWillMount () {
    global.window.addEventListener('keydown', this.handleKeydown);
  }

  componentWillUnmount () {
    global.window.removeEventListener('keydown', this.handleKeydown);
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

  onSelectChange = (evt) => {
    this.setState({
      theme: evt.target.value,
    });

    return false;
  }

  onRunQuery = () => {
    console.log('query run');
    return false;
  }

  onToggleSaveQuery = () => {
    const { saveQuery } = this.state;

    this.setState({
      saveQuery: !saveQuery,
    });

    return false;
  }

  handleKeydown = (evt) => {
    const { metaKey, code } = evt;

    if (metaKey && code === 'Enter') {
      return this.onRunQuery();
    }

    return false;
  };

  renderTextEditorThemeDropdown = () => {
    const { themeDropdownStyles } = componentStyles;
    const { theme } = this.state;
    const { onSelectChange } = this;

    return (
      <div style={themeDropdownStyles}>
        <span style={{ fontSize: '10px' }}>Editor Theme:</span>
        <select onChange={onSelectChange} style={themeDropdownStyles} value={theme}>
          <option value="kolide">Kolide</option>
          <option value="dreamweaver">Dreamweaver</option>
          <option value="cobalt">Cobalt</option>
          <option value="eclipse">Eclipse</option>
          <option value="github">Github</option>
          <option value="idle_fingers">Idle Fingers</option>
          <option value="iplastic">Iplastic</option>
          <option value="katzenmilch">Katzenmilch</option>
          <option value="kr_theme">KR Theme</option>
          <option value="kuroir">Kuroir</option>
          <option value="merbivore">Merbivore</option>
          <option value="merbivore_soft">Merbivore Soft</option>
          <option value="mono_industrial">Mono Industrial</option>
          <option value="monokai">Monokai</option>
          <option value="solarized_light">Solarized Light</option>
          <option value="sqlserver">SQL Server</option>
          <option value="tomorrow">Tomorrow</option>
        </select>
      </div>
    );
  }

  render () {
    const {
      containerStyles,
      runQueryButtonStyles,
      runQuerySectionStyles,
      runQueryTipStyles,
      saveResultsWrapper,
      saveQuerySection,
      saveWrapper,
      selectTargetsHeaderStyles,
      sliderText,
      targetsInputStyle,
      titleStyles,
    } = componentStyles;
    const { onTextEditorInputChange, textEditorText } = this.props;
    const { saveQuery, theme } = this.state;
    const {
      onBeforeLoad,
      onLoad,
      onRunQuery,
      onToggleSaveQuery,
      renderTextEditorThemeDropdown,
    } = this;

    return (
      <div style={containerStyles}>
        <p style={titleStyles}>
          New Query Page
        </p>
        {renderTextEditorThemeDropdown()}
        <div style={{ marginTop: '20px' }}>
          <AceEditor
            enableBasicAutocompletion
            enableLiveAutocompletion
            editorProps={{ $blockScrolling: Infinity }}
            mode="kolide"
            minLines={4}
            maxLines={4}
            name="query-editor"
            onBeforeLoad={onBeforeLoad}
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
          <input type="text" style={targetsInputStyle} />
        </div>
        <section style={saveQuerySection}>
          <div style={saveResultsWrapper}>
            <p>Save Query & Results For Later?</p>
            <small>For certain types of queries, like one that targets many hosts or one you plan to reuse frequently, we suggest saving the query & results. This allows you to set some advanced options, view the results later, and share with other users</small>
          </div>
          <div style={saveWrapper}>
            <span style={sliderText(!saveQuery)}>Dont save</span>
            <Slider onClick={onToggleSaveQuery} engaged={saveQuery} />
            <span style={sliderText(saveQuery)}>Save</span>
          </div>
        </section>
        <section style={runQuerySectionStyles}>
          <span style={runQueryTipStyles}>&#8984; + Enter</span>
          <GradientButton
            onClick={onRunQuery}
            style={runQueryButtonStyles}
            text="Run Query"
          />
        </section>
      </div>
    );
  }
}

export default radium(NewQuery);
