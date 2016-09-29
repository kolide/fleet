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
  selectTargetsHeaderStyles: {
    fontSize: font.base,
    color: color.textMedium,
  },
  targetsInputStyle: {
    width: '100%',
  },
  titleStyles: {
    color: color.textMedium,
    display: 'inline-block',
    fontSize: font.large,
  },
};

