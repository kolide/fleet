import styles from '../../../../styles';

const { color, font, padding } = styles;

export default {
  containerStyles: {
    marginTop: padding.base,
    position: 'relative',
    width: '100%',
  },
  inputStyles: (value, type) => {
    const baseStyles = {
      borderLeft: 'none',
      borderRight: 'none',
      borderTop: 'none',
      borderBottomWidth: '1px',
      fontSize: font.small,
      borderBottomStyle: 'solid',
      borderBottomColor: color.brand,
      color: color.accentText,
      paddingRight: '30px',
      opacity: '1',
      textIndent: '2px',
      position: 'relative',
      width: '100%',
      boxSizing: 'border-box',
      ':focus': {
        outline: 'none',
      },
    };

    if (type === 'password' && value) {
      return {
        ...baseStyles,
        letterSpacing: '7px',
        color: color.textUltradark,
      };
    }

    if (value) {
      return {
        ...baseStyles,
        color: color.textUltradark,
      };
    }

    return baseStyles;
  },
  labelStyles: {
    color: color.textLight,
    textTransform: 'uppercase',
    fontSize: font.mini,
  },
};
