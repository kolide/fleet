import React, { Component, PropTypes } from 'react';
import { isEqual, last } from 'lodash';
import componentStylesFunc from './styles';
import kolideLogo from '../../../assets/images/kolide-logo.svg';
import debounce from '../../utilities/debounce';
import navItems from './navItems';

const NAV_BREAKPOINT = 760;
export const NAV_STYLES = {
  FULL: 'full',
  SKINNY: 'skinny',
};

class SidePanel extends Component {
  static propTypes = {
    user: PropTypes.object,
  };

  constructor (props) {
    super(props);

    const { FULL, SKINNY } = NAV_STYLES;
    const { innerWidth } = global.window;
    const navStyle = innerWidth <= NAV_BREAKPOINT ? SKINNY : FULL;

    this.componentStyles = componentStylesFunc(navStyle);

    this.state = {
      activeTab: 'Hosts',
      activeSubItem: 'Add Hosts',
      navStyle,
      subItemsExpanded: false,
    };
  }

  componentDidMount () {
    global.window.addEventListener('resize', this.handleResize);
  }

  componentWillUnmount () {
    global.window.removeEventListener('resize', this.handleResize);
  }

  setActiveSubItem = (activeSubItem) => {
    return (evt) => {
      evt.preventDefault();

      this.setState({ activeSubItem });
      return false;
    };
  }

  setActiveTab = (activeTab) => {
    return (evt) => {
      evt.preventDefault();

      this.setState({
        activeTab,
      });

      return false;
    };
  }

  handleResize = debounce(() => {
    const { FULL, SKINNY } = NAV_STYLES;
    const { innerWidth } = global.window;
    const { navStyle } = this.state;

    if (innerWidth <= NAV_BREAKPOINT && navStyle !== SKINNY) {
      this.componentStyles = componentStylesFunc(SKINNY);
      this.setState({ navStyle: SKINNY });
    }

    if (innerWidth > NAV_BREAKPOINT && navStyle !== FULL) {
      this.componentStyles = componentStylesFunc(FULL);
      this.setState({ navStyle: FULL });
    }
  }, { leading: false, trailing: true, timeout: 300 })

  toggleSubItemsCollapse = (showSubItems) => {
    return (evt) => {
      evt.preventDefault();

      this.setState({
        subItemsExpanded: showSubItems,
      });

      return false;
    };
  }

  renderHeader = () => {
    const {
      user: {
        enabled,
        username,
      },
    } = this.props;
    const {
      companyLogoStyles,
      headerStyles,
      orgNameStyles,
      usernameStyles,
      userStatusStyles,
    } = this.componentStyles;

    return (
      <header style={headerStyles}>
        <img
          alt="Company logo"
          src={kolideLogo}
          style={companyLogoStyles()}
        />
        <h1 style={orgNameStyles()}>Kolide, Inc.</h1>
        <div style={userStatusStyles(enabled)} />
        <h2 style={usernameStyles()}>{username}</h2>
      </header>
    );
  }

  renderNavItem = (navItem, lastChild) => {
    const { activeTab } = this.state;
    const { icon, name, subItems } = navItem;
    const active = activeTab === name;
    const {
      iconStyles,
      navItemBeforeStyles,
      navItemNameStyles,
      navItemStyles,
      navItemWrapperStyles,
    } = this.componentStyles;
    const { renderSubItems, setActiveTab } = this;

    return (
      <div style={navItemWrapperStyles(lastChild)} key={`nav-item-${name}`}>
        {active && <div style={navItemBeforeStyles()} />}
        <li
          onClick={setActiveTab(name)}
          style={navItemStyles(active)}
        >
          <div style={{ position: 'relative' }}>
            <i className={icon} style={iconStyles()} />
            <span style={navItemNameStyles()}>
              {name}
            </span>
          </div>
          {active && renderSubItems(subItems)}
        </li>
      </div>
    );
  }

  renderNavItems = () => {
    const { renderNavItem } = this;
    const { navItemListStyles } = this.componentStyles;
    const { user: { admin } } = this.props;
    const userNavItems = navItems(admin);

    return (
      <ul style={navItemListStyles}>
        {userNavItems.map((navItem, index, collection) => {
          const lastChild = admin && isEqual(navItem, last(collection));
          return renderNavItem(navItem, lastChild);
        })}
      </ul>
    );
  }

  renderSubItem = (subItem) => {
    const { activeSubItem } = this.state;
    const { name, path } = subItem;
    const active = activeSubItem === name;
    const { setActiveSubItem } = this;
    const {
      subItemBeforeStyles,
      subItemStyles,
      subItemLinkStyles,
    } = this.componentStyles;

    return (
      <div
        key={`sub-item-${name}`}
        style={{ position: 'relative' }}
      >
        {active && <div style={subItemBeforeStyles} />}
        <li
          onClick={setActiveSubItem(name)}
          style={subItemStyles(active)}
        >
          <span to={path} style={subItemLinkStyles(active)}>{name}</span>
        </li>
      </div>
    );
  }

  renderSubItems = (subItems) => {
    const { subItemListStyles, subItemsStyles } = this.componentStyles;
    const { renderCollapseSubItems, renderSubItem } = this;
    const { subItemsExpanded } = this.state;

    if (!subItems.length) return false;

    return (
      <div style={subItemsStyles(subItemsExpanded)}>
        <ul style={subItemListStyles(subItemsExpanded)}>
          {subItems.map(subItem => {
            return renderSubItem(subItem);
          })}
        </ul>
        {renderCollapseSubItems()}
      </div>
    );
  }

  renderCollapseSubItems = () => {
    const { navStyle } = this.state;
    const { FULL } = NAV_STYLES;
    const { toggleSubItemsCollapse } = this;
    const { subItemsExpanded } = this.state;
    const { collapseSubItemsWrapper } = this.componentStyles;
    const iconName = subItemsExpanded ? 'kolidecon-chevron-bold-left' : 'kolidecon-chevron-bold-right';
    if (navStyle === FULL) return false;

    return (
      <div style={collapseSubItemsWrapper}>
        <i className={iconName} style={{ color: '#FFF' }} onClick={toggleSubItemsCollapse(!subItemsExpanded)} />
      </div>
    );
  }

  render () {
    const { navStyles } = this.componentStyles;
    const { renderHeader, renderNavItems } = this;

    return (
      <nav style={navStyles()}>
        {renderHeader()}
        {renderNavItems()}
      </nav>
    );
  }
}

export default SidePanel;
