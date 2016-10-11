import md5 from 'js-md5';

export const addGravatarUrlToResource = (resource) => {
  const { email } = resource;

  const emailHash = md5(email.toLowerCase());
  const gravatarURL = `https://www.gravatar.com/avatar/${emailHash}`;

  return {
    ...resource,
    gravatarURL,
  };
};

export default { addGravatarUrlToResource };
