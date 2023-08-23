const playwright = require('playwright');
const serviceLootTv = require("../service/loottv")
const utilsHelper = require('../utils/helper');
const utilsConstant = require('../utils/constant');

class LootTvEnv {
  constructor(
    playwright
    , playwrightBrowser
    , elementCheckAttemptNumber
    , elementCheckAttemptInterval
    , headless
    , lootTvUsername
    , lootTvPassword
    , videoMaxDuration
    , videoMinDuration
    , browserWidth
    , browserHeight
    ) {
    this.Playwright = playwright;
    this.PlaywrightBrowser = playwrightBrowser;
    this.ElementCheckAttemptNumber = elementCheckAttemptNumber;
    this.ElementCheckAttemptInterval = elementCheckAttemptInterval;
    this.Headless = headless;
    this.LootTvUsername = lootTvUsername;
    this.LootTvPassword = lootTvPassword;
    this.VideoMaxDuration = videoMaxDuration;
    this.VideoMinDuration = videoMinDuration;
    this.BrowserWidth = browserWidth;
    this.BrowserHeight = browserHeight;
  }

  async startLootTvProcess() {
    let page = null;
    
    try {
      let context = null;
      context = await this.PlaywrightBrowser.newContext({
        viewport: {
          width: this.BrowserWidth, 
          height: this.BrowserHeight,
        }
      });
      console.log("Open new page.");
      page = await context.newPage();
      
      console.log("Go to loottv page.");
      await serviceLootTv.accessLootTvLoginPage(page, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
      console.log("Login to loottv page.");
      await serviceLootTv.loginLootTv(page, this.LootTvUsername, this.LootTvPassword, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);      
      
      let videoItemXpath = "xpath=//div[@class='styles_recommendedVideosWrapper__5t_Oq']//a[contains(@class, 'VideoInfiniteLoader_videoItem__PqM6S')]"
      let videoItemDurationXpath = "//div[@class='VideoItem_videoLengthWrapper__bjMkx']"
      await serviceLootTv.getListAndClickVideo(page, videoItemXpath, videoItemDurationXpath, null, this.VideoMaxDuration, this.VideoMinDuration, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
      
      // aspetta che il video finisca
      await utilsHelper.waitPageIsChanged(page);

      // ora entra un loop in cui prendi i video laterali
      let videoItemSideXpath = "xpath=//div[@class='RecommendedVideo_recommendedVideoWrapper__j98_7']"
      
      let videoItemDurationSideXpath = "//div[@class='RecommendedVideo_videoLengthWrapper__K5vuC']"
      let videoImgXpath = "//img"
      while(true){
        await serviceLootTv.getListAndClickVideo(page, videoItemSideXpath, videoItemDurationSideXpath, videoImgXpath, this.VideoMaxDuration, this.VideoMinDuration, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        // aspetta che il video finisca
        await utilsHelper.waitPageIsChanged(page);
      }      
    }
    catch (err) {
      console.log(err);
      return err;
    }

    return err;
  }
}

module.exports = LootTvEnv;