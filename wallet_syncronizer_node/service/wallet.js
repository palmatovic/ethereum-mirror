module.exports = (sequelize,logger) => {
    const wallet = require('../database/wallet')(sequelize)

    const get = async (walletId) => {
        let walletDb;
        try {
            walletDb = wallet.findByPk(walletId);
            return walletDb
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error
        }
    };

    const list = async () => {
        let walletsDb
        try {
            walletsDb = await wallet.findAll();
            return walletsDb
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error
        }
    };

    const create = async (walletId) => {
        let newWallet
        try {
            newWallet = await wallet.create({ wallet_id: walletId });
            return newWallet
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error
        }
    };

    return {
        get,
        list,
        create,
    };
};