import Styles from '../../../../styles';

const { color, font } = Styles;

export default {
  saveQuerySection: {
    alignItems: 'flex-end',
    borderBottom: '1px solid #eaeaea',
    display: 'flex',
    justifyContent: 'space-between',
    paddingBottom: padding.base,
  },
  saveTextWrapper: {
    display: 'inline-block',
    width: '440px',
    '@media (max-width: 911px)': {
      width: '400px',
    },
  },
  saveWrapper: {
    alignItems: 'center',
    display: 'flex',
  },
  sliderTextDontSave: (saveQuery) => {
    return {
      color: saveQuery ? color.textDark : color.textUltradark,
      textTransform: 'uppercase',
      fontSize: font.small,
      marginRight: '5px',
    };
  },
  sliderTextSave: (saveQuery) => {
    return {
      color: saveQuery ? color.brand : color.textMedium,
      textTransform: 'uppercase',
      fontSize: font.small,
      marginLeft: '5px',
    };
  },
};
