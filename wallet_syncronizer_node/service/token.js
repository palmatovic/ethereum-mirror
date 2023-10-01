module.exports = (sequelize,logger) => {
    const get = async (tokenId) => {
        let token
        try {
            const Token = require('../database/token')(sequelize)
            token = await Token.findByPk(tokenId);
        } catch (error) {
            logger.error('Internal server error:', error);
        }
        return token
    };

    return {
        get,
    };
};