const { GoPlus} = require('goplus-sdk-js');

module.exports = (logger) => {
    const scamCheck = async (tokenAddress) => {
            try {
                const chainId = '1';
                const contractAddresses = [tokenAddress];

                const data = await GoPlus.tokenSecurity(chainId, [contractAddresses], 30);

                if (data.payload.code !== 'SUCCESS') {
                    logger.error(`ScamCheck failed with code: ${data.payload.code}`);
                    return null;
                }

                const tokenResponse = data.payload.result[tokenAddress];

                if (tokenResponse) {
                    return tokenResponse;
                } else {
                    logger.error('Result does not contain token address');
                    return null;
                }
            } catch (error) {
                logger.error('ScamCheck failed:', error);
                return null;
            }

    };

    return {
        scamCheck,
    };
};

