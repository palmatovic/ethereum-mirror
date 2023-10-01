// main.js

require('dotenv').config();
const express = require('express');
const { Sequelize } = require('sequelize');
const helmet = require('helmet'); // Security middleware
const winston = require('winston');
const uuid = require("uuid"); // For advanced logging

// Logger configuration
const logger = winston.createLogger({
    level: 'info',
    format: winston.format.json(),
    transports: [
        new winston.transports.Console(),
        new winston.transports.File({ filename: 'error.log', level: 'error' }),
    ],
});

async function startServer() {
    const app = express();

    // Configure Sequelize for the database
    const sequelize = new Sequelize({
        dialect: 'sqlite',
        storage: './wallet_synchronize.db',
        logging: false,
    });

    try {
        await sequelize.authenticate();
        logger.info('Database connection established successfully.');
        await sequelize.sync({alter: true})
        logger.info('Database synchronized successfully')
    } catch (error) {
        logger.error('Unable to setup database:', error);
        process.exit(1);
    }

    // Initialize API with sequelize and logger
    const tokenApi = require('./api/token')(sequelize, logger);
    const walletApi = require('./api/wallet')(sequelize, logger);

    // Middleware
    app.use(express.json());
    app.use(helmet()); // Use Helmet security middleware
    app.use((req, res, next) => {
        req.requestId = uuid.v4();
        next();
    });

    // Routes
    app.get('/token/list', tokenApi.list);
    app.get('/token/get/:token_id',tokenApi.get)
    app.get('/wallet/list', walletApi.list);
    app.get('/wallet/get/:wallet_id',walletApi.get)
    app.post('/wallet',walletApi.create)
    app.put('/wallet',walletApi.update)
    app.delete('/wallet/get/:wallet_id',walletApi.deleteWallet)





    const PORT = process.env.PORT || 3000;
    app.listen(PORT, () => {
        logger.info(`Server is running on port ${PORT}`);
    });
}

startServer().catch((error) => {
    logger.error('Error starting the server:', error);
});
