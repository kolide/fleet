const convertSeconds = (totalSeconds) => {
  const currentSeconds = totalSeconds / 1000;

  const hours = Math.floor(currentSeconds / 3600);
  const minutes = Math.floor((currentSeconds - (hours * 3600)) / 60);
  let seconds = currentSeconds - (hours * 3600) - (minutes * 60);
  seconds = Math.round((seconds * 100) / 100);

  const resultHrs = hours < 10 ? `0${hours}` : hours;
  const resultMins = minutes < 10 ? `0${minutes}` : minutes;
  const resultSecs = seconds < 10 ? `0${seconds}` : seconds;

  return `${resultHrs}:${resultMins}:${resultSecs}`;
}

export default convertSeconds;
