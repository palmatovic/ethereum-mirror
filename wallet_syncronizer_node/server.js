// main.js

require('dotenv').config();
const express = require('express');
const { Sequelize } = require('sequelize');
const helmet = require('helmet'); // Security middleware
const winston = require('winston'); // For advanced logging

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
    } catch (error) {
        logger.error('Unable to connect to the database:', error);
        process.exit(1);
    }

    // Initialize API with sequelize and logger
    const tokenApi = require('./token_api')(sequelize, logger);

    // Middleware
    app.use(express.json());
    app.use(helmet()); // Use Helmet security middleware

    // Routes
    app.get('/tokens/:tokenId', tokenApi.get);

    const PORT = process.env.PORT || 3000;
    app.listen(PORT, () => {
        logger.info(`Server is running on port ${PORT}`);
    });
}

startServer().catch((error) => {
    logger.error('Error starting the server:', error);
});
