const utilsHelper = require('../utils/helper');
const utilsConstant = require('../utils/constant');
const { chromium } = require('playwright');
const log = require('log');

async function accessZoombucksLoginPage(page, elementCheckAttemptNumber, elementCheckAttemptInterval) {
  await page.goto(utilsConstant.ZOOMBUCKS_URL);
  // log('Wait page loaded');
  await page.waitForTimeout(5000);
  // log('Page loaded');

  let loginButton = await utilsHelper.waitAndGetElement(page, utilsConstant.ZOOMBUCKS_LOGIN_BUTTON_XPATH, elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!loginButton) {
    log('Could not get login button');
    throw new Error('Could not get login button');
  }

  await loginButton.click();
  log('Login button clicked');
}

async function loginZoombucks(page, zoombucksUsername, zoombucksPassword, elementCheckAttemptNumber, elementCheckAttemptInterval) {
  let inputEmail = await utilsHelper.waitAndGetElement(page, utilsConstant.ZOOMBUCKS_EMAIL_INPUT_XPATH, elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!inputEmail) {
    log('Could not get email input');
    throw new Error('Could not get email input');
  }
  await inputEmail.fill(zoombucksUsername);

  let inputPassword = await utilsHelper.waitAndGetElement(page, utilsConstant.ZOOMBUCKS_PASSWORD_INPUT_XPATH, elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!inputPassword) {
    log('Could not get password input');
    throw new Error('Could not get password input');
  }
  await inputPassword.fill(zoombucksPassword);

  let signInButton = await utilsHelper.waitAndGetElement(page, utilsConstant.ZOOMBUCKS_SIGNIN_BUTTON_XPATH, elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!signInButton) {
    log('Could not get signin button');
    throw new Error('Could not get signin button');
  }
  await signInButton.click();

  // log('Wait page loaded');
  await page.waitForTimeout(5000);
  // log('Page loaded');
}


async function selectZoombucksService(page, pb, elementCheckAttemptNumber, elementCheckAttemptInterval, headless) {
    if (!headless) {
      try {
        await utilsHelper.removeSlideown(page, elementCheckAttemptNumber, elementCheckAttemptInterval);
      } catch (err) {
        console.log(`Could not click on cancel onesignal slidedown button: ${err}`);
        throw err;
      }
    }
  
    let watchSection;
    try {
      watchSection = await utilsHelper.waitAndGetElement(page, "//ul[contains(@class, 'nav-main')]", elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get watch section: ${err}`);
      throw err;
    }
  
    let watchLink;
    try {
      watchLink = await utilsHelper.waitAndGetElement(watchSection, '//span[text()="Watch"]', elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get watch link: ${err}`);
      throw err;
    }
  
    try {
      await watchLink.click();
    } catch (err) {
      console.log(`Could not click on watch link: ${err}`);
      throw err;
    }
  
    let innerCardDiv;
    try {
      innerCardDiv = await utilsHelper.waitAndGetElement(page, "//div[contains(@class, 'card-inner')]//*[contains(text(), 'Hideout')]/parent::div/parent::div/parent::div/parent::div/parent::div", elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get hideout inner card div: ${err}`);
      throw err;
    }
  
    let hideoutLink;
    try {
      hideoutLink = await utilsHelper.waitAndGetElement(innerCardDiv, '//a[text()="Start"]', elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get hideout link: ${err}`);
      throw err;
    }
  
    let href;
    try {
      href = await hideoutLink.getAttribute('href');
    } catch (err) {
      console.log(`Could not get hideout link href: ${err}`);
      throw err;
    }
  
    try {
      await page.goto(href);
    } catch (err) {
      console.log(`Could not go to hideout page: ${err}`);
      throw err;
    }
  
    // console.log('Wait page loaded');
    await page.waitForTimeout(5000);
    // console.log('Page loaded');
  }
  
  module.exports = {
    accessZoombucksLoginPage,
    loginZoombucks,
    selectZoombucksService,
  };
