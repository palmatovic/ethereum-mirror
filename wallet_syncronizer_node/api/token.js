// token.js

module.exports = (sequelize, logger) => {
    const tokenService = require('../service/token')(sequelize);

    const get = async (req, res) => {
        const tokenId = req.params.token_id;

        try {
            if (!tokenId || typeof tokenId !== 'string' || tokenId.trim().length === 0) {
                return res.status(400).json({ error: 'token_id must be a non-empty string' });
            }

            const token = await tokenService.get(tokenId);

            if (token) {
                res.json(token);
            } else {
                res.status(404).json({ error: 'Token not found' });
            }
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    };

    const list = async (req, res) => {
        try {
            const tokens = await tokenService.list();
            res.json(tokens);
        } catch (error) {
            logger.error('Internal server error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    };

    return {
        get,
        list
    };
};
