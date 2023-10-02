const axios = require('axios');

module.exports = (sequelize, logger) => {
    const baseURL = 'https://eth-mainnet.g.alchemy.com/v2/';

    const getTokenMetadata = async (apiKey, contractAddress) => {
        try {
            const data = {
                jsonrpc: '2.0',
                method: 'alchemy_getTokenMetadata',
                headers: {
                    'Content-Type': 'application/json',
                },
                params: [contractAddress],
                id: 1,
            };

            const response = await axios.post(`${baseURL}${apiKey}`, data);
            if (response.status !== 200) {
                logger.error(`Received unexpected response status code: ${response.status}`);
                return null;
            }

            return response.data;
        } catch (error) {
            logger.error('Cannot connect to server:', error.message);
            return null;
        }
    };

    const getTokenBalances = async (wallet, apiKey) => {
        try {
            const data = {
                jsonrpc: '2.0',
                method: 'alchemy_getTokenBalances',
                headers: {
                    'Content-Type': 'application/json',
                },
                params: [wallet, 'erc20'],
                id: 42,
            };

            const response = await axios.post(`${baseURL}${apiKey}`, data);
            if (response.status !== 200) {
                logger.error(`Received unexpected response status code: ${response.status}`);
                return null;
            }

            return response.data;
        } catch (error) {
            logger.error('Cannot connect to server:', error.message);
            return null;
        }
    };


    return {
        getTokenMetadata,
        getTokenBalances,
    };
};
