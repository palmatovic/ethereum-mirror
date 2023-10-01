// walletService.js

module.exports = (sequelize, logger) => {
    const wallet = require('../database/wallet')(sequelize);

    const get = async (walletId) => {
        let walletDb;
        try {
            walletDb = await wallet.findByPk(walletId);
            return walletDb;
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error;
        }
    };

    const list = async () => {
        let walletsDb;
        try {
            walletsDb = await wallet.findAll();
            return walletsDb;
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error;
        }
    };

    const create = async (walletId) => {
        let newWallet;
        try {
            newWallet = await wallet.create({ wallet_id: walletId });
            return newWallet;
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error;
        }
    };

    const update = async (walletId, newData) => {
        try {
            const walletToUpdate = await wallet.findByPk(walletId);
            if (walletToUpdate) {
                await walletToUpdate.update(newData);
                return walletToUpdate;
            }
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error;
        }
    };

    const deleteWallet = async (walletId) => {
        try {
            const walletToDelete = await wallet.findByPk(walletId);
            await walletToDelete.destroy();
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error;
        }
    };

    return {
        get,
        list,
        create,
        update,
        deleteWallet,
    };
};
