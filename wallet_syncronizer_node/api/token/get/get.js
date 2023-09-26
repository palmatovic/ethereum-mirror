const GetService = require('../../../service/token/get/get');
const logger = require('log4js').getLogger();
const json = require('../../../util/json');

class GetApi {
    constructor(sequelize, tokenId) {
        this.sequelize = sequelize;
        this.tokenId = tokenId;
    }

    async get() {
        logger.info('get token started', 'token_id', this.tokenId);

        if (this.tokenId.length === 0) {
            return { statusCode: 400, body: json.NewErrorResponse(400, 'empty token_id') };
        }

        const tokenGetService = new GetService(this.sequelize, this.tokenId);
        const { httpStatus, token, err } = await tokenGetService.get();

        if (err) {
            logger.error('get token terminated with failure', 'token_id', this.tokenId, 'error', err.message);
            return { statusCode: httpStatus, body: json.NewErrorResponse(httpStatus, err.message) };
        }

        logger.error('get token terminated with failure', 'token_id', this.tokenId, 'error', err.message);

        return { statusCode: 200, body: new json.NewSuccessResponse({ data: { token } }) };
    }
}

module.exports = GetApi;
