import React, { Component } from 'react';
import { map } from 'lodash';
import componentStyles from './styles';

const HOST_TABS = {
  FIRST: 'What Does This Script Do?',
  SECOND: 'Additional Script Options',
};

class NewHostPage extends Component {
  constructor (props) {
    super(props);

    this.state = {
      method1Text: 'curl https://kolide.acme.com/install/osquery.sh | sudo sh',
      method2Text: 'osqueryd --conﬁg_endpoint="https://kolide.acme.com/api/v1/…',
      selectedTab: HOST_TABS.FIRST,
    };
  }

  onSetActiveTab = (selectedTab) => {
    return (evt) => {
      evt.preventDefault();

      this.setState({ selectedTab });

      return false;
    };
  }

  renderHostTabContent = () => {
    const { selectedTab } = this.state;

    if (selectedTab === HOST_TABS.FIRST) {
      return (
        <div>
          <p style={{ marginTop: 0 }}>This script does the following:</p>
          <ol className="kolide-ol">
            <li>Detects operating system.</li>
            <li>Checks for any existing osqueryd installation.</li>
            <li>Installs osqueryd and ships your config to communicate with Kolide.</li>
          </ol>
        </div>
      );
    }

    return false;
  }

  renderHostTabHeaders = () => {
    const { hostTabHeaderStyles } = componentStyles;
    const { selectedTab } = this.state;
    const { onSetActiveTab } = this;

    return map(HOST_TABS, tab => {
      const selected = selectedTab === tab;

      return <span onClick={onSetActiveTab(tab)} key={tab} style={hostTabHeaderStyles(selected)}>{tab}</span>;
    });
  }

  render () {
    const {
      headerStyles,
      inputStyles,
      textStyles,
      scriptInfoWrapperStyles,
      selectedTabContentStyles,
      sectionWrapperStyles,
    } = componentStyles;
    const { method1Text, method2Text } = this.state;
    const { renderHostTabContent, renderHostTabHeaders } = this;

    return (
      <div>
        <div style={sectionWrapperStyles}>
          <p style={headerStyles}>Method 1 - One Liner</p>
          <input style={inputStyles} value={method1Text} readOnly />
          <div style={scriptInfoWrapperStyles}>
            {renderHostTabHeaders()}
            <div style={selectedTabContentStyles}>
              {renderHostTabContent()}
            </div>
          </div>
        </div>
        <div style={sectionWrapperStyles}>
          <p style={[headerStyles, { width: '626px' }]}>Method 2 - Your osqueryd with Kolid config</p>
          <input style={inputStyles} value={method2Text} readOnly />
          <p style={textStyles}>This method allows you to configure an existing osqueryd installation to work with Kolide. The <span style={{ color: '#AE6DDf', fontFamily: 'SourceCodePro, Oxygen' }}>--config_endpoints</span> flag allows us to point your osqueryd installation to your Kolide configuration.</p>
        </div>
        <div style={sectionWrapperStyles}>
          <p style={headerStyles}>Method 3 - Need More Methods?</p>
          <p style={textStyles}>Many IT automation frameworks offer direct recipes and scripts for deploying osquery. Choose a method below to learn more.</p>
        </div>
      </div>
    );
  }
}

export default NewHostPage;
