import Styles from '../../../styles';

const { color, font, padding } = Styles;

export default {
  containerStyles: {
    backgroundColor: color.white,
    padding: padding.base,
  },
  runQueryButtonStyles: {
    backgroundImage: 'none',
    backgroundColor: color.brandDark,
    boxShadow: '0 3px 0 #C38DEC',
    fontSize: font.medium,
    letterSpacing: '1px',
    paddingTop: padding.xSmall,
    paddingBottom: padding.xSmall,
    width: '200px',
  },
  runQuerySectionStyles: {
    borderTopColor: color.accentLight,
    borderTopStyle: 'solid',
    borderTopWidth: '1px',
    marginTop: padding.base,
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

