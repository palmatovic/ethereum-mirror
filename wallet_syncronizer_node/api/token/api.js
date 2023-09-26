const GetApi = require('../token/get/get');

class Token {
    constructor(sequelize) {
        this.sequelize = sequelize;
    }

    async get(req, res) {
        try {
            const result = await new GetApi(this.sequelize, req.params.tokenId).get();
            res.status(result.statusCode).json(result.body);
        } catch (error) {
            console.error('error trying GetApi.get():', error);
            res.status(500).json({ error: 'internal server error' });
        }
    }
}

module.exports = Token;
