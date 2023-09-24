const { DataTypes } = require('sequelize');
const Token = require('./token');
const Wallet = require('./wallet');

module.exports = (sequelize) => {
    const WalletToken = sequelize.define('WalletToken', {
        WalletId: {
            type: DataTypes.STRING,
            primaryKey: true,
        },
        TokenId: {
            type: DataTypes.STRING,
            primaryKey: true,
            allowNull: false,
        },
        TokenAmount: {
            type: DataTypes.FLOAT,
            allowNull: false,
        },
        TokenAmountHex: {
            type: DataTypes.STRING,
            allowNull: false,
        },
        CreatedAt: {
            type: DataTypes.DATE,
            defaultValue: DataTypes.NOW,
        },
        UpdatedAt: {
            type: DataTypes.DATE,
            defaultValue: DataTypes.NOW,
        },
    }, {
        tableName: 'WalletToken',
        timestamps: false, // Puoi impostare questo su true se desideri timestamps
    });

    // Definisci le relazioni con i modelli Token e Wallet
    WalletToken.belongsTo(Token(sequelize), {
        foreignKey: 'TokenId',
        as: 'Token',
    });
    WalletToken.belongsTo(Wallet(sequelize), {
        foreignKey: 'WalletId',
        as: 'Wallet',
    });

    return WalletToken;
};
