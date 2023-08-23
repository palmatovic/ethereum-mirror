const utilsHelper = require('../utils/helper');
const utilsConstant = require('../utils/constant');
const { error } = require('log');

async function accessHideoutLoginPage(page, elementCheckAttemptNumber, elementCheckAttemptInterval, headless) {
  // console.log('Wait seconds');
  await page.waitForTimeout(5000);

  try {
    await utilsHelper.clickElement(page, "//div[@role='dialog']//button//span[contains(text(), 'AGREE')]", elementCheckAttemptNumber, elementCheckAttemptInterval);
  } catch (err) {
    console.log(`Could not click on agree cookies button: ${err}`);
    throw err;
  }

  console.log('Wait seconds');
  await page.waitForTimeout(5000);

  if (!headless) {
    try {
      await utilsHelper.removeSlideown(page, elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not click on cancel onesignal slidedown button: ${err}`);
      throw err;
    }
  }

  let loginHideoutLink;
  try {
    loginHideoutLink = await utilsHelper.waitAndGetElement(page, "//div[@class='offerWallNotice']//a[text()='Log In']", elementCheckAttemptNumber, elementCheckAttemptInterval);
  } catch (err) {
    console.log(`Could not get login hideout link: ${err}`);
    throw err;
  }

  try {
    await loginHideoutLink.click();
  } catch (err) {
    console.log(`Could not click on login hideout link: ${err}`);
    throw err;
  }

  // console.log('Wait page loaded');
  await page.waitForTimeout(10000);
  // console.log('Page loaded');
}

async function loginHideout(page, hideoutUsername, hideoutPassword, elementCheckAttemptNumber, elementCheckAttemptInterval) {
    let inputUsername;
    try {
      inputUsername = await utilsHelper.waitAndGetElement(page, "//input[@name='username']", elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get username input: ${err}`);
      throw err;
    }
  
    try {
      await inputUsername.fill(hideoutUsername);
    } catch (err) {
      console.log(`Could not fill username: ${err}`);
      throw err;
    }
  
    let inputPassword;
    try {
      inputPassword = await utilsHelper.waitAndGetElement(page, "//input[@name='password']", elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get password input: ${err}`);
      throw err;
    }
  
    try {
      await inputPassword.fill(hideoutPassword);
    } catch (err) {
      console.log(`Could not fill password: ${err}`);
      throw err;
    }
  
    let logInButton;
    try {
      logInButton = await utilsHelper.waitAndGetElement(page, "//form[@id='login']//button[text()='LOG IN']", elementCheckAttemptNumber, elementCheckAttemptInterval);
    } catch (err) {
      console.log(`Could not get button: ${err}`);
      throw err;
    }
  
    try {
      await logInButton.click();
    } catch (err) {
      console.log(`Could not click on login button: ${err}`);
      throw err;
    }
  
    // console.log('Wait page loaded');
    await page.waitForTimeout(5000);
    // console.log('Page loaded');
  
    // try {
    //   await utilsHelper.waitAndGetElement(page, "//div[@class='rewards_wrapper']", elementCheckAttemptNumber, elementCheckAttemptInterval);
    // } catch (err) {
    //   console.log(`Could not get reward wrapper: ${err}`);
    //   throw err;
    // }
  }


  async function checkAndAddHideoutVoucher(page, elementCheckAttemptNumber, elementCheckAttemptInterval) {
    // controlla se ci sono rewards
    try {
      await utilsHelper.clickElement(page, "//button[@id='rewards']", elementCheckAttemptNumber, elementCheckAttemptInterval);
      await page.waitForTimeout(5000);
      let shareRewardsModal ;
      shareRewardsModal  = await utilsHelper.waitAndGetElement(page, "//div[@class='share-rewards-modal-box']", elementCheckAttemptNumber, elementCheckAttemptInterval);
    
      // Seleziona l'input con attributo name='promo_code' all'interno del div shareRewardsModal
      let promoCodeInput = null;
      promoCodeInput  = await utilsHelper.waitAndGetElement(shareRewardsModal, "//input[@name='promo_code']", elementCheckAttemptNumber, elementCheckAttemptInterval);

      // Verifica se l'input è valorizzato
      const promoCodeValue = await promoCodeInput.evaluate((element) => element.value);
      if (promoCodeValue !== '') {
        console.log("A voucher reward to add.")
        await utilsHelper.clickElement(shareRewardsModal, "//button[contains(@class, 'enter-promo-btn')]", elementCheckAttemptNumber, elementCheckAttemptInterval);
        console.log("Voucher reward to added.")
      } else {
        console.log("No voucher reward to add.")
      }
    } catch (err) {
      console.log(`Could not check voucher reward: ${err}`);
      throw err;
    }
  }

  async function searchHideoutVideos(page, search, elementCheckAttemptNumber, elementCheckAttemptInterval) {

    try {
      let inputSearch;
      inputSearch = await utilsHelper.waitAndGetElement(page, "//input[@id='searchTerm']", elementCheckAttemptNumber, elementCheckAttemptInterval);
      
      if (!inputSearch) {
        console.log(`Could not get search input`);
        new Error('Element not get search element')
      }
      
      await inputSearch.fill(search);
      await page.waitForTimeout(2000);
      
      await utilsHelper.clickElement(page, "//input[@id='headerSearchSubmit']", elementCheckAttemptNumber, elementCheckAttemptInterval);
            
    } catch (err) {
      console.log(`Could not search : ${err}`);
      throw err;
    }
  }
  
  module.exports = {
    accessHideoutLoginPage,
    loginHideout,
    checkAndAddHideoutVoucher,
    searchHideoutVideos,
  };