const { Token } = require('../../../database/token');

class GetService {
    constructor(sequelize, tokenId) {
        this.sequelize = sequelize;
        this.tokenId = tokenId;
    }

    async get() {
        const token = await Token.findOne({
            where: {
                tokenId: this.tokenId,
            },
        });

        if (!token) {
            return {
                status: 404,
                token: null,
                err: new Error('Token not found'),
            };
        }

        return {
            status: 200,
            token,
            err: null,
        };
    }
}

module.exports = GetService;