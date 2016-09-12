import styles from '../../../styles';

const { border, color, font, padding } = styles;

export default (disabled) => {
  const bgColor = {
    start: disabled ? '#B2B2B2' : '#7166D9',
    end: disabled ? '#C7B7C9' : '#C86DD7',
  };

  return {
    backgroundImage: `linear-gradient(to bottom right, ${bgColor.start}, ${bgColor.end})`,
    border: 'none',
    borderBottomLeftRadius: border.radius.base,
    borderBottomRightRadius: border.radius.base,
    boxSizing: 'border-box',
    color: color.white,
    cursor: disabled ? 'not-allowed' : 'pointer',
    fontSize: font.large,
    letterSpacing: '4px',
    padding: padding.base,
    textTransform: 'uppercase',
    width: '100%',
    ':focus': {
      outline: 'none',
    },
  };
};
