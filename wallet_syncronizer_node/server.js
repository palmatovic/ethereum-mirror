const express = require('express');
const { v4: uuid_v4 } = require('uuid');
const playwright = require('playwright');
const log = require('log4js').getLogger();
const { Sequelize} = require('sequelize');
const cron = require('node-cron');

const app = express();
const port = process.env.REST_PORT || 3000;

const sequelize = new Sequelize('sqlite::memory:', {
    logging: false,
});

const Token = require('./database/token')(sequelize);

const Wallet = require('./database/wallet')(sequelize);

const WalletToken = require('./database/wallet_token')(sequelize);

const WalletTransaction = require('./database/wallet_transaction')(sequelize);


app.use(express.json());
app.use(express.urlencoded({ extended: false }));

app.use((req, res, next) => {
    req.uuid = uuid_v4();
    next();
});

app.use((req, res, next) => {
    res.setTimeout(300000); // 5 minutes
    next();
});

const tokenApi = require('./api/token')(sequelize, Token);
const walletApi = require('./api/wallet')(sequelize, Wallet);
const walletTokenApi = require('./api/wallet_token')(sequelize, WalletToken);
const walletTransactionApi = require('./api/wallet_transaction')(sequelize, WalletTransaction);

app.get('/api/v1/token/get', tokenApi.get);
app.get('/api/v1/token/list', tokenApi.list);
app.get('/api/v1/wallet/get', walletApi.get);
app.get('/api/v1/wallet/list', walletApi.list);
app.post('/api/v1/wallet/create', walletApi.create);
app.put('/api/v1/wallet/update', walletApi.update);
app.delete('/api/v1/wallet/delete', walletApi.delete);
app.get('/api/v1/wallet_token/get', walletTokenApi.get);
app.get('/api/v1/wallet_token/list', walletTokenApi.list);
app.get('/api/v1/wallet_transaction/get', walletTransactionApi.get);
app.get('/api/v1/wallet_transaction/list', walletTransactionApi.list);

const startRestServer = () => {
    app.listen(port, () => {
        log.info(`Rest server listening on port ${port}`);
    });
};

const startCronJob = async () => {
    const browser = await initializeBrowser();

    cron.schedule(`*/${process.env.SCRAPE_INTERVAL_MINUTES} * * * *`, async () => {
        await syncData(browser);
    });

    log.info('Cron job started.');
};

const initializeBrowser = async () => {
    return await playwright.chromium.launch({
        headless: process.env.PLAYWRIGHT_HEADLESS === true,
    });
};

const syncData = async (browser) => {
    // Implement your data synchronization logic here
};

const initializeDatabase = async () => {
    await sequelize.sync({ force: true });
    log.info('Database synchronized.');
};

initializeDatabase()
    .then(() => startCronJob())
    .then(() => startRestServer())
    .catch((err) => {
        log.error('Error during initialization:', err);
        process.exit(1);
    });
