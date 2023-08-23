const utilsHelper = require('../utils/helper');
const utilsConstant = require('../utils/constant');
const { chromium } = require('playwright');
const log = require('log');

async function accessLootTvLoginPage(page, elementCheckAttemptNumber, elementCheckAttemptInterval) {
  await page.goto(utilsConstant.LOOTTV_URL);
  await page.waitForTimeout(5000);

  let dialogXpath = "//div[@role='dialog']";
  let consentCuttonXpath = "//button[contains(@class, 'fc-cta-consent')]";
  await utilsHelper.clickNestedElement(page, dialogXpath, consentCuttonXpath, elementCheckAttemptNumber, elementCheckAttemptInterval);

  await page.waitForTimeout(2000);

  let headerXpath = "//div[@id='header']";
  let loginButtonXpath = "//button[text()='Log In']";
  await utilsHelper.clickNestedElement(page, headerXpath, loginButtonXpath, elementCheckAttemptNumber, elementCheckAttemptInterval);
}

async function loginLootTv(page, lootTvUsername, lootTvPassword, elementCheckAttemptNumber, elementCheckAttemptInterval) {
  let inputEmail = await utilsHelper.waitAndGetElement(page, "//input[@type='email']", elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!inputEmail) {
    log('Could not get email input');
    throw new Error('Could not get email input');
  }
  await inputEmail.fill(lootTvUsername);

  let inputPassword = await utilsHelper.waitAndGetElement(page, "//input[@type='password']", elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!inputPassword) {
    log('Could not get password input');
    throw new Error('Could not get password input');
  }
  await inputPassword.fill(lootTvPassword);

  let signInButton = await utilsHelper.waitAndGetElement(page, "//button[text()='Sign In']", elementCheckAttemptNumber, elementCheckAttemptInterval);
  if (!signInButton) {
    log('Could not get signin button');
    throw new Error('Could not get signin button');
  }
  await signInButton.click();

  // log('Wait page loaded');
  await page.waitForTimeout(5000);
  // log('Page loaded');
}

async function searchLootTvVideos(page, search, elementCheckAttemptNumber, elementCheckAttemptInterval) {

  try {
    let inputSearch;
    inputSearch = await utilsHelper.waitAndGetElement(page, "//div[@class='Topnav_inputWrapper__cIwcM']//input[@placeholder='Search']", elementCheckAttemptNumber, elementCheckAttemptInterval);
    
    if (!inputSearch) {
      console.log(`Could not get search input`);
      new Error('Element not get search element')
    }
    
    await inputSearch.fill(search);
    await page.waitForTimeout(2000);
    
    await utilsHelper.clickElement(page, "//div[@class='Topnav_searchIconWrapper__EXEXI']", elementCheckAttemptNumber, elementCheckAttemptInterval);
          
  } catch (err) {
    console.log(`Could not search : ${err}`);
    throw err;
  }
}

async function getListAndClickVideo(page, listElementXpath, listElementDurationXpath, subElementToClickXpath, maximumMinutes, minusMinutes, elementCheckAttemptNumber, elementCheckAttemptInterval) {

  let videoPageUrl = null;
  let videoFinded = false;
  try {
    let maxMinutes = maximumMinutes;
    let minMinutes = minusMinutes;

    let videoRecommendedItemDivs;
    videoRecommendedItemDivs = await page.$$(listElementXpath);       
    if (!videoRecommendedItemDivs) {
      console.log("Could not fetch all videos on the first page");
      throw new Error("");
    }

    while(!videoFinded) {
      for (let video of videoRecommendedItemDivs) {
        let durationElementDiv;
        console.log("Cycle to get speficied duration video");
        durationElementDiv = await utilsHelper.waitAndGetElement(video, listElementDurationXpath, elementCheckAttemptNumber, elementCheckAttemptInterval);
        
        if (!durationElementDiv) {
          console.log("Could not get duration div");
          continue;
        }
        
        let durationString;
        durationString = await durationElementDiv.textContent();
        
        if (!durationString) {
          console.log("Could not get duration text");
          continue;
        }
        
        let durationStringParts = durationString.split(":");
        if (durationStringParts.length === 2) {
          let minutes = parseInt(durationStringParts[0].trim());
          if (minutes >= maxMinutes || minutes < minMinutes) {
            continue;
          }

          // clicca sull'elemento
          if(subElementToClickXpath != null){
            imageElement = await utilsHelper.waitAndGetElement(video, subElementToClickXpath, elementCheckAttemptNumber, elementCheckAttemptInterval);
        
            if (!durationElementDiv) {
              console.log("Could not get duration div");
              continue;
            }
            await imageElement.click();
          } else {
            await video.click();
          }
          
          await page.waitForTimeout(10000);
          videoPageUrl = page.url();
          videoFinded = true;
          break;          
        }
      }
      if(!videoFinded){
        maxMinutes = maxMinutes + 1;
      }
    }
    return videoFinded;

      
  } catch (err) {
    console.log(`Could not search : ${err}`);
    throw err;
  }
}


  
module.exports = {
    accessLootTvLoginPage,
    loginLootTv,
    searchLootTvVideos,
    getListAndClickVideo,
};
