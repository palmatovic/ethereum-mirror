require('dotenv').config();
// const playwright = require('playwright');
const playwright = require('playwright');
const utilsConstant = require('./utils/constant');
const HideoutEnv = require('./crawler/hideout');
const LootTvEnv = require('./crawler/loottv');

// (async () => {
//   // Avvia il browser Chromium
//   const browser = await chromium.launch();

//   // Crea una nuova pagina
//   const page = await browser.newPage();

//   // Naviga a una determinata URL
//   await page.goto('https://example.com');

//   // Esegui azioni sulla pagina
//   await page.fill('input[name="username"]', 'esempio@dominio.com');
//   await page.fill('input[name="password"]', 'passwordsegreta');
//   await page.click('button[type="submit"]');

//   // Attendi che la pagina si carichi completamente
//   await page.waitForLoadState('networkidle');

//   // Esegui altre azioni sulla pagina

//   // Chiudi il browser
//   await browser.close();
// })();


async function main() {

  let headlessEnabled 
  if (process.env.HEADLESSMODEENABLED.toLocaleLowerCase === "true") {
    headlessEnabled = true;
  } else {
    headlessEnabled = false;
  }

  const chromeChannel = 'chrome';
  const browser = await playwright.chromium.launch({
    headless: headlessEnabled,
    channel: chromeChannel,
  });

  

  try {
    let rewardPlatformUsername="";
    let rewardPlatformPassword="";
    if(process.env.REWARDPLATFORM === utilsConstant.ZOOMBUCKS_REWARDS){
      rewardPlatformUsername=process.env.ZOOMBUCKSUSERNAME;
      rewardPlatformPassword=process.env.ZOOMBUCKSPASSWORD;
    } else if (process.env.REWARDPLATFORM === utilsConstant.FREE_CRYPTO_REWARDS){
      rewardPlatformUsername=process.env.FREECRYPTOREWARDSUSERNAME;
      rewardPlatformPassword=process.env.FREECRYPTOREWARDSPASSWORD;
    }

    switch(process.env.PROCESSENABLED) {
      

      case utilsConstant.PROCESS_HIDEOUT_VIDEO:
        const hideoutEnv = new HideoutEnv(
          playwright,
          browser,
          parseInt(process.env.ELEMENTCHECKATTEMPTNUMBER),
          parseInt(process.env.ELEMENTCHECKINTERVAL),
          headlessEnabled,
          rewardPlatformUsername,
          rewardPlatformPassword,
          process.env.HIDEOUTUSERNAME,
          process.env.HIDEOUTPASSWORD,
          parseInt(process.env.VIDEOPERPAGE),
          parseInt(process.env.VIDEOMAXDURATION),
          parseInt(process.env.VIDEOMINDURATION),
          parseInt(process.env.BROWSERWIDTH),
          parseInt(process.env.BROWSERHEIGHT),
          process.env.REWARDPLATFORM,
        );
    
        await hideoutEnv.startHideoutProcess();
        break;
      case utilsConstant.PROCESS_LOOTTV_VIDEO:
        const lootTvEnv = new LootTvEnv(
          playwright,
          browser,
          parseInt(process.env.ELEMENTCHECKATTEMPTNUMBER),
          parseInt(process.env.ELEMENTCHECKINTERVAL),
          headlessEnabled,
          process.env.LOOTTVUSERNAME,
          process.env.LOOTTVPASSWORD,
          parseInt(process.env.VIDEOMAXDURATION),
          parseInt(process.env.VIDEOMINDURATION),
          parseInt(process.env.BROWSERWIDTH),
          parseInt(process.env.BROWSERHEIGHT),
        );
    
        await lootTvEnv.startLootTvProcess();
        break;
      default:
        // code block
    }

    
  } finally {
    // Chiudi il browser
    await browser.close();
  }
}

main().catch((error) => {
  console.error(error);
});
