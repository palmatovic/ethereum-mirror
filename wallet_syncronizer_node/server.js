const express = require('express');
const { v4: uuidv4 } = require('uuid');
const log4js = require('log4js');
const { Sequelize } = require('sequelize');

const app = express();
const port = process.env.REST_PORT || 3000;

const sequelize = new Sequelize('sqlite::memory:', {
    logging: true,
    storage: './wallet_synchronize.db'
});

// Definizione dei modelli Sequelize e delle rotte Express qui...

app.use(express.json());
app.use(express.urlencoded({ extended: false }));

app.use((req, res, next) => {
    req.uuid = uuidv4();
    next();
});

app.use((req, res, next) => {
    res.setTimeout(300000); // 5 minutes
    next();
});

require('./database/token')(sequelize)
const TokenApi = require('./api/token/api');
const token = new TokenApi(sequelize);

app.get('/api/v1/token/get', token.get);

const startRestServer = () => {
    app.listen(port, () => {
        log4js.getLogger().info(`Rest server listening on port ${port}`);
    });
};

const initializeDatabase = async () => {
    await sequelize.sync({ alter: true });
    log4js.getLogger().info('Database synchronized.');
};

initializeDatabase()
    .then(() => startRestServer())
    .catch((err) => {
        log4js.getLogger().error('Error during initialization:', err);
        process.exit(1);
    });
