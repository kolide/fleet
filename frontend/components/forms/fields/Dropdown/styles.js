import Styles from '../../../../styles';

const { border, color, font, padding } = Styles;

const DROPDOWN_BORDER = '1px solid rgba(176,183,210,0.5)';

export default {
  chevronStyles: {
    lineHeight: '38px',
  },
  chevronWrapperStyles: {
    backgroundColor: color.brand,
    color: color.white,
    cursor: 'pointer',
    float: 'right',
    height: '38px',
    textAlign: 'center',
    width: '32px',
  },
  optionWrapperStyles: {
    color: color.textMedium,
    cursor: 'pointer',
    padding: padding.half,
    ':hover': {
      background: '#F9F0FF',
      color: color.textUltradark,
    },
  },
  optionsWrapperStyles: (expanded) => {
    return {
      backgroundColor: color.white,
      borderBottom: DROPDOWN_BORDER,
      borderLeft: DROPDOWN_BORDER,
      borderRight: DROPDOWN_BORDER,
      borderTop: 'none',
      boxShadow: border.shadow.slight,
      display: expanded ? 'block' : 'none',
      fontSize: font.small,
      position: 'absolute',
      width: '250px',
    };
  },
  selectedOptionStyles: {
    cursor: 'pointer',
    borderColor: color.brand,
    borderRadius: border.radius.base,
    borderStyle: 'solid',
    borderWidth: '1px',
    fontSize: font.mini,
    height: '38px',
    paddingLeft: padding.half,
    textTransform: 'uppercase',
  },
  selectedTextStyles: {
    fontSize: font.mini,
    lineHeight: '38px',
  },
};
