// wallet.js

module.exports = (sequelize, logger) => {
    const walletService = require('../service/wallet')(sequelize);

    const get = async (req, res) => {
        const walletId = req.params.wallet_id;

        try {
            if (!walletId || typeof walletId !== 'string' || walletId.trim().length === 0) {
                return res.status(400).json({ error: 'wallet_id must be a non-empty string' });
            }

            const wallet = await walletService.get(walletId);

            if (wallet) {
                res.json(wallet);
            } else {
                res.status(404).json({ error: 'wallet not found' });
            }
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    };

    const list = async (req, res) => {
        try {
            const wallets = await walletService.list();
            res.json(wallets);
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    };

    const create = async (req, res) => {
        try {
            const walletId = req.body.wallet_id;

            if (!walletId || typeof walletId !== 'string' || walletId.trim().length === 0) {
                return res.status(400).json({ error: 'wallet_id is a mandatory field and must be a non-empty string' });
            }

            // Check if the wallet with the given ID already exists
            const existingWallet = await walletService.get(walletId);

            if (existingWallet) {
                return res.status(409).json({ error: 'wallet already exists' }); // 409 Conflict
            }

            // If the wallet doesn't exist, create it
            const wallet = await walletService.create(walletId);
            res.json(wallet);
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    };

    const update = async (req, res) => {
        let updatedWallet;
        try {
            const walletId = req.body.wallet_id;

            if (!walletId || typeof walletId !== 'string' || walletId.trim().length === 0) {
                return res.status(400).json({error: 'wallet_id is a mandatory field and must be a non-empty string'});
            }

            // Check if the wallet with the given ID already exists
            const existingWallet = await walletService.get(walletId);

            if (existingWallet) {
                updatedWallet = await walletService.update(existingWallet, req.body);
                res.json(updatedWallet);
            } else {
                res.status(404).json({error: 'wallet not found'});
            }
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({error: 'Internal server error'});
        }
    };

    const deleteWallet = async (req, res) => {
        try {
            const walletId = req.params.wallet_id;

            if (!walletId || typeof walletId !== 'string' || walletId.trim().length === 0) {
                return res.status(400).json({ error: 'wallet_id must be a non-empty string' });
            }

            // Check if the wallet with the given ID exists
            const existingWallet = await walletService.get(walletId);

            if (existingWallet) {
                await walletService.deleteWallet(existingWallet);
                res.json({ message: 'wallet deleted successfully' });
            } else {
                res.status(404).json({ error: 'wallet not found' });
            }
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    };

    return {
        create,
        update,
        deleteWallet,
        get,
        list
    };
};
