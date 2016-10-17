export const iconClassForLabel = (label) => {
  const lowerType = label.type && label.type.toLowerCase();
  const lowerTitle = label.title && label.title.toLowerCase();

  if (lowerType === 'all') return 'kolidecon-hosts';

  switch (lowerTitle) {
    case 'offline': return 'kolidecon-hosts';
    case 'online': return 'kolidecon-hosts';
    case 'mac os': return 'kolidecon-apple';
    case 'centos': return 'kolidecon-centos';
    case 'ubuntu': return 'kolidecon-ubuntu';
    case 'windows': return 'kolidecon-windows';
    default: return 'kolidecon-tag';
  }
};

export default { iconClassForLabel };
