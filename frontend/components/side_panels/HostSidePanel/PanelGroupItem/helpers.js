export const iconClassForLabel = (label) => {
  const lowerType = label.type && label.type.toLowerCase();
  const lowerDisplayText = label.display_text && label.display_text.toLowerCase();

  if (lowerType === 'all') return 'kolidecon-hosts';

  switch (lowerDisplayText) {
    case 'offline': return 'kolidecon-success-check';
    case 'online': return 'kolidecon-offline';
    case 'mac os': return 'kolidecon-apple';
    case 'centos': return 'kolidecon-centos';
    case 'ubuntu': return 'kolidecon-ubuntu';
    case 'windows': return 'kolidecon-windows';
    default: return 'kolidecon-label';
  }
};

export default { iconClassForLabel };
