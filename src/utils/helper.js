
const log = require('log');

async function waitAndGetElement(pe, xpath, attemptNumber, attemptInterval) {
  for (let i = 1; i <= attemptNumber; i++) {
    let element, err = null;
    
    if (pe.constructor.name !== "Page" && pe.constructor.name !==  "ElementHandle") {
      throw new Error('Invalid type');
    }

    try {
      await pe.waitForSelector(xpath);
      console.log(`xpath for: ${xpath}`);
      element = await pe.$(`xpath=${xpath}`);
    } catch (error) {
      err = error;
    }

    if (err === null && element !== null) {
      log(`element found ${xpath}`);
      return element;
    } else if (err === null && element === null) {
      log(`element not found ${xpath}`);
    } else if (!err.message.includes('failed to find element matching selector')) {
      throw err;
    }

    await new Promise((resolve) => setTimeout(resolve, attemptInterval * 1000));
  }

  log(`element not found ${xpath}. End of retry`);
  throw new Error('Element not found after the retries');
}

async function removeSlideown(page, attemptNumber, attemptInterval) {
  let onesignalHandle, cancelHandle;

  try {
    onesignalHandle = await waitAndGetElement(page, '//*[@id="onesignal-slidedown-container"]', attemptNumber, attemptInterval);
  } catch (error) {
    log(`could get slidown container button: ${error}`);
    throw error;
  }

  try {
    cancelHandle = await waitAndGetElement(onesignalHandle, '//*[@id="onesignal-slidedown-cancel-button"]', attemptNumber, attemptInterval);
  } catch (error) {
    log(`could get later button: ${error}`);
    throw error;
  }

  try {
    await cancelHandle.click();
  } catch (error) {
    log(`could not click on cancel onesignal slidedown button: ${error}`);
    throw error;
  }
}

async function clickElement(pe, xpath, attemptNumber, attemptInterval) {
  for (let i = 1; i <= attemptNumber; i++) {
    let element, err = null;

    try {
      await pe.waitForSelector(xpath);
      element = await pe.$(`xpath=${xpath}`);
    } catch (error) {
      err = error;
    }

    if (pe.constructor.name !== "Page" && pe.constructor.name !== "ElementHandle") {
      throw new Error('Invalid type');
    }

    if (err === null && element !== null) {
      log(`element found ${xpath}`);

      try {
        await element.click();
      
        return;
      } catch (error) {
        log(`could not click on button: ${error}. Retry`);
      }

      log(`element clicked ${xpath}`);
    } else if (err === null && element === null) {
      log(`element not found ${xpath}`);
    } else if (!err.message.includes('failed to find element matching selector')) {
      log(`some error occurred ${err}`);
    }

    await new Promise((resolve) => setTimeout(resolve, attemptInterval * 1000));
  }

  log(`element not click on element ${xpath}. End of retry`);
  throw new Error('Element not click element after the retries');
}

async function clickNestedElement(pe, parentXpath, childXpath, attemptNumber, attemptInterval) {

  let parentElement = null;
  parentElement  = await waitAndGetElement(pe, parentXpath, attemptNumber, attemptInterval);

  if (parentElement == null) {
    log(`cannot get parent element ${parentXpath}. End of retry`);
    throw new Error('Element not click element after the retries');
  }
  await clickElement(parentElement, childXpath, attemptNumber, attemptInterval);
}

async function waitPageIsChanged(page) {
  let startPageUrl = page.url();
  let endPageUrl = page.url();
  
  while (startPageUrl === endPageUrl) {
    await page.waitForTimeout(20000);
    endPageUrl = page.url();
  }
  console.log(`Out of while`)
  console.log(`Video has changed. startPageUrl ${startPageUrl}, endPageUrl ${endPageUrl}`);
}

module.exports = {
  waitAndGetElement,
  removeSlideown,
  clickElement,
  clickNestedElement,
  waitPageIsChanged,
};