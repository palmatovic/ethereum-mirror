const express = require('express');
const { v4: uuidv4 } = require('uuid');
const { Sequelize } = require('sequelize');

const app = express();
const port = process.env.REST_PORT || 3000;

const sequelize = new Sequelize({
    dialect: 'sqlite',
    storage: './wallet_synchronize.db',
    logging: false,
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
        console.info(`Rest server listening on port ${port}`);
    });
};

const initializeDatabase = async () => {
    await sequelize.sync({ alter: true });
    console.info('Database synchronized.');
};

initializeDatabase()
    .then(() => startRestServer())
    .catch((err) => {
        console.error('Error during initialization:', err);
        process.exit(1);
    });
