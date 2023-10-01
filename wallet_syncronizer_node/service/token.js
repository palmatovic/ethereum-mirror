module.exports = (sequelize,logger) => {
    const Token = require('../database/token')(sequelize)

    const get = async (tokenId) => {
        let token;
        try {
            token = Token.findByPk(tokenId);
            return token
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error
        }
    };

    const list = async () => {
        let tokens
        try {
            tokens = await Token.findAll();
            return tokens
        } catch (error) {
            logger.error('Internal server error:', error);
            throw error
        }
    };

    return {
        get,
        list,
    };
};