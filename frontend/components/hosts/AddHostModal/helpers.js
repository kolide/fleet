import select from 'select';

const removeSelectedText = () => {
  return global.window.getSelection().removeAllRanges();
};

export const copyText = (elementSelector) => {
  const { document } = global;

  const element = document.querySelector(elementSelector);
  element.type = 'text';

  select(element);

  const canCopy = document.queryCommandEnabled('copy');

  if (!canCopy) {
    return false;
  }

  document.execCommand('copy');
  element.type = 'password';
  removeSelectedText();
  return true;
};

export default { copyText };
