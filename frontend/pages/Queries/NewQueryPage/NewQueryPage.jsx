import React, { Component } from 'react';
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
import './mode';
import NewQuery from '../../../components/Queries/NewQuery';

class NewQueryPage extends Component {
  constructor (props) {
    super(props);

    this.state = {
      osqueryTable: 'users',
      textEditorText: 'SELECT * FROM users u JOIN groups g WHERE u.gid = g.gid',
      theme: 'tomorrow',
    };
  }

  onLoad = (editor) => {
    editor.setOptions({
      enableLinking: true,
    });

    editor.on('linkClick', (data) => {
      const { type, value } = data.token;

      if (type === 'osquery-token') {
        this.setState({
          osqueryTable: value,
        });
      }
    });
  }

  onSelectChange = (evt) => {
    this.setState({
      theme: evt.target.value,
    });

    return false;
  }

  onTextEditorInputChange = (textEditorText) => {
    this.setState({ textEditorText });

    return false;
  }

  renderTextEditorThemeDropdown = () => {
    const { theme } = this.state;
    const { onSelectChange } = this;

    return (
      <select onChange={onSelectChange} value={theme}>
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
    );
  }

  render () {
    const { osqueryTable, textEditorText, theme } = this.state;
    const {
      onBeforeLoad,
      onLoad,
      onTextEditorInputChange,
      renderTextEditorThemeDropdown
    } = this;

    return (
      <div>
        <h1>New Query Page</h1>
        <h1>OsqueryTable: {osqueryTable}</h1>
        {renderTextEditorThemeDropdown()}
        <div style={{ marginTop: '20px' }}>
          <AceEditor
            enableBasicAutocompletion
            enableLiveAutocompletion
            editorProps={{$blockScrolling: Infinity}}
            height="200px"
            mode="kolide"
            name="query-editor"
            onBeforeLoad={this.onBeforeLoad}
            onLoad={this.onLoad}
            onChange={onTextEditorInputChange}
            setOptions={{enableLinking: true}}
            showGutter
            showPrintMargin={false}
            theme={theme}
            value={textEditorText}
            width="100%"
          />
        </div>
      </div>
    );
  }
}

export default NewQueryPage;
