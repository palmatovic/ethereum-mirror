// token.js

module.exports = (sequelize, logger) => {
    const tokenService = require('../service/token')(sequelize);

    const get = async (req, res) => {
        const tokenId = req.params.tokenId;

        try {
            // Validate tokenId as a string with a minimum length of 1 character
            if (typeof tokenId !== 'string' || tokenId.length < 1) {
                return res.status(400).json({ error: 'TokenId must be a non-empty string' });
            }

            const token = await tokenService.get(tokenId)

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

    return {
        get,
    };
};
