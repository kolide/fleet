import Styles from '../../styles';
import { NAV_STYLES } from './SidePanel';

const { border, color, font, padding } = Styles;

const componentStyles = (navStyle) => {
  const { FULL, SKINNY } = NAV_STYLES;

  return {
    companyLogoStyles: () => {
      return {
        position: 'absolute',
        left: navStyle === FULL ? '16px' : '4px',
        height: '44px',
        marginRight: '10px',
      };
    },
    headerStyles: {
      borderBottomColor: color.accentLight,
      borderBottomStyle: 'solid',
      borderBottomWidth: '1px',
      height: '67px',
      marginBottom: padding.half,
      marginRight: padding.medium,
      paddingLeft: '54px',
    },
    iconStyles: () => {
      return {
        position: 'relative',
        fontSize: '22px',
        marginRight: '16px',
        top: '4px',
        left: navStyle === FULL ? 0 : '5px',
      };
    },
    navItemBeforeStyles: () => {
      return {
        content: '',
        width: '6px',
        height: '50px',
        position: 'absolute',
        left: navStyle === FULL ? '-24px' : 0,
        top: '2px',
        bottom: 0,
        backgroundColor: '#9a61c6',
      };
    },
    navItemListStyles: {
      listStyle: 'none',
      margin: 0,
      padding: 0,
    },
    navItemNameStyles: () => {
      if (navStyle === SKINNY) {
        return {
          display: 'none',
        };
      }

      return {
        display: 'inline-block',
        textDecoration: 'none',
      };
    },
    navItemStyles: (active) => {
      const activeStyles = {
        color: color.brand,
        borderBottom: navStyle === FULL ? 'none' : '8px solid #9a61c6',
        transition: 'none',
      };

      const baseStyles = {
        minHeight: '40px',
        position: 'relative',
        color: color.textLight,
        cursor: 'pointer',
        fontSize: font.small,
        textTransform: 'uppercase',
        paddingTop: padding.half,
        WebkitTransition: 'all 0.2s ease-in-out',
        MozTransition: 'all 0.2s ease-in-out',
        textAlign: navStyle === FULL ? 'inherit' : 'center',
        transition: 'all 0.2s ease-in-out',
      };

      if (active) {
        return {
          ...baseStyles,
          ...activeStyles,
        };
      }

      return baseStyles;
    },
    navItemWrapperStyles: (lastChild) => {
      const baseStyles = {
        position: 'relative',
      };
      const lastChildStyles = {
        borderTopColor: color.accentLight,
        borderTopStyle: 'solid',
        borderTopWidth: '1px',
        marginRight: navStyle === FULL ? '16px' : 0,
        marginTop: '5px',
      };

      if (lastChild) {
        return {
          ...baseStyles,
          ...lastChildStyles,
        };
      }

      return baseStyles;
    },
    navStyles: () => {
      return {
        backgroundColor: color.white,
        borderRightColor: color.borderMedium,
        borderRightStyle: 'solid',
        borderRightWidth: '1px',
        bottom: 0,
        boxShadow: '2px 0 8px 0 rgba(0, 0, 0, 0.1)',
        left: 0,
        paddingLeft: navStyle === FULL ? padding.large : 0,
        paddingTop: padding.large,
        position: 'fixed',
        top: 0,
        width: navStyle === FULL ? '216px' : '54px',
      };
    },
    orgNameStyles: () => {
      if (navStyle === SKINNY) {
        return { display: 'none' };
      }

      return {
        fontSize: font.medium,
        letterSpacing: '0.5px',
        margin: 0,
        overFlow: 'hidden',
        padding: 0,
        position: 'relative',
        textOverflow: 'ellipsis',
        top: '3px',
        whiteSpace: 'nowrap',
      };
    },
    subItemBeforeStyles: {
      backgroundColor: color.white,
      borderRadius: border.radius.circle,
      content: '',
      display: 'block',
      height: '7px',
      left: '24px',
      position: 'absolute',
      top: '15px',
      width: '7px',
    },
    subItemLinkStyles: (active) => {
      const activeStyles = {
        textDecoration: 'none',
        textTransform: 'none',
      };

      return active ? activeStyles : {};
    },
    subItemStyles: (active) => {
      const activeStyles = {
        fontWeight: font.weight.bold,
        opacity: '1',
      };

      const baseStyles = {
        color: color.white,
        marginBottom: '5px',
        marginLeft: 0,
        marginRight: 0,
        marginTop: '5px',
        opacity: '0.5',
        paddingBottom: padding.xSmall,
        paddingLeft: padding.most,
        paddingTop: padding.xSmall,
        position: 'relative',
        WebkitTransition: 'all 0.2s ease-in-out',
        MozTransition: 'all 0.2s ease-in-out',
        textTransform: 'none',
        transition: 'all 0.2s ease-in-out',
      };

      if (active) {
        return {
          ...baseStyles,
          ...activeStyles,
        };
      }

      return baseStyles;
    },
    subItemsStyles: (expanded) => {
      const baseStyles = {
        backgroundColor: color.brand,
        boxShadow: 'inset 0 5px 8px 0 rgba(0, 0, 0, 0.12), inset 0 -5px 8px 0 rgba(0, 0, 0, 0.12)',
        marginBottom: 0,
        marginRight: 0,
        minHeight: '87px',
        paddingBottom: padding.half,
        paddingTop: padding.half,
      };

      const fullNavStyles = {
        marginLeft: '-24px',
        marginTop: padding.medium,
      };

      const skinnyNavStyles = {
        bottom: '-8px',
        left: '54px',
        position: 'absolute',
        width: expanded ? '251px' : '18px',
      };

      if (navStyle === FULL) {
        return {
          ...baseStyles,
          ...fullNavStyles,
        };
      }

      return {
        ...baseStyles,
        ...skinnyNavStyles,
      };
    },
    subItemListStyles: (expanded) => {
      const skinnyNavStyles = {
        borderRight: '1px solid #eaeaea',
        display: 'inline-block',
        padding: 0,
        textAlign: 'left',
        width: '84%',
      };

      const baseStyles = {
        listStyle: 'none',
      };

      if (navStyle === SKINNY) {
        if (!expanded) return { display: 'none' };

        return {
          ...baseStyles,
          ...skinnyNavStyles,
        };
      }

      return baseStyles;
    },
    collapseSubItemsWrapper: {
      position: 'absolute',
      right: '3px',
      top: '41%',
    },
    usernameStyles: () => {
      if (navStyle === SKINNY) {
        return { display: 'none' };
      }

      return {
        position: 'relative',
        top: '3px',
        display: 'inline-block',
        margin: 0,
        padding: 0,
        fontSize: font.small,
        textTransform: 'uppercase',
      };
    },
    userStatusStyles: (enabled) => {
      if (navStyle === SKINNY) {
        return { display: 'none' };
      }

      const backgroundColor = enabled ? color.success : color.warning;
      const size = '16px';

      return {
        backgroundColor,
        borderRadius: border.radius.circle,
        display: 'inline-block',
        height: size,
        left: '1px',
        marginRight: '6px',
        position: 'relative',
        top: '6px',
        width: size,
      };
    },
  };
};

export default componentStyles;
