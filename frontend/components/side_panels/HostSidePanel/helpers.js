export const iconClassForLabel = (label) => {
  const lowerName = label.name && label.name.toLowerCase();

  switch (lowerName) {
    case 'all': return 'kolidecon-hosts';
    case 'offline': return 'kolidecon-hosts';
    case 'online': return 'kolidecon-hosts';
    case 'macs': return 'kolidecon-apple';
    case 'centos': return 'kolidecon-centos';
    case 'ubuntu': return 'kolidecon-ubuntu';
    case 'windows': return 'kolidecon-windows';
    default: return 'kolidecon-tag';
  }
};

export default { iconClassForLabel };
