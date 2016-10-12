import Styles from '../../../styles';

const { color, font, padding } = Styles;

export default {
  containerStyles: (status) => {
    const baseStyles = {
      backgroundColor: color.white,
      borderBottom: `solid 1px ${color.silver}`,
      borderLeft: `solid 1px ${color.silver}`,
      borderRight: `solid 1px ${color.silver}`,
      borderTop: `solid 1px ${color.silver}`,
      borderRadius: '3px',
      boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.05)',
      boxSizing: 'border-box',
      display: 'inline-block',
      height: '286px',
      marginLeft: padding.base,
      marginTop: padding.base,
      paddingBottom: padding.half,
      paddingLeft: padding.half,
      paddingRight: padding.half,
      paddingTop: padding.half,
      textAlign: 'center',
      width: '240px',
    };
    const statusStyles = {
      ONLINE: {
        borderTop: `6px solid ${color.success}`,
      },
      OFFLINE: {
        borderTop: `6px solid ${color.alert}`,
      },
      NEEDS_UPGRADE: {
        borderTop: `6px solid ${color.alert}`,
      },
    };

    return {
      ...baseStyles,
      ...statusStyles[status],
    };
  },
  contentSeparatorStyles: {
    borderTop: `1px solid ${color.accentLight}`,
    marginTop: padding.half,
  },
  hostContentItemStyles: {
    color: color.textUltradark,
    fontSize: font.small,
    marginLeft: '3px',
    marginRight: '3px',
  },
  hostnameStyles: {
    color: color.link,
    fontSize: font.mini,
    fontWeight: 'bold',
    marginTop: 0,
    marginBottom: 0,
  },
  statusStyles: (status) => {
    const baseStyles = {
      fontSize: font.medium,
      textAlign: 'left',
    };
    const statusStyles = {
      ONLINE: {
        color: color.success,
      },
      OFFLINE: {
        color: color.alert,
      },
      NEEDS_UPGRADE: {
        color: color.warning,
      },
    };

    return {
      ...baseStyles,
      ...statusStyles[status],
    };
  },
};
