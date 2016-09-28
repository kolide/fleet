import Styles from '../../../styles';

const { color, font, padding } = Styles;

export default {
  containerStyles: {
    backgroundColor: color.white,
    padding: padding.base,
  },
  runQueryButtonStyles: {
    backgroundImage: 'none',
    backgroundColor: color.brand,
    fontSize: font.medium,
    paddingTop: padding.half,
    paddingBottom: padding.half,
    width: '200px',
  },
  runQuerySectionStyles: {
    paddingBottom: padding.base,
    paddingTop: padding.base,
    textAlign: 'right',
  },
  runQueryTipStyles: {
    color: color.textLight,
    fontSize: font.small,
    marginRight: padding.half,
  },
  saveResultsWrapper: {
    display: 'inline-block',
    width: '440px',
    '@media (max-width: 911px)': {
      width: '400px',
    },
  },
  saveQuerySection: {
    alignItems: 'flex-end',
    borderBottom: '1px solid #eaeaea',
    display: 'flex',
    justifyContent: 'space-between',
    paddingBottom: padding.base,
  },
  saveWrapper: {
    alignItems: 'center',
    display: 'flex',
  },
  selectTargetsHeaderStyles: {
    fontSize: font.base,
    color: color.textMedium,
  },
  sliderText: (saveQuery) => {
    return {
      color: saveQuery ? color.textDark : color.textMedium,
      textTransform: 'uppercase',
      fontSize: font.small,
      marginLeft: '5px',
      marginRight: '5px',
    };
  },
  targetsInputStyle: {
    width: '100%',
  },
  themeDropdownStyles: {
    display: 'inline-block',
    marginLeft: padding.half,
  },
  titleStyles: {
    color: color.textMedium,
    display: 'inline-block',
    fontSize: font.large,
  },
};

