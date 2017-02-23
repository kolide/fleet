import selenium from 'selenium-webdriver';

const { By, until } = selenium;

export default (driver) => {
  const url = 'https://localhost:8080/login';
  const elements = {
    loginForm: By.css('.login-form'),
    loginSuccess: By.css('.login-success'),
    passwordInput: By.name('password'),
    submitBtn: By.css('.login-form__submit-btn'),
    usernameInput: By.name('username'),
  };

  return {
    url,
    elements,
    navigate: () => {
      return driver.navigate().to(url);
    },
    enterUsername: (username) => {
      return driver.findElement(elements.usernameInput).sendKeys(username);
    },
    enterPassword: (password) => {
      return driver.findElement(elements.passwordInput).sendKeys(password);
    },
    getSuccessText: () => {
      const form = driver.findElement(elements.loginForm);
      driver.wait(until.elementIsNotVisible(form));

      return driver.findElement(elements.loginSuccess).getText();
    },
    submit: () => {
      return driver.findElement(elements.submitBtn).click();
    },
  };
};
