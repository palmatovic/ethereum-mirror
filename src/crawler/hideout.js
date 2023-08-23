const playwright = require('playwright');
const serviceZoombucks = require("../service/zoombucks")
const serviceFreeCryptoRewards = require("../service/freecryptorewards")
const serviceHideout = require("../service/hideout")
const utilsHelper = require('../utils/helper');
const utilsConstant = require('../utils/constant');

class HideoutEnv {
  constructor(
    playwright
    , playwrightBrowser
    , elementCheckAttemptNumber
    , elementCheckAttemptInterval
    , headless
    , rewardsPlatformUsername
    , rewardsPlatformPassword
    , hideoutUsername
    , hideoutPassword
    , videoPerPage
    , videoMaxDuration
    , videoMinDuration
    , browserWidth
    , browserHeight
    , rewardPlatform
    ) {
    this.Playwright = playwright;
    this.PlaywrightBrowser = playwrightBrowser;
    this.ElementCheckAttemptNumber = elementCheckAttemptNumber;
    this.ElementCheckAttemptInterval = elementCheckAttemptInterval;
    this.Headless = headless;
    this.RewardsPlatformUsername = rewardsPlatformUsername;
    this.RewardsPlatformPassword = rewardsPlatformPassword;
    this.HideoutUsername = hideoutUsername;
    this.HideoutPassword = hideoutPassword;
    this.VideoPerPage = videoPerPage;
    this.VideoMaxDuration = videoMaxDuration;
    this.VideoMinDuration = videoMinDuration;
    this.BrowserWidth = browserWidth;
    this.BrowserHeight = browserHeight;
    this.RewardPlatform = rewardPlatform;
  }

  async startHideoutProcess() {
    let page = null;
    
    try {
      let context = null;
      context = await this.PlaywrightBrowser.newContext({
        viewport: {
          width: this.BrowserWidth, 
          height: this.BrowserHeight,
        }
      })
      console.log("Open new page.")
      page = await context.newPage();

      if(this.RewardPlatform === utilsConstant.ZOOMBUCKS_REWARDS){
        console.log("Go to zoombucks page.")
        await serviceZoombucks.accessZoombucksLoginPage(page, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        console.log("Login to zoombucks page.")
        await serviceZoombucks.loginZoombucks(page, this.RewardsPlatformUsername, this.RewardsPlatformPassword, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);        
        console.log("Select Hideout partner.")
        await serviceZoombucks.selectZoombucksService(page, this.PlaywrightBrowser, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval, this.HeadlessModeEnabled);
      
      } else if (this.RewardPlatform === utilsConstant.FREE_CRYPTO_REWARDS){
        console.log("Go to free crypto rewards page.")
        await serviceFreeCryptoRewards.accessFreeCryptoRewardsLoginPage(page, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        console.log("Login to free crypto rewards page.")
        await serviceFreeCryptoRewards.loginFreeCryptoRewards(page, this.RewardsPlatformUsername, this.RewardsPlatformPassword, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);        
        console.log("Select Hideout partner.")
        await serviceFreeCryptoRewards.selectFreeCryptoRewardsService(page, this.PlaywrightBrowser, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval, this.HeadlessModeEnabled);
      }

      console.log("Go to unlogged hideout page.")
      await serviceHideout.accessHideoutLoginPage(page, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval, this.HeadlessModeEnabled);
      console.log("Login to hideout.")
      await serviceHideout.loginHideout(page, this.HideoutUsername, this.HideoutPassword, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
      await page.waitForTimeout(10000);
      try{
        await serviceHideout.checkAndAddHideoutVoucher(page, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);     
      } catch (err) {
        console.log("Cannot get rewards voucher")
        console.log(err);
      }
      
      // Prendi dal DB le ricerche da effettuare
      let searches = [
        'gaming', 'web', 'minecraft', 'mining', 'chef', 
        'amazon', 'lol',
        'crawler', 'husky', 'dog', 'plants', 'apple', 
        'castle', 'paris','london', 'rome', 'berlin', 
        'africa', 'america', 'cat','wifi', 'gardening', 
        'living', 'house'];
      
      for (const search of searches) {
        console.log(`Search for: ${search}`);
        
        await serviceHideout.searchHideoutVideos(page, search, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        // let inputSearch;
        // inputSearch = await utilsHelper.waitAndGetElement(page, "//input[@id='searchTerm']", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        
        // if (!inputSearch) {
        //   console.log(`Could not get search input`);
        //   continue;
        // }
        
        // await inputSearch.fill(search);
        // await page.waitForTimeout(2000);
        
        // const searchInputButtonXPath = "//input[@id='headerSearchSubmit']";
        // await utilsHelper.clickElement(page, searchInputButtonXPath, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
                
        await page.waitForTimeout(10000);
        
        let numberOfVideoTextDiv;
        numberOfVideoTextDiv = await utilsHelper.waitAndGetElement(page, "//div[@class='homeHeading']", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        
        if (!numberOfVideoTextDiv) {
          console.log(`Could not get div with number of videos`);
          continue;
        }
        
        let numberOfVideoText;
        numberOfVideoText = await numberOfVideoTextDiv.textContent();
        
        if (!numberOfVideoText) {
          console.log(`Could not get text with number of videos`);
          continue;
        }
        
        console.log(`Text number of videos: ${numberOfVideoText}`) 

        let re = /\d+/;
        const numberMatches = numberOfVideoText.match(re);

        let numberOfVideo = 0;
        if (numberMatches && numberMatches.length > 0) {
          numberOfVideo = parseInt(numberMatches[0]);

          if (numberOfVideo < 1 || numberOfVideo > 100000) {
            console.log("Number of videos converted but not valid");
            continue;
          }
        } else {
          console.log("No videos found");
          continue;
        }

        let startPage = 1;
        let totalNumberPages = Math.ceil(numberOfVideo / this.VideoPerPage);
        let endPage = (totalNumberPages >= 5) ? 5 : totalNumberPages;
        let indexPage = 0;
        for (indexPage = startPage; indexPage <= endPage; indexPage++){
          if(indexPage > 1){
            // quando finisci i video della pagina sei nel dettaglio del video, 
            // quindi devi tornare alla search
            await serviceHideout.searchHideoutVideos(page, search, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        
            // let inputSearch;
            // inputSearch = await utilsHelper.waitAndGetElement(page, "//input[@id='searchTerm']", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
            
            // if (!inputSearch) {
            //   console.log(`Could not get search input`);
            //   continue;
            // }
            
            // await inputSearch.fill(search);
            // await page.waitForTimeout(2000);
            
            // const searchInputButtonXPath = "//input[@id='headerSearchSubmit']";
            // await utilsHelper.clickElement(page, searchInputButtonXPath, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
                    
            await page.waitForTimeout(10000);

            // per pagine maggiori di uno prima devi cliccare sulla pagina desiderata
            console.log(`Change to page ${indexPage}`);
            let pageButtonXPath = `//div[normalize-space()="${indexPage}"]`;
            await utilsHelper.clickNestedElement(page, "//div[@class='pagination']", pageButtonXPath, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval)
            await page.waitForTimeout(10000);
          }

          console.log("Get list of video from search");
          let videoItemDivs;
          videoItemDivs = await page.$$("xpath=//div[@class='search_video_wrapper']//div[@class='videoItem']");
          
          if (!videoItemDivs) {
            console.log("Could not fetch all videos on the first page");
            continue;
          }
          
          let maxMinutes = this.VideoMaxDuration;
          let minMinutes = this.VideoMinDuration;
          let videosToSee = [];
          for (let video of videoItemDivs) {
            let durationElementDiv;
            console.log("Cycle to get speficied duration video");
            durationElementDiv = await utilsHelper.waitAndGetElement(video, "//div[@class='duration']", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
            
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
              
              let linkElement;
              linkElement = await utilsHelper.waitAndGetElement(video, "//a[@data-value]", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
              
              if (!linkElement) {
                console.log("Could not get single link element div");
                continue;
              }
              
              let href;
              href = await linkElement.getAttribute("href");
              
              if (!href) {
                console.log("Could not get single link href");
                continue;
              }

              videosToSee.push(href);
            }
          }

          await page.waitForTimeout(10000);

          for (let videoHref of videosToSee) {
            console.log("Process new video");
              
            await serviceHideout.searchHideoutVideos(page, search, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
        
            // let inputSearch;
            // inputSearch = await utilsHelper.waitAndGetElement(page, "//input[@id='searchTerm']", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
            
            // if (!inputSearch) {
            //   console.log(`Could not get search input`);
            //   continue;
            // }
            
            // await inputSearch.fill(search);
            // await page.waitForTimeout(2000);
            
            // await utilsHelper.clickElement(page, "//input[@id='headerSearchSubmit']", this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
            
            await page.waitForTimeout(10000);

            if(indexPage > 1){
              // per pagine maggiori di uno prima devi cliccare sulla pagina desiderata
              console.log(`Change to page ${indexPage}`);
              let pageButtonXPath = `//div[normalize-space()="${indexPage}"]`;
              await utilsHelper.clickNestedElement(page, "//div[@class='pagination']", pageButtonXPath, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval)
              await page.waitForTimeout(10000);
            }
            
            const videoLinkXPath = `//div[@class='search_video_wrapper']//div[@class='videoItem']//a[@href='${videoHref}']`;
            const videoLinkError = await utilsHelper.clickElement(page, videoLinkXPath, this.ElementCheckAttemptNumber, this.ElementCheckAttemptInterval);
            
            if (videoLinkError) {
              console.log(`Could not get single link element div`);
              continue;
            }
            
            await page.waitForTimeout(10000);
            
            let startPageUrl = page.url();
            let endPageUrl = page.url();
            
            while (startPageUrl === endPageUrl) {
                      
              // await page.waitForTimeout(1000);
              // endPageUrl = startPage;

              await page.waitForTimeout(20000);
              endPageUrl = page.url();
            }
            console.log(`Out of while`)
            if (startPageUrl === endPageUrl) {
              console.log(`Same video. startPage`)
              console.log(`Same video. startPageUrl ${startPageUrl}, endPageUrl ${endPageUrl}`);
            } else {
              console.log(`Video has changed. startPageUrl ${startPageUrl}, endPageUrl ${endPageUrl}`);
            }
          }
        }          
      }
    }
    catch (err) {
      console.log(err);
      return err;
    }
    return err;
  }
}

module.exports = HideoutEnv;



